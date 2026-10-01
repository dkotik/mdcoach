package epub

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/dkotik/mdcoach"
	"github.com/yuin/goldmark/v2/ast"
)

func TestImageLoaderLoadLocalImage(t *testing.T) {
	imageData := encodePNG(t, loaderTestImage(20, 10))
	filesystem := fstest.MapFS{
		"images/sample.png": &fstest.MapFile{Data: imageData},
	}
	cache := NewImageCache()
	loader := NewImageLoader(cache, MediaOptions{
		FS:          filesystem,
		WidthLimit:  8,
		HeightLimit: 6,
	})

	got, err := loader.LoadImage(context.Background(), "images/sample.png")
	if err != nil {
		t.Fatal(err)
	}
	if got.Location != "images/sample.png" {
		t.Errorf("image location = %q, want %q", got.Location, "images/sample.png")
	}
	if got.Width != 8 || got.Height != 4 {
		t.Errorf("image dimensions = %dx%d, want 8x4", got.Width, got.Height)
	}
	if !bytes.HasPrefix(got.Data, []byte("\x89PNG\r\n\x1a\n")) {
		t.Errorf("cached data is not raw PNG: %v", got.Data)
	}

	cached, ok := cache.Get("images/sample.png")
	if !ok || cached.Hash != got.Hash {
		t.Fatalf("cached image = %#v, want image with hash %q", cached, got.Hash)
	}
}

func TestImageLoaderLoadRemoteImage(t *testing.T) {
	imageData := encodePNG(t, loaderTestImage(3, 2))
	requests := 0
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests++
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(bytes.NewReader(imageData)),
			Header:     make(http.Header),
			Request:    request,
		}, nil
	})}

	cache := NewImageCache()
	loader := NewImageLoader(cache, MediaOptions{HTTPClient: client})
	for i := 0; i < 2; i++ {
		got, err := loader.LoadImage(context.Background(), "https://images.example/sample.png")
		if err != nil {
			t.Fatal(err)
		}
		if got.Width != 3 || got.Height != 2 {
			t.Errorf("image dimensions = %dx%d, want 3x2", got.Width, got.Height)
		}
	}
	if requests != 1 {
		t.Errorf("remote requests = %d, want 1 due to cache", requests)
	}
}

func TestImageLoaderLoadRemoteErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{
			name:       "non-success status",
			statusCode: http.StatusNotFound,
			body:       "not found",
			wantError:  "unexpected HTTP status",
		},
		{
			name:       "invalid image data",
			statusCode: http.StatusOK,
			body:       "not an image",
			wantError:  "decode remote image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: tt.statusCode,
					Status:     http.StatusText(tt.statusCode),
					Body:       io.NopCloser(strings.NewReader(tt.body)),
					Header:     make(http.Header),
					Request:    request,
				}, nil
			})}
			loader := NewImageLoader(NewImageCache(), MediaOptions{HTTPClient: client})
			_, err := loader.LoadImage(context.Background(), "https://images.example/sample.png")
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("LoadImage() error = %v, want it to contain %q", err, tt.wantError)
			}
		})
	}
}

func TestImageLoaderLoadImagesAnnotatesAST(t *testing.T) {
	source := []byte("![sample](images/sample.png)")
	tree := mdcoach.NewParser().Parse(source)
	filesystem := fstest.MapFS{
		"images/sample.png": &fstest.MapFile{Data: encodePNG(t, loaderTestImage(3, 2))},
	}
	loader := NewImageLoader(NewImageCache(), MediaOptions{FS: filesystem})
	if err := loader.LoadImages(context.Background(), source, tree); err != nil {
		t.Fatal(err)
	}

	var found bool
	if err := ast.Walk(tree, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || node.Kind() != ast.KindImage {
			return ast.WalkContinue, nil
		}
		found = true
		imageNode := node.(*ast.Image)
		for _, name := range []string{"data-hash", "data-width", "data-height", "data-aspect-ratio"} {
			var hasAttribute bool
			for _, attribute := range imageNode.Attributes() {
				if attribute.Name == name {
					hasAttribute = true
					break
				}
			}
			if !hasAttribute {
				t.Errorf("image node is missing %q attribute", name)
			}
		}
		return ast.WalkSkipChildren, nil
	}); err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("Markdown AST contains no image node")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func encodePNG(t *testing.T, source image.Image) []byte {
	t.Helper()

	var encoded bytes.Buffer
	if err := png.Encode(&encoded, source); err != nil {
		t.Fatal(err)
	}
	return encoded.Bytes()
}

func loaderTestImage(width, height int) image.Image {
	result := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			result.Set(x, y, color.White)
		}
	}
	return result
}
