package renderer

import (
	"bytes"
	"context"
	"errors"
	"strconv"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Image)
	sourceSet, err := r.pictureProvider.GetSourceSet(
		context.TODO(),
		string(n.Destination),
	)
	if err != nil {
		return ast.WalkStop, err
	}
	originalIndex := len(sourceSet) - 1
	if originalIndex < 0 {
		return ast.WalkStop, errors.New("picture provider returned an empty source set")
	}
	original := sourceSet[originalIndex]

	_, _ = w.WriteString("<picture>")
	for _, source := range sourceSet[:originalIndex] {
		_, _ = w.WriteString("<source media=\"(max-width: ")
		_, _ = w.WriteString(strconv.Itoa(source.Width))
		_, _ = w.WriteString("px) or (max-height: ")
		_, _ = w.WriteString(strconv.Itoa(source.Height))
		_, _ = w.WriteString("px)\" srcset=\"")
		_, _ = w.Write(util.EscapeHTML([]byte(source.Location)))
		_, _ = w.WriteString("\" />")
	}
	_, _ = w.WriteString("<source srcset=\"")
	_, _ = w.Write(util.EscapeHTML([]byte(original.Location)))
	_, _ = w.WriteString("\" />")

	alt := nodeToHTMLText(n, source)
	_, _ = w.WriteString("<img src=\"")
	// if r.Unsafe || !html.IsDangerousURL(n.Destination) {
	_, _ = w.Write(util.EscapeHTML(util.URLEscape([]byte(original.Location), true)))
	// }
	_, _ = w.WriteString(`" width="`)
	_, _ = w.WriteString(strconv.Itoa(original.Width))
	_, _ = w.WriteString(`" height="`)
	_, _ = w.WriteString(strconv.Itoa(original.Height))
	_, _ = w.WriteString(`" alt="`)
	_, _ = w.Write(alt)
	_ = w.WriteByte('"')
	if n.Title != nil {
		_, _ = w.WriteString(` title="`)
		// r.Writer.Write(w, n.Title)
		_, _ = w.Write(util.EscapeHTML(n.Title))
		_ = w.WriteByte('"')
	}
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ImageAttributeFilter)
	}
	_, _ = w.WriteString(" />")
	_, _ = w.WriteString("</picture>")
	if n.Title != nil {
		_, _ = w.WriteString(`<figcaption>`)
		_, _ = w.Write(util.EscapeHTML(n.Title))
		_, _ = w.WriteString(`</figcaption>`)
	}
	if len(alt) > 0 {
		_, _ = w.WriteString(`<figcaption class="source">`)
		_, _ = w.Write(util.EscapeHTML(alt))
		_, _ = w.WriteString(`</figcaption>`)
	}
	return ast.WalkSkipChildren, nil
}

func nodeToHTMLText(n ast.Node, source []byte) []byte {
	var buf bytes.Buffer
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if s, ok := c.(*ast.String); ok && s.IsCode() {
			buf.Write(s.Text(source))
		} else if !c.HasChildren() {
			buf.Write(util.EscapeHTML(c.Text(source)))
		} else {
			buf.Write(nodeToHTMLText(c, source))
		}
	}
	return buf.Bytes()
}
