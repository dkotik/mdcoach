package test

import (
	"fmt"
	"testing"

	"github.com/dkotik/mdcoach"
	"github.com/dkotik/mdcoach/parser"
	"github.com/yuin/goldmark/text"
)

func TestSlideCut(t *testing.T) {
	p, err := parser.New()
	if err != nil {
		t.Fatal(err)
	}
	markdown := []byte("# test\n\n\n## first\n\n\n12\n\n12\n\n***\nsome notes\n\n---\n\n sdfsdf sdf \n\n# final slide")
	// r := text.NewReader(
	// 	[]byte("\n\n\n12\n\n12\n\n***\n---\n"),
	// 	// testPresentation,
	// )
	// tree := mdcoach.DefaultParser().Parse(r)
	// renderer := renderer.NewRenderer()
	mdcoach.Walk(
		p.Parse(text.NewReader(markdown)),
		markdown,
		NewTestRenderer(t),
		func(s, n, ft []byte) error {
			fmt.Println("-------------------------------")
			fmt.Println("slide:", string(s))
			fmt.Println("notes:", string(n))
			fmt.Println("footnotes:", string(ft))
			return nil
		},
	)
	// iter := mdcoach.NewIterator(
	//
	// )
	//
	// for ; err := iter.Next(); err != io.EOF {
	// 	fmt.Println("one  ################")
	// 	fmt.Println(iter.Slide.String())
	// }

	// t.Fatal("impl")
}
