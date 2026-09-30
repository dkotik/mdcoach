package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dkotik/mdcoach/presentation"
	"github.com/skratchdot/open-golang/open"

	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

func compileMarkdownToHTML(
	ctx context.Context,
	output string,
	sources []string,
	force bool,
) (err error) {
	if err = confirmOverwrite(output, force); err != nil {
		if errors.Is(err, errSkip) {
			return nil // decided to skip file
		}
		return err
	}

	w, err := os.Create(output)
	if err != nil {
		return err
	}
	defer w.Close()

	if err := presentation.New(
		ctx,
		w,
		sources,
	); err != nil {
		return fmt.Errorf("compile presentation: %w", err)
	}
	return nil
}

func compileMarkdownToEPUB(
	ctx context.Context,
	output string,
	sources []string,
	force bool,
) (err error) {
	if err = confirmOverwrite(output, force); err != nil {
		if errors.Is(err, errSkip) {
			return nil
		}
		return err
	}

	var htmlPage bytes.Buffer
	if err := presentation.New(ctx, &htmlPage, sources); err != nil {
		return fmt.Errorf("compile presentation: %w", err)
	}

	w, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("create EPUB output %q: %w", output, err)
	}
	defer w.Close()

	if err := writeEPUB(w, htmlPage.Bytes()); err != nil {
		return fmt.Errorf("write EPUB: %w", err)
	}
	return nil
}

func compileCmd() *cli.Command {
	return &cli.Command{
		Name:  "compile",
		Usage: "convert Markdown to an HTML presentation or EPUB book",
		Flags: []cli.Flag{
			outputFlag,
			openFlag,
			overwriteFlag,
			silentFlag,
		},
		Action: func(ctx context.Context, c *cli.Command) (err error) {
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
			output := c.String("output")
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

				if filepath.Ext(output) == ".epub" {
					err = compileMarkdownToEPUB(
						ctx,
						output,
						args,
						c.Bool("force"),
					)
				} else {
					if !strings.HasSuffix(output, ".html") {
						output = output + ".html"
					}
					err = compileMarkdownToHTML(
						// TODO: add notify context to respond to Ctrl+C signal and others.
						ctx,
						output,
						args,
						c.Bool("force"),
					)
				}
				if err != nil {
					return err
				}
				if c.IsSet("open") {
					return open.Run("file://" + output)
				}
				return nil
			}

			// TODO: add notify context to respond to Ctrl+C signal and others.
			g, ctx := errgroup.WithContext(ctx)
			for _, p := range args {
				p := p // golang.org/doc/faq#closures_and_goroutines
				// if len(p) > 0 && p[0] != filepath.Separator {
				if filepath.IsLocal(p) {
					p = filepath.Join(cwd, p)
				}
				g.Go(func() (err error) {
					destination := filepath.Join(output, strings.TrimSuffix(filepath.Base(p), ".md")+".html")
					if err = compileMarkdownToHTML(
						ctx,
						destination,
						[]string{p},
						c.Bool("force"),
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
