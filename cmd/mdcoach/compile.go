package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/dkotik/mdcoach/epub"
	"github.com/dkotik/mdcoach/presentation"
	presentationindex "github.com/dkotik/mdcoach/presentation/index"
)

var presentationIndexUpdateMutex sync.Mutex

func compileMarkdownToHTML(
	ctx context.Context,
	output string,
	sources []string,
	force bool,
	addToIndex bool,
) (err error) {
	if err = confirmOverwrite(output, force); err != nil {
		if errors.Is(err, errSkip) {
			return nil // decided to skip file
		}
		return err
	}

	var parsedSources []presentation.Source
	var metadata presentation.Frontmatter
	if addToIndex {
		parsedSources, metadata, err = presentation.Parse(sources)
		if err != nil {
			return fmt.Errorf("parse presentation: %w", err)
		}
	}

	w, err := os.Create(output)
	if err != nil {
		return err
	}

	var renderErr error
	if addToIndex {
		renderErr = presentation.Render(ctx, w, parsedSources, metadata)
	} else {
		renderErr = presentation.New(ctx, w, sources)
	}
	closeErr := w.Close()
	if renderErr != nil {
		return fmt.Errorf("compile presentation: %w", renderErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close presentation output %q: %w", output, closeErr)
	}
	if addToIndex {
		if err := addPresentationToIndex(output, metadata); err != nil {
			return err
		}
	}
	return nil
}

func compileMarkdownToEPUB(
	ctx context.Context,
	output string,
	sources []string,
	force bool,
	addToIndex bool,
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
	closeErr := w.Close()
	if renderErr != nil {
		return fmt.Errorf("render EPUB: %w", renderErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close EPUB output %q: %w", output, closeErr)
	}
	if addToIndex {
		if err := addPresentationToIndex(output, metadata); err != nil {
			return err
		}
	}
	return nil
}

func addPresentationToIndex(output string, metadata presentation.Frontmatter) error {
	indexPath := filepath.Join(filepath.Dir(output), "index.html")
	if filepath.Clean(output) == filepath.Clean(indexPath) {
		return fmt.Errorf("presentation output %q conflicts with its index path", output)
	}

	presentationIndexUpdateMutex.Lock()
	defer presentationIndexUpdateMutex.Unlock()

	if _, err := os.Stat(indexPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("check presentation index %q: %w", indexPath, err)
		}
		if err := presentationindex.New(indexPath); err != nil {
			return fmt.Errorf("create presentation index %q: %w", indexPath, err)
		}
	}
	if err := presentationindex.Add(indexPath, output, metadata); err != nil {
		return fmt.Errorf("add presentation %q to index: %w", output, err)
	}
	return nil
}
