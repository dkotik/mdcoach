package renderer

import (
	"bytes"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestImageEncoding(t *testing.T) {
	var buf bytes.Buffer
	err := testMarkdown.Convert([]byte(`

![img](../picture/testdata/notfound.jpg "title")

  `), &buf)

	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(buf.String())
	t.Fatal(`impl`)
}
