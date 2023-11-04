package test

import (
	_ "embed"
	"testing"

	"github.com/dkotik/mdcoach/picture"
	"github.com/dkotik/mdcoach/renderer"
	gr "github.com/yuin/goldmark/renderer"
)

//go:embed testdata/demo-presentation.md
var testPresentation []byte

func NewTestRenderer(t *testing.T) gr.Renderer {
	r, err := renderer.New(
		renderer.WithPictureProviderOptions(
			picture.WithDestinationPath(t.TempDir()),
		),
	)
	if err != nil {
		t.Fatal("cannot create test renderer:", err)
	}
	return r
}
