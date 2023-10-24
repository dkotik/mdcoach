package main

import (
	"context"
	"errors"
	"fmt"
	"html"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dkotik/mdcoach/document"

	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func compileMarkdownToHTML(ctx context.Context, p, output string) error {
	output = filepath.Join(output, "demo.html")
	w, err := os.Create(output)
	if err != nil {
		return err
	}
	defer w.Close()

	if err = document.WriteHeader(w, &document.HeadMeta{
		Title: "demo",
	}); err != nil {
		return err
	}

	// header, err := templates.NewHeader(&templates.Header{
	// 	Title:      "demo",
	// 	StyleSheet: templates.StyleSheet,
	// 	JavascriptModulesES6: []string{
	// 		templates.Javascript,
	// 	},
	// })
	// if err != nil {
	// 	return err
	// }
	// if _, err = io.Copy(w, bytes.NewReader(header)); err != nil {
	// 	return err
	// }
	if _, err = io.Copy(w, strings.NewReader(html.EscapeString(`"wooo<ul><li>1</li><li>2</li></ul>", "wooNotes", "wooFT", "1", "1n", "1ft"`))); err != nil {
		return err
	}
	// fmt.Println("compiling:", p)
	return document.WriteFooter(w)
}

func compileCmd() *cli.Command {
	return &cli.Command{
		Name:  "compile",
		Usage: "display text editor autocompletion snippets that can accelerate presentation composition",
		// TODO: add silent flag.
		Action: func(c *cli.Context) error {
			directory, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot locate working directory: %w", err)
			}
			args := c.Args().Slice()
			if len(args) == 0 {
				return errors.New("compile command requires a file path to a Markdown file")
			}

			// TODO: add notify context to respond to Ctrl+C signal and others.
			g, ctx := errgroup.WithContext(context.TODO())
			for _, p := range args {
				p := p // golang.org/doc/faq#closures_and_goroutines
				if len(p) > 0 && p[0] != filepath.Separator {
					p = filepath.Join(directory, p)
				}
				g.Go(func() error {
					return compileMarkdownToHTML(ctx, p, directory)
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
