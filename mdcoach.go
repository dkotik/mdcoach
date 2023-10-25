/*
Package mdcoach renders Markdown files into HTML presentations with separate presentation slides and notes.
*/
package mdcoach

import (
	"bytes"
	"io"
	"strings"

	"github.com/dkotik/mdcoach/document"
)

func Compile(w io.Writer, markdown []byte) error {
	return Walk(markdown, func(slide, notes, footnotes []byte) (err error) {
		slide = bytes.ReplaceAll(slide, []byte("\n"), []byte("&#10;"))
		notes = bytes.ReplaceAll(notes, []byte("\n"), []byte("&#10;"))
		footnotes = bytes.ReplaceAll(footnotes, []byte("\n"), []byte("&#10;"))

		if _, err = io.Copy(w, strings.NewReader("\"")); err != nil {
			return err
		}
		if err = document.WriteEscapedJSON(w, string(slide)); err != nil {
			return err
		}
		if _, err = io.Copy(w, strings.NewReader("\",\"")); err != nil {
			return err
		}
		if err = document.WriteEscapedJSON(w, string(notes)); err != nil {
			return err
		}
		if _, err = io.Copy(w, strings.NewReader("\",\"")); err != nil {
			return err
		}
		if err = document.WriteEscapedJSON(w, string(footnotes)); err != nil {
			return err
		}
		_, err = io.Copy(w, strings.NewReader("\","))
		return err
	})
}
