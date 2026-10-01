package epub

import (
	"bytes"
	"fmt"
	stdImage "image"
	"image/png"
	"sort"
	"sync"

	"github.com/OneOfOne/xxhash"
)

// Image stores a PNG-encoded image and its source metadata.
type Image struct {
	Location string
	Hash     string
	Data     []byte
	Width    int
	Height   int
}

// NewImage encodes source as PNG and records its dimensions and content hash.
func NewImage(location string, source stdImage.Image) (*Image, error) {
	if source == nil {
		return nil, fmt.Errorf("encode image as PNG: nil image")
	}

	var encoded bytes.Buffer
	if err := png.Encode(&encoded, source); err != nil {
		return nil, fmt.Errorf("encode image as PNG: %w", err)
	}

	data := encoded.Bytes()
	hash := xxhash.New64()
	if _, err := hash.Write(data); err != nil {
		return nil, fmt.Errorf("hash encoded image: %w", err)
	}

	bounds := source.Bounds()
	return &Image{
		Location: location,
		Hash:     fmt.Sprintf("%x", hash.Sum(nil)),
		Data:     data,
		Width:    bounds.Dx(),
		Height:   bounds.Dy(),
	}, nil
}

// ImageCache stores PNG images by source location and content hash.
type ImageCache struct {
	mu     sync.Mutex
	hashes map[string]string
	images map[string]*Image
}

// NewImageCache creates an empty, concurrency-safe image cache.
func NewImageCache() *ImageCache {
	return &ImageCache{
		hashes: make(map[string]string),
		images: make(map[string]*Image),
	}
}

// Get returns the cached image at location, if present.
func (c *ImageCache) Get(location string) (*Image, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	hash, ok := c.hashes[location]
	if !ok {
		return nil, false
	}
	image, ok := c.images[hash]
	return image, ok
}

// Set adds image to the cache by its content hash and source location.
func (c *ImageCache) Set(image *Image) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.images[image.Hash] = image
	c.hashes[image.Location] = image.Hash
}

func (c *ImageCache) all() []*Image {
	c.mu.Lock()
	defer c.mu.Unlock()

	images := make([]*Image, 0, len(c.images))
	for _, image := range c.images {
		images = append(images, image)
	}
	sort.Slice(images, func(i, j int) bool {
		return images[i].Hash < images[j].Hash
	})
	return images
}
