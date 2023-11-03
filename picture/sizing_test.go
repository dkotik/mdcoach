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
	hash, image, err := LoadImage("./testdata/notfound.jpg")
	if err != nil {
		t.Fatal("unable to load test image:", err)
	}

	provider := newTestLocalProvider(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err = provider.resize(ctx, image, []Source{
		{
			Sizing: Sizing{
				Width:   240,
				Height:  240,
				Quality: 50,
			},
			Location: filepath.Join(provider.destinationPath, hash+"medium.jpg"),
		},
		{
			Sizing: Sizing{
				Width:   120,
				Height:  120,
				Quality: 0,
			},
			Location: filepath.Join(provider.destinationPath, hash+"small.jpg"),
		},
	}); err != nil {
		t.Fatal("could not resize the image:", err)
	}
}
