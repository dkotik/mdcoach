package mdcoach

import (
	"context"
	"fmt"
	stdImage "image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"

	"github.com/nfnt/resize"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/text"
	"golang.org/x/sync/errgroup"
)

type ImageLoader struct {
	fs          fs.FS
	widthLimit  int
	heightLimit int
	quality     int
	cache       *ImageCache
}

type MediaOptions struct {
	FS          fs.FS
	Path        string
	WidthLimit  int
	HeightLimit int
	Quality     int
}

func NewImageLoader(cache *ImageCache, ic MediaOptions) *ImageLoader {
	if cache == nil {
		panic("nil cache")
	}
	if ic.WidthLimit == 0 {
		ic.WidthLimit = 800
	}
	if ic.HeightLimit == 0 {
		ic.HeightLimit = 600
	}
	if ic.Quality == 0 {
		ic.Quality = 80
	}
	if ic.FS == nil {
		ic.FS = os.DirFS(ic.Path)
	}
	return &ImageLoader{
		fs:          ic.FS,
		widthLimit:  ic.WidthLimit,
		heightLimit: ic.HeightLimit,
		quality:     ic.Quality,
		cache:       cache,
	}
}

func (l *ImageLoader) decodeImage(r io.Reader) (*imageWebp, error) {
	decoded, _, err := stdImage.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	if decoded.Bounds().Dx() > l.widthLimit || decoded.Bounds().Dy() > l.heightLimit {
		decoded = resize.Thumbnail(
			uint(l.widthLimit),
			uint(l.heightLimit),
			decoded,
			resize.Lanczos3,
		)
	}
	return newImageFromImage(decoded)
}

func (l *ImageLoader) loadImage(
	ctx context.Context,
	location string,
) (img *imageWebp, err error) {
	url := newURLFromLocation(location)
	location = url.String()
	img, ok := l.cache.Get(location)
	if ok {
		return img, nil
	}

	if url.Host == "" {
		file, err := l.fs.Open(path.Clean(location))
		if err != nil {
			return nil, fmt.Errorf("open local image %q: %w", location, err)
		}
		defer file.Close()
		img, err = l.decodeImage(file)
		if err != nil {
			return nil, fmt.Errorf("decode local image %q: %w", location, err)
		}
	} else {
		client := &http.Client{}
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, location, nil)
		if err != nil {
			return nil, fmt.Errorf("create image request for %q: %w", location, err)
		}
		response, err := client.Do(request)
		if err != nil {
			return nil, fmt.Errorf("download image %q: %w", location, err)
		}
		defer response.Body.Close()
		if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
			return nil, fmt.Errorf("download image %q: unexpected HTTP status %s", location, response.Status)
		}
		img, err = l.decodeImage(response.Body)
		if err != nil {
			return nil, fmt.Errorf("decode remote image %q: %w", location, err)
		}
	}

	img.Location = location
	l.cache.Set(img)
	return img, nil
}

// LoadImages downloads image destinations that are owned by the AST nodes.
// Use LoadImagesFromSource when the AST was parsed from Markdown source.
func (l *ImageLoader) LoadImages(ctx context.Context, source []byte, tree ast.Node) error {
	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(min(4, runtime.NumCPU()))
	if err := ast.Walk(tree, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || node.Kind() != ast.KindImage {
			return ast.WalkContinue, nil
		}
		imageNode, ok := node.(*ast.Image)
		if !ok {
			return ast.WalkContinue, nil
		}
		location := imageNode.Destination.Value(source)
		group.Go(func() error {
			image, err := l.loadImage(ctx, location)
			if err != nil {
				return err
			}
			imageNode.SetAttribute(
				"data-hash",
				text.NewMultiLineValue(image.Hash, text.IdentityDecoder),
			)
			imageNode.SetAttribute(
				"data-width",
				text.NewMultiLineValue(fmt.Sprintf("%d", image.Width), text.IdentityDecoder),
			)
			imageNode.SetAttribute(
				"data-height",
				text.NewMultiLineValue(fmt.Sprintf("%d", image.Height), text.IdentityDecoder),
			)
			imageNode.SetAttribute(
				"data-aspect-ratio",
				text.NewMultiLineValue(fmt.Sprintf("%.4f", float32(image.Width)/float32(image.Height)), text.IdentityDecoder),
			)
			return nil
		})
		return ast.WalkSkipChildren, nil
	}); err != nil {
		return err
	}

	return group.Wait()
}

func newURLFromLocation(s string) *url.URL {
	parsed, err := url.Parse(s)
	if err != nil || parsed.Host == "" {
		// s = path.Clean(root, s)
		return &url.URL{
			RawPath: s,
			Path:    s,
		}
	}
	parsed.Path = path.Clean(parsed.Path)
	return parsed
}
