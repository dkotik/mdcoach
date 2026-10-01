package epub

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dkotik/mdcoach/internal"
	"github.com/dkotik/mdcoach/presentation"
	"github.com/yuin/goldmark/v2/ast"
)

const epubMimetype = "application/epub+zip"

func writePublication(w io.Writer, htmlPage []byte, metadata presentation.Frontmatter, stylesheet []byte, images []*Image) error {
	if metadata.Title == "" {
		metadata.Title = "Presentation"
	}

	archive := zip.NewWriter(w)
	mimetype := []byte(epubMimetype)
	header := &zip.FileHeader{
		Name:               "mimetype",
		Method:             zip.Store,
		CRC32:              crc32.ChecksumIEEE(mimetype),
		CompressedSize64:   uint64(len(mimetype)),
		UncompressedSize64: uint64(len(mimetype)),
	}
	mimetypeWriter, err := archive.CreateRaw(header)
	if err != nil {
		return fmt.Errorf("create EPUB mimetype entry: %w", err)
	}
	if _, err := mimetypeWriter.Write(mimetype); err != nil {
		return fmt.Errorf("write EPUB mimetype entry: %w", err)
	}

	identifier := epubIdentifier(htmlPage)
	modified := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	creator := ""
	if metadata.Author != "" {
		creator = fmt.Sprintf("    <dc:creator>%s</dc:creator>\n", escapeXML(metadata.Author))
	}

	var manifest strings.Builder
	manifest.WriteString(`    <item id="nav" href="nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>` + "\n")
	manifest.WriteString(`    <item id="presentation" href="presentation.xhtml" media-type="application/xhtml+xml"/>` + "\n")
	if stylesheet != nil {
		manifest.WriteString(`    <item id="stylesheet" href="styles.css" media-type="text/css"/>` + "\n")
	}
	for i, image := range images {
		fmt.Fprintf(&manifest, "    <item id=\"image-%d\" href=\"images/%s.png\" media-type=\"image/png\"/>\n", i, escapeXML(image.Hash))
	}

	packageDocument := []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="3.0" unique-identifier="pub-id">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="pub-id">%s</dc:identifier>
    <dc:title>%s</dc:title>
    <dc:language>en</dc:language>
%s    <meta property="dcterms:modified">%s</meta>
  </metadata>
  <manifest>
%s  </manifest>
  <spine>
    <itemref idref="presentation"/>
  </spine>
</package>`, identifier, escapeXML(metadata.Title), creator, modified, manifest.String()))
	navigation := []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
  <head><title>Contents</title></head>
  <body>
    <nav epub:type="toc" id="toc">
      <h1>Contents</h1>
      <ol><li><a href="presentation.xhtml">%s</a></li></ol>
    </nav>
  </body>
</html>`, escapeXML(metadata.Title)))
	entries := []struct {
		name    string
		content []byte
	}{
		{
			name: "META-INF/container.xml",
			content: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="OEBPS/package.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>`),
		},
		{name: "OEBPS/package.opf", content: packageDocument},
		{name: "OEBPS/nav.xhtml", content: navigation},
		{name: "OEBPS/presentation.xhtml", content: htmlPage},
	}
	if stylesheet != nil {
		entries = append(entries, struct {
			name    string
			content []byte
		}{name: "OEBPS/styles.css", content: stylesheet})
	}
	for _, image := range images {
		entries = append(entries, struct {
			name    string
			content []byte
		}{name: fmt.Sprintf("OEBPS/images/%s.png", image.Hash), content: image.Data})
	}
	for _, entry := range entries {
		file, err := archive.Create(entry.name)
		if err != nil {
			return fmt.Errorf("create EPUB entry %q: %w", entry.name, err)
		}
		if _, err := file.Write(entry.content); err != nil {
			return fmt.Errorf("write EPUB entry %q: %w", entry.name, err)
		}
	}
	if err := archive.Close(); err != nil {
		return fmt.Errorf("finish EPUB archive: %w", err)
	}
	return nil
}

func escapeXML(value string) string {
	var escaped bytes.Buffer
	_ = xml.EscapeText(&escaped, []byte(value))
	return escaped.String()
}

func epubIdentifier(htmlPage []byte) string {
	identifier := sha256.Sum256(htmlPage)
	identifier[6] = identifier[6]&0x0f | 0x80
	identifier[8] = identifier[8]&0x3f | 0x80
	return fmt.Sprintf(
		"urn:uuid:%08x-%04x-%04x-%04x-%012x",
		binary.BigEndian.Uint32(identifier[0:4]),
		binary.BigEndian.Uint16(identifier[4:6]),
		binary.BigEndian.Uint16(identifier[6:8]),
		binary.BigEndian.Uint16(identifier[8:10]),
		identifier[10:16],
	)
}

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
