package mdcoach

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dkotik/mdcoach/internal"
	"github.com/sebdah/goldie/v2"
	"github.com/yuin/goldmark/v2/ast"
)

func TestFigureTransformerPreservesImage(t *testing.T) {
	source := []byte("![cat](media/cat_1.jpg)\n")
	tree := NewParser().Parse(source)

	var figure *Figure
	if err := ast.Walk(tree, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && node.Kind() == KindFigure {
			figure = node.(*Figure)
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	}); err != nil {
		t.Fatal(err)
	}
	if figure == nil {
		t.Fatal("figure node not found")
	}
	if _, ok := figure.FirstChild().(*ast.Image); !ok {
		t.Fatalf("figure first child = %T, want *ast.Image", figure.FirstChild())
	}

	var dump bytes.Buffer
	internal.WriteAST(&dump, tree, source)
	goldie.New(t).Assert(t, "figure_ast", dump.Bytes())
}

func TestFigureRenderer(t *testing.T) {
	source := []byte("![cat](media/cat_1.jpg)\n")
	tree := NewParser().Parse(source)

	var rendered bytes.Buffer
	if err := NewRenderer(NewImageCache()).Render(&rendered, source, tree); err != nil {
		t.Fatal(err)
	}

	goldie.New(t).Assert(t, "figure_html", rendered.Bytes())
}

func TestFigureRendererUsesImageTitleForCaption(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     string
		wantNot  string
	}{
		{
			name:     "renders the image title instead of alt text and escapes it",
			markdown: `![alternative text](media/cat.jpg "Title & <caption>")`,
			want:     `<figcaption>Title &amp; &lt;caption&gt;</figcaption>`,
			wantNot:  `<figcaption>alternative text</figcaption>`,
		},
		{
			name:     "does not render a caption without an image title",
			markdown: `![alternative text](media/cat.jpg)`,
			wantNot:  `<figcaption>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			source := []byte(tt.markdown)
			tree := NewParser().Parse(source)
			var rendered bytes.Buffer
			if err := NewRenderer(NewImageCache()).Render(&rendered, source, tree); err != nil {
				t.Fatal(err)
			}
			output := rendered.String()
			if tt.want != "" && !strings.Contains(output, tt.want) {
				t.Errorf("rendered output does not contain %q:\n%s", tt.want, output)
			}
			if tt.wantNot != "" && strings.Contains(output, tt.wantNot) {
				t.Errorf("rendered output unexpectedly contains %q:\n%s", tt.wantNot, output)
			}
		})
	}
}

func TestFigureRendererRecognizesVideoDestinations(t *testing.T) {
	tests := []struct {
		name        string
		markdown    string
		want        string
		wantNoVideo bool
	}{
		{
			name:     "youtube watch URL",
			markdown: "![clip](https://www.youtube.com/watch?v=dQw4w9WgXcQ)",
			want:     `src="https://www.youtube-nocookie.com/embed/dQw4w9WgXcQ"`,
		},
		{
			name:     "youtube short URL",
			markdown: "![clip](https://youtu.be/dQw4w9WgXcQ)",
			want:     `src="https://www.youtube-nocookie.com/embed/dQw4w9WgXcQ"`,
		},
		{
			name:     "vimeo URL",
			markdown: "![clip](https://vimeo.com/12345678)",
			want:     `src="https://player.vimeo.com/video/12345678"`,
		},
		{
			name:     "dailymotion URL",
			markdown: "![clip](https://www.dailymotion.com/video/x8abcde)",
			want:     `src="https://www.dailymotion.com/embed/video/x8abcde"`,
		},
		{
			name:     "direct mp4 URL",
			markdown: "![clip](https://cdn.example/videos/clip.mp4?download=1&lang=en)",
			want:     `<video controls="controls" preload="metadata"><source src="https://cdn.example/videos/clip.mp4?download=1&amp;lang=en" type="video/mp4" /></video>`,
		},
		{
			name:        "ordinary image falls back to figure rendering",
			markdown:    "![cat](media/cat.jpg)",
			want:        `<figure><div data-src="media/cat.jpg" data-alt="cat"></div></figure>`,
			wantNoVideo: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := []byte(tt.markdown)
			tree := NewParser().Parse(source)
			var rendered bytes.Buffer
			if err := NewRenderer(NewImageCache()).Render(&rendered, source, tree); err != nil {
				t.Fatal(err)
			}
			output := rendered.String()
			if !strings.Contains(output, tt.want) {
				t.Errorf("rendered output does not contain %q:\n%s", tt.want, output)
			}
			if tt.wantNoVideo && (strings.Contains(output, "<video") || strings.Contains(output, "<iframe")) {
				t.Errorf("non-video destination rendered a player:\n%s", output)
			}
		})
	}
}
