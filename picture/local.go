package picture

import (
	"context"
	"errors"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nfnt/resize"
)

type LocalProvider struct {
	sourcePath             string
	destinationPath        string
	destinationPathCreated bool
	encoder                Encoder
	sizings                []Sizing
	wg                     *sync.WaitGroup
}

func (p *LocalProvider) encode(to string, m image.Image, quality int) error {
	// TODO: confirm overwrite.
	w, err := os.Create(to)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, w.Close())
	}()
	return p.encoder.EncodeImage(w, m, quality)
}

func (p *LocalProvider) resize(
	ctx context.Context,
	m image.Image,
	targets []Source,
) (err error) {
	for _, target := range targets {
		if err = isContextAlive(ctx); err != nil {
			return err
		}
		m = resize.Thumbnail(
			uint(target.Width),
			uint(target.Height),
			m, resize.Lanczos3,
		)
		if err = p.encode(target.Location, m, target.Quality); err != nil {
			return err
		}
	}
	return nil
}

func (p *LocalProvider) matchSizings(w, h int) (smaller []Sizing) {
	for _, sizing := range p.sizings {
		if sizing.Height < h || sizing.Width < w {
			smaller = append(smaller, sizing)
		}
	}
	return
}

func (p *LocalProvider) ensureDestinationPathExists() (err error) {
	if !p.destinationPathCreated {
		if err = os.MkdirAll(p.destinationPath, 0700); err != nil {
			return fmt.Errorf("unable to create directory %q: %w", p.destinationPath, err)
		}
		p.destinationPathCreated = true
	}
	return nil
}

func (p *LocalProvider) GetSourceSet(
	ctx context.Context,
	location string,
) (set []Source, err error) {
	if location == "" {
		return nil, errors.New("empty picture reference")
	}
	if location[0] != filepath.Separator {
		location = filepath.Join(p.sourcePath, location)
	}
	if err = isContextAlive(ctx); err != nil {
		return nil, err
	}
	hash, m, err := LoadImage(location)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	if err = isContextAlive(ctx); err != nil {
		return nil, err
	}
	if err = p.ensureDestinationPathExists(); err != nil {
		return nil, err
	}

	baseName, ext, _ := strings.Cut(filepath.Base(location), ".")
	original := filepath.Join(p.destinationPath, baseName+"-"+hash+"."+ext)

	p.wg.Add(1)
	go func(ctx context.Context, destination, source string) {
		// TODO: check if already exists.
		if err := copyFile(ctx, destination, source); err != nil {
			panic(err) // TODO: graceful error capture!
		}
		p.wg.Done()
	}(ctx, original, location)

	bounds := m.Bounds().Size()
	sizings := p.matchSizings(bounds.X, bounds.Y)
	for _, sizing := range sizings {
		set = append(set, Source{
			Sizing: sizing,
			Location: filepath.Join(
				p.destinationPath,
				baseName+"-"+hash+fmt.Sprintf(`-%d-%d.jpg`, sizing.Width, sizing.Height),
			),
		})
	}

	p.wg.Add(1)
	go func(ctx context.Context, set []Source) {
		err := p.resize(ctx, m, set)
		if err != nil {
			panic(err) // TODO: add graceful error handling.
		}
		p.wg.Done()
	}(ctx, set[:]) // copy set to prevent overwriting original

	return append(set, Source{
		Sizing: Sizing{
			Width:   bounds.X,
			Height:  bounds.Y,
			Quality: 100,
		},
		Location: original,
	}), nil
}

func (p *LocalProvider) FinishScaling() {
	p.wg.Wait()
}

func NewLocalProvider(withOptions ...Option) (p *LocalProvider, err error) {
	o := &options{}
	for _, option := range append(
		withOptions,
		withDefaultSourcePath(),
		withDefaultDestinationPath(),
		withDefaultEncoder(),
		withDefaultSizing(),
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("unable to initialize local picture provider: %w", err)
		}
	}

	SortSizingsInDescendingOrder(o.sizings)
	return &LocalProvider{
		sourcePath:      o.sourcePath,
		destinationPath: o.destinationPath,
		encoder:         o.encoder,
		sizings:         o.sizings,
		wg:              &sync.WaitGroup{},
	}, nil
}
