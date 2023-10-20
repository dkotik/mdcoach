package test

import (
	"fmt"
	"testing"

	"github.com/dkotik/mdcoach/v1"

	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
)

func TestSlideCut(t *testing.T) {
	r := text.NewReader(
		[]byte("\n\n\n12\n\n12\n\n***\n---\n"),
		// testPresentation,
	)
	tree := mdcoach.DefaultParser().Parse(r)
	// spew.Dump(tree)
	//
	// // for c := tree.FirstChild(); c != nil; c = c.NextSibling() {
	// for c := tree; c != nil; c = c.NextSibling() {
	// 	spew.Dump(c)
	// 	for a := c.FirstChild(); a != nil; a = a.NextSibling() {
	// 		fmt.Println("one  ################")
	// 	}
	// }
	//
	// if err := goldmark.Convert([]byte("\n\n\n12\n\n12\n\n***\n"), os.Stdout); err != nil {
	// 	panic(err)
	// }
	//
	renderer := renderer.NewRenderer()
	iter := mdcoach.NewIterator(tree, renderer)

	for iter.Next() {
		fmt.Println("one  ################")
	}

	t.Fatal("impl")
}
