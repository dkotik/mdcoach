package main

import (
	"bytes"
	_ "embed"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

//go:embed assets/snippets.cson
var pulsarSnippets []byte

func snippetsCmd() *cli.Command {
	return &cli.Command{
		Name:  "snippets",
		Usage: "display text editor autocompletion snippets that can accelerate presentation composition",
		Subcommands: []*cli.Command{
			{
				Name:  "pulsar",
				Usage: "snippets.cson for Pulsar or Atom",
				Action: func(_ *cli.Context) error {
					_, err := io.Copy(
						os.Stdout,
						bytes.NewReader(pulsarSnippets),
					)
					return err
				},
			},
		},
	}
}
