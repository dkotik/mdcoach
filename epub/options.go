package epub

import (
	"fmt"
	"io/fs"

	"github.com/dkotik/mdcoach/presentation"
	"github.com/yuin/goldmark/v2/ast"
)

// Option configures AST-to-EPUB rendering.
type Option func(*newOptions) error

type newOptions struct {
	metadata presentation.Frontmatter
	media    MediaOptions
	emoteFS  fs.FS
}

// WithMetadata sets the EPUB publication metadata.
func WithMetadata(metadata presentation.Frontmatter) Option {
	return func(options *newOptions) error {
		options.metadata = metadata
		return nil
	}
}

// WithTitle sets the EPUB title.
func WithTitle(title string) Option {
	return func(options *newOptions) error {
		options.metadata.Title = title
		return nil
	}
}

// WithAuthor sets the EPUB author.
func WithAuthor(author string) Option {
	return func(options *newOptions) error {
		options.metadata.Author = author
		return nil
	}
}

// WithMediaOptions configures local and remote image loading.
func WithMediaOptions(media MediaOptions) Option {
	return func(options *newOptions) error {
		options.media = media
		return nil
	}
}

// WithEmoteFS sets the filesystem used to resolve emote PNGs.
func WithEmoteFS(assets fs.FS) Option {
	return func(options *newOptions) error {
		if assets == nil {
			return fmt.Errorf("nil emote asset filesystem")
		}
		options.emoteFS = assets
		return nil
	}
}

func metadataFromNode(node ast.Node) presentation.Frontmatter {
	document, ok := node.(*ast.Document)
	if !ok {
		return presentation.Frontmatter{}
	}
	metadata := document.Metadata()
	return presentation.Frontmatter{
		Title:  metadataText(metadata["title"]),
		Author: metadataText(metadata["author"]),
	}
}

func metadataText(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}
