package mdcoach

import (
	"bytes"
	"testing"

	"github.com/dkotik/mdcoach/internal"
	"github.com/sebdah/goldie/v2"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"

	"github.com/yuin/goldmark/v2/util"
)

func TestEmoteParser(t *testing.T) {
	source := []byte("Try :smile: and :cat-2:; leave :: and :not an emote: as text.\n")
	markdownParser := parser.New(parser.WithInlineParsers(
		util.Prioritized(NewEmoteParser(), 500),
	))
	tree := markdownParser.Parse(source)

	var names []string
	if err := ast.Walk(tree, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if emote, ok := node.(*Emote); ok {
				names = append(names, emote.Name.Value(source))
			}
		}
		return ast.WalkContinue, nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 || names[0] != "smile" || names[1] != "cat-2" {
		t.Fatalf("parsed emote names = %v, want [smile cat-2]", names)
	}

	var dump bytes.Buffer
	internal.WriteAST(&dump, tree, source)
	goldie.New(t).Assert(t, "emote_ast", dump.Bytes())
}
