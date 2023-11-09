/*
Package mdcoach renders Markdown files into HTML presentations with separate presentation slides and notes.
*/
package mdcoach

import (
	"io"
	"strings"

	"github.com/dkotik/mdcoach/document"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

func Compile(w io.Writer, tree ast.Node, markdown []byte, r renderer.Renderer) error {
	return Walk(tree, markdown, r, func(slide, notes, footnotes []byte) (err error) {
		// slide = bytes.ReplaceAll(slide, []byte("\n"), []byte("&#10;"))
		// notes = bytes.ReplaceAll(notes, []byte("\n"), []byte("&#10;"))
		// footnotes = bytes.ReplaceAll(footnotes, []byte("\n"), []byte("&#10;"))

		if _, err = io.Copy(w, strings.NewReader("\"")); err != nil {
			return err
		}
		if err = document.WriteEscapedJSON(w, slide); err != nil {
			return err
		}
		if _, err = io.Copy(w, strings.NewReader("\",\"")); err != nil {
			return err
		}
		if err = document.WriteEscapedJSON(w, notes); err != nil {
			return err
		}
		if _, err = io.Copy(w, strings.NewReader("\",\"")); err != nil {
			return err
		}
		if err = document.WriteEscapedJSON(w, footnotes); err != nil {
			return err
		}
		_, err = io.Copy(w, strings.NewReader("\","))
		return err
	})
}
