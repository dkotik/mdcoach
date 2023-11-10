package parser

import (
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"

	"github.com/yuin/goldmark/text"
)

var reFrontMatterMarker = regexp.MustCompile(`^[\-\+]{3,}\n$`)

func locateFrontMatterEnd(r text.Reader) int {
	segment, _, ok := r.SkipBlankLines()
	if !ok {
		return 0
	}
	line := r.Value(segment)
	if !reFrontMatterMarker.Match(line) {
		// panic(string(line))
		return 0
	}
	// spew.Dump(segment, pos)
	r.Advance(segment.Stop)
	for i := 0; i < 2000; i++ {
		matchingLine, _ := r.PeekLine()
		if bytes.Equal(line, matchingLine) {
			_, segment = r.Position()
			// spew.Dump(segment.Stop)
			return segment.Stop
		}
		r.AdvanceLine()
	}
	return 0
}

// Glue builds a single Markdown document out of a set of files. Metadata is taken from the first file only.
func Glue(files ...string) ([]byte, error) {
	if len(files) < 1 {
		return nil, errors.New("at least one Markdown file is required")
	}
	b := &bytes.Buffer{}
	first, err := os.ReadFile(files[0])
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(b, bytes.NewReader(first)); err != nil {
		return nil, err
	}
	for _, file := range files[1:] {
		additional, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		if _, err = b.Write([]byte("\n\n")); err != nil {
			return nil, err
		}
		cutOff := locateFrontMatterEnd(text.NewReader(additional))
		if _, err = io.Copy(b, bytes.NewReader(additional[cutOff:])); err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}
