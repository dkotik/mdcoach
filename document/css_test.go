package document

import (
	"bytes"
	"testing"
)

func TestEmbedCascadingStyleSheet(t *testing.T) {
	t.Skip("proper encoding is not implemented yet") // TODO: fix.

	b := &bytes.Buffer{}
	err := WriteCascadingStyleSheet(b, `body{color:"red";}`)
	if err != nil {
		t.Fatal(err)
	}
	if b.String() != `<style type="text/css">body{color:&#34;red&#34;;}</style>` {
		t.Fatal("encoding mismatch:", b.String())
	}
}
