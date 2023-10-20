package mdcoach

import (
	"testing"
)

func TestHandoutGenerator(t *testing.T) {
	e := NewEnvironment(`/tmp/mdcoachtest/handout.html`)
	defer e.Close()
	// Debug settings.
	e.Debug = true
	e.Overwrite = true
	// // case `syllabus`:
	// Paper(e, doc, []string{`html/handout.html`, `html/meta.tmpl`}, []string{`html/foot.html`})
	// // case `handout`:
	// Paper(e, doc, []string{`html/handout.html`, `html/meta.tmpl`}, []string{`html/review.html`, `html/foot.html`})
	// // case `book`:
	// // case `assessment`:
	// default:
	// e.Output = `test/handout.pdf`
	// Paper(e, `handout.css`, []string{`html/book-f.tmpl`}, `assets/test-demo.md`, `assets/test-another.md`, `assets/test-handout.md`)
	// Assessment(e, `assets/test-demo.md`, `assets/test-another.md`, `assets/test-handout.md`)
	Presentation(e, `assets/test-demo.md`, `assets/test-another.md`, `assets/test-handout.md`)
	// e.MakeIndex(`test`)
	// t.Fail() // stops Makefile recipe execution
}
