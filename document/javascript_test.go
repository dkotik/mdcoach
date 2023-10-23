package document

import (
	"bytes"
	"testing"
)

func TestEmbedJavascriptModuleES6(t *testing.T) {
	b := &bytes.Buffer{}
	err := WriteJavascriptModuleES6(b, `alert("!");`)
	if err != nil {
		t.Fatal(err)
	}
	if b.String() != `<script type="module">alert(&#34;!&#34;);</script>` {
		t.Fatal("encoding mismatch:", b.String())
	}
}
