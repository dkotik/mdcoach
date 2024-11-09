package aside

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type AsideBlockRenderer struct{}

func (r AsideBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(AsideBlockKind, r.render)
}

func (r AsideBlockRenderer) render(
	w util.BufWriter, source []byte, node ast.Node, entering bool,
) (ast.WalkStatus, error) {
	n := node.(*AsideBlockNode)
	if entering {
		_, _ = w.WriteString("<aside>")
		if len(n.Title) > 0 {
			_, _ = w.WriteString("<h3>")
			_, _ = w.Write(util.EscapeHTML(n.Title))
			_, _ = w.WriteString("</h3>\n")
		}
	} else {
		_, _ = w.WriteString("</aside>\n")
	}
	return ast.WalkContinue, nil
}
