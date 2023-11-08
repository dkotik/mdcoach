package review

import (
	"errors"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
)

type options struct {
	renderer renderer.Renderer
}

type Option func(*options) error

func WithRenderer(r renderer.Renderer) Option {
	return func(o *options) error {
		if r == nil {
			return errors.New("cannot use a <nil> renderer")
		}
		if o.renderer != nil {
			return errors.New("renderer is already set")
		}
		o.renderer = r
		return nil
	}
}

func withDefaultRenderer() Option {
	return func(o *options) error {
		if o.renderer != nil {
			return nil
		}
		return WithRenderer(goldmark.DefaultRenderer())(o)
	}
}
