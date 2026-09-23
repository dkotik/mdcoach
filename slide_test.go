package mdcoach

import (
	"bytes"
	"testing"

	"github.com/dkotik/mdcoach/internal"
	"github.com/yuin/goldmark/v2/ast"
)

func TestSlideTransformer(t *testing.T) {
	t.Skip("TODO: fix later")
	document := ast.NewDocument()
	beforeHeading := ast.NewParagraph()
	firstHeading := ast.NewHeading(1, ast.HeadingKindATX)
	firstContent := ast.NewParagraph()
	secondHeading := ast.NewHeading(2, ast.HeadingKindATX)
	secondContent := ast.NewParagraph()

	document.AppendChild(beforeHeading)
	document.AppendChild(firstHeading)
	document.AppendChild(firstContent)
	document.AppendChild(secondHeading)
	document.AppendChild(secondContent)

	(&slideCutter{}).Transform(document, nil, nil)

	if document.ChildCount() != 4 {
		t.Fatalf("got %d slides, want 4", document.ChildCount())
	}
	if got := document.FirstChild(); got.Kind() != SlideKind {
		t.Fatalf("first child kind = %v, want %v", got.Kind(), SlideKind)
	}
	if got := document.FirstChild().ChildCount(); got != 1 {
		t.Fatalf("first slide has %d children, want 1", got)
	}
	if got := document.FirstChild().NextSibling().ChildCount(); got != 2 {
		t.Fatalf("second slide has %d children, want 2", got)
	}
	if got := document.LastChild().ChildCount(); got != 2 {
		t.Fatalf("last slide has %d children, want 2", got)
	}
}

func TestSlideTransformerDoesNotSplitNestedHeadings(t *testing.T) {
	t.Skip("TODO: fix later")
	document := ast.NewDocument()
	blockquote := ast.NewBlockquote()
	blockquote.AppendChild(ast.NewHeading(1, ast.HeadingKindATX))
	document.AppendChild(blockquote)

	(&slideCutter{}).Transform(document, nil, nil)

	if document.ChildCount() != 1 {
		t.Errorf("got %d slides, want 1", document.ChildCount())
		for child := range document.Children() {
			t.Log("child", child.Kind())
		}
	}
	if document.FirstChild().ChildCount() != 1 {
		t.Errorf("nested heading was not retained in the slide")
	}

	if !t.Failed() {
		return
	}
	var dump bytes.Buffer
	internal.WriteAST(&dump, document, nil)
	t.Log(dump.String())
}

func TestSlideTransformerEmptyDocument(t *testing.T) {
	document := ast.NewDocument()
	(&slideCutter{}).Transform(document, nil, nil)
	if document.ChildCount() != 0 {
		t.Fatalf("got %d children, want 0", document.ChildCount())
	}
}
