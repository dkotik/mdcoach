package renderer

import (
	"bytes"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestAsideRendering(t *testing.T) {
	var buf bytes.Buffer
	err := testMarkdown.Convert([]byte(`

> >
p

> > > > p

  `), &buf)

	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(buf.String())
	// t.Fatal(`impl`)
}
