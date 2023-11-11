package test

import (
	"bytes"
	"testing"

	mdcParser "github.com/dkotik/mdcoach/parser"
	"github.com/dkotik/mdcoach/picture"
	"github.com/dkotik/mdcoach/renderer"
	"github.com/yuin/goldmark/ast"
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
	p, err := mdcParser.New()
	if err != nil {
		t.Fatal(err)
	}
	p.Parse(text.NewReader(source), parser.WithContext(ctx))

	questions := mdcParser.QuestionListFromContext(ctx)
	if len(questions) != 2 {
		t.Fatalf("unexpected number of questions detected, wanted %d, but found %d instead", 2, len(questions))
	}

	b := &bytes.Buffer{}
	r, err := renderer.New(
		renderer.WithPictureProviderOptions(
			picture.WithSourcePath(`../test/testdata`),
			picture.WithDestinationPath(t.TempDir()),
		),
	)
	for _, q := range questions {
		if err := r.Render(b, source, q); err != nil {
			t.Fatal(err)
		}
		t.Log(b.String())
		b.Reset()
	}

	p, err = mdcParser.New(
		mdcParser.WithQuestionExtractor(mdcParser.QuestionExtractor(
			func(pc parser.Context, question ast.Node) {
				if err := r.Render(b, source, question); err != nil {
					t.Fatal(err)
				}
				t.Log(b.String())
				b.Reset()
			},
		)),
	)
	p.Parse(text.NewReader(source), parser.WithContext(ctx))
	if err != nil {
		t.Fatal(err)
	}
	// spew.Dump(questions)
	// t.Fatal("impl")
}
