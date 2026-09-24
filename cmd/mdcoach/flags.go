package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/urfave/cli/v3"
)

var (
	outputFlagValue *string
	outputFlag      = &cli.StringFlag{
		Destination: outputFlagValue,
		Name:        "output",
		Aliases:     []string{"o"},
		Value:       os.TempDir(),
		Usage:       "destination directory to safe presentation files to",
		DefaultText: os.TempDir(),
	}

	openFlagValue *bool
	openFlag      = &cli.BoolFlag{
		Destination: openFlagValue,
		Name:        "open",
		Aliases:     []string{"p"},
		Usage:       "open created files in system browser",
	}

	confirmOverwriteMutex = &sync.Mutex{}
	errSkip               = errors.New("skip file, do not overwrite")
	overwriteFlag         = &cli.BoolFlag{
		Name:    "force",
		Aliases: []string{"f"},
		Usage:   "overwrite files without requesting confirmation",
	}

	silentFlagValue *bool
	silentFlag      = &cli.BoolFlag{
		Destination: silentFlagValue,
		Name:        "silent",
		Aliases:     []string{"s"},
		Usage:       "hide all log messages unless they report errors",
	}
)

var ( // Document flags
	titleFlag = &cli.StringFlag{
		Name:    "title",
		Aliases: []string{"t"},
		Value:   "",
		Usage:   "title of the generated document",
	}
)

func confirmOverwrite(destination string, force bool) error {
	confirmOverwriteMutex.Lock()
	defer confirmOverwriteMutex.Unlock()

	stat, err := os.Stat(destination)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("check output file %q: %w", destination, err)
	}
	if stat.IsDir() {
		return fmt.Errorf("target %q cannot be overwritten because it is a directory", destination)
	}
	if force {
		return nil
	}

	if _, err := fmt.Fprintf(os.Stderr, "File %q already exists. Overwrite? [y/N] ", destination); err != nil {
		return fmt.Errorf("write overwrite prompt: %w", err)
	}
	var answer string
	if _, err := fmt.Fscan(os.Stdin, &answer); err != nil {
		return fmt.Errorf("read overwrite confirmation: %w", err)
	}
	if answer != "y" && answer != "Y" && strings.ToLower(answer) != "yes" {
		return errSkip
	}
	return nil
}
