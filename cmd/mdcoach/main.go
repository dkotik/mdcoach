/*
Package main provides command line interface to [mdcoach.Iterator].
*/
package main

import (
	"bytes"
	_ "embed"
	"io"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

//go:embed assets/snippets.cson
var pulsarSnippets []byte

func main() {
	app := &cli.App{
		Name:  "mdcoach",
		Usage: "convert markdown documents to HTML slide presentations with notes",
		Commands: []*cli.Command{
			{
				Name:  "snippets",
				Usage: "display text editor autocompletion snippets that can accelerate presentation composition",
				Subcommands: []*cli.Command{
					{
						Name:  "pulsar",
						Usage: "snippets.cson for Pulsar or Atom",
						Action: func(cCtx *cli.Context) error {
							_, err := io.Copy(
								os.Stdout,
								bytes.NewReader(pulsarSnippets),
							)
							return err
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
