package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/dkotik/mdcoach/presentation"
)

const epubMimetype = "application/epub+zip"

var (
	scriptElementPattern  = regexp.MustCompile(`(?is)<script\b(?:[^>"']|"[^"]*"|'[^']*')*>.*?</script\s*>`)
	styleElementPattern   = regexp.MustCompile(`(?is)<style\b(?:[^>"']|"[^"]*"|'[^']*')*>.*?</style\s*>`)
	openingElementPattern = regexp.MustCompile(`(?is)^<(?:script|style)\b(?:[^>"']|"[^"]*"|'[^']*')*>`)
	elementNamePattern    = regexp.MustCompile(`(?is)^<(?:script|style)\b`)
	attributePattern      = regexp.MustCompile(`(?is)([^\s=/>]+)(?:\s*=\s*(?:"[^"]*"|'[^']*'|[^\s>]+))?`)
)

func writeEPUB(w io.Writer, htmlPage []byte, metadata presentation.Frontmatter) error {
	htmlPage = stripUnmarkedScriptAndStyleTags(htmlPage)
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
		var escapedAuthor bytes.Buffer
		if err := xml.EscapeText(&escapedAuthor, []byte(metadata.Author)); err != nil {
			return fmt.Errorf("escape EPUB author: %w", err)
		}
		creator = fmt.Sprintf("    <dc:creator>%s</dc:creator>\n", escapedAuthor.String())
	}
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
		{
			name: "OEBPS/package.opf",
			content: []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="3.0" unique-identifier="pub-id">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="pub-id">%s</dc:identifier>
    <dc:title>Presentation</dc:title>
    <dc:language>en</dc:language>
%s    <meta property="dcterms:modified">%s</meta>
  </metadata>
  <manifest>
    <item id="nav" href="nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>
    <item id="presentation" href="presentation.xhtml" media-type="application/xhtml+xml"/>
  </manifest>
  <spine>
    <itemref idref="presentation"/>
  </spine>
</package>`, identifier, creator, modified)),
		},
		{
			name: "OEBPS/nav.xhtml",
			content: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
  <head><title>Contents</title></head>
  <body>
    <nav epub:type="toc" id="toc">
      <h1>Contents</h1>
      <ol><li><a href="presentation.xhtml">Presentation</a></li></ol>
    </nav>
  </body>
</html>`),
		},
		{
			name:    "OEBPS/presentation.xhtml",
			content: htmlPage,
		},
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

func stripUnmarkedScriptAndStyleTags(htmlPage []byte) []byte {
	stripElements := func(page []byte, pattern *regexp.Regexp) []byte {
		return pattern.ReplaceAllFunc(page, func(element []byte) []byte {
			openingTag := openingElementPattern.Find(element)
			if hasRoleAttribute(openingTag) {
				return element
			}
			return nil
		})
	}

	htmlPage = stripElements(htmlPage, scriptElementPattern)
	return stripElements(htmlPage, styleElementPattern)
}

func hasRoleAttribute(openingTag []byte) bool {
	nameEnd := elementNamePattern.FindIndex(openingTag)
	tagEnd := bytes.LastIndexByte(openingTag, '>')
	if nameEnd == nil || tagEnd < nameEnd[1] {
		return false
	}

	for _, attribute := range attributePattern.FindAllSubmatch(openingTag[nameEnd[1]:tagEnd], -1) {
		if strings.EqualFold(string(attribute[1]), "role") {
			return true
		}
	}
	return false
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
