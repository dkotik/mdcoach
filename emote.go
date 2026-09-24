package mdcoach

import (
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/text"
)

var KindEmote = ast.NewNodeKind("Emote")

// Emote is an inline text emoji such as :smile:.
type Emote struct {
	ast.BaseInline
	Name text.MultiLineValue
}

var _ ast.Node = (*Emote)(nil)
var _ parser.InlineParser = (*emoteParser)(nil)

// Kind returns EmoteKind.
func (*Emote) Kind() ast.NodeKind {
	return KindEmote
}

// Dump includes the emote name in AST dumps.
func (e *Emote) Dump(source []byte) *ast.NodeDump {
	return ast.NewNodeDump(e, map[string]any{"Name": e.Name.Str(source)})
}

// NewEmoteParser returns a Goldmark v2 inline parser for :text: emotes.
func NewEmoteParser() parser.InlineParser {
	return &emoteParser{}
}

type emoteParser struct{}

func (*emoteParser) Trigger() []byte {
	return []byte{':'}
}

func (*emoteParser) Parse(_ ast.Node, reader text.Reader, _ parser.Context) ast.Node {
	line, segment := reader.PeekLine()
	if len(line) < 3 || line[0] != ':' {
		return nil
	}

	end := 1
	for end < len(line) && isEmoteCharacter(line[end]) {
		end++
	}
	if end == 1 || end >= len(line) || line[end] != ':' {
		return nil
	}

	name := reader.ValueBetween(segment.Start+1, segment.Start+end)
	reader.Advance(end + 1)
	emote := &Emote{Name: name}
	emote.Init(emote)
	return emote
}

func isEmoteCharacter(character byte) bool {
	return character >= 'a' && character <= 'z' ||
		character >= 'A' && character <= 'Z' ||
		character >= '0' && character <= '9' ||
		character == '_' || character == '+' || character == '-'
}
