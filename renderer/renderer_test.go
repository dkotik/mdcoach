package renderer

import (
	"testing"

	"github.com/dkotik/mdcoach/picture"
	"github.com/yuin/goldmark"
)

// TODO: replace with test.NewTestRenderer version!
func newTestRenderer(t *testing.T) goldmark.Markdown {
	r, err := New(
		WithPictureProviderOptions(
			picture.WithDestinationPath(t.TempDir()),
		),
	)
	if err != nil {
		t.Fatal("cannot create test renderer:", err)
	}
	return goldmark.New(goldmark.WithRenderer(r))
}

func TestRenderer(t *testing.T) {
	newTestRenderer(t)
	// t.Fatal("renderer/slideMedia folder is created wrongly, because destination should be a set to a temporary directory; create option set for renderer")
}
