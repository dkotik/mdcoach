package mdcoach

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	stdImage "image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/OneOfOne/xxhash"
	"github.com/nfnt/resize"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/text"
	"golang.org/x/sync/errgroup"
)

type imageJPG struct {
	Hash       string
	DataBase64 []byte
	Width      int
	Height     int
}

func newImageFromImage(i stdImage.Image) (*imageJPG, error) {
	if i == nil {
		return nil, fmt.Errorf("encode image as JPEG: nil image")
	}

	var encoded bytes.Buffer
	if err := jpeg.Encode(&encoded, i, nil); err != nil {
		return nil, fmt.Errorf("encode image as JPEG: %w", err)
	}

	encodedBytes := encoded.Bytes()
	hash := xxhash.New64()
	if _, err := hash.Write(encodedBytes); err != nil {
		return nil, fmt.Errorf("hash encoded image: %w", err)
	}

	bounds := i.Bounds()
	return &imageJPG{
		Hash:       fmt.Sprintf("%x", hash.Sum(nil)),
		DataBase64: []byte(base64.StdEncoding.EncodeToString(encodedBytes)),
		Width:      bounds.Dx(),
		Height:     bounds.Dy(),
	}, nil
}

type ImageLoader struct {
	widthLimit  int
	heightLimit int
	quality     int

	mu     *sync.Mutex
	images map[string]*imageJPG
}

type ImageConstraints struct {
	WidthLimit  int
	HeightLimit int
	Quality     int
}

func NewImageLoader(ic ImageConstraints) *ImageLoader {
	if ic.WidthLimit == 0 {
		ic.WidthLimit = 800
	}
	if ic.HeightLimit == 0 {
		ic.HeightLimit = 600
	}
	if ic.Quality == 0 {
		ic.Quality = 80
	}
	return &ImageLoader{
		widthLimit:  ic.WidthLimit,
		heightLimit: ic.HeightLimit,
		quality:     ic.Quality,
		mu:          &sync.Mutex{},
		images:      make(map[string]*imageJPG),
	}
}

func (l *ImageLoader) decodeImage(r io.Reader) (*imageJPG, error) {
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
) (img *imageJPG, err error) {
	url := newURLFromLocation(location)
	location = url.String()
	l.mu.Lock()
	img, ok := l.images[location]
	if ok {
		l.mu.Unlock()
		return img, nil
	}
	l.mu.Unlock()

	if url.Host == "" {
		file, err := os.Open(location)
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

	l.mu.Lock()
	l.images[location] = img
	l.mu.Unlock()
	return img, nil
}

// LoadImages downloads image destinations that are owned by the AST nodes.
// Use LoadImagesFromSource when the AST was parsed from Markdown source.
func (l *ImageLoader) LoadImages(ctx context.Context, source []byte, tree *ast.Document) error {
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
