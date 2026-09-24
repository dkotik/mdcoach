package mdcoach

import (
	_ "embed"
	"fmt"
	"io"
	"path"
	"regexp"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

var KindEmote = ast.NewNodeKind("Emote")

// Emote is an inline text emoji such as :smile:.
type Emote struct {
	ast.BaseInline
	Name text.SingleLineValue
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

	name := text.NewSingleLineValueFromIndex(
		text.NewIndex(segment.Start+1, segment.Start+end),
		reader.Decoder(),
	)
	reader.Advance(end + 1)
	emote := &Emote{Name: name}
	emote.Init(emote)
	emote.AppendChild(ast.NewText(name))
	return emote
}

func isEmoteCharacter(character byte) bool {
	return character >= 'a' && character <= 'z' ||
		character >= 'A' && character <= 'Z' ||
		character >= '0' && character <= '9' ||
		character == '_' || character == '+' || character == '-'
}

var emojiCodepointPattern = regexp.MustCompile(`'([^']+)'\s*:\s*'([^']+)'`)

//go:embed internal/assets/sass/emote.sass
var emoteStyles string

// NewEmoteRenderer returns a Goldmark renderer that resolves emote aliases to
// cached images from internal/assets/emojis.
func NewEmoteRenderer(cache *ImageCache) html.NodeRenderer {
	if cache == nil {
		panic("nil image cache")
	}
	return html.NodeRendererFunc((&emoteRenderer{
		cache:      cache,
		codepoints: emojiCodepoints(),
	}).render)
}

type emoteRenderer struct {
	cache      *ImageCache
	codepoints map[string]string
}

func (r *emoteRenderer) render(
	writer io.Writer,
	source []byte,
	node ast.Node,
	entering bool,
	rc renderer.Context,
) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	w, ok := writer.(util.BufWriter)
	if !ok {
		w = util.NewErrorBufWriter(writer)
	}
	emote := node.(*Emote)
	name := emote.Name.Value(source)
	class := name
	codepoint, ok := r.codepoints[class]
	if !ok {
		_, _ = html.ContextTextWriter(rc).WriteString(":" + name + ":")
		return ast.WalkSkipChildren, nil
	}

	location := path.Join("internal", "assets", "emojis", codepoint+".png")
	image, ok := r.cache.Get(location)
	if !ok {
		_, _ = html.ContextTextWriter(rc).WriteString(":" + name + ":")
		return ast.WalkSkipChildren, nil
	}

	emote.SetAttribute(
		"data-hash",
		text.NewMultiLineValueFromString(image.Hash, text.IdentityDecoder),
	)
	_, _ = fmt.Fprintf(w, `<span class="emote emote-%s" data-hash="%s">`, class, image.Hash)
	_, _ = emote.Name.WriteTo(html.ContextTextWriter(rc), source)
	_, _ = w.WriteString("</span>")
	return ast.WalkSkipChildren, nil
}

func emojiCodepoints() map[string]string {
	codepoints := make(map[string]string)
	for _, match := range emojiCodepointPattern.FindAllStringSubmatch(emoteStyles, -1) {
		codepoints[match[2]] = match[1]
	}
	return codepoints
}
