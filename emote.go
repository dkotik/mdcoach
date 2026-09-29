package mdcoach

import (
	"bytes"
	"fmt"
	stdImage "image"
	"io"
	"path"
	"sort"

	"github.com/dkotik/mdcoach/internal"
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
var _ parser.InlineParser = (*emoteColonCodeParser)(nil)
var _ parser.InlineParser = (*emoteParser)(nil)

// Kind returns EmoteKind.
func (*Emote) Kind() ast.NodeKind {
	return KindEmote
}

// Dump includes the emote name in AST dumps.
func (e *Emote) Dump(source []byte) *ast.NodeDump {
	return ast.NewNodeDump(e, map[string]any{"Name": e.Name.Str(source)})
}

type emoteParser struct {
	codes    map[string]string
	keys     []string
	triggers []byte
}

// NewEmoteParser returns a Goldmark v2 inline parser for configured emoticons.
func NewEmoteParser(emoticons map[string]string) parser.InlineParser {
	p := &emoteParser{
		codes:    make(map[string]string, len(emoticons)),
		keys:     make([]string, 0, len(emoticons)),
		triggers: make([]byte, 0, len(emoticons)),
	}
	triggerSet := make(map[byte]struct{}, len(emoticons))
	for emoticon, code := range emoticons {
		if emoticon == "" {
			continue
		}
		p.codes[emoticon] = code
		p.keys = append(p.keys, emoticon)
		triggerSet[emoticon[0]] = struct{}{}
	}
	sort.Slice(p.keys, func(i, j int) bool {
		if len(p.keys[i]) == len(p.keys[j]) {
			return p.keys[i] < p.keys[j]
		}
		return len(p.keys[i]) > len(p.keys[j])
	})
	for trigger := range triggerSet {
		p.triggers = append(p.triggers, trigger)
	}
	sort.Slice(p.triggers, func(i, j int) bool {
		return p.triggers[i] < p.triggers[j]
	})
	return p
}

func (p *emoteParser) Trigger() []byte {
	return append([]byte(nil), p.triggers...)
}

func (p *emoteParser) Parse(_ ast.Node, reader text.Reader, _ parser.Context) ast.Node {
	line, _ := reader.PeekLine()
	for _, emoticon := range p.keys {
		if !bytes.HasPrefix(line, []byte(emoticon)) {
			continue
		}

		reader.Advance(len(emoticon))
		name := text.NewSingleLineValueFromString(p.codes[emoticon], text.IdentityDecoder)
		emote := &Emote{Name: name}
		emote.Init(emote)
		emote.AppendChild(ast.NewText(name))
		emote.SetAttribute("class", text.NewMultiLineValueFromString("icon", text.IdentityDecoder))
		return emote
	}
	return nil
}

type emoteColonCodeParser struct{}

// NewEmoteColonCodeParser returns a Goldmark v2 inline parser for :text: emotes.
func NewEmoteColonCodeParser() parser.InlineParser {
	return &emoteColonCodeParser{}
}

func (*emoteColonCodeParser) Trigger() []byte {
	return []byte{':'}
}

func (*emoteColonCodeParser) Parse(_ ast.Node, reader text.Reader, _ parser.Context) ast.Node {
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
	emote.SetAttribute("class", text.NewMultiLineValueFromString("icon", text.IdentityDecoder))
	return emote
}

func isEmoteCharacter(character byte) bool {
	return character >= 'a' && character <= 'z' ||
		character >= 'A' && character <= 'Z' ||
		character >= '0' && character <= '9' ||
		character == '_' || character == '+' || character == '-'
}

// NewEmoteRenderer returns a Goldmark renderer that uses cached images for
// emote names found in internal/assets/emojis.
func NewEmoteRenderer(cache *ImageCache) html.NodeRenderer {
	if cache == nil {
		panic("nil image cache")
	}
	return html.NodeRendererFunc((&emoteRenderer{cache: cache}).render)
}

type emoteRenderer struct {
	cache *ImageCache
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
	_, ok = emoteMap[name]
	if !ok {
		_, _ = html.ContextTextWriter(rc).WriteString(":" + name + ":")
		return ast.WalkSkipChildren, nil
	}

	assetPath := path.Join("assets", "emojis", name+".png")
	location := path.Join("internal", assetPath)
	image, ok := r.cache.Get(location)
	if !ok {
		data, err := internal.Assets.ReadFile(assetPath)
		if err != nil {
			return ast.WalkStop, fmt.Errorf("read emote asset %q: %w", assetPath, err)
		}
		decoded, _, err := stdImage.Decode(bytes.NewReader(data))
		if err != nil {
			return ast.WalkStop, fmt.Errorf("decode emote asset %q: %w", assetPath, err)
		}
		image, err = newImageFromImage(decoded)
		if err != nil {
			return ast.WalkStop, fmt.Errorf("cache emote asset %q: %w", assetPath, err)
		}
		image.Location = location
		r.cache.Set(image)
	}

	emote.SetAttribute(
		"data-hash",
		text.NewMultiLineValueFromString(image.Hash, text.IdentityDecoder),
	)
	// class := name
	// emote.SetAttribute(
	// 	"class",
	// 	text.NewMultiLineValueFromString("emote emote-"+class, text.IdentityDecoder),
	// )
	_, _ = w.WriteString("<span")
	renderImageAttributes(w, source, emote)
	// _, _ =
	// _, _ = emote.Name.WriteTo(html.a(rc), source)
	_, _ = w.WriteString("title=\"")
	_ = internal.WriteEscapedHTML(w, emote.Name.Str(source))
	_, _ = w.WriteString("\">")
	_, _ = w.WriteString("</span>")
	return ast.WalkSkipChildren, nil
}
