package picture

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/OneOfOne/xxhash"
)

var reDetectURL = regexp.MustCompile(`^(https?\:\/\/|(www\.)?\w+\.\w+\/)`)

type InternetProvider struct {
	client *http.Client
	local  *LocalProvider
}

func NewInternetProvider(withOptions ...Option) (*InternetProvider, error) {
	local, err := NewLocalProvider(withOptions...)
	if err != nil {
		return nil, err
	}
	return &InternetProvider{
		client: http.DefaultClient,
		local:  local,
	}, nil
}

func (p *InternetProvider) download(ctx context.Context, URL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong response status code: %s", http.StatusText(resp.StatusCode))
	}
	defer func() {
		err = errors.Join(err, resp.Body.Close())
	}()

	b := &bytes.Buffer{}
	if _, err = io.Copy(b, resp.Body); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (p *InternetProvider) GetSourceSet(
	ctx context.Context,
	location string,
) (set []Source, err error) {
	if !reDetectURL.MatchString(location) {
		return p.local.GetSourceSet(ctx, location)
	}
	raw, err := p.download(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("unable to download %q image: %w", location, err)
	}

	h := xxhash.New64()
	if _, err = io.CopyN(h, bytes.NewReader(raw), 4096); err != nil {
		return nil, err
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))
	m, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	if err = isContextAlive(ctx); err != nil {
		return nil, err
	}
	if err = p.local.ensureDestinationPathExists(); err != nil {
		return nil, err
	}

	baseName, ext, _ := strings.Cut(path.Base(location), ".")
	original := filepath.Join(p.local.destinationPath, baseName+"-"+hash+"."+ext)

	p.local.wg.Add(1)
	go func(ctx context.Context, destination string, r io.Reader) {
		w, err := os.Create(destination)
		if err != nil {
			panic(err) // TODO: graceful error capture!
		}
		defer w.Close() // TODO: err capture.
		// TODO: check if already exists.
		if _, err = io.Copy(w, r); err != nil {
			panic(err) // TODO: graceful error capture!
		}
		p.local.wg.Done()
	}(ctx, original, bytes.NewReader(raw))

	bounds := m.Bounds().Size()
	sizings := p.local.matchSizings(bounds.X, bounds.Y)
	for _, sizing := range sizings {
		set = append(set, Source{
			Sizing: sizing,
			Location: filepath.Join(
				p.local.destinationPath,
				baseName+"-"+hash+fmt.Sprintf(`-%d-%d.jpg`, sizing.Width, sizing.Height),
			),
		})
	}
	set = append(set, Source{
		Sizing: Sizing{
			Width:   bounds.X,
			Height:  bounds.Y,
			Quality: 100,
		},
		Location: original,
	})

	p.local.wg.Add(1)
	go func(ctx context.Context) {
		err := p.local.resize(ctx, m, set)
		if err != nil {
			panic(err) // TODO: add graceful error handling.
		}
		p.local.wg.Done()
	}(ctx)
	return set, nil
}
