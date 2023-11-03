/*
Package picture loads, modifies, and saves images.
*/
package picture

import (
	"context"
	_ "image/gif" // add GIF format to image reader
	_ "image/png" // add PNG format to image reader
)

type Source struct {
	Sizing
	Location string
}

type Provider interface {
	GetSourceSet(ctx context.Context, location string) ([]Source, error)
}

func isContextAlive(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
