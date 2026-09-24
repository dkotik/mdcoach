```go
	return parser.NewParser(
		parser.WithBlockParsers(
			util.Prioritized(parser.NewSetextHeadingParser(), 100),
			util.Prioritized(extension.NewDefinitionListParser(), 101),
			util.Prioritized(extension.NewDefinitionDescriptionParser(), 102),
			util.Prioritized(NewNotesBreakParser(), 180),
			util.Prioritized(parser.NewThematicBreakParser(), 200),
			util.Prioritized(parser.NewListParser(), 300),
			// util.Prioritized(parser.NewListItemParser(), 400),
			util.Prioritized(o.listItemParser, 400),
			util.Prioritized(o.listItemParser, 400),
			util.Prioritized(parser.NewCodeBlockParser(), 500),
			util.Prioritized(parser.NewATXHeadingParser(), 600),
			util.Prioritized(parser.NewFencedCodeBlockParser(), 700),
			util.Prioritized(parser.NewBlockquoteParser(), 800),
			// util.Prioritized(NewBlockquoteParser(), 800),
			util.Prioritized(parser.NewHTMLBlockParser(), 900),
			util.Prioritized(extension.NewFootnoteBlockParser(), 999),
			util.Prioritized(parser.NewParagraphParser(), 1000),
		),
		parser.WithInlineParsers(append(
			parser.DefaultInlineParsers(),
			util.Prioritized(extension.NewTaskCheckBoxParser(), 0),
			util.Prioritized(extension.NewFootnoteParser(), 101),
			util.Prioritized(extension.NewStrikethroughParser(), 500),
			// Linkify extension
			util.Prioritized(extension.NewLinkifyParser(
			// opts ...LinkifyOption
			), 999),
			util.Prioritized(extension.NewTypographerParser(
			// opts ...TypographerOption
			), 9999),
		)...),
		parser.WithParagraphTransformers(append(
			parser.DefaultParagraphTransformers(),
			util.Prioritized(extension.NewTableParagraphTransformer(), 200),
		)...),
		parser.WithASTTransformers(
			util.Prioritized(extension.NewTableASTTransformer(), 0),
			// TODO: Footnote collect footnotes per slide, separately.
			// util.Prioritized(extension.NewFootnoteASTTransformer(), 999),
		),
	), nil
}
```
