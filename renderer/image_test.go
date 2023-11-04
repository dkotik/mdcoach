package renderer

import (
	"bytes"
	"testing"
)

func TestImageEncoding(t *testing.T) {
	var buf bytes.Buffer
	err := newTestRenderer(t).Convert([]byte(`

![img](../picture/testdata/notfound.jpg "title")

  `), &buf)

	if err != nil {
		t.Fatal(err)
	}

	// TODO: add proper test cases.
	// spew.Dump(buf.String())
	// t.Fatal(`impl`)
}
