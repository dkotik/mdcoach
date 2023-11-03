/*
Package renderer converts parsed [ast.Node] Markdown to HTML. Images are resized and saved to a directory.
*/
package renderer

import (
	"fmt"

	"github.com/OneOfOne/xxhash"
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

// type Renderer interface {
// 	Render(w io.Writer, source []byte, n ast.Node) error
//
// 	// AddOptions adds given option to this renderer.
// 	AddOptions(...Option)
// }

func New() (renderer.Renderer, error) {
	provider, err := picture.NewLocalProvider()
	if err != nil {
		return nil, err
	}

	r := &Renderer{
		pictureProvider: provider,
	}

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
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindParagraph, r.renderParagraph)
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
