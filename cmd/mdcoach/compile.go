package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"

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

	var htmlPage bytes.Buffer
	if err := presentation.Render(ctx, &htmlPage, parsedSources, metadata); err != nil {
		return fmt.Errorf("render presentation: %w", err)
	}

	w, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("create EPUB output %q: %w", output, err)
	}
	defer w.Close()

	if err := epub.Write(w, htmlPage.Bytes(), metadata); err != nil {
		return fmt.Errorf("write EPUB: %w", err)
	}
	return nil
}
