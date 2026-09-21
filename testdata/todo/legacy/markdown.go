package mdcoach

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"path"
	"regexp"
	"strings"

	blackfriday "github.com/russross/blackfriday/v2"
)

// GatherQuestions fill a stack with renderer questions.
func GatherQuestions(m map[string]interface{}, stack []template.HTML,
	l func(string) ([]byte, error), r blackfriday.Renderer) []template.HTML {
	if v, ok := m[`questions`].([]interface{}); ok {
		for _, q := range v {
			d := NewDocument()
			d.Parse([]byte(fmt.Sprintf(`%v`, q)), path.Dir(m[`source`].(string)), l)
			b := bytes.NewBuffer(nil)
			d.Render(b, r)
			stack = append(stack, template.HTML(IOcleanHTML(b.Bytes())))
			// log.Println(`caught question`, b.String())
		}
	}
	return stack
}

// ArticleRenderer produces HTML from markdown tree nodes.
type ArticleRenderer struct {
	blackfriday.Renderer
}

// RenderNode implements sugar over standard markdown HTML rendering.
func (r *ArticleRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	// if entering {
	// 	fmt.Fprintf(w, "\n<!-- [%s] -->", node.Type)
	// } else {
	// 	defer fmt.Fprintf(w, "<!-- [/%s] -->\n", node.Type)
	// }
	if entering {
		switch node.Type {
		case blackfriday.Text: // TODO: I did the Text node wrong in the past!
			if node.Parent.Type == blackfriday.Link {
				escLink(w, node.Literal)
			} else {
				// fmt.Printf(":: %s\n", node.Literal)
				renderTextNode(w, node.Literal)
			}
			return blackfriday.GoToNext
		// case blackfriday.Emph:
		case blackfriday.Paragraph:
			if node.FirstChild.Type == blackfriday.Text &&
				len(node.FirstChild.Literal) == 0 &&
				node.FirstChild.Next.Type == blackfriday.Image &&
				node.FirstChild.Next.Next == nil {
				io.WriteString(w, `<figure>`)
				r.Walk(w, node.FirstChild.Next)
				if len(node.FirstChild.Next.LinkData.Title) > 0 {
					fmt.Fprintf(w, `<figcaption>%s</figcaption>`, node.FirstChild.Next.LinkData.Title)
				}
				io.WriteString(w, `</figure>`)
				return blackfriday.SkipChildren
			}
		case blackfriday.List:
			todoList := false
			reTask := regexp.MustCompile(`^\s*\[( |x|X|\-)\]\s+`)
			current := node.FirstChild // go through children
			for current != nil {
				if current.Type == blackfriday.Item && current.FirstChild.Type == blackfriday.Paragraph && current.FirstChild.FirstChild.Type == blackfriday.Text {
					if m := reTask.FindSubmatch(current.FirstChild.FirstChild.Literal); m != nil {
						todoList = true
						current.FirstChild.FirstChild.Literal = current.FirstChild.FirstChild.Literal[len(m[0]):]
						if bytes.Equal(m[1], []byte{' '}) {
							current.FirstChild.FirstChild.InsertBefore(&blackfriday.Node{
								Type: blackfriday.HTMLSpan, Literal: []byte(`<span class="emote emote-checkbox">&nbsp;</span>`)})
						} else {
							current.FirstChild.FirstChild.InsertBefore(&blackfriday.Node{
								Type: blackfriday.HTMLSpan, Literal: []byte(`<span class="emote emote-checkboxchecked">&nbsp;</span> `)})
						}
					}
				}
				current = current.Next
			}
			if todoList {
				if node.ListFlags&blackfriday.ListTypeOrdered != 0 {
					io.WriteString(w, `<ol class="todo">`)
				} else {
					io.WriteString(w, `<ul class="todo">`)
				}
				return blackfriday.GoToNext
			}
		case blackfriday.CodeBlock:
			switch strings.ToLower(string(node.CodeBlockData.Info)) {
			case `greek`, `hebrew`, `bible`:
				fmt.Fprintf(w, `<blockquote class="language-%s">`, node.CodeBlockData.Info)
				d := NewDocument()
				d.Parse(regexp.MustCompile(`\s*([0-9][0-9\:\.]*)\.? \s*`).ReplaceAll(node.Literal, []byte(`<sup>$1</sup>`)), ``, nil)
				var tmp bytes.Buffer
				d.Render(&tmp, r)
				catchCitation(w, tmp.Bytes())
				w.Write([]byte(`</blockquote>`))
				return blackfriday.GoToNext
			}
		case blackfriday.BlockQuote:
			if node.FirstChild.Type == blackfriday.BlockQuote && node.FirstChild.Next == nil {
				// detected nested block quote -> must be aside block
				io.WriteString(w, `<aside>`)
				current := node.FirstChild.FirstChild
				for current != nil {
					r.Walk(w, current)
					current = current.Next
				}
				io.WriteString(w, `</aside>`)
				return blackfriday.SkipChildren
			}
			io.WriteString(w, `<blockquote>`)
			current := node.FirstChild
			for current != nil {
				if current.Type == blackfriday.Paragraph { // fish for citations
					var tmp bytes.Buffer
					r.Walk(&tmp, current)
					catchCitation(w, tmp.Bytes())
				} else {
					r.Walk(w, current)
				}
				current = current.Next
			}
			io.WriteString(w, `</blockquote>`)
			return blackfriday.SkipChildren
		case blackfriday.HorizontalRule:
			if node.Prev != nil && node.Prev.Type == blackfriday.HorizontalRule {
				node.Prev.Unlink() // is this wise? modifying the tree while rendering
				io.WriteString(w, `<hr class="double" />`)
				return blackfriday.SkipChildren
			}
		}
	}
	return r.Renderer.RenderNode(w, node, entering)
}

// Walk runs down a node tree, rendering each node.
func (r *ArticleRenderer) Walk(w io.Writer, node *blackfriday.Node) {
	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return r.RenderNode(w, node, entering)
	})
}
