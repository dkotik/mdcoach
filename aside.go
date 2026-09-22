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

var AsideKind = ast.NewNodeKind("Aside")

var _ html.NodeRenderer = (*asideRenderer)(nil)

// A Aside struct represents a marginal note.
type Aside struct {
	ast.BaseBlock
}

// Dump implements Node.Dump .
func (n *Aside) Dump(_ []byte) *ast.NodeDump {
	return ast.NewNodeDump(n, nil)
}

// Kind implements Node.Kind.
func (n *Aside) Kind() ast.NodeKind {
	return AsideKind
}

// NewAside returns a new [Aside] node.
func NewAside() *Aside {
	n := &Aside{}
	n.Init(n)
	return n
}

type asideRenderer struct{}

// NewAsideRenderer returns a Goldmark v2 HTML renderer for Aside nodes.
func NewAsideRenderer() html.NodeRenderer {
	return &asideRenderer{}
}

func (*asideRenderer) Render(
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
		_, _ = w.WriteString("<aside>")
	} else {
		_, _ = w.WriteString("</aside>")
	}
	return ast.WalkContinue, nil
}

type asideParser struct {
}

var defaultAsideParser = &asideParser{}

// NewAsideParser returns a new BlockParser that
// parses marginal notes.
func NewAsideParser() parser.BlockParser {
	return defaultAsideParser
}

func (a *asideParser) Trigger() []byte {
	return []byte{'*'}
}

func (a *asideParser) Open(_ ast.Node, reader text.Reader, _ parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	if isAsideBreak(line, reader.LineOffset()) {
		reader.AdvanceToEOL()
		return NewAside(), parser.HasChildren
	}
	return nil, parser.NoChildren
}

func (a *asideParser) Continue(_ ast.Node, reader text.Reader, _ parser.Context) parser.State {
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

func (a *asideParser) Close(_ ast.Node, _ text.Reader, _ parser.Context) {
	// nothing to do
}

func (a *asideParser) CanInterruptParagraph() bool {
	return true
}

func (a *asideParser) CanAcceptIndentedLine() bool {
	return true
}

func isAsideBreak(line []byte, offset int) bool {
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
