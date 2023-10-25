package picture

import (
	"context"
	"fmt"
	"image"
	"sort"

	"github.com/nfnt/resize"
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
	Width   uint
	Height  uint
	Quality uint
}

type SizingTarget struct {
	Sizing
	Destination string
}

type SizingRequest struct {
	Image   image.Image
	Outputs []SizingTarget
}

type Resizer struct {
	// TODO: confirm overwrite.
	encoder Encoder
}

func (r *Resizer) Execute(ctx context.Context, s *SizingRequest) error {
	for _, target := range s.Outputs {
		select { // check if context is alive
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		smaller := resize.Thumbnail(
			target.Width,
			target.Height,
			s.Image,
			resize.Lanczos3,
		)
		select { // check if context is alive
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := r.encoder.EncodeImage(&target, smaller); err != nil {
			return fmt.Errorf("could not save image %q: %w", target.Destination, err)
		}
		s.Image = smaller
	}
	return nil
}

func SortSizingsInDescendingOrder(sizings []Sizing) {
	area := make([]uint, len(sizings))
	for i, s := range sizings {
		area[i] = s.Width * s.Height
	}
	sort.Slice(sizings, func(i, j int) bool {
		return area[i] > area[j]
	})
}
