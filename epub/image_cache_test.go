package epub

import (
	"bytes"
	stdImage "image"
	"image/color"
	"image/png"
	"testing"
)

func TestNewImage(t *testing.T) {
	tests := []struct {
		name       string
		location   string
		source     stdImage.Image
		wantErr    bool
		wantWidth  int
		wantHeight int
	}{
		{
			name:    "nil image",
			wantErr: true,
		},
		{
			name:       "encodes raw PNG data",
			location:   "images/sample.png",
			source:     solidImage(3, 2),
			wantWidth:  3,
			wantHeight: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			image, err := NewImage(tt.location, tt.source)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewImage() error = %v, wantErr %t", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if image.Location != tt.location {
				t.Errorf("image location = %q, want %q", image.Location, tt.location)
			}
			if image.Hash == "" {
				t.Error("image hash is empty")
			}
			if !bytes.HasPrefix(image.Data, []byte("\x89PNG\r\n\x1a\n")) {
				t.Fatalf("image data is not a raw PNG: %v", image.Data)
			}
			if image.Width != tt.wantWidth || image.Height != tt.wantHeight {
				t.Errorf("image dimensions = %dx%d, want %dx%d", image.Width, image.Height, tt.wantWidth, tt.wantHeight)
			}
			decoded, err := png.Decode(bytes.NewReader(image.Data))
			if err != nil {
				t.Fatalf("decode cached PNG data: %v", err)
			}
			if got := decoded.Bounds().Size(); got != (stdImage.Point{X: tt.wantWidth, Y: tt.wantHeight}) {
				t.Errorf("decoded dimensions = %v, want (%d, %d)", got, tt.wantWidth, tt.wantHeight)
			}
		})
	}
}

func TestImageCacheGet(t *testing.T) {
	cachedImage, err := NewImage("images/sample.png", solidImage(1, 1))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		set     bool
		wantHit bool
	}{
		{
			name: "missing location",
		},
		{
			name:    "cached location",
			set:     true,
			wantHit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache := NewImageCache()
			if tt.set {
				cache.Set(cachedImage)
			}

			got, ok := cache.Get("images/sample.png")
			if ok != tt.wantHit {
				t.Fatalf("Get() hit = %t, want %t", ok, tt.wantHit)
			}
			if tt.wantHit && (got == nil || got.Hash != cachedImage.Hash || !bytes.Equal(got.Data, cachedImage.Data)) {
				t.Errorf("Get() = %#v, want cached image with hash %q", got, cachedImage.Hash)
			}
		})
	}
}

func solidImage(width, height int) stdImage.Image {
	image := stdImage.NewRGBA(stdImage.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			image.Set(x, y, color.RGBA{R: 255, A: 255})
		}
	}
	return image
}
