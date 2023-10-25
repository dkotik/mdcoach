package picture

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestResize(t *testing.T) {
	if testing.Short() {
		t.Skip("image processing takes a lot of resources")
	}
	image, err := LoadImage("./testdata/notfound.jpg")
	if err != nil {
		t.Fatal("unable to load test image:", err)
	}
	r := &Resizer{
		encoder: &JPEGEncoder{},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	directory := t.TempDir()
	if err = r.Execute(ctx, &SizingRequest{
		Image: image,
		Outputs: []SizingTarget{
			{
				Sizing: Sizing{
					Width:   240,
					Height:  240,
					Quality: 50,
				},
				Destination: filepath.Join(directory, "medium"),
			},
			{
				Sizing: Sizing{
					Width:   120,
					Height:  120,
					Quality: 0,
				},
				Destination: filepath.Join(directory, "smaller"),
			},
		},
	}); err != nil {
		t.Fatal("could not resize the image:", err)
	}
}
