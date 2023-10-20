package mdcoach

import (
	"bytes"
	"io"

	"github.com/davecgh/go-spew/spew"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

// Iterator yields one slide for each [Iterator.Next] call with all of its notes and footnes.
type Iterator struct {
	Slide     *bytes.Buffer
	Notes     *bytes.Buffer
	Footnotes *bytes.Buffer
	renderer  renderer.Renderer
	cursor    ast.Node
	writer    io.Writer
}

func NewIterator(from ast.Node, r renderer.Renderer) *Iterator {
	if from.Kind() == ast.KindDocument {
		from = from.FirstChild()
	}
	return &Iterator{
		Slide:     &bytes.Buffer{},
		Notes:     &bytes.Buffer{},
		Footnotes: &bytes.Buffer{},
		renderer:  r,
		cursor:    from,
	}
}

func (i *Iterator) Next() bool {
	if i.cursor == nil {
		return false // iteration complete
	}

	i.Slide.Reset()
	i.Notes.Reset()
	i.Footnotes.Reset()
	i.writer = i.Slide

	switch i.cursor.Kind() {
	case ast.KindThematicBreak:
		// panic(i.cursor.Lines())
		spew.Dump(IsNotesThematicBreak(i.cursor), i.cursor.Attributes)
		fallthrough
	default:
		// spew.Dump(i.cursor)
		// spew.Dump(i.cursor.Lines())
		// log.Fatal(i.cursor.Kind(), i.cursor.FirstChild().Kind())
	}
	i.cursor = i.cursor.NextSibling()
	// for c := i.cursor.FirstChild(); c != nil; c = c.NextSibling() {
	//
	// }
	return true
}
