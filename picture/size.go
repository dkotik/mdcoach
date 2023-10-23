package picture

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

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

// Write saves provided image to disk.
func Write(path string, m *image.Image, s *Sizing) error {
	handle, err := os.Create(path)
	if err != nil {
		return err
	}
	defer handle.Close()
	t := resize.Thumbnail(s.Width, s.Height, *m, resize.Lanczos3)
	switch filepath.Ext(path) {
	// case `.webp`:
	// 	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, float32(s.Quality))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return webp.Encode(handle, *m, options)
	// 	// return webp.Encode(handle, *m, &webp.Options{
	// 	// 	Lossless: false,
	// 	// 	Quality:  float32(s.Quality),
	// 	// })
	default:
		return jpeg.Encode(handle, t, &jpeg.Options{
			Quality: int(s.Quality),
		})
	}
}
