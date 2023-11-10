package renderer

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// TODO: kull this in favor of the one in parser. All operations of this kind are parser operations.
func hasOnlyOneChildOfKind(n ast.Node, k ast.NodeKind) bool {
	if n.ChildCount() != 1 {
		return false
	}
	return n.FirstChild().Kind() == k
}

func (r *Renderer) renderParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if hasOnlyOneChildOfKind(n, ast.KindImage) {
		if entering {
			_, _ = w.WriteString("<figure>")
		} else {
			_, _ = w.WriteString("</figure>")
		}
		return ast.WalkContinue, nil
	}

	if entering {
		if n.Attributes() != nil {
			_, _ = w.WriteString("<p")
			html.RenderAttributes(w, n, html.ParagraphAttributeFilter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString("<p>")
		}
	} else {
		_, _ = w.WriteString("</p>\n")
	}
	return ast.WalkContinue, nil
}
