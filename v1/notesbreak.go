package mdcoach

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

const notesBreakAttributeName = "isNotesThematicBreak"

func IsNotesThematicBreak(n ast.Node) (ok bool) {
	value, ok := n.AttributeString(notesBreakAttributeName)
	if !ok {
		return
	}
	ok, _ = value.(bool)
	return
}

type notesBreakParser struct{}

var defaultNotesBreakParser = &notesBreakParser{}

// NewNotesBreakParser returns a new BlockParser that
// parses only thematic breaks composed of asterisks.
// It is based on [parser.NewThematicBreakParser],
// so it should be placed higher in parser priority to work
func NewNotesBreakParser() parser.BlockParser {
	return defaultNotesBreakParser
}

func isThematicBreak(line []byte, offset int) bool {
	w, pos := util.IndentWidth(line, offset)
	if w > 3 {
		return false
	}
	mark := byte(0)
	count := 0
	for i := pos; i < len(line); i++ {
		c := line[i]
		if util.IsSpace(c) {
			continue
		}
		if mark == 0 {
			mark = c
			count = 1
			if mark == '*' { // || mark == '-' || mark == '_'
				continue
			}
			return false
		}
		if c != mark {
			return false
		}
		count++
	}
	return count > 2
}

func (b *notesBreakParser) Trigger() []byte {
	// return []byte{'-', '*', '_'}
	return []byte{'*'}
}

func (b *notesBreakParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, segment := reader.PeekLine()
	if isThematicBreak(line, reader.LineOffset()) {
		reader.Advance(segment.Len() - 1)
		n := ast.NewThematicBreak()
		// make this break identifiable
		n.SetAttributeString(notesBreakAttributeName, true)
		return n, parser.NoChildren
	}
	return nil, parser.NoChildren
}

func (b *notesBreakParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	return parser.Close
}

func (b *notesBreakParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// nothing to do
}

func (b *notesBreakParser) CanInterruptParagraph() bool {
	return true
}

func (b *notesBreakParser) CanAcceptIndentedLine() bool {
	return false
}
