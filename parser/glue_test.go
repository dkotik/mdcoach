package parser

import (
	"os"
	"testing"

	"github.com/yuin/goldmark/text"
)

func TestFrontMatterDetection(t *testing.T) {
	markdown, err := os.ReadFile("../test/testdata/another-presentation.md")
	if err != nil {
		t.Fatal(err)
	}
	frontMatterEnds := locateFrontMatterEnd(text.NewReader(markdown))
	if frontMatterEnds == 0 {
		t.Fatal("frontmatter was not detected")
	}
	// spew.Dump(string(markdown[frontMatterEnds:]))
	// t.Fatal("implt")
}
