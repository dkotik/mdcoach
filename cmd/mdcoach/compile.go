package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dkotik/mdcoach"
	"github.com/dkotik/mdcoach/document"
	mdcParser "github.com/dkotik/mdcoach/parser"
	"github.com/dkotik/mdcoach/picture"
	"github.com/dkotik/mdcoach/renderer"
	"github.com/skratchdot/open-golang/open"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func compileMarkdownToHTML(
	ctx context.Context,
	output string,
	markdownContent []byte,
) (err error) {
	if err = confirmOverwrite(output); err != nil {
		if errors.Is(err, errSkip) {
			return nil // decided to skip file
		}
		return err
	}

	prsr, err := mdcParser.New()
	if err != nil {
		return err
	}
	pc := parser.NewContext()
	tree := prsr.Parse(
		text.NewReader(markdownContent),
		parser.WithContext(pc))
	meta, err := document.NewMetadata(pc)
	if err != nil {
		return fmt.Errorf("cannot accept document metadata: %w", err)
	}

	pictureProvider, err := picture.NewInternetProvider(
		picture.WithDestinationPath(filepath.Join(
			filepath.Dir(output),
			"presentationMedia",
		)),
	)
	if err != nil {
		return err
	}

	r, err := renderer.New(
		renderer.WithPictureProvider(&picture.SourceFilter{
			Provider: pictureProvider,
			IsAllowed: func(source *picture.Source) (bool, error) {
				// trim output path from the source set
				source.Location = strings.TrimPrefix(source.Location, filepath.Dir(output)+"/")
				return true, nil
			},
		}),
	)
	if err != nil {
		return err
	}

	w, err := os.Create(output)
	if err != nil {
		return err
	}
	defer w.Close()

	if err = document.WriteHeader(w, meta); err != nil {
		return err
	}

	if err = mdcoach.Compile(w, tree, markdownContent, r); err != nil {
		return err
	}
	if err = document.WriteFooter(w); err != nil {
		return err
	}
	pictureProvider.FinishScaling()
	return nil
}

func compileCmd() *cli.Command {
	return &cli.Command{
		Name:  "compile",
		Usage: "convert Markdown to an HTML presentation",
		Flags: []cli.Flag{
			outputFlag,
			openFlag,
			overwriteFlag,
			silentFlag,
		},
		Action: func(c *cli.Context) (err error) {
			// TODO: use c.IsSet("open") instead of output value!
			// if outputFlagValue == nil {
			// 	return errors.New("output flag is required")
			// }
			// output := *outputFlagValue
			// TODO: add silent flag.
			cwd, err := os.Getwd() // TODO: should be flag -C
			if err != nil {
				return fmt.Errorf("cannot locate working directory: %w", err)
			}
			output := c.Value("output").(string)
			args := c.Args().Slice()
			if len(args) == 0 {
				return errors.New("compile command requires a file path to at least one Markdown file")
			}

			isDir, err := isDirectory(output)
			if err != nil {
				return err
			}
			if !isDir {
				for i, p := range args {
					if filepath.IsLocal(p) {
						args[i] = filepath.Join(cwd, p)
					}
				}
				markdownContent, err := mdcParser.Glue(args...)
				if err != nil {
					return err
				}
				if !strings.HasSuffix(output, ".html") {
					output = output + ".html"
				}
				if err = compileMarkdownToHTML(
					// TODO: add notify context to respond to Ctrl+C signal and others.
					context.TODO(),
					output,
					markdownContent,
				); err != nil {
					return err
				}
				if c.IsSet("open") {
					return open.Run("file://" + output)
				}
				return nil
			}

			// TODO: add notify context to respond to Ctrl+C signal and others.
			g, ctx := errgroup.WithContext(context.TODO())
			for _, p := range args {
				p := p // golang.org/doc/faq#closures_and_goroutines
				// if len(p) > 0 && p[0] != filepath.Separator {
				if filepath.IsLocal(p) {
					p = filepath.Join(cwd, p)
				}
				g.Go(func() (err error) {
					markdownContent, err := os.ReadFile(p)
					if err != nil {
						return fmt.Errorf("cannot read file %q: %w", p, err)
					}
					destination := filepath.Join(output, strings.TrimSuffix(filepath.Base(p), ".md")+".html")
					if err = compileMarkdownToHTML(
						ctx,
						destination,
						markdownContent,
					); err != nil {
						return err
					}
					if c.IsSet("open") {
						return open.Run("file://" + destination)
					}
					return nil
				})
			}

			// searches := []Search{Web, Image, Video}
			// results := make([]Result, len(searches))
			// for i, search := range searches {
			// 	i, search := i, search // https://golang.org/doc/faq#closures_and_goroutines
			// 	g.Go(func() error {
			// 		result, err := search(ctx, query)
			// 		if err == nil {
			// 			results[i] = result
			// 		}
			// 		return err
			// 	})
			// }
			// if err := g.Wait(); err != nil {
			// 	return nil, err
			// }
			return g.Wait()
		},
	}
}
