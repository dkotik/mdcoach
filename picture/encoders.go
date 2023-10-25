package picture

import (
	"errors"
	"image"
	"image/jpeg"
	"os"
)

// TODO: encoding should be customizeable for the resizer via an option.
type Encoder interface {
	EncodeImage(*SizingTarget, image.Image) error
}

type JPEGEncoder struct {
	// TODO: overwrite confirm.
}

func (j *JPEGEncoder) EncodeImage(t *SizingTarget, m image.Image) (err error) {
	// if strings.ToLower(filepath.Ext(p)) != "jpg"
	// TODO: confirm overwrite.
	w, err := os.Create(t.Destination + ".jpg")
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, w.Close())
	}()
	return jpeg.Encode(w, m, &jpeg.Options{Quality: int(t.Quality)})
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
