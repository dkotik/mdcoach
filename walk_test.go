package mdcoach

import (
	"os"
	"testing"

	"github.com/dkotik/mdcoach/parser"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func TestWalk(t *testing.T) {
	demo, err := os.ReadFile("./test/testdata/demo-presentation.md")
	if err != nil {
		t.Fatal("could not load demo presentation:", err)
	}
	p, err := parser.New()
	if err != nil {
		t.Fatal(err)
	}
	if err = Walk(p.Parse(text.NewReader(demo)), demo, goldmark.DefaultRenderer(), func(slide, notes, footnotes []byte) error {
		t.Logf("slide: %s", slide)
		t.Logf("notes: %s", notes)
		t.Logf("footnotes: %s", footnotes)
		return nil
	}); err != nil {
		t.Fatal("unable to walk Markdown tree:", err)
	}
}
