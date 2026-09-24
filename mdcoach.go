/*
Package mdcoach converts Markdown files to HTML presentations with notes.
*/
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
			extension.NewTableParser(),
		),
		parser.WithInlineParsers(
			util.Prioritized(NewEmoteParser(), 500),
		),
		parser.WithBlockParsers(
			util.Prioritized(NewAsideParser(), 10),
		),
		parser.WithASTTransformers(
			// Run after extension AST transformers so slides contain their
			// final block structure.
			util.Prioritized(NewFigureTransformer(), 4),
			util.Prioritized(NewSlideTransformer(2), 900),
			util.Prioritized(NewDocumentIDInjector(), 1400),
		),
	)
}

type imageRendererExtension struct{}

func (*imageRendererExtension) RendererOptions(*html.Config) []html.Option {
	return []html.Option{
		html.WithNodeRenderer(ast.KindImage, NewImageRenderer()),
		html.WithNodeRenderer(KindEmote, NewEmoteRenderer(NewImageCache())),
		html.WithNodeRenderer(KindAside, NewAsideRenderer()),
		html.WithNodeRenderer(KindFigure, NewFigureRenderer()),
		html.WithNodeRenderer(SlideKind, NewSlideRenderer(
			NewFigureRenderer(),
		)),
	}
}

func NewRenderer() html.Renderer {
	return html.New(
		html.WithExtensions(&imageRendererExtension{}),
	)
}
