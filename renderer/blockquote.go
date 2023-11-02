package renderer

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

func isAsideBlockquote(n ast.Node) bool {
	return hasOnlyOneChildOfKind(n, ast.KindBlockquote)
}

func (r *Renderer) renderAside(
	w util.BufWriter,
	source []byte,
	n ast.Node,
	entering bool,
) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<aside>\n")
	} else {
		_, _ = w.WriteString("</aside>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderBlockquote(
	w util.BufWriter,
	source []byte,
	n ast.Node,
	entering bool,
) (ast.WalkStatus, error) {
	if isAsideBlockquote(n.Parent()) {
		// never render blockquote children of aside block quote
		return ast.WalkContinue, nil
	}

	if isAsideBlockquote(n) {
		return r.renderAside(w, source, n, entering)
	}

	if entering {
		if n.Attributes() != nil {
			_, _ = w.WriteString("<blockquote")
			html.RenderAttributes(w, n, html.BlockquoteAttributeFilter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString("<blockquote>\n")
		}
	} else {
		_, _ = w.WriteString("</blockquote>\n")
	}
	return ast.WalkContinue, nil
}
