package test

import (
	"bytes"
	"os"
	"testing"

	"github.com/dkotik/mdcoach/document/review"
	mdcParser "github.com/dkotik/mdcoach/parser"
	"github.com/dkotik/mdcoach/picture"
	"github.com/dkotik/mdcoach/renderer"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func TestReviewQuestionExtraction(t *testing.T) {
	pictureProvider, err := picture.NewInternetProvider(
		picture.WithDestinationPath(t.TempDir()),
	)
	if err != nil {
		t.Fatal(err)
	}

	r, err := renderer.New(
		renderer.WithPictureProvider(&picture.SourceFilter{
			Provider: pictureProvider,
			// IsAllowed: func(source *picture.Source) (bool, error) {
			//   // trim output path from the source set
			//   source.Location = strings.TrimPrefix(source.Location, filepath.Dir(output)+"/")
			//   return true, nil
			// },
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	questions, err := review.New(
		review.WithRenderer(r),
	)
	if err != nil {
		t.Fatal(err)
	}

	source, err := os.ReadFile(`testdata/review-questions.md`)
	if err != nil {
		t.Fatal(err)
	}

	if err = questions.AddSource(source); err != nil {
		t.Fatal(err)
	}

	// spew.Dump(source)
	// panic("11")

	if questions.Len() != 13 { // one is duplicate
		t.Fatalf("unexpected number of questions detected, wanted %d, but found %d instead", 13, questions.Len())
	}
}

func TestReviewQuestionParsing(t *testing.T) {
	source, err := os.ReadFile(`testdata/review-questions.md`)
	if err != nil {
		t.Fatal(err)
	}

	ctx := parser.NewContext()
	p, err := mdcParser.New()
	if err != nil {
		t.Fatal(err)
	}
	p.Parse(text.NewReader(source), parser.WithContext(ctx))

	questions := mdcParser.QuestionListFromContext(ctx)
	if len(questions) != 14 {
		t.Fatalf("unexpected number of questions detected, wanted %d, but found %d instead", 14, len(questions))
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
