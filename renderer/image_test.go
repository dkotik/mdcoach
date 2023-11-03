package renderer

import (
	"bytes"
	"testing"
)

func TestImageEncoding(t *testing.T) {
	var buf bytes.Buffer
	err := testMarkdown.Convert([]byte(`

![img](url "title")

  `), &buf)

	if err != nil {
		t.Fatal(err)
	}

	// spew.Dump(buf.String())
	// t.Fatal(`impl`)
}
