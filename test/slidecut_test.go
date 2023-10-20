package test

import (
	"fmt"
	"testing"

	"github.com/dkotik/mdcoach/v1"
)

func TestSlideCut(t *testing.T) {
	// r := text.NewReader(
	// 	[]byte("\n\n\n12\n\n12\n\n***\n---\n"),
	// 	// testPresentation,
	// )
	// tree := mdcoach.DefaultParser().Parse(r)
	// renderer := renderer.NewRenderer()
	mdcoach.Walk(
		[]byte("# test\n\n\n## first\n\n\n12\n\n12\n\n***\nsome notes\n\n---\n\n sdfsdf sdf \n\n# final slide"),
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

	t.Fatal("impl")
}
