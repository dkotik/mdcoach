/*
Package picture loads, modifies, and saves images.
*/
package picture

import (
	_ "image/gif" // add GIF format to image reader
	_ "image/png" // add PNG format to image reader
)

type ResizedImage struct {
	FilePath string
	size     *Sizing
}

func (i *ResizedImage) Width() uint {
	return i.size.Width
}

func (i *ResizedImage) Height() uint {
	return i.size.Height
}

type Provider interface {
	GetResizedImages(p string) ([]ResizedImage, error)
}

type LocalProvider struct {
	directory string
	sizings   []Sizing
}

func NewLocalProvider(p string, sizings ...Sizing) *LocalProvider {
	SortSizingsInDescendingOrder(sizings)
	return &LocalProvider{
		directory: p,
		sizings:   sizings,
	}
}

// func (p *LocalProvider) CopyOriginal(p string) (ResizedImage, error) {
//   handle, err := os.Open(p)
//   if err != nil {
//     return m, err
//   }
//   defer handle.Close()
//   m, _, err = image.Decode(handle)
//   return m, err
// }

// func (p *LocalProvider) GetResizedImages(source string) ([]ResizedImage, error) {
// 	image, err := Load(source)
// 	if err != nil {
// 		return nil, err
// 	}
// }
