package mdcoach

import (
	"os"
	"testing"

	"github.com/dkotik/mdcoach/parser"
	"github.com/dkotik/mdcoach/picture"
	"github.com/dkotik/mdcoach/renderer"
	"github.com/yuin/goldmark/text"
)

func TestWalk(t *testing.T) {
	t.Skip("some demo features must be enabled")
	demo, err := os.ReadFile("./test/testdata/demo-presentation.md")
	if err != nil {
		t.Fatal("could not load demo presentation:", err)
	}
	p, err := parser.New()
	if err != nil {
		t.Fatal(err)
	}
	r, err := renderer.New(
		renderer.WithPictureProviderOptions(
			picture.WithSourcePath(`test/testdata`),
			picture.WithDestinationPath(t.TempDir()),
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err = Walk(p.Parse(text.NewReader(demo)), demo, r, func(slide, notes, footnotes []byte) error {
		t.Logf("slide: %s", slide)
		t.Logf("notes: %s", notes)
		t.Logf("footnotes: %s", footnotes)
		return nil
	}); err != nil {
		t.Fatal("unable to walk Markdown tree:", err)
	}
}
