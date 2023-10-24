package main

import (
	"bytes"
	_ "embed"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

//go:embed assets/demo.md
var demoPresentation []byte

func demoCmd() *cli.Command {
	return &cli.Command{
		Name:  "demo",
		Usage: "print a demo file",
		Action: func(_ *cli.Context) error {
			// TODO: add convert demo file to an HTML presentation.
			_, err := io.Copy(
				os.Stdout,
				bytes.NewReader(demoPresentation),
			)
			return err
		},
	}
}
