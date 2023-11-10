/*
Package renderer converts parsed [ast.Node] Markdown to HTML. Images are resized and saved to a directory.
*/
package renderer

import (
	"errors"
	"fmt"

	"github.com/OneOfOne/xxhash"
	"github.com/dkotik/mdcoach/parser"
	"github.com/dkotik/mdcoach/picture"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type Renderer struct {
	sizings         []*picture.Sizing
	pictureProvider picture.Provider
	sourceDirectory string
	outputDirectory string
}

func New(withOptions ...Option) (_ renderer.Renderer, err error) {
	o := &options{}
	for _, option := range append(
		withOptions,
		func(o *options) error {
			if o.PictureProvider == nil {
				return errors.New("a picture provider is required")
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot create a Markdown renderer: %w", err)
		}
	}

	return renderer.NewRenderer(renderer.WithNodeRenderers(
		util.Prioritized(&Renderer{
			pictureProvider: o.PictureProvider,
		}, 900),
		util.Prioritized(html.NewRenderer(), 1000),
	)), nil
}

func Must(r renderer.Renderer, err error) renderer.Renderer {
	if err != nil {
		panic(err)
	}
	return r
}

func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(parser.KindNotesBreak, r.renderNotesBreak)
	// reg.Register(ast.KindLink, r.renderLink)
	// reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	// reg.Register(ast.KindFencedCodeBlock, r.renderCodeBlock)
}

func (r *Renderer) PathHash(p string) string {
	h := xxhash.New64()
	h.WriteString(r.sourceDirectory)
	h.WriteString("^")
	h.WriteString(p)
	return fmt.Sprintf("%x", h.Sum64())
}
