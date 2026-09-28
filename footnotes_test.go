package mdcoach

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark/v2/ast"
	footnoteast "github.com/yuin/goldmark/v2/extension/ast"
)

func TestFootnoteAsideTransformer(t *testing.T) {
	testCases := []struct {
		name                string
		source              string
		wantAsideLabels     map[string]bool
		wantRemainingLabels map[string]bool
	}{
		{
			name: "move definitions to their referring slides",
			source: "# First slide\n\nReference[^first].\n\n***\nFirst slide notes.\n---\n" +
				"# Second slide\n\nReference[^second].\n\n***\nSecond slide notes.\n---\n" +
				"[^first]: First footnote.\n[^second]: Second footnote.\n[^unused]: Unreferenced footnote.\n",
			wantAsideLabels:     map[string]bool{"first": true, "second": true},
			wantRemainingLabels: map[string]bool{"unused": true},
		},
		{
			name: "leave definition when slide has no aside",
			source: "# Slide without aside\n\nReference[^note].\n\n---\n" +
				"[^note]: Footnote body.\n",
			wantAsideLabels:     map[string]bool{},
			wantRemainingLabels: map[string]bool{"note": true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			source := []byte(testCase.source)
			document := NewParser().Parse(source).(*ast.Document)

			foundAsideLabels := make(map[string]bool)
			foundRemainingLabels := make(map[string]bool)
			for slide := document.FirstChild(); slide != nil; slide = slide.NextSibling() {
				if slide.Kind() != SlideKind {
					continue
				}

				references := make(map[string]bool)
				var aside ast.Node
				_ = ast.Walk(slide, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
					if !entering {
						return ast.WalkContinue, nil
					}
					if node.Kind() == KindAside && aside == nil {
						aside = node
					}
					if node.Kind() == footnoteast.KindFootnoteDefinition {
						return ast.WalkSkipChildren, nil
					}
					if reference, ok := node.(*footnoteast.FootnoteReference); ok {
						references[reference.Label.Str(source)] = true
					}
					return ast.WalkContinue, nil
				})

				if aside != nil {
					for child := aside.FirstChild(); child != nil; child = child.NextSibling() {
						definition, ok := child.(*footnoteast.FootnoteDefinition)
						if !ok {
							continue
						}
						label := definition.Label.Str(source)
						if !references[label] {
							t.Errorf("aside contains footnote %q without a reference in its slide (references: %v)", label, references)
						}
						foundAsideLabels[label] = true
					}
				}
			}

			_ = ast.Walk(document, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				if !entering {
					return ast.WalkContinue, nil
				}
				definition, ok := node.(*footnoteast.FootnoteDefinition)
				if !ok {
					return ast.WalkContinue, nil
				}
				if definition.Parent() == nil || definition.Parent().Kind() != KindAside {
					foundRemainingLabels[definition.Label.Str(source)] = true
				}
				return ast.WalkSkipChildren, nil
			})

			if len(foundAsideLabels) != len(testCase.wantAsideLabels) {
				t.Errorf("footnotes in slide asides = %v, want %v", foundAsideLabels, testCase.wantAsideLabels)
			}
			for label := range testCase.wantAsideLabels {
				if !foundAsideLabels[label] {
					t.Errorf("footnote %q was not moved into its slide aside", label)
				}
			}
			if len(foundRemainingLabels) != len(testCase.wantRemainingLabels) {
				t.Errorf("remaining footnotes = %v, want %v", foundRemainingLabels, testCase.wantRemainingLabels)
			}
			for label := range testCase.wantRemainingLabels {
				if !foundRemainingLabels[label] {
					t.Errorf("footnote %q should remain outside an aside", label)
				}
			}
		})
	}
}

func TestFootnoteExtensionRendering(t *testing.T) {
	testCases := []struct {
		name   string
		source string
		want   []string
	}{
		{
			name:   "single footnote",
			source: "A reference[^note].\n\n[^note]: Footnote body.\n",
			want:   []string{`href="#fn:1"`, `id="fn:1"`, "Footnote body."},
		},
		{
			name:   "multiple footnotes",
			source: "First[^first] and second[^second].\n\n[^first]: First note.\n[^second]: Second note.\n",
			want:   []string{`href="#fn:2"`, "First note.", "Second note."},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			source := []byte(testCase.source)
			tree := NewParser().Parse(source)

			var output bytes.Buffer
			if err := NewRenderer(nil).Render(&output, source, tree); err != nil {
				t.Fatal(err)
			}

			rendered := output.String()
			for _, expected := range testCase.want {
				if !strings.Contains(rendered, expected) {
					t.Errorf("rendered HTML does not contain %q: %s", expected, rendered)
				}
			}
		})
	}
}
