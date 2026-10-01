package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dkotik/mdcoach/epub"
	"github.com/dkotik/mdcoach/presentation"
)

func compileMarkdownToHTML(
	ctx context.Context,
	output string,
	sources []string,
	force bool,
) (err error) {
	if err = confirmOverwrite(output, force); err != nil {
		if errors.Is(err, errSkip) {
			return nil // decided to skip file
		}
		return err
	}

	w, err := os.Create(output)
	if err != nil {
		return err
	}
	defer w.Close()

	if err := presentation.New(
		ctx,
		w,
		sources,
	); err != nil {
		return fmt.Errorf("compile presentation: %w", err)
	}
	return nil
}

func compileMarkdownToEPUB(
	ctx context.Context,
	output string,
	sources []string,
	force bool,
) (err error) {
	if err = confirmOverwrite(output, force); err != nil {
		if errors.Is(err, errSkip) {
			return nil
		}
		return err
	}

	parsedSources, metadata, err := presentation.Parse(sources)
	if err != nil {
		return fmt.Errorf("parse presentation: %w", err)
	}

	w, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("create EPUB output %q: %w", output, err)
	}
	defer w.Close()

	var renderErr error
	if len(parsedSources) == 1 {
		source := parsedSources[0]
		renderErr = epub.New(
			ctx,
			w,
			source.Source,
			source.Presentation,
			epub.WithMetadata(metadata),
			epub.WithMediaOptions(epub.MediaOptions{Path: filepath.Dir(source.Path)}),
		)
	} else {
		renderErr = epub.NewSources(ctx, w, parsedSources, epub.WithMetadata(metadata))
	}
	if renderErr != nil {
		return fmt.Errorf("render EPUB: %w", renderErr)
	}
	return nil
}
