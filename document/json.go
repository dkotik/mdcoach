package document

import (
	"io"
)

var escapeDataJSON = [256][]byte{
	'&':  []byte("&amp;"),
	'<':  []byte("&lt;"),
	'>':  []byte("&gt;"),
	'\'': []byte("&#39;"),
	'"':  []byte("\\&#34;"),
	'\n': []byte("&amp;#10;"),
	'\r': []byte("&amp;#13;"),
}

// WriteEscapedJSON writes a string a JSON value with HTML escape codes applied. Double quotes, new lines, and carrier returns are escaped twice. This preserves JSON objects embedded into HTML element body.
func WriteEscapedJSON(w io.Writer, s []byte) (err error) {
	var start, end int
	for end < len(s) {
		escSeq := escapeDataJSON[s[end]]
		if escSeq != nil {
			_, err = w.Write(s[start:end])
			if err != nil {
				return err
			}
			_, err = w.Write(escSeq)
			if err != nil {
				return err
			}
			start = end + 1
		}
		end++
	}
	if start < len(s) && end <= len(s) {
		_, err = w.Write(s[start:end])
	}
	return err
}
