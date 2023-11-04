package renderer

import (
	"errors"

	"github.com/dkotik/mdcoach/picture"
)

type options struct {
	PictureProvider picture.Provider
}

type Option func(*options) error

func WithPictureProvider(p picture.Provider) Option {
	return func(o *options) error {
		if o.PictureProvider != nil {
			return errors.New("picture provider is already set")
		}
		if p == nil {
			return errors.New("cannot use a <nil> picture provider")
		}
		o.PictureProvider = p
		return nil
	}
}

func WithPictureProviderOptions(withOptions ...picture.Option) Option {
	return func(o *options) error {
		provider, err := picture.NewLocalProvider(withOptions...)
		if err != nil {
			return err
		}
		return WithPictureProvider(provider)(o)
	}
}
