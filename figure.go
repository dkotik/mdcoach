package mdcoach

import (
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/text"
)

// FigureKind is the node kind for image-only figures.
var FigureKind = ast.NewNodeKind("Figure")

// Figure is a figure containing a single image.
type Figure struct {
	ast.BaseBlock
}

var _ ast.Node = (*Figure)(nil)
var _ parser.ASTTransformer = (*figureTransformer)(nil)

// Kind returns FigureKind.
func (*Figure) Kind() ast.NodeKind {
	return FigureKind
}

// Dump dumps the figure and its children.
func (f *Figure) Dump(_ []byte) *ast.NodeDump {
	return ast.NewNodeDump(f, nil)
}

// figureTransformer replaces top-level paragraphs containing only one image
// with Figure nodes.
type figureTransformer struct{}

// NewFigureTransformer returns an AST transformer that turns image-only
// document children into Figure nodes.
func NewFigureTransformer() parser.ASTTransformer {
	return &figureTransformer{}
}

// Transform replaces only direct Document children that are paragraphs with
// exactly one image child.
func (*figureTransformer) Transform(document *ast.Document, _ text.Reader, _ parser.Context) {
	for child := document.FirstChild(); child != nil; {
		next := child.NextSibling()
		paragraph, ok := child.(*ast.Paragraph)
		if !ok || paragraph.ChildCount() != 1 || paragraph.FirstChild().Kind() != ast.KindImage {
			child = next
			continue
		}

		image := paragraph.FirstChild()
		paragraph.RemoveChild(image)
		figure := &Figure{}
		figure.AppendChild(image)
		document.ReplaceChild(paragraph, figure)
		child = next
	}
}
