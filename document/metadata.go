package document

import (
	"html/template"
	"io"
	"sync"
	"time"

	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

type Metadata struct {
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
