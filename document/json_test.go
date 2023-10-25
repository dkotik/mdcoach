package document

// func TestJSONEncoding(t *testing.T) {
// 	b := &bytes.Buffer{}
// 	if err := WriteEscapedJSON(b, `"&lt;h1&gt;Top Heading&lt;/h1&gt;
// "`); err != nil {
// 		t.Fatal("could not encode:", err)
// 	}
// 	for _, c := range b.Bytes() {
// 		if c == '\n' {
// 			t.Fatal("found a newline char in JSON encoded string")
// 		}
// 	}
// }
