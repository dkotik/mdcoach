/*
Package renderer converts parsed [ast.Node] Markdown to HTML. Images are resized and saved to a directory.
*/
package renderer

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type Renderer struct {
	sourceDirectory string
	outputDirectory string
}

// type Renderer interface {
// 	Render(w io.Writer, source []byte, n ast.Node) error
//
// 	// AddOptions adds given option to this renderer.
// 	AddOptions(...Option)
// }

func New() (renderer.Renderer, error) {
	r := &Renderer{}

	return renderer.NewRenderer(renderer.WithNodeRenderers(
		util.Prioritized(r, 900),
		util.Prioritized(html.NewRenderer(), 1000),
	)), nil
}

func Must() renderer.Renderer {
	r, err := New()
	if err != nil {
		panic(err)
	}
	return r
}

func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	// reg.Register(ast.KindLink, r.renderLink)
	// reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	// reg.Register(ast.KindFencedCodeBlock, r.renderCodeBlock)
}
