package mdcoach

import (
	"bytes"
	"os"
	"testing"

	"github.com/dkotik/mdcoach/internal"
	"github.com/sebdah/goldie/v2"
)

func TestPresentationAST(t *testing.T) {
	source, err := os.ReadFile("testdata/presentation-1.md")
	if err != nil {
		t.Fatal(err)
	}

	tree := NewParser().Parse(source)
	var dump bytes.Buffer
	internal.WriteAST(&dump, tree, source)

	for child := range tree.Children() {
		switch child.Kind() {
		case SlideKind:
			for _, attr := range child.Attributes() {
				switch attr.Name {
				case "class":
					t.Error("slide contains class attribute:", attr.Value.Str(source))
				}
			}
		}
	}

	goldie.New(t).Assert(t, "presentation_ast", dump.Bytes())
}
