package mdcoach

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	blackfriday "github.com/russross/blackfriday/v2"
)

// PDFManual generates reading booklet with a table of contents.
func PDFManual(source string, target string) error {
	if strings.HasSuffix(target, `.html`) {
		target = target[:len(target)-5] + `.pdf`
	}
	// -p is important for ol! https://github.com/Kozea/WeasyPrint/issues/398
	return Exec(`weasyprint`, source, target, `-p`)
}

// Assessment generates a test from questions contained in sources meta.
func Assessment(e *Environment, sources ...string) error {
	r := &ArticleRenderer{&ImageRenderer{blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		FootnoteReturnLinkContents: `⮭`,
		Flags:                      blackfriday.CommonHTMLFlags | blackfriday.FootnoteReturnLinks,
		// Flags:                      blackfriday.TOC,
	}), e}}
	questions := make([]template.HTML, 0)
	meta := make(map[string]interface{})
	cb := func(m map[string]interface{}) {
		questions = GatherQuestions(m, questions, e.Loader(func(m map[string]interface{}) {}), r)
	}
	d := NewDocument()
	d.Load(e.Loader(func(m map[string]interface{}) {
		meta = m
		cb(m)
	}), sources[0])
	for _, s := range sources[1:] {
		sub := NewDocument()
		sub.Load(e.Loader(cb), s)
	}

	temp, handle := e.Create(e.Output+`.html`, true)
	defer PDFManual(temp, filepath.Clean(e.Output))
	defer handle.Close()

	meta[`questions`] = questions
	if v, ok := meta[`shuffle`]; ok && fmt.Sprintf(`%s`, v) == `true` {
		log.Println(`Shuffling test questions.`)
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(questions), func(i, j int) {
			questions[i], questions[j] = questions[j], questions[i]
		})
	}
	max := len(questions)
	if max == 0 {
		log.Fatal(`No questions were found in specified markdown documents.`)
	}
	switch t := meta[`eliminate`].(type) {
	case int:
		if t >= max {
			max = 0
		} else {
			max = t
		}
		log.Printf(`Eliminating %d questions. %d questions remain.`, t, max)
	case string:
		if ok, _ := regexp.MatchString(`^\d+\%$`, t); ok {
			el, _ := strconv.Atoi(t[:len(t)-1])
			el = el/max + 1
			if el > max {
				max = 0
			} else {
				max -= el
			}
			log.Printf(`Eliminating %s of questions. %d questions remain.`, t, max)
		} else {
			log.Printf(`Eliminate value "%s" is neither a number nor a percentage.`, t)
		}
	case nil:
	default:
		log.Printf(`Eliminate value "%s" is neither a number nor a percentage.`, t)
	}
	questions = questions[:max]
	meta[`missed`] = []int{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	grades := make([]int, 0)
	max *= 3 // max 3 points per question
	for _, v := range meta[`missed`].([]int) {
		grades = append(grades, (max-v)*100/max)
	}
	meta[`grades`] = grades
	meta[`stylesheets`] = []string{`assessment.css`}

	// Render body of the assessment.
	w := bytes.NewBuffer(nil)
	d.Render(w, r)
	meta[`content`] = template.HTML(IOcleanHTML(w.Bytes()))
	w.Reset()
	err := template.Must(template.New("head").Parse(tmplHead)).Execute(handle, meta)
	if err != nil {
		return err
	}
	err = template.Must(template.New("assessment").Parse(tmplAssessment)).Execute(handle, meta)
	if err != nil {
		return err
	}
	io.WriteString(handle, tmplFoot)
	return nil
}

// Paper generates a printable document.
func Paper(e *Environment, styleSheet string, sources ...string) error {
	meta := make(map[string]interface{})
	d := NewDocument()
	d.Load(e.Loader(func(m map[string]interface{}) { meta = m }), sources[0])
	if len(sources) > 1 {
		d.Load(e.Loader(func(m map[string]interface{}) {}), sources[1:]...)
	}

	tmp, handle := e.Create(e.Output+`.html`, true)
	e.SaveMeta(filepath.Base(e.Output), ``, meta)
	defer PDFManual(tmp, filepath.Clean(e.Output))
	meta[`stylesheets`] = []string{styleSheet, `pygments.css`, `emote.css`}
	err := template.Must(template.New("head").Parse(tmplHead)).Execute(handle, meta)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(nil)
	r := &ArticleRenderer{&ImageRenderer{blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		FootnoteReturnLinkContents: `⮭`,
		Flags:                      blackfriday.CommonHTMLFlags | blackfriday.FootnoteReturnLinks,
		// Flags:                    | blackfriday.TOC,
	}), e}}
	d.Render(b, r)
	handle.Write(IOcleanHTML(b.Bytes()))

	err = template.Must(template.New("review").Parse(tmplReview)).Execute(handle, meta)
	if err != nil {
		return err
	}
	io.WriteString(handle, tmplFoot)
	return nil
}
