package mdcoach

import (
	"os"
	"testing"
)

func TestWalk(t *testing.T) {
	demo, err := os.ReadFile("./test/testdata/demo-presentation.md")
	if err != nil {
		t.Fatal("could not load demo presentation:", err)
	}
	if err = Walk(demo, func(slide, notes, footnotes []byte) error {
		t.Logf("slide: %s", slide)
		t.Logf("notes: %s", notes)
		t.Logf("footnotes: %s", footnotes)
		return nil
	}); err != nil {
		t.Fatal("unable to walk Markdown tree:", err)
	}
}
