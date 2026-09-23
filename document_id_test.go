package mdcoach

import (
	"testing"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/text"
)

func TestASTTransformerSetsDocumentID(t *testing.T) {
	document := ast.NewDocument()
	document.AppendChild(ast.NewParagraph())

	(&documentIDInjector{}).Transform(document, text.NewReader(nil, text.IdentityDecoder), nil)

	value, ok := document.Attribute("id")
	if !ok {
		t.Fatal("document id attribute was not set")
	}
	id := value.Value(nil)
	if len(id) == 0 {
		t.Fatal("document id is empty")
	}
}
