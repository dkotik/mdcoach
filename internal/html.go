package internal

import (
	"io"
	"strings"
)

const escapedChars = "&'<>\"\r"

// WriteEscapedHTML writes s to w while escaping HTML-sensitive characters.
// Its escaping logic is simplified from the standard library's html package.
func WriteEscapedHTML(w io.Writer, s string) error {
	i := strings.IndexAny(s, escapedChars)
	for i != -1 {
		if _, err := io.WriteString(w, s[:i]); err != nil {
			return err
		}
		var esc string
		switch s[i] {
		case '&':
			esc = "&amp;"
		case '\'':
			// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
			esc = "&#39;"
		case '<':
			esc = "&lt;"
		case '>':
			esc = "&gt;"
		case '"':
			// "&#34;" is shorter than "&quot;".
			esc = "&#34;"
		case '\r':
			esc = "&#13;"
		default:
			panic("unrecognized escape character")
		}
		s = s[i+1:]
		if _, err := io.WriteString(w, esc); err != nil {
			return err
		}
		i = strings.IndexAny(s, escapedChars)
	}
	_, err := io.WriteString(w, s)
	return err
}

func WriteCascadingStyleSheet(w io.Writer, css string) (err error) {
	if _, err = io.WriteString(w, `<style type="text/css">`); err != nil {
		return err
	}
	// TODO: proper escaping, &gt; `>` breaks css.
	// if err = WriteEscapedHTML(w, css); err != nil {
	// 	return err
	// }
	if _, err = io.Copy(w, strings.NewReader(css)); err != nil {
		return err
	}
	if _, err = io.WriteString(w, `</style>`); err != nil {
		return err
	}
	return nil
}
