package renderer

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderNotesBreak(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString("<hr class=\"notes\"")
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ThematicAttributeFilter)
	}
	_, _ = w.WriteString(" />\n")
	return ast.WalkContinue, nil
}
