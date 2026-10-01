package epub

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"path"
	"path/filepath"

	"github.com/dkotik/mdcoach/internal"
	"github.com/dkotik/mdcoach/presentation"
	"github.com/yuin/goldmark/v2/ast"
)

const epubStylesheet = `body {
  font-family: serif;
  line-height: 1.5;
  margin: 5%;
}
section {
  break-after: page;
  page-break-after: always;
}
img {
  height: auto;
  max-width: 100%;
}
figure {
  margin: 1em 0;
  text-align: center;
}
figcaption {
  font-size: 0.9em;
  margin-top: 0.5em;
}
table {
  border-collapse: collapse;
  width: 100%;
}
th, td {
  border: 1px solid currentColor;
  padding: 0.25em 0.5em;
}
pre {
  white-space: pre-wrap;
}
.speaker-notes, .footnote {
  border-top: 1px solid currentColor;
  margin-top: 1em;
  padding-top: 0.5em;
}
`

// New renders a parsed Markdown AST as EPUB 3 and writes the archive to w.
// source must be the byte slice used to parse node with mdcoach.NewParser.
func New(ctx context.Context, w io.Writer, source []byte, node ast.Node, withOptions ...Option) error {
	return newFromSources(ctx, w, []presentation.Source{{
		Source:       source,
		Presentation: node,
	}}, withOptions...)
}

// NewSources renders parsed Markdown sources together into one EPUB 3 archive.
// Each source should be parsed with mdcoach.NewParser; source paths are used to
// resolve relative images unless media options provide a filesystem or path.
func NewSources(
	ctx context.Context,
	w io.Writer,
	sources []presentation.Source,
	withOptions ...Option,
) error {
	return newFromSources(ctx, w, sources, withOptions...)
}

func newFromSources(
	ctx context.Context,
	w io.Writer,
	sources []presentation.Source,
	withOptions ...Option,
) error {
	if w == nil {
		return fmt.Errorf("nil EPUB writer")
	}
	for i, source := range sources {
		if source.Presentation == nil {
			return fmt.Errorf("nil Markdown AST node for source %d", i)
		}
	}
	if ctx == nil {
		ctx = context.Background()
	}

	options := newOptions{}
	if len(sources) > 0 {
		options.metadata = metadataFromNode(sources[0].Presentation)
	}
	for _, option := range withOptions {
		if option == nil {
			return fmt.Errorf("nil EPUB option")
		}
		if err := option(&options); err != nil {
			return fmt.Errorf("configure EPUB: %w", err)
		}
	}
	if options.metadata.Title == "" {
		options.metadata.Title = "Presentation"
	}

	emoteFS := options.emoteFS
	if emoteFS == nil {
		var err error
		emoteFS, err = fs.Sub(internal.Assets, path.Join("assets", "emojis"))
		if err != nil {
			return fmt.Errorf("open embedded emote assets: %w", err)
		}
	}

	cache := NewImageCache()
	renderer := newEPUBHTMLRenderer(cache, emoteFS)
	var body bytes.Buffer
	for _, source := range sources {
		media := options.media
		if media.FS == nil && media.Path == "" {
			media.Path = "."
			if source.Path != "" {
				media.Path = filepath.Dir(source.Path)
			}
		}
		loader := NewImageLoader(cache, media)
		if err := loader.LoadImages(ctx, source.Source, source.Presentation); err != nil {
			return fmt.Errorf("load EPUB images from %q: %w", source.Path, err)
		}
		if err := renderer.Render(&body, source.Source, source.Presentation); err != nil {
			return fmt.Errorf("render EPUB XHTML from %q: %w", source.Path, err)
		}
	}

	author := ""
	if options.metadata.Author != "" {
		author = fmt.Sprintf("  <meta name=\"author\" content=\"%s\" />\n", escapeXML(options.metadata.Author))
	}
	page := []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <title>%s</title>
%s  <link rel="stylesheet" type="text/css" href="styles.css" />
</head>
<body>
%s
</body>
</html>`, escapeXML(options.metadata.Title), author, body.String()))
	stylesheet := append([]byte(epubStylesheet), []byte(options.metadata.Stylesheet)...)
	if err := writePublication(w, page, options.metadata, stylesheet, cache.all()); err != nil {
		return fmt.Errorf("write EPUB publication: %w", err)
	}
	return nil
}
