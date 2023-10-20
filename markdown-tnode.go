package mdcoach

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	blackfriday "github.com/russross/blackfriday/v2"
)

// Experimental handler for the blackfriday.Text node.

var smartyPants = blackfriday.NewSmartypantsRenderer(
	blackfriday.SmartypantsFractions |
		blackfriday.SmartypantsDashes |
		blackfriday.SmartypantsLatexDashes)

var reCitation = regexp.MustCompile(`\s*\(([^\)]{2,128})\)\<\/p\>\s*$`)

func catchCitation(w io.Writer, renderedQuote []byte) {
	if i := reCitation.FindSubmatchIndex(renderedQuote); i != nil {
		w.Write(renderedQuote[:i[0]])
		w.Write(renderedQuote[i[3]+1:])
		io.WriteString(w, `<cite>`)
		w.Write(renderedQuote[i[2]:i[3]])
		io.WriteString(w, `</cite>`)
		return
	}
	w.Write(renderedQuote)
}

func renderTextNode(w io.Writer, literal []byte) {
	var tmp bytes.Buffer
	// TODO: rewrite it as for loop to avoid unnecessary function nesting.
	// for start, l := 0, len(literal); start < l; {
	// 	i := emojiFilter.FindIndex(literal[start:])
	// 	if i == nil {
	// 		start = l
	// 	} else {
	// 		start = i[1]
	// 	}
	// 	escapeHTML(&tmp, literal[start:])
	// 	smartyPants.Process(w, tmp.Bytes())
	// }

	if i := emojiFilter.FindIndex(literal); i != nil {
		escapeHTML(&tmp, literal[:i[0]])
		smartyPants.Process(w, tmp.Bytes())
		class, text := emojiClass(string(literal[i[0]:i[1]]))
		fmt.Fprintf(w, `<span class="emote emote-%s">%s</span>`, class, text)
		renderTextNode(w, literal[i[1]:]) // keep going
	} else {
		escapeHTML(&tmp, literal)
		smartyPants.Process(w, tmp.Bytes())
	}
}
