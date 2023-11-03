package renderer

import (
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

var reCaptureCitation = regexp.MustCompile(`\s+\(([^\)]+)\)$`)

const citationAttribute = `citeSource`

func isAsideBlockquote(n ast.Node) bool {
	if n.Kind() != ast.KindBlockquote {
		return false
	}
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
		// panic(n.Parent().Kind())
		// n.Parent().Kind() == ast.KindBlockquote &&
		// never render the tag itself of aside block quote
		return ast.WalkContinue, nil
	}

	if isAsideBlockquote(n) {
		return r.renderAside(w, source, n, entering)
	}

	if entering {
		if n.HasChildren() && n.LastChild().Kind() == ast.KindParagraph {
			lines := n.LastChild().Lines()
			if index := lines.Len() - 1; index >= 0 {
				lastLine := lines.At(index)
				m := reCaptureCitation.FindSubmatch(lastLine.Value(source))
				if len(m) == 2 {
					cutoff := len(m[0])
					// spew.Dump(lastLine.Stop)
					lastLine.Stop -= cutoff //
					lines.Set(index, lastLine)

					t := n.LastChild().LastChild().(*ast.Text)
					t.Segment.Stop -= cutoff
					n.SetAttributeString(citationAttribute, m[1])
					// n.LastChild().SetLines(lines)
					// spew.Dump(t.Text(source))
					// panic(string(n.LastChild().Text(source)))
				}
			}
		}

		if n.Attributes() != nil {
			_, _ = w.WriteString("<blockquote")
			html.RenderAttributes(w, n, html.BlockquoteAttributeFilter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString("<blockquote>\n")
		}
	} else {
		citationValue, ok := n.AttributeString(citationAttribute)
		if ok {
			if citation, ok := citationValue.([]byte); ok {
				_, _ = w.WriteString("<cite>")
				_, _ = w.Write(util.EscapeHTML(citation))
				_, _ = w.WriteString("</cite>\n")
			}
		}
		_, _ = w.WriteString("</blockquote>\n")
	}
	return ast.WalkContinue, nil
}
