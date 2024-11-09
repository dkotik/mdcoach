package aside

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

var _ parser.BlockParser = AsideBlockParser{}

// https://youtu.be/MNkqyCybWeM?feature=shared

type AsideBlockExtension struct {
}

func (ext AsideBlockExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(AsideBlockParser{}, 500),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(AsideBlockRenderer{}, 500)),
	)
}
