package mdcomb

type Renderer interface {
	Render(*blackfriday.Node)
}

type renderer struct {
	w  io.Writer
	re blackfriday.Renderer
}

func (r *renderer) Render(node *blackfriday.Node) {
	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return r.re.RenderNode(r.w, node, entering)
	})
}

func NewRenderer(w io.Writer, re blackfriday.Render) Renderer {
	return &renderer{w: w, re: re}
}
