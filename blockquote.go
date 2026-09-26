package mdcoach

import (
	"html"
	"io"
	"strings"
	"unicode"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	htmlrenderer "github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

var _ parser.ASTTransformer = (*blockquoteTransformer)(nil)
var _ htmlrenderer.NodeRenderer = (*blockquoteFooterRenderer)(nil)

// NewBlockquoteTransformer returns an AST transformer that moves trailing
// parenthetical citations in blockquotes into footer elements.
func NewBlockquoteTransformer() parser.ASTTransformer {
	return &blockquoteTransformer{}
}

type blockquoteTransformer struct{}

func (*blockquoteTransformer) Transform(document *ast.Document, reader text.Reader, _ parser.Context) {
	if reader == nil {
		return
	}

	source := reader.Source()
	_ = ast.Walk(document, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || node.Kind() != ast.KindBlockquote {
			return ast.WalkContinue, nil
		}

		paragraph, ok := node.LastChild().(*ast.Paragraph)
		if !ok {
			return ast.WalkContinue, nil
		}
		moveBlockquoteCitation(node, paragraph, source)
		return ast.WalkContinue, nil
	})
}

func moveBlockquoteCitation(blockquote ast.Node, paragraph *ast.Paragraph, source []byte) {
	var contents strings.Builder
	appendBlockText(&contents, paragraph, source)
	textContents := contents.String()
	if !strings.HasSuffix(textContents, ")") {
		return
	}

	opening := strings.LastIndex(textContents, "(")
	if opening < 0 {
		return
	}

	citation := strings.TrimSpace(textContents[opening+1 : len(textContents)-1])
	if citation == "" {
		return
	}

	prefix := strings.TrimRightFunc(textContents[:opening], unicode.IsSpace)
	removedLength := len(textContents) - len(prefix)
	textNodes := blockTextNodes(paragraph)
	var plainText strings.Builder
	for _, textNode := range textNodes {
		plainText.WriteString(textNode.Value.Value(source))
	}
	plainTextValue := plainText.String()
	if removedLength > len(plainTextValue) || !strings.HasSuffix(
		plainTextValue,
		textContents[len(textContents)-removedLength:],
	) {
		return
	}
	removeTrailingText(textNodes, paragraph, source, removedLength)

	footer := ast.NewRawHTML(text.NewMultiLineValueFromString(
		"<footer>"+html.EscapeString(citation)+"</footer>",
		text.IdentityDecoder,
	))
	if paragraph.ChildCount() == 0 {
		blockquote.ReplaceChild(paragraph, footer)
		return
	}
	blockquote.InsertAfter(paragraph, footer)
}

func appendBlockText(contents *strings.Builder, node ast.Node, source []byte) {
	switch node := node.(type) {
	case *ast.Text:
		contents.WriteString(node.Value.Value(source))
		return
	case *ast.CodeSpan:
		contents.WriteString(node.Value.Value(source))
		return
	case *ast.RawHTML:
		contents.WriteString(node.Value.Value(source))
		return
	}

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		appendBlockText(contents, child, source)
	}
}

func blockTextNodes(node ast.Node) []*ast.Text {
	if textNode, ok := node.(*ast.Text); ok {
		return []*ast.Text{textNode}
	}

	var nodes []*ast.Text
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		nodes = append(nodes, blockTextNodes(child)...)
	}
	return nodes
}

func removeTrailingText(nodes []*ast.Text, stop ast.Node, source []byte, length int) {
	for i := len(nodes) - 1; i >= 0 && length > 0; i-- {
		node := nodes[i]
		value := node.Value.Value(source)
		if len(value) <= length {
			length -= len(value)
			removeTextNode(node, stop)
			continue
		}

		node.Value = text.NewSingleLineValueFromString(
			value[:len(value)-length],
			text.IdentityDecoder,
		)
		length = 0
	}
}

func removeTextNode(node *ast.Text, stop ast.Node) {
	parent := node.Parent()
	parent.RemoveChild(node)
	for parent != stop && parent.ChildCount() == 0 {
		grandparent := parent.Parent()
		grandparent.RemoveChild(parent)
		parent = grandparent
	}
}

// blockquoteFooterRenderer allows only generated, owned footer HTML through the
// renderer's otherwise-safe raw HTML handling.
type blockquoteFooterRenderer struct {
	next htmlrenderer.NodeRenderer
}

func (r *blockquoteFooterRenderer) Render(
	writer io.Writer,
	source []byte,
	node ast.Node,
	entering bool,
	context renderer.Context,
) (ast.WalkStatus, error) {
	if rawHTML, ok := node.(*ast.RawHTML); ok {
		value := rawHTML.Value.Value(source)
		if rawHTML.Value.IsOwned() && strings.HasPrefix(value, "<footer>") && strings.HasSuffix(value, "</footer>") {
			w, ok := writer.(util.BufWriter)
			if !ok {
				w = util.NewErrorBufWriter(writer)
			}
			_, err := w.WriteString(value)
			return ast.WalkSkipChildren, err
		}
	}
	return r.next.Render(writer, source, node, entering, context)
}
