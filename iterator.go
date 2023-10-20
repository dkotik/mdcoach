package mdcoach

import (
	"bytes"
	"errors"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var EndWalk = errors.New("iteration interrupted")

type WalkFunc func(slide, notes, footnotes []byte) error

func Walk(source []byte, walk WalkFunc) (err error) {
	slides := &bytes.Buffer{}
	notes := &bytes.Buffer{}
	footnotes := &bytes.Buffer{}
	renderer := goldmark.DefaultRenderer()
	tree := DefaultParser().Parse(text.NewReader(source))
	w := slides

	flush := func() (err error) {
		err = walk(slides.Bytes(), notes.Bytes(), footnotes.Bytes())
		slides.Reset()
		notes.Reset()
		footnotes.Reset()
		w = slides
		if errors.Is(err, EndWalk) {
			return nil
		}
		return err
	}

	for c := tree.FirstChild(); c != nil; c = c.NextSibling() {
		switch c.Kind() {
		case ast.KindThematicBreak:
			if IsNotesThematicBreak(c) {
				if w != notes { // render repeated notes HR
					if err = renderer.Render(w, source, c); err != nil {
						return err
					}
				} else {
					w = notes
				}
				continue
			}
			if err = flush(); err != nil {
				return err
			}
		case ast.KindHeading:
			heading := c.(*ast.Heading)
			switch heading.Level {
			case 1, 2:
				if err = flush(); err != nil {
					return err
				}
			}
			fallthrough
		default:
			if err = renderer.Render(w, source, c); err != nil {
				return err
			}
		}
	}
	// if slides.Len() == 0 && notes.Len() == 0 &&
	return flush()
}

// // Iterator yields one slide for each [Iterator.Next] call with all of its notes and footnes.
// type Iterator struct {
// 	Slide     *bytes.Buffer
// 	Notes     *bytes.Buffer
// 	Footnotes *bytes.Buffer
// 	renderer  renderer.Renderer
// 	source    []byte
// 	cursor    ast.Node
// 	writer    io.Writer
// }
//
// func NewIterator(source []byte) *Iterator {
// 	return &Iterator{
// 		Slide:     &bytes.Buffer{},
// 		Notes:     &bytes.Buffer{},
// 		Footnotes: &bytes.Buffer{},
// 		source:    source,
// 		renderer:  goldmark.DefaultRenderer(),
// 		cursor:    DefaultParser().Parse(text.NewReader(source)).FirstChild(),
// 	}
// }
//
// func (i *Iterator) Next() error {
// 	if i.cursor == nil {
// 		fmt.Println("--- done ---")
// 		return io.EOF // iteration complete
// 	}
//
// 	i.Slide.Reset()
// 	i.Notes.Reset()
// 	i.Footnotes.Reset()
// 	i.writer = i.Slide
//
// 	// i.cursor = i.cursor.NextSibling()
// 	for c := i.cursor; c != nil; c = c.NextSibling() {
// 		i.cursor = c
// 		spew.Dump(string(c.Text(i.source)))
// 		switch c.Kind() {
// 		case ast.KindThematicBreak:
// 			// spew.Dump(IsNotesThematicBreak(c))
//
// 			// fallthrough
// 		default:
// 			if err := i.renderer.Render(i.writer, i.source, c); err != nil {
// 				panic(err) // TODO: handle error.
// 			}
// 			// for j := 0; j < c.Lines().Len(); j++ {
// 			// 	line := c.Lines().At(j)
// 			// 	fmt.Printf("$$$ %s\n", line.Value(i.source))
// 			// }
//
// 			// spew.Dump(i.cursor)
// 			// spew.Dump(i.cursor.Lines())
// 			// log.Fatal(i.cursor.Kind(), i.cursor.FirstChild().Kind())
// 		}
// 	}
// 	return io.EOF
// }
