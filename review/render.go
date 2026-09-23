package review

import (
	"fmt"

	"github.com/gpdf-dev/gpdf"
	gpdfdocument "github.com/gpdf-dev/gpdf/document"
	gpdftemplate "github.com/gpdf-dev/gpdf/template"
)

func New(questions []string, title string) ([]byte, error) {
	document := gpdf.NewDocument(
		gpdf.WithPageSize(gpdf.A4),
		gpdf.WithMargins(gpdfdocument.UniformEdges(gpdfdocument.Mm(20))),
		gpdf.WithMetadata(gpdfdocument.DocumentMetadata{Title: title}),
	)
	page := document.AddPage()
	for _, question := range questions {
		page.AutoRow(func(row *gpdftemplate.RowBuilder) {
			row.Col(12, func(column *gpdftemplate.ColBuilder) {
				column.Text(question, gpdftemplate.FontSize(14))
			})
		})
	}
	data, err := document.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate questions PDF: %w", err)
	}
	return data, nil
}
