package mdcoach

import (
	"bytes"
	"errors"
	"io"

	"github.com/dkotik/mdcoach/parser"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
)

var EndWalk = errors.New("iteration interrupted")

type WalkFunc func(slide, notes, footnotes []byte) error

type iterator struct {
	callback  WalkFunc
	slide     *bytes.Buffer
	notes     *bytes.Buffer
	footnotes *bytes.Buffer
	w         io.Writer
	renderer  renderer.Renderer
}

func newIterator(callback WalkFunc, r renderer.Renderer) *iterator {
	i := &iterator{
		callback:  callback,
		slide:     &bytes.Buffer{},
		notes:     &bytes.Buffer{},
		footnotes: &bytes.Buffer{},
		renderer:  r,
	}
	i.w = i.slide
	return i
}

func (i *iterator) Render(tree ast.Node, source []byte) (err error) {
	for n := tree.FirstChild(); n != nil; n = n.NextSibling() {
		switch n.Kind() {
		case ast.KindThematicBreak:
			if parser.IsNotesThematicBreak(n) {
				if i.w == i.notes { // render repeated notes HR
					if err = i.renderer.Render(i.w, source, n); err != nil {
						return err
					}
				} else {
					i.w = i.notes
				}
				continue
			}
			if err = i.Flush(); err != nil {
				return err
			}
		case ast.KindHeading:
			heading := n.(*ast.Heading)
			switch heading.Level {
			case 1, 2:
				if err = i.Flush(); err != nil {
					return err
				}
			}
			fallthrough
		default:
			if err = i.renderer.Render(i.w, source, n); err != nil {
				return err
			}
		}
	}
	return i.Flush()
}

func (i *iterator) Flush() error {
	if i.slide.Len()+i.notes.Len()+i.footnotes.Len() == 0 {
		return nil // skip entirely empty slides
	}
	// TODO: render footnotes.
	err := i.callback(
		i.slide.Bytes(),
		i.notes.Bytes(),
		i.footnotes.Bytes(),
	)
	i.slide.Reset()
	i.notes.Reset()
	i.footnotes.Reset()
	i.w = i.slide
	if !errors.Is(err, EndWalk) {
		return err
	}
	return nil
}

func Walk(source []byte, renderer renderer.Renderer, walk WalkFunc) (err error) {
	p, err := parser.New()
	if err != nil {
		return err
	}
	return newIterator(walk, renderer).Render(
		p.Parse(text.NewReader(source)),
		source,
	)
}
