package epub

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

// ImageLoader loads, resizes, and caches images as PNG files.
type ImageLoader struct {
	fs          fs.FS
	widthLimit  int
	heightLimit int
	httpClient  *http.Client
	cache       *ImageCache
}

// MediaOptions configures an ImageLoader's filesystem and output dimensions.
type MediaOptions struct {
	FS          fs.FS
	Path        string
	WidthLimit  int
	HeightLimit int
	HTTPClient  *http.Client
}

// NewImageLoader creates a loader that caches PNG images.
func NewImageLoader(cache *ImageCache, options MediaOptions) *ImageLoader {
	if cache == nil {
		panic("nil image cache")
	}
	if options.WidthLimit == 0 {
		options.WidthLimit = 800
	}
	if options.HeightLimit == 0 {
		options.HeightLimit = 600
	}
	if options.FS == nil {
		options.FS = os.DirFS(options.Path)
	}
	if options.HTTPClient == nil {
		options.HTTPClient = http.DefaultClient
	}
	return &ImageLoader{
		fs:          options.FS,
		widthLimit:  options.WidthLimit,
		heightLimit: options.HeightLimit,
		httpClient:  options.HTTPClient,
		cache:       cache,
	}
}

func (l *ImageLoader) decodeImage(location string, reader io.Reader) (*Image, error) {
	decoded, _, err := stdImage.Decode(reader)
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
	return NewImage(location, decoded)
}

// LoadImage loads location, returning a cached image when available.
func (l *ImageLoader) LoadImage(ctx context.Context, location string) (*Image, error) {
	url := imageURLFromLocation(location)
	location = url.String()
	if image, ok := l.cache.Get(location); ok {
		return image, nil
	}

	if url.Host == "" {
		file, err := l.fs.Open(path.Clean(location))
		if err != nil {
			return nil, fmt.Errorf("open local image %q: %w", location, err)
		}
		defer file.Close()

		image, err := l.decodeImage(location, file)
		if err != nil {
			return nil, fmt.Errorf("decode local image %q: %w", location, err)
		}
		l.cache.Set(image)
		return image, nil
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, location, nil)
	if err != nil {
		return nil, fmt.Errorf("create image request for %q: %w", location, err)
	}
	response, err := l.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("download image %q: %w", location, err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("download image %q: unexpected HTTP status %s", location, response.Status)
	}

	image, err := l.decodeImage(location, response.Body)
	if err != nil {
		return nil, fmt.Errorf("decode remote image %q: %w", location, err)
	}
	l.cache.Set(image)
	return image, nil
}

// LoadImages downloads image destinations owned by nodes in the Markdown AST.
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
			image, err := l.LoadImage(ctx, location)
			if err != nil {
				return err
			}
			imageNode.SetAttribute(
				"data-hash",
				text.NewMultiLineValueFromString(image.Hash, text.IdentityDecoder),
			)
			imageNode.SetAttribute(
				"data-width",
				text.NewMultiLineValueFromString(fmt.Sprintf("%d", image.Width), text.IdentityDecoder),
			)
			imageNode.SetAttribute(
				"data-height",
				text.NewMultiLineValueFromString(fmt.Sprintf("%d", image.Height), text.IdentityDecoder),
			)
			imageNode.SetAttribute(
				"data-aspect-ratio",
				text.NewMultiLineValueFromString(fmt.Sprintf("%.4f", float32(image.Width)/float32(image.Height)), text.IdentityDecoder),
			)
			return nil
		})
		return ast.WalkSkipChildren, nil
	}); err != nil {
		return err
	}
	return group.Wait()
}

func imageURLFromLocation(location string) *url.URL {
	parsed, err := url.Parse(location)
	if err != nil || parsed.Host == "" {
		return &url.URL{
			RawPath: location,
			Path:    location,
		}
	}
	parsed.Path = path.Clean(parsed.Path)
	return parsed
}
