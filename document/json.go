package document

import (
	"io"
	"strings"
)

const escapedCharsJSON = "&'<>\"\r\n"

// WriteEscapedJSON writes a string a JSON value with HTML escape codes applied. Single and double quotes are escaped twice. This preserves JSON objects embedded into HTML element body.
//
// code borrowed from `html` standard library package, simplified
func WriteEscapedJSON(w io.Writer, s string) error {
	i := strings.IndexAny(s, escapedCharsJSON)
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
			esc = "&amp;#39;"
		case '<':
			esc = "&lt;"
		case '>':
			esc = "&gt;"
		case '"':
			// "&#34;" is shorter than "&quot;".
			esc = "&amp;#34;" // TODO: I am adding a slash here - needs to be factored out.
		case '\r':
			esc = "&#13;"
		case '\n':
			esc = "&#10;"
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
