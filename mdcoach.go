package mdcoach

import (
	meta "github.com/yuin/goldmark-meta/v2"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/util"
)

// NewParser returns a Goldmark v2 parser configured for presentations.
func NewParser() parser.Parser {
	return parser.New(
		parser.WithExtensions(
			meta.Parser,
			extension.NewGFMParser(),
			extension.NewDefinitionListParser(),
			extension.NewFootnoteParser(),
			extension.NewTypographerParser(),
		),
		parser.WithBlockParsers(
			util.Prioritized(NewAsideParser(), 10),
		),
		parser.WithASTTransformers(
			// Run after extension AST transformers so slides contain their
			// final block structure.
			util.Prioritized(NewSlideTransformer(2), 999),
		),
	)
}

func NewRenderer() html.Renderer {
	return html.New(
		html.WithNodeRenderer(ast.KindImage, NewImageRenderer()),
	)
}
