/*
Package main provides command line interface to [mdcoach.Iterator].
*/
package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "mdcoach",
		Usage: "convert markdown documents to HTML slide presentations with notes",
		Commands: []*cli.Command{
			snippetsCmd(),
			compileCmd(),
			demoCmd(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
