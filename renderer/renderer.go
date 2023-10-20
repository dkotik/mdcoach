package renderer

import "github.com/yuin/goldmark/renderer"

type Renderer struct{}

// type Renderer interface {
// 	Render(w io.Writer, source []byte, n ast.Node) error
//
// 	// AddOptions adds given option to this renderer.
// 	AddOptions(...Option)
// }

func New() (*Renderer, error) {
	return &Renderer{}, nil
}

func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// reg.Register(ast.KindParagraph, r.renderParagraph)
	// reg.Register(ast.KindLink, r.renderLink)
	// reg.Register(ast.KindBlockquote, r.renderBlockquote)
	// reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	// reg.Register(ast.KindFencedCodeBlock, r.renderCodeBlock)
}
