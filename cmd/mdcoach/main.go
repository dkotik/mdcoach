/*
Package main provides command line interface to [mdcoach.Iterator].
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

func runtimeVersion() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok || buildInfo.Main.Version == "" {
		return "unknown"
	}
	return buildInfo.Main.Version
}

func isCapableOfPDF() bool {
	p, _ := exec.LookPath("weasyprint")
	return p != ""
}

func isDirectory(p string) (bool, error) {
	info, err := os.Stat(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

func main() {
	app := &cli.Command{
		Name:           "mdcoach",
		Version:        runtimeVersion(),
		Usage:          "convert markdown documents to HTML slide presentations with notes",
		DefaultCommand: "compile",
		Commands: []*cli.Command{
			reviewCmd(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "directory",
				Aliases: []string{"C"},
				Usage:   "change working directory before processing input files",
			},
			outputFlag,
			openFlag,
			overwriteFlag,
			silentFlag,
		},
		Action: func(ctx context.Context, c *cli.Command) (err error) {
			args := c.Args().Slice()
			if len(args) == 0 {
				return cli.ShowRootCommandHelp(c.Root())
			}
			args = args[1:] // the first one is the command name

			if directory := c.String("directory"); directory != "" {
				if err := os.Chdir(directory); err != nil {
					return fmt.Errorf("failed to change working directory to %q: %w", directory, err)
				}
			}

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to locate working directory: %w", err)
			}
			output := c.String("output")

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
			return g.Wait()
		},
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	err := app.Run(ctx, os.Args)
	stop()
	if err != nil {
		log.Fatal(err)
	}
}
