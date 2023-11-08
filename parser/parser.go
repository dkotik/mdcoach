/*
Package parser extends default Goldmark Markdown parser with additional functionality, necessary for building slide presentations.
*/
package parser

import (
	"fmt"

	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/frontmatter"
)

// New returns a new Parser that is configured by default values.
func New(withOptions ...Option) (_ parser.Parser, err error) {
	o := &options{}
	for _, option := range append(
		withOptions,
		withDefaultListItemParser(),
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot create parser: %w", err)
		}
	}

	return parser.NewParser(
		parser.WithBlockParsers(
			util.Prioritized(&frontmatter.Parser{
				Formats: frontmatter.DefaultFormats,
			}, 0),
			util.Prioritized(parser.NewSetextHeadingParser(), 100),
			util.Prioritized(NewNotesBreakParser(), 180),
			util.Prioritized(parser.NewThematicBreakParser(), 200),
			util.Prioritized(parser.NewListParser(), 300),
			// util.Prioritizparser.Nd(NewListItemParser(), 400),
			// util.Prioritized(NewReviewQuestionParser(), 400),
			util.Prioritized(o.listItemParser, 400),
			util.Prioritized(parser.NewCodeBlockParser(), 500),
			util.Prioritized(parser.NewATXHeadingParser(), 600),
			util.Prioritized(parser.NewFencedCodeBlockParser(), 700),
			util.Prioritized(parser.NewBlockquoteParser(), 800),
			util.Prioritized(parser.NewHTMLBlockParser(), 900),
			util.Prioritized(parser.NewParagraphParser(), 1000),
		),
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
		parser.WithParagraphTransformers(parser.DefaultParagraphTransformers()...),
		parser.WithASTTransformers(
			util.Prioritized(&frontmatter.MetaTransformer{}, 0),
		),
	), nil
}
