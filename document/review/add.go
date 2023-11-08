package review

import (
	"bytes"

	mdcParser "github.com/dkotik/mdcoach/parser"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func (r *Review) AddSource(markdown []byte) (err error) {
	questions := make([]ast.Node, 0, preAllocateQuestions)
	p, err := mdcParser.New(
		mdcParser.WithQuestionExtractor(
			func(_ parser.Context, question ast.Node) {
				questions = append(questions, question)
			},
		),
	)
	if err != nil {
		return err
	}
	_ = p.Parse(text.NewReader(markdown), parser.WithContext(parser.NewContext()))

	b := &bytes.Buffer{}
	for _, q := range questions {
		for n := q.FirstChild(); n != nil; n = n.NextSibling() {
			if err = r.renderer.Render(b, markdown, n); err != nil {
				return err
			}
		}
		if r.AddQuestion(b.String()) {
			r.logger.Info(string(q.FirstChild().Text(markdown)))
		}
		b.Reset()
	}
	return nil
}

func (r *Review) AddQuestion(question string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.known[question]; ok {
		return false
	}
	r.known[question] = struct{}{}
	r.questions = append(r.questions, question)
	return true
}
