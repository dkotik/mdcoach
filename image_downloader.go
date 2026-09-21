package mdcoach

import (
	"context"
	"io"
	"net/url"
	"path"
)

type ImageDownloader interface {
	ImageRepository
	DownloadContent(context.Context) error
	WriteStyleClasses(io.Writer) error
}

type imageRepository struct {
}

func NewImageRepository() ImageDownloader {
	return &imageRepository{}
}

func (r *imageRepository) GetStyleClass(string) string {
	return ""
}

func (r *imageRepository) DownloadContent(ctx context.Context) error {
	return nil
}

func (r *imageRepository) WriteStyleClasses(w io.Writer) error {

	return nil
}

func normalizeURL(s string) string {
	url, err := url.Parse(s)
	if err != nil {
		return s
	}
	url.Path = path.Clean(url.Path)
	return url.String()
}
