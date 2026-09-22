package mdcoach

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/sebdah/goldie/v2"
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

	tree, ok := NewParser().Parse(source).(*ast.Document)
	if !ok {
		t.Fatal("parser returned a non-document AST")
	}
	cache := NewImageCache()
	loader := NewImageLoader(cache, MediaOptions{Path: "testdata"})
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

	var css bytes.Buffer
	_, _ = css.WriteString("<!DOCTYPE html>\n<html>\n<body>\n")
	_, _ = css.WriteString("<style>\n")
	_, _ = css.WriteString(ImageStyle)
	_, _ = css.WriteString("</style>\n")

	if err := renderImagesOnly(&css, source, tree); err != nil {
		t.Fatal(err)
	}
	if err := cache.WriteImageDataCSS(&css); err != nil {
		t.Fatal(err)
	}
	_, _ = css.WriteString("</body>\n</html>")
	goldie.New(t).Assert(t, "presentation_icss", css.Bytes())
}
