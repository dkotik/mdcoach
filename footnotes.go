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

const FootnoteIndexAttribute = "data-slide-footnote-index"

func getFootnoteIndex(node ast.Node, source []byte) (int, bool) {
	index, ok := node.Attribute(FootnoteIndexAttribute)
	if !ok {
		return 0, false
	}
	i, err := strconv.Atoi(index.Str(source))
	if err != nil {
		return 0, false
	}
	return i, true
}

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
) (_ ast.WalkStatus, err error) {
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
		index, ok := getFootnoteIndex(node, source)
		if !ok {
			index = node.Index
		}
		_, err = fmt.Fprintf(
			w,
			`<sup id="fnref-%s-%d"><a data-ref="#%s" class="footnote-ref" role="doc-noteref">%d</a></sup>`,
			label,
			node.RefIndex,
			label,
			index,
		)
		return ast.WalkSkipChildren, err
	case *footnoteast.FootnoteDefinition:
		if entering {
			label := footnoteID(node.Label.Str(source))
			if _, err = fmt.Fprintf(w, `<div class="footnote-definition" id="%s" role="doc-footnote">`, label); err != nil {
				return ast.WalkContinue, err
			}
			if index, ok := node.Attribute(FootnoteIndexAttribute); ok {
				_, err = fmt.Fprintf(w, `<a class="footnote-ref">%s</a>`, index.Value(source))
			}
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
