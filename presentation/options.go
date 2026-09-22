package presentation

import (
	"errors"

	"github.com/dkotik/mdcoach"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

type options struct {
	Parser           parser.Parser
	Renderer         html.Renderer
	ImageHeightLimit int
	ImageWidthLimit  int
	ImageQuality     int
	MediaOptions     *mdcoach.MediaOptions
}

type Option func(*options) error

func WithParser(p parser.Parser) Option {
	return func(o *options) error {
		if p == nil {
			return errors.New("nil parser")
		}
		if o.Parser != nil {
			return errors.New("parser already set")
		}
		o.Parser = p
		return nil
	}
}
func WithRenderer(r html.Renderer) Option {
	return func(o *options) error {
		if r == nil {
			return errors.New("nil renderer")
		}
		if o.Renderer != nil {
			return errors.New("renderer already set")
		}
		o.Renderer = r
		return nil
	}
}

func WithImageSizeLimit(w, h int) Option {
	return func(o *options) error {
		if o.ImageHeightLimit != 0 || o.ImageWidthLimit != 0 {
			return errors.New("image size limit already set")
		}
		o.ImageHeightLimit = h
		o.ImageWidthLimit = w
		return nil
	}
}

func WithImageQuality(q int) Option {
	return func(o *options) error {
		o.ImageQuality = q
		return nil
	}
}
