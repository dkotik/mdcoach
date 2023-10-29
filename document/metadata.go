package document

import (
	"html/template"
	"io"
	"sync"
	"time"
)

type Metadata struct {
	Title       string
	Description string
	Author      string
	SlideCount  int
	Created     time.Time
}

// func NewMetadata(n ast.Node) (*Metadata)

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
