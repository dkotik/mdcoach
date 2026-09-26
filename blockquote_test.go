package mdcoach

import (
	"bytes"
	"strings"
	"testing"
)

func TestBlockquoteTransformer(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		wantFooter string
		wantText   string
	}{
		{
			name:       "citation at end of quote paragraph",
			source:     "> Quote text (Citation)\n",
			wantFooter: "<footer>Citation</footer>",
			wantText:   "<p>Quote text</p>",
		},
		{
			name:       "citation in final blockquote paragraph",
			source:     "> Quote text\n>\n> (Citation)\n",
			wantFooter: "<footer>Citation</footer>",
			wantText:   "<p>Quote text</p>",
		},
		{
			name:     "no trailing parenthetical",
			source:   "> Quote text (not at the end) more text\n",
			wantText: "Quote text (not at the end) more text",
		},
		{
			name:       "citation is escaped as text",
			source:     "> Quote text (A & B)\n",
			wantFooter: "<footer>A &amp; B</footer>",
			wantText:   "<p>Quote text</p>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			source := []byte(test.source)
			tree := NewParser().Parse(source)
			var rendered bytes.Buffer
			if err := NewRenderer(nil).Render(&rendered, source, tree); err != nil {
				t.Fatal(err)
			}
			output := rendered.String()

			if test.wantFooter != "" && !strings.Contains(output, test.wantFooter) {
				t.Errorf("rendered HTML %q does not contain footer %q", output, test.wantFooter)
			}
			if !strings.Contains(output, test.wantText) {
				t.Errorf("rendered HTML %q does not contain expected quote text %q", output, test.wantText)
			}
			if test.wantFooter == "" && strings.Contains(output, "<footer>") {
				t.Errorf("rendered HTML %q unexpectedly contains a footer", output)
			}
		})
	}
}

func TestBlockquoteFooterRendererOmitsSourceRawHTML(t *testing.T) {
	source := []byte("> <footer onclick=\"alert(1)\">unsafe</footer>\n")
	tree := NewParser().Parse(source)

	var rendered bytes.Buffer
	if err := NewRenderer(nil).Render(&rendered, source, tree); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(rendered.String(), `<footer onclick=`) {
		t.Fatalf("source raw HTML was rendered: %s", rendered.String())
	}
}
