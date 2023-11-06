package document

import (
	"fmt"
	"html/template"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/OneOfOne/xxhash"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

type Metadata struct {
	ID          string
	Title       string
	Description string
	Author      string
	Year        int
	Created     time.Time
}

func NewMetadata(ctx parser.Context) (*Metadata, error) {
	raw := frontmatter.Get(ctx)
	var metadata Metadata
	if err := raw.Decode(&metadata); err != nil {
		return nil, err
	}
	if metadata.ID == "" {
		h := xxhash.New64()
		_, _ = io.Copy(h, strings.NewReader(metadata.Title))
		_, _ = io.Copy(h, strings.NewReader(metadata.Description))
		_, _ = io.Copy(h, strings.NewReader(metadata.Author))
		metadata.ID = fmt.Sprintf("%x", h.Sum(nil))
	}
	return &metadata, nil
}

var metadataTemplate *template.Template

func (m *Metadata) ToHTML(w io.Writer) (err error) {
	sync.OnceFunc(func() {
		metadataTemplate = template.Must(
			// TODO: fill out the template.
			template.New("metadata").Parse(``),
		)
	})()
	return metadataTemplate.Execute(w, m)
}
