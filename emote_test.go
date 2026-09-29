package mdcoach

import (
	"slices"
	"testing"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"

	"github.com/yuin/goldmark/v2/util"
)

func TestEmoteParser(t *testing.T) {
	tests := []struct {
		name      string
		emoticons map[string]string
		source    string
		wantNames []string
		wantBytes []byte
	}{
		{
			name:      "parses configured emoticons and unique triggers",
			emoticons: map[string]string{":-)": "smile", ";-)": "wink"},
			source:    "Hello :-) ;-)!\n",
			wantNames: []string{"smile", "wink"},
			wantBytes: []byte{':', ';'},
		},
		{
			name:      "prefers the longest matching emoticon",
			emoticons: map[string]string{":)": "slight_smile", ":-)": "smile"},
			source:    ":-) :)\n",
			wantNames: []string{"smile", "slight_smile"},
			wantBytes: []byte{':'},
		},
		{
			name:      "empty mapping has no triggers or matches",
			emoticons: map[string]string{},
			source:    ":-)\n",
			wantNames: nil,
			wantBytes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emotes := NewEmoteParser(tt.emoticons)
			if got := emotes.Trigger(); !slices.Equal(got, tt.wantBytes) {
				t.Fatalf("trigger bytes = %q, want %q", got, tt.wantBytes)
			}

			source := []byte(tt.source)
			markdownParser := parser.New(parser.WithInlineParsers(
				util.Prioritized(emotes, 500),
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
			if !slices.Equal(names, tt.wantNames) {
				t.Fatalf("parsed emote names = %v, want %v", names, tt.wantNames)
			}
		})
	}
}

func TestEmoteColonCodeParser(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		wantNames []string
	}{
		{
			name:      "valid and invalid shortcodes",
			source:    "Try :smile: and :cat-2:; leave :: and :not an emote: as text. Escaped \\:smile: will be ignored.\n",
			wantNames: []string{"smile", "cat-2"},
		},
		{
			name:      "rejects incomplete and spaced forms",
			source:    ":: and :not an emote: as text\n",
			wantNames: nil,
		},
		{
			name:      "recognizes a single shortcode",
			source:    ":+1:\n",
			wantNames: []string{"+1"},
		},
		{
			name:      "does not accept:+1: unless there is space before it",
			source:    "does not accept:+1: unless there is space before it\n",
			wantNames: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := []byte(tt.source)
			markdownParser := parser.New(parser.WithInlineParsers(
				util.Prioritized(NewEmoteColonCodeParser(), 500),
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
			if !slices.Equal(names, tt.wantNames) {
				t.Fatalf("parsed emote names = %v, want %v", names, tt.wantNames)
			}
		})
	}
}
