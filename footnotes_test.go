package mdcoach

import (
	"bytes"
	"strings"
	"testing"
)

func TestFootnoteRendererRendersDefinitionsInSlideNotes(t *testing.T) {
	source := []byte("# Slide\n\nReference[^note].\n\n***\nSpeaker notes.\n[^note]: Footnote body.\n")
	tree := NewParser().Parse(source)

	var output bytes.Buffer
	if err := NewRenderer(nil).Render(&output, source, tree); err != nil {
		t.Fatal(err)
	}

	rendered := output.String()
	for _, expected := range []string{
		`href="#fn-bm90ZQ"`,
		`<div class="footnote-definition" id="fn-bm90ZQ" role="doc-footnote">`,
		"Footnote body.",
	} {
		if !strings.Contains(rendered, expected) {
			t.Errorf("rendered HTML does not contain %q: %s", expected, rendered)
		}
	}
	if strings.Contains(rendered, `<div class="footnotes"`) {
		t.Errorf("rendered HTML contains an end-of-document footnotes list: %s", rendered)
	}
	if got := strings.Count(rendered, `class="footnote-definition"`); got != 1 {
		t.Errorf("rendered %d footnote definitions, want 1: %s", got, rendered)
	}

	notesStart := strings.Index(rendered, "<slide-notes>")
	definitionStart := strings.Index(rendered, `<div class="footnote-definition"`)
	notesEnd := strings.Index(rendered, "</slide-notes>")
	if notesStart < 0 || definitionStart < notesStart || notesEnd < definitionStart {
		t.Errorf("footnote definition was not rendered inside slide notes: %s", rendered)
	}
}
