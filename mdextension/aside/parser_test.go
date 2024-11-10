package aside_test

import (
	"bytes"
	"testing"

	"github.com/dkotik/mdcoach/mdextension/aside"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func TestAsideParserSuccess(t *testing.T) {
	cases := [...][]byte{
		[]byte(`::: title
body
:::`),
		[]byte(`::::::::::::::::::::::::: title
body
::::::::::::::::::::::::::::::::: `),
	}

	p := parser.NewParser(parser.WithBlockParsers(util.Prioritized(aside.AsideBlockParser{}, 500)))

	for i, cs := range cases {
		tree := p.Parse(text.NewReader(cs))
		if tree == nil {
			t.Fatalf("could not parse case %d: parse tree is empty", i)
		}
		if tree.ChildCount() != 1 {
			t.Fatal("too many children parsed")
		}

		first := tree.FirstChild()
		if first.Kind() != aside.AsideBlockKind {
			t.Fatal("parsed node is not an aside block")
		}

		n, ok := first.(*aside.AsideBlockNode)
		if !ok {
			t.Fatal("parsed node is not an aside block")
		}

		if !bytes.Equal(n.Title, []byte("title")) {
			t.Fatal("unexpected parsed title:", string(n.Title))
		}
	}
}

func TestAsideParserFailure(t *testing.T) {
	cases := [...][]byte{
		[]byte(`:: title
body
:::`),
		[]byte(`::::::::::::::::::::::::: title
body
::`),
	}

	p := parser.NewParser(parser.WithBlockParsers(util.Prioritized(aside.AsideBlockParser{}, 500)))

	for i, cs := range cases {
		tree := p.Parse(text.NewReader(cs))
		if tree == nil {
			t.Fatalf("could not parse case %d: parse tree is empty", i)
		}

		first := tree.FirstChild()
		if first == nil {
			return
		}
		if first.Kind() == aside.AsideBlockKind {
			t.Fatal("parsed node is an aside block")
		}

		_, ok := first.(*aside.AsideBlockNode)
		if ok {
			t.Fatal("parsed node is an aside block")
		}
	}
}
