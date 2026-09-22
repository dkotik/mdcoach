package mdcoach

import (
	"context"
	"os"
	"testing"

	"github.com/yuin/goldmark/v2/ast"
)

func TestLocalURLPath(t *testing.T) {
	url := newURLFromLocation(".")
	if url == nil {
		t.Fatal("expected non-nil URL")
	}
	if url.Host != "" {
		t.Fatal("expected local URL")
	}
	if url.Path != "." {
		t.Fatal("expected path to be '.'")
	}
}

func TestImageLoader(t *testing.T) {
	source, err := os.ReadFile("testdata/presentation.md")
	if err != nil {
		t.Fatal(err)
	}

	t.Chdir("testdata")
	tree, ok := NewParser().Parse(source).(*ast.Document)
	if !ok {
		t.Fatal("parser returned a non-document AST")
	}
	cache := NewImageCache()
	loader := NewImageLoader(cache, ImageConstraints{})
	if err := loader.LoadImages(context.Background(), source, tree); err != nil {
		t.Fatal(err)
	}

	if len(cache.images) < 3 {
		t.Fatalf("got %d cached images, want at least 3", len(cache.images))
	}
	for hash, image := range cache.images {
		if hash == "" || image.Hash == "" {
			t.Fatalf("cached image has no hash: key=%q image=%+v", hash, image)
		}
		if image.Width <= 0 || image.Height <= 0 {
			t.Fatalf("cached image %q has invalid size %dx%d", hash, image.Width, image.Height)
		}
	}
}
