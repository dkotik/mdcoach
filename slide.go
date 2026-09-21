package mdcoach

import (
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/text"
)

// SlideKind is the node kind for presentation slides.
var SlideKind = ast.NewNodeKind("Slide")

// Slide is a presentation slide containing a section of the document.
//
// The heading that starts a section is part of that section and is therefore
// the first child of the corresponding Slide node.
type Slide struct {
	ast.BaseBlock
}

var _ ast.Node = (*Slide)(nil)
var _ parser.ASTTransformer = (*SlideTransformer)(nil)

// Kind returns SlideKind.
func (*Slide) Kind() ast.NodeKind {
	return SlideKind
}

// Dump dumps the slide and its children.
func (s *Slide) Dump(_ []byte) *ast.NodeDump {
	return ast.NewNodeDump(s, map[string]any{
		"attributes": s.Attributes(),
		"children":   s.Children(),
	})
}

// SlideTransformer groups the document's top-level blocks into Slide nodes.
// A top-level heading starts a new slide; headings nested in other blocks do
// not affect the grouping.
type SlideTransformer struct{}

// NewSlideTransformer returns an AST transformer that groups sections into
// Slide nodes.
func NewSlideTransformer() parser.ASTTransformer {
	return &SlideTransformer{}
}

// NewSlideASTTransformer is an alias using Goldmark's AST transformer naming.
func NewSlideASTTransformer() parser.ASTTransformer {
	return NewSlideTransformer()
}

// Transform groups the document's top-level blocks into sections. Content
// before the first heading is retained in the first slide. An empty document
// remains empty.
func (*SlideTransformer) Transform(document *ast.Document, _ text.Reader, _ parser.Context) {
	if document.ChildCount() == 0 {
		return
	}

	children := make([]ast.Node, 0, document.ChildCount())
	for child := document.FirstChild(); child != nil; child = child.NextSibling() {
		children = append(children, child)
	}
	document.RemoveChildren()

	var slide *Slide
	for _, child := range children {
		if slide == nil || child.Kind() == ast.KindHeading {
			slide = &Slide{}
			document.AppendChild(slide)
		}
		slide.AppendChild(child)
	}
}
