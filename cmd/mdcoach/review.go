package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dkotik/mdcoach/document"
	"github.com/dkotik/mdcoach/document/review"
	"github.com/dkotik/mdcoach/picture"
	"github.com/dkotik/mdcoach/renderer"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func reviewCmd() *cli.Command {
	return &cli.Command{
		Name:  "review",
		Usage: "generate a review sheet using questions found in Markdown list items",
		Flags: []cli.Flag{
			outputFlag,
			openFlag,
			overwriteFlag,
			silentFlag,
			titleFlag,
			// TODO: add -C flag.
		},
		Action: func(c *cli.Context) (err error) {
			cwd, err := os.Getwd() // TODO: should be flag -C
			if err != nil {
				return fmt.Errorf("cannot locate working directory: %w", err)
			}
			output := c.Value("output").(string)
			if filepath.IsLocal(output) {
				output = filepath.Join(cwd, output)
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

			args := c.Args().Slice()
			if len(args) == 0 {
				return errors.New("compile command requires a file path to a Markdown file")
			}

			questions, err := review.New(
				review.WithRenderer(r),
			)
			if err != nil {
				return err
			}

			// TODO: add notify context to respond to Ctrl+C signal and others.
			g, _ := errgroup.WithContext(context.TODO())
			for _, filePath := range args {
				g.Go(func() (err error) {
					// TODO: make sure file is markdown file!
					// prevent loading huge files!
					markdown, err := os.ReadFile(filepath.Join(cwd, filePath))
					if err != nil {
						return fmt.Errorf("unable to read file %q: %w", filePath, err)
					}
					return questions.AddSource(markdown)
				})
			}
			if err = g.Wait(); err != nil {
				return err
			}
			if questions.Len() == 0 {
				return errors.New("no questions were found in list items of given files")
			}
			questions.Shuffle()

			// if err = questions.Render(os.Stdout); err != nil {
			// 	return err
			// }
			// time.Now().Format(`2006-01-02`)
			pdf := isCapableOfPDF()
			switch ext := filepath.Ext(output); ext {
			case ".html":
			case ".pdf":
				if !pdf {
					return errors.New("PDF generator weasyprint is not installed")
				}
				output = strings.TrimSuffix(output, ".pdf") + ".html"
			case "": // directory
				output = filepath.Join(
					output,
					"review"+time.Now().Format(`2006-01-02`)+".html",
				)
			default:
				return fmt.Errorf("output format %q is not supported", ext)
			}
			if err = questions.RenderToFile(output, &document.Metadata{
				Title: c.Value("title").(string),
			}); err != nil {
				return err
			}

			if pdf {
				renderered := output
				output = strings.TrimSuffix(output, ".html") + ".pdf"
				if err = Exec(
					`weasyprint`, renderered, output,
					`-p`, // -p is important for ol! https://github.com/Kozea/WeasyPrint/issues/398
				); err != nil {
					return err
				}
			}
			fmt.Println(output)
			if c.IsSet("open") {
				return open.Run("file://" + output)
			}
			return nil
		},
	}
}

func uniqueOnly(set []string) (unique []string) {
	known := make(map[string]struct{})
	for _, check := range set {
		if _, ok := known[check]; !ok {
			known[check] = struct{}{}
			unique = append(unique, check)
		}
	}
	return unique
}
