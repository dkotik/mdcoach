/*
Package main provides command line interface to [mdcoach.Iterator].
*/
package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func isCapableOfPDF() bool {
	p, _ := exec.LookPath("weasyprint")
	return p != ""
}

func isDirectory(p string) (bool, error) {
	info, err := os.Stat(p)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func main() {
	app := &cli.App{
		Name:  "mdcoach",
		Usage: "convert markdown documents to HTML slide presentations with notes",
		Commands: []*cli.Command{
			snippetsCmd(),
			compileCmd(),
			reviewCmd(),
			demoCmd(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
