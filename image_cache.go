package mdcoach

import (
	"bytes"
	"encoding/base64"
	"fmt"
	stdImage "image"
	"image/jpeg"
	"sync"

	"github.com/OneOfOne/xxhash"
)

type imageJPG struct {
	Hash       string
	DataBase64 []byte
	Width      int
	Height     int
}

func newImageFromImage(i stdImage.Image) (*imageJPG, error) {
	if i == nil {
		return nil, fmt.Errorf("encode image as JPEG: nil image")
	}

	var encoded bytes.Buffer
	if err := jpeg.Encode(&encoded, i, nil); err != nil {
		return nil, fmt.Errorf("encode image as JPEG: %w", err)
	}

	encodedBytes := encoded.Bytes()
	hash := xxhash.New64()
	if _, err := hash.Write(encodedBytes); err != nil {
		return nil, fmt.Errorf("hash encoded image: %w", err)
	}

	bounds := i.Bounds()
	return &imageJPG{
		Hash:       fmt.Sprintf("%x", hash.Sum(nil)),
		DataBase64: []byte(base64.StdEncoding.EncodeToString(encodedBytes)),
		Width:      bounds.Dx(),
		Height:     bounds.Dy(),
	}, nil
}

type ImageCache struct {
	mu     *sync.Mutex
	hashes map[string]string
	images map[string]*imageJPG
}

func NewImageCache() *ImageCache {
	return &ImageCache{
		mu:     &sync.Mutex{},
		hashes: make(map[string]string),
		images: make(map[string]*imageJPG),
	}
}

func (c *ImageCache) Get(location string) (*imageJPG, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	hash, ok := c.hashes[location]
	if !ok {
		return nil, false
	}
	img, ok := c.images[hash]
	return img, ok
}

func (c *ImageCache) Set(location string, img *imageJPG) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.images[img.Hash] = img
	c.hashes[location] = img.Hash
}
