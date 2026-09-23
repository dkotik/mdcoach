package mdcoach

import (
	"bytes"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/yuin/goldmark/v2/ast"
)

func TestFigureTransformerPreservesImage(t *testing.T) {
	source := []byte("![cat](media/cat_1.jpg)\n")
	tree := NewParser().Parse(source)

	var figure *Figure
	if err := ast.Walk(tree, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && node.Kind() == KindFigure {
			figure = node.(*Figure)
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	}); err != nil {
		t.Fatal(err)
	}
	if figure == nil {
		t.Fatal("figure node not found")
	}
	if _, ok := figure.FirstChild().(*ast.Image); !ok {
		t.Fatalf("figure first child = %T, want *ast.Image", figure.FirstChild())
	}

	var dump bytes.Buffer
	dumpAST(tree, source, &dump)
	goldie.New(t).Assert(t, "figure_ast", dump.Bytes())
}
