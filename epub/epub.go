package epub

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"io"
	"strings"
	"time"

	"github.com/dkotik/mdcoach/presentation"
)

const epubMimetype = "application/epub+zip"

// Write creates an EPUB 3 archive containing the provided rendered HTML page.
// It is the low-level archive writer; use New to render an AST and package its assets.
func Write(w io.Writer, htmlPage []byte, metadata presentation.Frontmatter) error {
	return writePublication(w, htmlPage, metadata, nil, nil)
}

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
