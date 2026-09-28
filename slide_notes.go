package mdcoach

import (
	"io"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

var KindSlideNotes = ast.NewNodeKind("SlideNotes")

var _ html.NodeRenderer = (*slideNotesRenderer)(nil)

// SlideNotes represents notes associated with a slide.
type SlideNotes struct {
	ast.BaseBlock
}

// Dump implements Node.Dump .
func (n *SlideNotes) Dump(_ []byte) *ast.NodeDump {
	return ast.NewNodeDump(n, nil)
}

// Kind implements Node.Kind.
func (n *SlideNotes) Kind() ast.NodeKind {
	return KindSlideNotes
}

// NewSlideNotes returns a new [SlideNotes] node.
func NewSlideNotes() *SlideNotes {
	n := &SlideNotes{}
	n.Init(n)
	return n
}

type slideNotesRenderer struct{}

// NewSlideNotesRenderer returns a Goldmark v2 HTML renderer for SlideNotes nodes.
func NewSlideNotesRenderer() html.NodeRenderer {
	return &slideNotesRenderer{}
}

func (*slideNotesRenderer) Render(
	writer io.Writer,
	_ []byte,
	node ast.Node,
	entering bool,
	_ renderer.Context,
) (ast.WalkStatus, error) {
	w, ok := writer.(util.BufWriter)
	if !ok {
		w = util.NewErrorBufWriter(writer)
	}
	if entering {
		_, _ = w.WriteString("<slide-notes>")
	} else {
		_, _ = w.WriteString("</slide-notes>")
	}
	return ast.WalkContinue, nil
}

type slideNotesParser struct{}

var defaultSlideNotesParser = &slideNotesParser{}

// NewSlideNotesParser returns a BlockParser for slide notes.
func NewSlideNotesParser() parser.BlockParser {
	return defaultSlideNotesParser
}

func (p *slideNotesParser) Trigger() []byte {
	return []byte{'*'}
}

func (p *slideNotesParser) Open(_ ast.Node, reader text.Reader, _ parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	if isSlideNotesBreak(line, reader.LineOffset()) {
		reader.AdvanceToEOL()
		return NewSlideNotes(), parser.HasChildren
	}
	return nil, parser.NoChildren
}

func (p *slideNotesParser) Continue(_ ast.Node, reader text.Reader, _ parser.Context) parser.State {
	line, _ := reader.PeekLine()
	offset := reader.LineOffset()
	_, pos := util.IndentWidth(line, offset)
	switch line[pos] {
	case '#':
		return parser.Close
	case '*', '-', '_':
		if isThematicBreak(line, offset) {
			return parser.Close
		}
	}
	return parser.Continue | parser.HasChildren
}

func (p *slideNotesParser) Close(_ ast.Node, _ text.Reader, _ parser.Context) {
	// nothing to do
}

func (p *slideNotesParser) CanInterruptParagraph() bool {
	return true
}

func (p *slideNotesParser) CanAcceptIndentedLine() bool {
	return true
}

func isSlideNotesBreak(line []byte, offset int) bool {
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
			if mark == '*' {
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
			if mark == '*' || mark == '-' || mark == '_' {
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
