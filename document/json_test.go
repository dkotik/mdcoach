package document

import (
	"bytes"
	"testing"
)

func TestJSONEncoding(t *testing.T) {
	b := &bytes.Buffer{}
	// 	WriteEscapedJSON(b, []byte(`
	// `))

	// 	spew.Dump([]byte(`"&lt;h1&gt;Top Heading&lt;/h1&gt; dfs
	// newline"` + "\n"))
	// 	spew.Dump('\n')

	if err := WriteEscapedJSON(b, []byte(`"&lt;h1&gt;Top Heading&lt;/h1&gt; dfs
newline"`+"\n")); err != nil {
		t.Fatal("could not encode:", err)
	}
	for _, c := range b.Bytes() {
		if c == '\n' {
			t.Fatal("found a newline char in JSON encoded string")
		}
	}
}
