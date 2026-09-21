package mdcoach

import (
	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/util"
)

// NewParser returns a Goldmark v2 parser configured for presentations.
func NewParser() parser.Parser {
	return parser.New(
		parser.WithExtensions(
			extension.NewGFMParser(),
			extension.NewDefinitionListParser(),
			extension.NewFootnoteParser(),
			extension.NewTypographerParser(),
		),
		parser.WithASTTransformers(
			// Run after extension AST transformers so slides contain their
			// final block structure.
			util.Prioritized(NewSlideTransformer(), 999),
		),
	)
}
