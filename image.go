package mdcoach

import (
	"io"
	"strings"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/util"
)

var _ html.NodeRenderer = (*imageRenderer)(nil)

const (
	ImageCSSClass           = "mdcoachImage"
	ImageContentClassPrefix = ImageCSSClass + "Content"
)

const ImageStyle = `.` + ImageCSSClass + ` {
	min-height: 1.5em;
	min-width: 1.5em;
	background-repeat: no-repeat;
	background-position: center center;
	background-size: contain;
	aspect-ratio: attr(data-aspect-ratio type(<number>));
}

`

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
	_, _ = w.WriteString("<div src=\"")
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
		renderImageAttributes(w, source, n)
	}
	_, _ = w.WriteString("></div>")
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

func renderImageAttributes(writer io.Writer, source []byte, node ast.Node) {
	w, ok := writer.(util.BufWriter)
	if !ok {
		w = util.NewErrorBufWriter(w)
	}
	tw := &textWriter{w}
	classes := make([]string, 0, 2)
	classes = append(classes, ImageCSSClass)
	for _, attr := range node.Attributes() {
		if !html.ImageAttributeFilter.ContainsString(attr.Name) {
			if !strings.HasPrefix(attr.Name, "data-") {
				continue
			}
			if attr.Name == "data-hash" {
				classes = append(classes, ImageContentClassPrefix+attr.Value.Str(source))
				continue
			}
		}
		if attr.Name == "class" {
			classes = append(classes, attr.Value.Str(source))
			continue
		}
		_, _ = w.WriteString(" ")
		_, _ = w.WriteString(attr.Name)
		_, _ = w.WriteString(`="`)
		_, _ = attr.Value.WriteTo(tw, source)
		_ = w.WriteByte('"')
	}

	_, _ = w.WriteString(` class="`)
	_, _ = tw.WriteString(strings.Join(classes, " "))
	_ = w.WriteByte('"')
}
