package mdcoach

import (
	"encoding/base64"
	"fmt"
	"io"
	"strconv"

	"github.com/yuin/goldmark/v2/ast"
	footnoteast "github.com/yuin/goldmark/v2/extension/ast"
	"github.com/yuin/goldmark/v2/renderer"
	htmlrenderer "github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/util"
)

var _ htmlrenderer.NodeRenderer = (*footnoteRenderer)(nil)

type footnoteRenderer struct{}

// NewFootnoteRenderer returns an HTML renderer that renders footnotes at their
// position in the AST instead of collecting them at the end of the document.
func NewFootnoteRenderer() htmlrenderer.NodeRenderer {
	return &footnoteRenderer{}
}

func (*footnoteRenderer) Render(
	writer io.Writer,
	source []byte,
	node ast.Node,
	entering bool,
	_ renderer.Context,
) (ast.WalkStatus, error) {
	w, ok := writer.(util.BufWriter)
	if !ok {
		w = util.NewErrorBufWriter(writer)
	}

	switch node := node.(type) {
	case *footnoteast.FootnoteReference:
		if !entering {
			return ast.WalkSkipChildren, nil
		}
		label := footnoteID(node.Label.Str(source))
		_, err := fmt.Fprintf(
			w,
			`<sup id="fnref-%s-%d"><a href="#%s" class="footnote-ref" role="doc-noteref">%s</a></sup>`,
			label,
			node.RefIndex,
			label,
			strconv.Itoa(node.Index),
		)
		return ast.WalkSkipChildren, err
	case *footnoteast.FootnoteDefinition:
		label := footnoteID(node.Label.Str(source))
		if entering {
			_, err := fmt.Fprintf(w, `<div class="footnote-definition" id="%s" role="doc-footnote">`, label)
			return ast.WalkContinue, err
		}
		_, err := io.WriteString(w, `</div>`)
		return ast.WalkContinue, err
	default:
		return ast.WalkContinue, nil
	}
}

func footnoteID(label string) string {
	return "fn-" + base64.RawURLEncoding.EncodeToString([]byte(label))
}
