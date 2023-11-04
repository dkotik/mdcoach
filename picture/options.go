package picture

import (
	"errors"
	"os"
)

type options struct {
	sourcePath      string
	destinationPath string
	encoder         Encoder
	sizings         []Sizing
}

type Option func(*options) error

func withDefaultSourcePath() Option {
	return func(o *options) error {
		if o.sourcePath != "" {
			return nil
		}
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		o.sourcePath = cwd
		return nil
	}
}

func WithSourcePath(p string) Option {
	return func(o *options) error {
		if p == "" {
			return errors.New("cannot use an empty source path")
		}
		if o.sourcePath != "" {
			return errors.New("source path is already set")
		}
		o.sourcePath = p
		return nil
	}
}

func withDefaultDestinationPath() Option {
	return func(o *options) error {
		if o.destinationPath != "" {
			return nil
		}
		o.destinationPath = `slideMedia`
		return nil
	}
}

func WithDestinationPath(p string) Option {
	return func(o *options) error {
		if p == "" {
			return errors.New("cannot use an empty destination path")
		}
		if o.destinationPath != "" {
			return errors.New("destination path is already set")
		}
		o.destinationPath = p
		return nil
	}
}

func withDefaultEncoder() Option {
	return func(o *options) error {
		if o.encoder != nil {
			return nil
		}
		o.encoder = &JPEGEncoder{}
		return nil
	}
}

func withDefaultSizing() Option {
	return func(o *options) error {
		if len(o.sizings) > 0 {
			return nil
		}
		return WithSizing(128, 128, 30)(o)
	}
}

func WithSizing(w, h, quality int) Option {
	return func(o *options) error {
		if w < 64 {
			return errors.New("width cannot be less than 64 pixels")
		}
		if w > 4096 {
			return errors.New("width cannot be greater than 4096 pixels")
		}
		if h < 64 {
			return errors.New("height cannot be less than 64 pixels")
		}
		if h > 4096 {
			return errors.New("height cannot be greater than 4096 pixels")
		}
		if quality < 0 {
			return errors.New("quality cannot be less than 0")
		}
		if quality > 100 {
			return errors.New("quality cannot be greater than 100")
		}
		o.sizings = append(o.sizings, Sizing{
			Width:   w,
			Height:  h,
			Quality: quality,
		})
		return nil
	}
}
