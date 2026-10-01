package epub

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	stdhtml "html"
	stdImage "image"
	"io"
	"io/fs"
	"path"
	"strconv"
	"strings"

	"github.com/dkotik/mdcoach"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/extension"
	footnoteast "github.com/yuin/goldmark/v2/extension/ast"
	"github.com/yuin/goldmark/v2/renderer"
	htmlrenderer "github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/util"
)

type epubNodeRenderer struct {
	cache   *ImageCache
	emoteFS fs.FS
}

func newEPUBHTMLRenderer(cache *ImageCache, emoteFS fs.FS) htmlrenderer.Renderer {
	nodes := &epubNodeRenderer{cache: cache, emoteFS: emoteFS}
	return htmlrenderer.New(
		htmlrenderer.WithXHTML(),
		htmlrenderer.WithExtensions(
			extension.NewStrikethroughHTMLRenderer(),
			extension.NewTableHTMLRenderer(),
		),
		htmlrenderer.WithNodeRendererDecorator(ast.KindParagraph, func(next htmlrenderer.NodeRenderer) htmlrenderer.NodeRenderer {
			return &taskParagraphRenderer{next: next}
		}),
		htmlrenderer.WithNodeRenderer(mdcoach.SlideKind, nodes),
		htmlrenderer.WithNodeRenderer(mdcoach.KindFigure, nodes),
		htmlrenderer.WithNodeRenderer(mdcoach.KindSlideNotes, nodes),
		htmlrenderer.WithNodeRenderer(mdcoach.KindEmote, nodes),
		htmlrenderer.WithNodeRendererDecorator(ast.KindImage, func(htmlrenderer.NodeRenderer) htmlrenderer.NodeRenderer {
			return nodes
		}),
		htmlrenderer.WithNodeRendererDecorator(ast.KindRawHTML, func(htmlrenderer.NodeRenderer) htmlrenderer.NodeRenderer {
			return nodes
		}),
		htmlrenderer.WithNodeRendererDecorator(ast.KindHTMLBlock, func(htmlrenderer.NodeRenderer) htmlrenderer.NodeRenderer {
			return nodes
		}),
		htmlrenderer.WithNodeRenderer(footnoteast.KindFootnoteReference, nodes),
		htmlrenderer.WithNodeRenderer(footnoteast.KindFootnoteDefinition, nodes),
	)
}

func (r *epubNodeRenderer) Render(
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
	case *mdcoach.Slide:
		if entering {
			_, err := fmt.Fprintf(w, "<section id=\"slide-%d\">\n", node.Index)
			return ast.WalkContinue, err
		}
		_, err := io.WriteString(w, "</section>\n")
		return ast.WalkContinue, err
	case *mdcoach.Figure:
		if entering {
			_, err := io.WriteString(w, "<figure>\n")
			return ast.WalkContinue, err
		}
		if imageNode, ok := node.FirstChild().(*ast.Image); ok {
			if alt := imageAltText(imageNode, source); alt != "" {
				if _, err := fmt.Fprintf(w, "<figcaption>%s</figcaption>\n", stdhtml.EscapeString(alt)); err != nil {
					return ast.WalkContinue, err
				}
			}
		}
		_, err := io.WriteString(w, "</figure>\n")
		return ast.WalkContinue, err
	case *ast.Image:
		if !entering {
			return ast.WalkSkipChildren, nil
		}
		if err := r.renderImage(w, source, node); err != nil {
			return ast.WalkStop, err
		}
		return ast.WalkSkipChildren, nil
	case *mdcoach.Emote:
		if !entering {
			return ast.WalkSkipChildren, nil
		}
		if err := r.renderEmote(w, source, node); err != nil {
			return ast.WalkStop, err
		}
		return ast.WalkSkipChildren, nil
	case *footnoteast.FootnoteReference:
		if !entering {
			return ast.WalkSkipChildren, nil
		}
		if err := renderFootnoteReference(w, source, node); err != nil {
			return ast.WalkStop, err
		}
		return ast.WalkSkipChildren, nil
	case *footnoteast.FootnoteDefinition:
		if err := renderFootnoteDefinition(w, source, node, entering); err != nil {
			return ast.WalkStop, err
		}
		return ast.WalkContinue, nil
	case *ast.RawHTML:
		if !entering {
			return ast.WalkSkipChildren, nil
		}
		if err := renderRawHTML(w, source, node); err != nil {
			return ast.WalkStop, err
		}
		return ast.WalkSkipChildren, nil
	case *ast.HTMLBlock:
		if entering {
			_, err := io.WriteString(w, stdhtml.EscapeString(node.Value.Str(source)))
			return ast.WalkContinue, err
		}
		return ast.WalkContinue, nil
	default:
		if node.Kind() == mdcoach.KindSlideNotes {
			if entering {
				_, err := io.WriteString(w, "<aside class=\"speaker-notes\"><h2>Speaker notes</h2>\n")
				return ast.WalkContinue, err
			}
			_, err := io.WriteString(w, "</aside>\n")
			return ast.WalkContinue, err
		}
		return ast.WalkContinue, nil
	}
}

func (r *epubNodeRenderer) renderImage(w io.Writer, source []byte, node *ast.Image) error {
	location := node.Destination.Value(source)
	cacheLocation := imageURLFromLocation(location).String()
	image, ok := r.cache.Get(cacheLocation)
	if !ok {
		return fmt.Errorf("image %q is missing from EPUB cache", location)
	}

	if _, err := fmt.Fprintf(
		w,
		`<img src="images/%s.png" alt="%s" width="%d" height="%d"`,
		stdhtml.EscapeString(image.Hash),
		stdhtml.EscapeString(imageAltText(node, source)),
		image.Width,
		image.Height,
	); err != nil {
		return err
	}
	if !node.Title.IsEmpty() {
		if _, err := fmt.Fprintf(w, ` title="%s"`, stdhtml.EscapeString(node.Title.Value(source))); err != nil {
			return err
		}
	}
	_, err := io.WriteString(w, " />")
	return err
}

func imageAltText(node ast.Node, source []byte) string {
	var result strings.Builder
	_ = ast.Walk(node, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch child := child.(type) {
		case *ast.Text:
			result.WriteString(child.Value.Value(source))
		case *ast.CodeSpan:
			result.WriteString(child.Value.Value(source))
		}
		return ast.WalkContinue, nil
	})
	return result.String()
}

func (r *epubNodeRenderer) renderEmote(w io.Writer, source []byte, node *mdcoach.Emote) error {
	name := node.Name.Value(source)
	image, ok, err := r.loadEmote(name)
	if err != nil {
		return fmt.Errorf("load emote %q: %w", name, err)
	}
	if !ok {
		_, err := io.WriteString(w, stdhtml.EscapeString(":"+name+":"))
		return err
	}
	_, err = fmt.Fprintf(
		w,
		`<img src="images/%s.png" alt="%s" title="%s" width="%d" height="%d" />`,
		stdhtml.EscapeString(image.Hash),
		stdhtml.EscapeString(":"+name+":"),
		stdhtml.EscapeString(name),
		image.Width,
		image.Height,
	)
	return err
}

func (r *epubNodeRenderer) loadEmote(name string) (*Image, bool, error) {
	assetPath := name + ".png"
	if !fs.ValidPath(assetPath) {
		return nil, false, nil
	}
	location := path.Join("emotes", assetPath)
	if image, ok := r.cache.Get(location); ok {
		return image, true, nil
	}
	data, err := fs.ReadFile(r.emoteFS, assetPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, err
	}
	decoded, _, err := stdImage.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, false, fmt.Errorf("decode PNG: %w", err)
	}
	image, err := NewImage(location, decoded)
	if err != nil {
		return nil, false, err
	}
	r.cache.Set(image)
	return image, true, nil
}

func renderFootnoteReference(w io.Writer, source []byte, node *footnoteast.FootnoteReference) error {
	label := footnoteID(node.Label.Str(source))
	index := footnoteIndex(node, source, node.Index)
	_, err := fmt.Fprintf(
		w,
		`<sup id="fnref-%s-%d"><a href="#%s">%d</a></sup>`,
		label,
		node.RefIndex,
		label,
		index,
	)
	return err
}

func renderFootnoteDefinition(w io.Writer, source []byte, node *footnoteast.FootnoteDefinition, entering bool) error {
	label := footnoteID(node.Label.Str(source))
	if entering {
		index := footnoteIndex(node, source, 1)
		_, err := fmt.Fprintf(
			w,
			`<aside class="footnote" id="%s"><p><a href="#fnref-%s-0">↩</a> <sup>%d</sup></p>`,
			label,
			label,
			index,
		)
		return err
	}
	_, err := io.WriteString(w, "</aside>")
	return err
}

func footnoteIndex(node ast.Node, source []byte, fallback int) int {
	attribute, ok := node.Attribute(mdcoach.FootnoteIndexAttribute)
	if !ok {
		return fallback
	}
	index, err := strconv.Atoi(attribute.Str(source))
	if err != nil {
		return fallback
	}
	return index
}

func footnoteID(label string) string {
	return "fn-" + base64.RawURLEncoding.EncodeToString([]byte(label))
}

func renderRawHTML(w io.Writer, source []byte, node *ast.RawHTML) error {
	value := node.Value.Value(source)
	if node.Value.IsOwned() && strings.HasPrefix(value, "<footer>") && strings.HasSuffix(value, "</footer>") {
		content := value[len("<footer>") : len(value)-len("</footer>")]
		_, err := fmt.Fprintf(w, "<footer>%s</footer>", stdhtml.EscapeString(stdhtml.UnescapeString(content)))
		return err
	}
	_, err := io.WriteString(w, stdhtml.EscapeString(value))
	return err
}

type taskParagraphRenderer struct {
	next htmlrenderer.NodeRenderer
}

func (r *taskParagraphRenderer) Render(
	writer io.Writer,
	source []byte,
	node ast.Node,
	entering bool,
	context renderer.Context,
) (ast.WalkStatus, error) {
	if entering {
		if paragraph, ok := node.(*ast.Paragraph); ok && paragraph.Parent() != nil && paragraph.Parent().FirstChild() == paragraph {
			if status, isTask := extension.TaskStatusOf(paragraph.Parent()); isTask {
				marker := "[ ] "
				if status == extension.TaskStatusCompleted {
					marker = "[x] "
				}
				if _, err := io.WriteString(writer, marker); err != nil {
					return ast.WalkStop, err
				}
			}
		}
	}
	return r.next.Render(writer, source, node, entering, context)
}
