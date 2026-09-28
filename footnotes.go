package mdcoach

import (
	"github.com/yuin/goldmark/v2/ast"
	footnoteast "github.com/yuin/goldmark/v2/extension/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/text"
)

var _ parser.ASTTransformer = (*footnoteAsideTransformer)(nil)

// NewFootnoteAsideTransformer moves referenced footnote definitions into the
// first aside in the slide containing each reference.
func NewFootnoteAsideTransformer() parser.ASTTransformer {
	return &footnoteAsideTransformer{}
}

type footnoteAsideTransformer struct{}

type footnoteSlide struct {
	aside      ast.Node
	references map[string]struct{}
}

func findASTParent(root, target ast.Node) ast.Node {
	var parent ast.Node
	_ = ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			if child == target {
				parent = node
				return ast.WalkStop, nil
			}
		}
		return ast.WalkContinue, nil
	})
	return parent
}

func (*footnoteAsideTransformer) Transform(
	document *ast.Document,
	reader text.Reader,
	_ parser.Context,
) {
	if reader == nil {
		return
	}

	source := reader.Source()
	var definitions []*footnoteast.FootnoteDefinition
	_ = ast.Walk(document, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if definition, ok := node.(*footnoteast.FootnoteDefinition); ok {
			definitions = append(definitions, definition)
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})

	var slides []footnoteSlide
	for slide := document.FirstChild(); slide != nil; slide = slide.NextSibling() {
		if slide.Kind() != SlideKind {
			continue
		}

		current := footnoteSlide{references: make(map[string]struct{})}
		_ = ast.Walk(slide, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
			if !entering {
				return ast.WalkContinue, nil
			}
			if node.Kind() == KindAside && current.aside == nil {
				current.aside = node
			}
			if node.Kind() == footnoteast.KindFootnoteDefinition {
				return ast.WalkSkipChildren, nil
			}
			if reference, ok := node.(*footnoteast.FootnoteReference); ok {
				current.references[reference.Label.Str(source)] = struct{}{}
			}
			return ast.WalkContinue, nil
		})
		slides = append(slides, current)
	}

	for _, definition := range definitions {
		label := definition.Label.Str(source)
		for _, slide := range slides {
			if slide.aside == nil {
				continue
			}
			if _, referenced := slide.references[label]; !referenced {
				continue
			}
			if definition.Parent() != slide.aside {
				parent := definition.Parent()
				if parent == nil {
					parent = findASTParent(document, definition)
				}
				if parent != nil {
					parent.RemoveChild(definition)
				} else {
					definition.SetPreviousSibling(nil)
					definition.SetNextSibling(nil)
				}
				slide.aside.AppendChild(definition)
			}
			break
		}
	}
}
