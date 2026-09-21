package mdcoach

import (
	"fmt"

	meta "github.com/yuin/goldmark-meta/v2"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/text"
)

var (
	// SlideKind is the node kind for presentation slides.
	SlideKind                       = ast.NewNodeKind("Slide")
	_         ast.Node              = (*Slide)(nil)
	_         parser.ASTTransformer = (*slideCutter)(nil)
)

// Slide is a presentation slide containing a section of the document.
//
// The heading that starts a section is part of that section and is therefore
// the first child of the corresponding Slide node.
type Slide struct {
	ast.BaseBlock
	HeadingLevel        int
	Image               *ast.Image
	IsImageRightAligned bool
}

func NewSlide(children ...ast.Node) (s *Slide) {
	s = &Slide{}
	s.withChildren(children...)
	return
}

func (s *Slide) withChildren(children ...ast.Node) {
	if len(children) == 0 {
		return
	}
	for _, child := range children {
		s.AppendChild(child)
	}
	ok := false
	s.Image, ok = s.FirstChild().(*ast.Image)
	if ok {
		s.RemoveChild(s.Image)
		return
	}

	var firstHeading *ast.Heading
	firstHeading, ok = s.FirstChild().(*ast.Heading)
	if ok {
		s.HeadingLevel = firstHeading.Level
		s.SetAttribute("data-heading-level", text.NewMultiLineValueFromString(
			fmt.Sprintf("%d", s.HeadingLevel),
			text.IdentityDecoder,
		))
	}

	s.Image, ok = s.FirstChild().NextSibling().(*ast.Image)
	if ok {
		s.RemoveChild(s.Image)
		return
	}
	s.Image, ok = s.LastChild().(*ast.Image)
	if ok {
		s.IsImageRightAligned = true
		s.RemoveChild(s.Image)
	}
}

// Kind returns SlideKind.
func (*Slide) Kind() ast.NodeKind {
	return SlideKind
}

// Dump dumps the slide and its children.
func (s *Slide) Dump(_ []byte) *ast.NodeDump {
	return ast.NewNodeDump(s, map[string]any{
		"aside":             s.Image,
		"asideRightAligned": s.IsImageRightAligned,
		// "headingLevel": s.HeadingLevel,
		// "children":   s.Children(),
	})
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
