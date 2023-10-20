package mdcoach

import (
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// DefaultParser returns a new Parser that is configured by default values.
func DefaultParser() parser.Parser {
	return parser.NewParser(parser.WithBlockParsers(
		append(parser.DefaultBlockParsers(),
			// Inject notes break detector above thematic break parser.
			util.Prioritized(NewNotesBreakParser(), 180),
		// util.Prioritized(NewThematicBreakParser(), 200),
		)...,
	),
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
		parser.WithParagraphTransformers(parser.DefaultParagraphTransformers()...),
	)
}
