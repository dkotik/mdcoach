package aside

import (
	"bytes"
	"unicode"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type AsideBlockParser struct{}

func (p AsideBlockParser) Trigger() []byte {
	return []byte{':'}
}

func (p AsideBlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	if len(line) < 3 {
		return nil, parser.NoChildren
	}

	for i, c := range line {
		if c == ':' {
			continue
		}
		if i <= 2 {
			return nil, parser.NoChildren
		}
		for _, c = range line[i:] {
			if !unicode.IsSpace(rune(c)) {
				reader.Advance(len(line))
				return &AsideBlockNode{
					Title: bytes.TrimSpace(line[i:]),
				}, parser.HasChildren
			}
		}
	}
	reader.Advance(len(line))
	return &AsideBlockNode{}, parser.HasChildren
}

func (p AsideBlockParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, _ := reader.PeekLine()
	if len(line) > 3 {
		for i, c := range line {
			if c == ':' {
				continue
			}
			if i <= 2 {
				return parser.Continue | parser.HasChildren
			}
			for i, c = range line[i:] {
				if !unicode.IsSpace(rune(c)) {
					reader.Advance(len(line))
					return parser.Close
				}
			}
		}
		reader.Advance(len(line))
		return parser.Close
	}
	return parser.Continue | parser.HasChildren
}

func (p AsideBlockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// no operation
}

func (p AsideBlockParser) CanInterruptParagraph() bool {
	return false
}

func (p AsideBlockParser) CanAcceptIndentedLine() bool {
	return false
}
