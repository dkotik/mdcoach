package mdcoach

import (
	"io"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

// KindFigure is the node kind for image-only figures.
var KindFigure = ast.NewNodeKind("Figure")

// Figure is a figure containing a single image.
type Figure struct {
	ast.BaseBlock
}

var _ ast.Node = (*Figure)(nil)
var _ html.NodeRenderer = (*figureRenderer)(nil)
var _ parser.ASTTransformer = (*figureTransformer)(nil)

// Kind returns FigureKind.
func (*Figure) Kind() ast.NodeKind {
	return KindFigure
}

// Dump dumps the figure and its children.
func (f *Figure) Dump(_ []byte) *ast.NodeDump {
	return ast.NewNodeDump(f, nil)
}

type figureRenderer struct{}

// NewFigureRenderer returns a Goldmark v2 HTML renderer for Figure nodes.
func NewFigureRenderer() html.NodeRenderer {
	return &figureRenderer{}
}

func (*figureRenderer) Render(
	writer io.Writer,
	source []byte,
	node ast.Node,
	entering bool,
	rc renderer.Context,
) (ast.WalkStatus, error) {
	w := writer.(util.BufWriter)
	if entering {
		_, _ = w.WriteString("<figure>")
		return ast.WalkContinue, nil
	}

	firstChild := node.FirstChild()
	if firstChild != nil {
		_, _ = w.WriteString("<figcaption>")
		renderFigureText(w, source, firstChild, rc)
		_, _ = w.WriteString("</figcaption>")
	}
	_, _ = w.WriteString("</figure>")
	return ast.WalkContinue, nil
}

func renderFigureText(w util.BufWriter, source []byte, node ast.Node, rc renderer.Context) {
	if textNode, ok := node.(*ast.Text); ok {
		_, _ = textNode.Value.WriteTo(html.ContextTextWriter(rc), source)
		return
	}
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		renderFigureText(w, source, child, rc)
	}
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
