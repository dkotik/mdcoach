package parser

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func TestReviewQuestionExtraction(t *testing.T) {
	source := []byte(`
# Heading 1

- one?
- two?

  paragraph providing more instructions to two
- three\?
`)

	ctx := parser.NewContext()
	_ = New().Parse(text.NewReader(source), parser.WithContext(ctx))

	questions := QuestionListFromContext(ctx)
	if len(questions) != 2 {
		t.Fatalf("unexpected number of questions detected, wanted %d, but found %d instead", 2, len(questions))
	}

	b := &bytes.Buffer{}
	r := goldmark.DefaultRenderer()
	for _, q := range questions {
		if err := r.Render(b, source, q); err != nil {
			t.Fatal(err)
		}
		t.Log(b.String())
		b.Reset()
	}
	// spew.Dump(questions)
	// t.Fatal("impl")
}
