package picture

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"io"
	"os"
	"path/filepath"
)

//go:embed testdata/notfound.jpg
var PictureNotFound []byte

type EmbeddedProvider struct {
	*LocalProvider
}

func (p *EmbeddedProvider) SaveAsset(name string, raw []byte) (_ string, err error) {
	if err = p.ensureDestinationPathExists(); err != nil {
		return "", err
	}
	target := filepath.Join(p.LocalProvider.destinationPath, name)
	w, err := os.Create(target)
	if err != nil {
		return "", err
	}
	defer func() {
		err = errors.Join(err, w.Close())
	}()

	if _, err = io.Copy(w, bytes.NewReader(raw)); err != nil {
		return "", err
	}
	return target, nil
}

func (p *EmbeddedProvider) GetSourceSet(
	ctx context.Context,
	location string,
) (set []Source, err error) {
	set, err = p.LocalProvider.GetSourceSet(ctx, location)
	if errors.Is(err, os.ErrNotExist) {
		switch location {
		case "notfound.jpg":
			target, err := p.SaveAsset(location, PictureNotFound)
			if err != nil {
				return nil, err
			}
			return p.LocalProvider.GetSourceSet(ctx, target)
		}
	}
	return
}
