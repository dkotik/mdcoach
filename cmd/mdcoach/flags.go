package main

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/urfave/cli/v2"
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

	overwriteFlagValue    *bool
	confirmOverwriteMutex = &sync.Mutex{}
	errSkip               = errors.New("skip file, do not overwrite")
	overwriteFlag         = &cli.BoolFlag{
		Destination: overwriteFlagValue,
		Name:        "force",
		Aliases:     []string{"f"},
		Usage:       "overwrite files without requesting confirmation",
	}

	silentFlagValue *bool
	silentFlag      = &cli.BoolFlag{
		Destination: silentFlagValue,
		Name:        "silent",
		Aliases:     []string{"s"},
		Usage:       "hide all log messages unless they report errors",
	}
)

func confirmOverwrite(destination string) error {
	confirmOverwriteMutex.Lock()
	defer confirmOverwriteMutex.Unlock()

	if overwriteFlagValue != nil && *overwriteFlagValue {
		return nil // always overwrite
	}
	stat, err := os.Stat(destination)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // conflict not possible
		}
		return err
	}

	if stat.IsDir() {
		return fmt.Errorf("target %q cannot be overwritten, because it is a directory", destination)
	}

	sync.OnceFunc(func() {
		fmt.Println("Detected file conflict. Type 'yes' or 'y' to confirm. Type 'all' to assume 'yes' answer for every other file conflict. You may also use --force command line flag to assume 'all' answer when the program runs.")
	})
	fmt.Printf("File %q already exists. Overwrite?", destination)
	var answer string
	if _, err = fmt.Scanf("%s", &answer); err != nil {
		return err
	}
	switch answer {
	case "all":
		all := true
		overwriteFlagValue = &all
		fallthrough
	case "y", "Y", "yes", "Yes":
		return nil
	default:
		return errSkip
	}

	return fmt.Errorf("file %s already exists", destination)
}
