package parser

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var KindNotesBreak = ast.NewNodeKind("NotesBreak")

type NotesBreak struct {
	*ast.ThematicBreak
}

func (n *NotesBreak) Kind() ast.NodeKind {
	return KindNotesBreak
}

type notesBreakParser struct{}

// NewNotesBreakParser returns a new BlockParser that
// parses only thematic breaks composed of asterisks.
// It is based on [parser.NewThematicBreakParser],
// so it should be placed higher in parser priority to work
func NewNotesBreakParser() parser.BlockParser {
	return &notesBreakParser{}
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
		return &NotesBreak{ast.NewThematicBreak()}, parser.NoChildren
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
