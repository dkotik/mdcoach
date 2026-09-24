package main

import (
	"bufio"
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

	stat, err := os.Lstat(destination)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("check output path %q: %w", destination, err)
	}
	if stat.IsDir() {
		return fmt.Errorf("target %q cannot be overwritten because it is a directory", destination)
	}
	if stat.Mode()&os.ModeSymlink != 0 {
		target, err := os.Stat(destination)
		if err == nil && target.IsDir() {
			return fmt.Errorf("target %q cannot be overwritten because it is a directory", destination)
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("check symlink target %q: %w", destination, err)
		}
	}
	if force {
		return nil
	}

	stdinInfo, err := os.Stdin.Stat()
	if err != nil {
		return fmt.Errorf("check standard input for overwrite confirmation: %w", err)
	}
	if stdinInfo.Mode()&os.ModeCharDevice == 0 {
		return fmt.Errorf("cannot confirm overwrite of %q from non-interactive input; use --force to overwrite", destination)
	}

	if _, err := fmt.Fprintf(os.Stderr, "Overwrite %q? [y/N] ", destination); err != nil {
		return fmt.Errorf("write overwrite prompt: %w", err)
	}
	answer, err := readOverwriteAnswer()
	if err != nil {
		return fmt.Errorf("read overwrite confirmation: %w", err)
	}
	answer = strings.TrimSpace(answer)
	if !strings.EqualFold(answer, "y") && !strings.EqualFold(answer, "yes") {
		return errSkip
	}
	return nil
}

func readOverwriteAnswer() (string, error) {
	// A one-byte buffer prevents a fresh reader from consuming subsequent
	// answers when multiple output files are confirmed in sequence.
	answer, err := bufio.NewReaderSize(os.Stdin, 1).ReadString('\n')
	if err != nil {
		return "", err
	}
	return answer, nil
}
