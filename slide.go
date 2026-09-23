package mdcoach

import (
	"fmt"
	"io"
	"strings"

	meta "github.com/yuin/goldmark-meta/v2"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

var (
	// SlideKind is the node kind for presentation slides.
	SlideKind                       = ast.NewNodeKind("Slide")
	_         ast.Node              = (*Slide)(nil)
	_         html.NodeRenderer     = (*slideRenderer)(nil)
	_         parser.ASTTransformer = (*slideCutter)(nil)
)

type SlideLayout string

const (
	SlideNormalLayout     SlideLayout = "normal"
	SlideSplashLayout     SlideLayout = "splash"
	SlideLeftAsideLayout  SlideLayout = "left-aside"
	SlideRightAsideLayout SlideLayout = "right-aside"
)

// Slide is a presentation slide containing a section of the document.
//
// The heading that starts a section is part of that section and is therefore
// the first child of the corresponding Slide node.
type Slide struct {
	ast.BaseBlock
	HeadingLevel int
	SlideLayout  SlideLayout
}

func NewSlide(children ...ast.Node) (s *Slide) {
	s = &Slide{}
	s.withChildren(children...)
	return
}

func (s *Slide) withChildren(children ...ast.Node) {
	s.SlideLayout = SlideNormalLayout
	defer func() {
		s.SetAttribute("data-layout", text.NewMultiLineValueFromString(string(s.SlideLayout), text.IdentityDecoder))
	}()
	if len(children) == 0 {
		return
	}
	var child ast.Node
	contentElementCount := 0
	for _, child = range children {
		// s.OwnerDocument().RemoveChild(child)
		s.AppendChild(child)
		switch child.Kind() {
		case ast.KindHeading:
		case KindFigure:
		case KindAside:
		default:
			contentElementCount++
		}
	}

	ok := false
	var firstHeading *ast.Heading
	firstHeading, ok = s.FirstChild().(*ast.Heading)
	if ok {
		s.HeadingLevel = firstHeading.Level
		s.SetAttribute("data-heading-level", text.NewMultiLineValueFromString(
			fmt.Sprintf("%d", s.HeadingLevel),
			text.IdentityDecoder,
		))
		if s.HeadingLevel == 1 {
			// && contentElementCount == 0
			// possible background image in the last element
			_, ok = s.LastChild().(*Figure)
			if ok {
				s.SlideLayout = SlideSplashLayout
				return
			}
		}
	}
	if contentElementCount == 0 {
		return
	}
	_, ok = s.LastChild().(*Figure)
	if ok {
		s.SlideLayout = SlideRightAsideLayout
		return
	}

	for _, child = range children {
		switch child.Kind() {
		case ast.KindHeading:
		case KindAside:
		case KindFigure:
			_, ok = child.(*Figure)
			if ok {
				s.SlideLayout = SlideLeftAsideLayout
				return
			}
			return
		default:
			return
		}
	}
}

// Kind returns SlideKind.
func (*Slide) Kind() ast.NodeKind {
	return SlideKind
}

// Dump dumps the slide and its children.
func (s *Slide) Dump(_ []byte) *ast.NodeDump {
	return ast.NewNodeDump(s, map[string]any{
		"layout": s.SlideLayout,
		// "headingLevel": s.HeadingLevel,
		// "children":   s.Children(),
	})
}

type slideRenderer struct {
	FigureRenderer html.NodeRenderer
}

// NewSlideRenderer returns a Goldmark v2 HTML renderer for Slide nodes.
func NewSlideRenderer(fr html.NodeRenderer) html.NodeRenderer {
	return &slideRenderer{
		FigureRenderer: fr,
	}
}

func (s *slideRenderer) Render(
	writer io.Writer,
	source []byte,
	node ast.Node,
	entering bool,
	rc renderer.Context,
) (ast.WalkStatus, error) {
	w := writer.(util.BufWriter)
	slide := node.(*Slide)
	if entering {
		_, _ = fmt.Fprintf(w, `<section data-heading-level="%d"`, slide.HeadingLevel)
		renderSlideAttributes(w, source, slide)
		_ = w.WriteByte('>')
		_, _ = w.WriteString(`<div class="grid"><div class="content">`)
		return ast.WalkContinue, nil
	}

	_, _ = w.WriteString("</div></div></section>")
	return ast.WalkContinue, nil
}

// slideCutter groups the document's top-level blocks into Slide nodes.
// A top-level heading starts a new slide; headings nested in other blocks do
// not affect the grouping.
type slideCutter struct {
	HeadingLevelLimit int
}

// NewSlideTransformer returns an AST transformer that groups sections into slides.
func NewSlideTransformer(headingLevelLimit int) parser.ASTTransformer {
	if headingLevelLimit < 1 {
		headingLevelLimit = 6
	}
	return &slideCutter{
		HeadingLevelLimit: headingLevelLimit,
	}
}

// Transform groups the document's top-level blocks into sections. Content
// before the first heading is retained in the first slide. An empty document
// remains empty.
func (s *slideCutter) Transform(document *ast.Document, _ text.Reader, _ parser.Context) {
	if document.ChildCount() == 0 {
		return
	}

	children := make([]ast.Node, 0, 12)
	lastChild := document.FirstChild()
	makeSlide := func() ast.Node {
		slide := &Slide{}
		document.InsertAfter(lastChild, slide)
		// for _, child := range children {
		// 	document.RemoveChild(child)
		// }
		slide.withChildren(children...)
		children = children[:0]
		return slide
	}
	var (
		heading       *ast.Heading
		thematicBreak *ast.ThematicBreak
	)
	for child := document.FirstChild(); child != nil; child = child.NextSibling() {
		switch child.Kind() {
		case meta.KindMetaBlock:
			continue // skip
		case ast.KindThematicBreak:
			lastChild = child
			thematicBreak = child.(*ast.ThematicBreak)
			child = makeSlide()
			document.RemoveChild(thematicBreak)
		case ast.KindHeading:
			heading = child.(*ast.Heading)
			if len(children) > 0 && heading.Level <= s.HeadingLevelLimit {
				_ = makeSlide()
			}
			fallthrough
		default:
			lastChild = child
			children = append(children, child)
		}
	}
	if len(children) > 0 {
		makeSlide()
	}
	// document.RemoveChildren()

	// var slide *Slide
	// for _, child := range children {
	// 	if slide == nil || child.Kind() == ast.KindHeading {
	// 		slide = &Slide{}
	// 		document.AppendChild(slide)
	// 	}
	// 	slide.AppendChild(child)
	// }
}

func renderSlideAttributes(writer io.Writer, source []byte, node ast.Node) {
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
		_, _ = w.WriteString(" ")
		_, _ = w.WriteString(attr.Name)
		_, _ = w.WriteString(`="`)
		_, _ = attr.Value.WriteTo(tw, source)
		_ = w.WriteByte('"')
	}
}
