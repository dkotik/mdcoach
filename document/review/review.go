/*
Package review presents a printable test constructed using a set of questions extracted from Markdown list items.
*/
package review

import (
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/yuin/goldmark/renderer"
)

const preAllocateQuestions = 24

// Review holds the collected questions. Safe for concurrent use.
type Review struct {
	renderer renderer.Renderer
	logger   *slog.Logger

	mu        *sync.Mutex
	questions []string
	known     map[string]struct{}
}

func New(withOptions ...Option) (_ *Review, err error) {
	o := &options{}
	for _, option := range append(
		withOptions,
		withDefaultRenderer(),
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot create presentation review: %w", err)
		}
	}
	return &Review{
		renderer: o.renderer,
		logger:   slog.Default(),

		mu:        &sync.Mutex{},
		questions: make([]string, 0, 24),
		known:     make(map[string]struct{}),
	}, nil
}

func (r *Review) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.questions)
}

func (r *Review) Shuffle() {
	r.mu.Lock()
	defer r.mu.Unlock()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(r.questions), func(i, j int) {
		r.questions[i], r.questions[j] = r.questions[j], r.questions[i]
	})
}

// Assessment generates a test from questions contained in sources meta.
// func Assessment(e *Environment, sources ...string) error {
// 	r := &ArticleRenderer{&ImageRenderer{blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
// 		FootnoteReturnLinkContents: `⮭`,
// 		Flags:                      blackfriday.CommonHTMLFlags | blackfriday.FootnoteReturnLinks,
// 		// Flags:                      blackfriday.TOC,
// 	}), e}}
// 	questions := make([]template.HTML, 0)
// 	meta := make(map[string]interface{})
// 	cb := func(m map[string]interface{}) {
// 		questions = GatherQuestions(m, questions, e.Loader(func(m map[string]interface{}) {}), r)
// 	}
// 	d := NewDocument()
// 	d.Load(e.Loader(func(m map[string]interface{}) {
// 		meta = m
// 		cb(m)
// 	}), sources[0])
// 	for _, s := range sources[1:] {
// 		sub := NewDocument()
// 		sub.Load(e.Loader(cb), s)
// 	}
//
// 	temp, handle := e.Create(e.Output+`.html`, true)
// 	defer PDFManual(temp, filepath.Clean(e.Output))
// 	defer handle.Close()
