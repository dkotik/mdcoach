package aside

import "github.com/yuin/goldmark/ast"

var AsideBlockKind = ast.NewNodeKind("AsideBlock")

var _ ast.Node = &AsideBlockNode{}

type AsideBlockNode struct {
	ast.BaseBlock

	Title []byte
}

func (n *AsideBlockNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"Title": string(n.Title),
	}, nil)
}

func (n *AsideBlockNode) Kind() ast.NodeKind {
	return AsideBlockKind
}
