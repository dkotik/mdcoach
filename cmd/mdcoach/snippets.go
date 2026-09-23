package main

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"os"

	"github.com/urfave/cli/v3"
)

//go:embed assets/snippets.cson
var pulsarSnippets []byte

func snippetsCmd() *cli.Command {
	return &cli.Command{
		Name:  "snippets",
		Usage: "display text editor autocompletion snippets that can accelerate presentation composition",
		Commands: []*cli.Command{
			{
				Name:  "pulsar",
				Usage: "snippets.cson for Pulsar or Atom",
				Action: func(_ context.Context, _ *cli.Command) error {
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
