package mdcoach

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark/v2/ast"
	footnoteast "github.com/yuin/goldmark/v2/extension/ast"
)

func TestFootnoteIndicesFollowReferenceOrder(t *testing.T) {
	source := []byte("# Slide\n\nSecond[^second], then first[^first], second again[^second].\n\n[^first]: First body.\n[^second]: Second body.\n")
	tree := NewParser().Parse(source)

	wantIndices := map[string]string{"first": "2", "second": "1"}
	definitionIndices := make(map[string]string)
	referenceIndices := make(map[string][]string)
	if err := ast.Walk(tree, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		var label string
		switch node := node.(type) {
		case *footnoteast.FootnoteDefinition:
			label = node.Label.Value(source)
			definitionIndices[label] = slideFootnoteIndex(node, source)
		case *footnoteast.FootnoteReference:
			label = node.Label.Value(source)
			referenceIndices[label] = append(referenceIndices[label], slideFootnoteIndex(node, source))
		}
		return ast.WalkContinue, nil
	}); err != nil {
		t.Fatal(err)
	}

	for label, want := range wantIndices {
		if got := definitionIndices[label]; got != want {
			t.Errorf("definition %q index = %q, want %q", label, got, want)
		}
		for _, got := range referenceIndices[label] {
			if got != want {
				t.Errorf("reference %q index = %q, want %q", label, got, want)
			}
		}
	}
	if got := len(referenceIndices["second"]); got != 2 {
		t.Errorf("found %d references to second footnote, want 2", got)
	}
}

func slideFootnoteIndex(node ast.Node, source []byte) string {
	for _, attribute := range node.Attributes() {
		if attribute.Name == FootnoteIndexAttribute {
			return attribute.Value.Str(source)
		}
	}
	return ""
}

func TestFootnoteRendererRendersDefinitionsInSlideNotes(t *testing.T) {
	source := []byte("# Slide\n\nReference[^note].\n\n***\nSpeaker notes.\n[^note]: Footnote body.\n")
	tree := NewParser().Parse(source)

	var output bytes.Buffer
	if err := NewRenderer(nil).Render(&output, source, tree); err != nil {
		t.Fatal(err)
	}

	rendered := output.String()
	for _, expected := range []string{
		`<a data-ref="#fn-bm90ZQ" class="footnote-ref" role="doc-noteref">1</a>`,
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

	notesStart := strings.Index(rendered, "<slide-notes")
	definitionStart := strings.Index(rendered, `<div class="footnote-definition"`)
	notesEnd := strings.Index(rendered, "</slide-notes>")
	if notesStart < 0 || definitionStart < notesStart || notesEnd < definitionStart {
		t.Errorf("footnote definition was not rendered inside slide notes: %s", rendered)
	}
}
