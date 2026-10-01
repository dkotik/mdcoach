package epub

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"path"

	"github.com/dkotik/mdcoach/internal"
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

// New renders node as EPUB 3 and writes the archive to w. The node should be
// produced by mdcoach.NewParser; source must be the byte slice used to parse it.
func New(ctx context.Context, w io.Writer, source []byte, node ast.Node, withOptions ...Option) error {
	if w == nil {
		return fmt.Errorf("nil EPUB writer")
	}
	if node == nil {
		return fmt.Errorf("nil Markdown AST node")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	options := newOptions{metadata: metadataFromNode(node)}
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
	if options.media.Path == "" && options.media.FS == nil {
		options.media.Path = "."
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
	loader := NewImageLoader(cache, options.media)
	if err := loader.LoadImages(ctx, source, node); err != nil {
		return fmt.Errorf("load EPUB images: %w", err)
	}

	renderer := newEPUBHTMLRenderer(cache, emoteFS)
	var body bytes.Buffer
	if err := renderer.Render(&body, source, node); err != nil {
		return fmt.Errorf("render EPUB XHTML: %w", err)
	}

	page := []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <title>%s</title>
  <link rel="stylesheet" type="text/css" href="styles.css" />
</head>
<body>
%s
</body>
</html>`, escapeXML(options.metadata.Title), body.String()))
	stylesheet := append([]byte(epubStylesheet), []byte(options.metadata.Stylesheet)...)
	if err := writePublication(w, page, options.metadata, stylesheet, cache.all()); err != nil {
		return fmt.Errorf("write EPUB publication: %w", err)
	}
	return nil
}
