package mdcoach

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	blackfriday "github.com/russross/blackfriday/v2"
)

var reImageAdjuster = regexp.MustCompile(`\<img src\=\"([^\"]+)\.thumb\.jpg\" alt\=\"([^\"]*)\" (title\=\"([^\"]*)\" )?\/\>`)

func isCutNode(node *blackfriday.Node) bool {
	switch node.Type {
	case blackfriday.Heading:
		return node.HeadingData.Level <= 2
	case blackfriday.Paragraph: // figure
		return node.FirstChild.Type == blackfriday.Text &&
			len(node.FirstChild.Literal) == 0 &&
			node.FirstChild.Next.Type == blackfriday.Image &&
			node.FirstChild.Next.Next == nil
	case blackfriday.BlockQuote: // regular BlockQuote, but not the aside element
		return !(node.FirstChild.Type == blackfriday.BlockQuote && node.FirstChild.Next == nil)
	case blackfriday.HorizontalRule, blackfriday.CodeBlock:
		return true
	}
	return false
}

// Presentation creates slides from document.
func Presentation(e *Environment, sources ...string) error {
	r := &ArticleRenderer{&ImageRenderer{blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		// HeadingLevelOffset: 1,
		// AbsolutePrefix:             `|urlprefix|`,
		FootnoteReturnLinkContents: `⮭`,
		Flags:                      blackfriday.CommonHTMLFlags | blackfriday.FootnoteReturnLinks,
		// Flags:                      blackfriday.CommonHTMLFlags | blackfriday.TOC | blackfriday.FootnoteReturnLinks,
	}), e}}

	init := true
	meta := make(map[string]interface{})
	questions := make([]template.HTML, 0)
	cb := func(m map[string]interface{}) {
		if init {
			meta = m
			init = false
		}
		questions = GatherQuestions(m, questions, e.Loader(func(m map[string]interface{}) {}), r)
	}
	d := NewDocument()
	d.Load(e.Loader(cb), sources...)

	temp, fnotes := e.Create(e.Output+`.html`, true)
	pdfpath := path.Join(path.Dir(e.Output), `notes-`+strings.TrimSuffix(path.Base(e.Output), `.html`)+`.pdf`)
	defer PDFManual(temp, pdfpath)
	defer fnotes.Close()
	_, fslides := e.Create(e.Output, false)
	e.SaveMeta(filepath.Base(e.Output), pdfpath, meta)

	// Rendering HEAD for slides and notes:
	head, err := template.New("head").Parse(tmplHead)
	if err != nil {
		return err
	}
	meta[`stylesheets`] = []string{`.cache/presentation.css`, `.cache/pygments.css`, `.cache/emote.css`}
	meta[`scripts`] = []string{`.cache/vue.min.js`}
	err = head.Execute(fslides, meta)
	if err != nil {
		return err
	}
	meta[`scripts`] = []string{}
	meta[`stylesheets`] = []string{`notes.css`, `pygments.css`, `emote.css`}
	fslides.Write([]byte(`<script type="text/javascript">var slides=['`))
	err = head.Execute(fnotes, meta)
	if err != nil {
		return err
	}

	total := 0
	totalFootnotes := 0
	ba, bm, bn := bytes.NewBuffer(nil), bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	flush := func() {
		if ba.Len() > 0 || bm.Len() > 0 {
			total++
			fmt.Fprintf(fnotes, `<div class="slidebr"><mark>%d</mark></div>`, total)
			defer io.WriteString(fslides, `','`)
		}
		fnotes.Write(IOcleanHTML(ba.Bytes()))
		fnotes.Write(IOcleanHTML(bm.Bytes()))
		if bn.Len() > 0 {
			fmt.Fprintf(fnotes, `<ol start="%d" class="footnotes">`, totalFootnotes)
			fnotes.Write(IOcleanHTML(bn.Bytes()))
			bn.Reset()
			io.WriteString(fnotes, `</ol>`)
		}
		template.JSEscape(fslides, reImageAdjuster.ReplaceAll(bm.Bytes(), []byte(`<img src=".cache/$1.webp" alt="$2" title="$3" />`)))
		template.JSEscape(fslides, reImageAdjuster.ReplaceAll(ba.Bytes(), []byte(`<img src=".cache/$1.webp" alt="$2" title="$3" />`)))
		ba.Reset()
		bm.Reset()
	}

	captureFootnotes := func(w io.Writer, node *blackfriday.Node) {
		node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
			if node.Type == blackfriday.Link && entering && node.Footnote != nil {
				totalFootnotes++
				r.Walk(bn, node.Footnote)
			}
			return r.RenderNode(w, node, entering)
		})
	}

	current := d.Tree.FirstChild
	for current != nil {
		if isCutNode(current) {
			flush()
		}
		switch current.Type {
		case blackfriday.BlockQuote:
			if current.FirstChild.Type == blackfriday.BlockQuote && current.FirstChild.Next == nil {
				// detected nested block quote -> must be aside block
				if current.FirstChild.FirstChild == nil {
					// detected an empty aside element, write the remaining elements to notes.
					io.WriteString(bn, `<aside class="notes">`)
					current = current.Next
					for current != nil && !isCutNode(current) {
						r.Walk(bn, current)
						current = current.Next
					}
					io.WriteString(bn, `</aside>`)
					continue
				}
				if current.Next == nil || isCutNode(current.Next) {
					if current.Prev == nil || isCutNode(current.Prev) {
						io.WriteString(bm, `<aside class="splash">`)
					} else {
						io.WriteString(bm, `<aside class="right">`)
					}
				} else {
					io.WriteString(bm, `<aside class="left">`)
				}
				cc := current.FirstChild.FirstChild
				for cc != nil {
					r.Walk(bm, cc)
					cc = cc.Next
				}
				io.WriteString(bm, `</aside>`)
			} else {
				captureFootnotes(ba, current)
			}
		case blackfriday.HorizontalRule: // drop
		case blackfriday.CodeBlock:
			if len(current.CodeBlockData.Info) == 0 {
				// Unmarked code block - taken as comment!
				fmt.Fprintf(bn, `<aside class="notes">%s</aside>`, current.Literal) // TODO: render this with all the guts?
				current = current.Next                                              // TODO: hacky, not ideal
				continue
			}
			fallthrough
		default:
			captureFootnotes(ba, current)
		}
		current = current.Next
	}

	flush()
	fmt.Fprint(fnotes, `<div class="slidebr"><mark>&#8943;</mark></div>`)
	err = template.Must(template.New("presentation").Parse(tmplPresentation)).Execute(fslides, nil)
	if err != nil {
		return err
	}
	io.WriteString(fslides, tmplFoot)
	err = template.Must(template.New("review").Parse(tmplReview)).Execute(fnotes, map[string]interface{}{"Questions": questions})
	if err != nil {
		return err
	}
	io.WriteString(fnotes, tmplFoot)
	return nil
}
