package mdcoach

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/yuin/goldmark/v2/ast"
)

func TestPresentationAST(t *testing.T) {
	source, err := os.ReadFile("testdata/presentation.md")
	if err != nil {
		t.Fatal(err)
	}

	tree := NewParser().Parse(source)
	var dump bytes.Buffer
	dumpAST(tree, &dump)

	goldie.New(t).Assert(t, "presentation_ast", dump.Bytes())
}

func dumpAST(tree ast.Node, dump *bytes.Buffer) {
	level := 0
	_ = ast.Walk(tree, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			fmt.Fprintf(dump, "%*s%s {\n", level*4, "", node.Kind())
			level++
		} else {
			level--
			fmt.Fprintf(dump, "%*s}\n", level*4, "")
		}
		return ast.WalkContinue, nil
	})
}
