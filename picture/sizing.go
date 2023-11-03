package picture

import (
	"sort"
	// "golang.org/x/image/webp"
	// "github.com/kolesa-team/go-webp/encoder"
	// "github.com/kolesa-team/go-webp/webp"
	// github.com/nickalie/go-webpbin
	// github.com/chai2010/webp
	// https://github.com/chai2010/webp/issues/51
	//
	// Alternative for resizing:
	// github.com/disintegration/imaging
)

// Sizing represets desired image parameters.
type Sizing struct {
	Width   int
	Height  int
	Quality int
}

func SortSizingsInDescendingOrder(sizings []Sizing) {
	area := make([]int, len(sizings))
	for i, s := range sizings {
		area[i] = s.Width * s.Height
	}
	sort.Slice(sizings, func(i, j int) bool {
		return area[i] > area[j]
	})
}
