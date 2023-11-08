package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
			// TODO: add -C flag.
		},
		Action: func(c *cli.Context) (err error) {
			output := c.Value("output").(string)
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
			cwd, err := os.Getwd() // TODO: should be flag -C
			if err != nil {
				return fmt.Errorf("cannot locate working directory: %w", err)
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
			output = filepath.Join(
				output,
				"review.html",
			)
			if err = questions.RenderToFile(output); err != nil {
				return err
			}

			wp, _ := exec.LookPath("weasyprint")
			if wp != "" {
				renderered := output
				if strings.HasSuffix(output, `.html`) {
					output = output[:len(output)-5] + `.pdf`
				}
				// -p is important for ol! https://github.com/Kozea/WeasyPrint/issues/398
				if err = Exec(`weasyprint`, renderered, output, `-p`); err != nil {
					return err
				}
			}
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
