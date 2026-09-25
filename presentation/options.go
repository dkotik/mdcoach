package presentation

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"

	"github.com/dkotik/mdcoach"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

type options struct {
	Parser           parser.Parser
	Renderer         html.Renderer
	HeaderTemplate   *template.Template
	FooterTemplate   *template.Template
	ImageCache       *mdcoach.ImageCache
	ImageHeightLimit int
	ImageWidthLimit  int
	ImageQuality     int
	MediaOptions     *mdcoach.MediaOptions
}

type Option func(*options) error

//go:embed html/header.html
var header []byte

//go:embed html/footer.html
var footer []byte

func withDefaultTemplates(o *options) error {
	if o.HeaderTemplate == nil {
		t, err := template.New("header").Parse(string(header))
		if err != nil {
			return fmt.Errorf("parse default header template: %w", err)
		}
		o.HeaderTemplate = t
	}
	if o.FooterTemplate == nil {
		t, err := template.New("footer").Parse(string(footer))
		if err != nil {
			return fmt.Errorf("parse default footer template: %w", err)
		}
		o.FooterTemplate = t
	}
	return nil
}

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

func WithHeaderTemplate(t *template.Template) Option {
	return func(o *options) error {
		if t == nil {
			return errors.New("nil header template")
		}
		if o.HeaderTemplate != nil {
			return errors.New("header template already set")
		}
		o.HeaderTemplate = t
		return nil
	}
}

func WithFooterTemplate(t *template.Template) Option {
	return func(o *options) error {
		if t == nil {
			return errors.New("nil footer template")
		}
		if o.FooterTemplate != nil {
			return errors.New("footer template already set")
		}
		o.FooterTemplate = t
		return nil
	}
}

func WithImageCache(cache *mdcoach.ImageCache) Option {
	return func(o *options) error {
		if cache == nil {
			return errors.New("nil image cache")
		}
		if o.ImageCache != nil {
			return errors.New("image cache already set")
		}
		o.ImageCache = cache
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
