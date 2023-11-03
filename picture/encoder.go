package picture

import (
	"image"
	"image/jpeg"
	"io"
)

var _ Encoder = (*JPEGEncoder)(nil)

type Encoder interface {
	EncodeImage(w io.Writer, m image.Image, quality int) error
}

type JPEGEncoder struct{}

func (j *JPEGEncoder) EncodeImage(
	w io.Writer,
	m image.Image,
	quality int,
) error {
	return jpeg.Encode(w, m, &jpeg.Options{Quality: quality})
}

// // Write saves provided image to disk.
// func Write(path string, m image.Image, s *Sizing) error {
// 	w, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer w.Close()
// 	t := resize.Thumbnail(s.Width, s.Height, m, resize.Lanczos3)
// 	return jpeg.Encode(w, t, &jpeg.Options{Quality: int(s.Quality)})
// 	// switch filepath.Ext(path) {
// 	// case `.webp`:
// 	// 	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, float32(s.Quality))
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	return webp.Encode(handle, *m, options)
// 	// 	// return webp.Encode(handle, *m, &webp.Options{
// 	// 	// 	Lossless: false,
// 	// 	// 	Quality:  float32(s.Quality),
// 	// 	// })
// 	// default:
// 	// 	return jpeg.Encode(handle, t, &jpeg.Options{
// 	// 		Quality: int(s.Quality),
// 	// 	})
// 	// }
// }
