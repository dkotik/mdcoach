package renderer

import (
	"bytes"
	"testing"
)

func TestCitationDetection(t *testing.T) {
	var buf bytes.Buffer
	err := newTestRenderer(t).Convert([]byte(`

> Some quote.
>
> Last part of it. (Author)

  `), &buf)

	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != `<blockquote><p>Some quote.</p>
<p>Last part of it.</p>
<cite>Author</cite>
</blockquote>
` {
		t.Log(buf.String())
		t.Fatal("expected result does not match")
	}
}

func TestAsideRendering(t *testing.T) {
	var buf bytes.Buffer
	err := newTestRenderer(t).Convert([]byte(`

> > p

> > > > p

  `), &buf)

	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != `<aside>
<p>p</p>
</aside>
<aside>
<p>p</p>
</aside>
` {
		t.Log(buf.String())
		t.Fatal("expected result does not match")
	}

	// spew.Dump(buf.String())
	// t.Fatal(`impl`)
}
