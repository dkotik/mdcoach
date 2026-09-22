package mdcoach

import (
	"io"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/util"
)

var _ html.NodeRenderer = (*imageRenderer)(nil)

type imageRenderer struct {
}

func NewImageRenderer() html.NodeRenderer {
	return &imageRenderer{}
}

func (r *imageRenderer) Render(writer io.Writer, source []byte, node ast.Node, entering bool, rc renderer.Context) (ast.WalkStatus, error) {
	if node.Kind() != ast.KindImage {
		return ast.WalkContinue, nil
	}
	if !entering {
		return ast.WalkContinue, nil
	}
	w := writer.(util.BufWriter)
	n := node.(*ast.Image)
	_, _ = w.WriteString("<img src=\"")
	dest := n.Destination.Value(source)
	if !html.IsDangerousURL(dest) {
		_, _ = html.ContextLinkURLWriter(rc).WriteString(dest)
	}
	_, _ = w.WriteString(`" alt="`)
	renderTexts(w, source, n, rc)
	_ = w.WriteByte('"')
	if !n.Title.IsEmpty() {
		_, _ = w.WriteString(` title="`)
		_, _ = n.Title.WriteTo(html.ContextTextWriter(rc), source)
		_ = w.WriteByte('"')
	}
	if n.Attributes() != nil {
		html.RenderAttributes(w, source, n, html.ImageAttributeFilter, rc)
	}
	_, _ = w.WriteString(">")
	return ast.WalkSkipChildren, nil
}

func renderTexts(w util.BufWriter, source []byte, node ast.Node, rc renderer.Context) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if text, ok := child.(*ast.Text); ok {
			_, _ = text.Value.WriteTo(html.ContextTextWriter(rc), source)
			continue
		}
		renderTexts(w, source, child, rc)
	}
}
