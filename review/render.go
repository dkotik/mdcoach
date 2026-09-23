package review

import (
	"fmt"

	"github.com/gpdf-dev/gpdf"
	gpdfdocument "github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/pdf"
	gpdftemplate "github.com/gpdf-dev/gpdf/template"
)

type Page struct {
	Title       string
	Description string
	Questions   []string
}

func New(page Page) ([]byte, error) {
	document := gpdf.NewDocument(
		gpdf.WithPageSize(gpdf.A4),
		gpdf.WithMargins(gpdfdocument.UniformEdges(gpdfdocument.Mm(20))),
		gpdf.WithMetadata(gpdfdocument.DocumentMetadata{Title: page.Title}),
	)
	pdfPage := document.AddPage()
	if page.Title != "" {
		pdfPage.AutoRow(func(row *gpdftemplate.RowBuilder) {
			row.Col(12, func(column *gpdftemplate.ColBuilder) {
				column.Text(page.Title, gpdftemplate.FontSize(18))
			})
		})
	}
	if page.Description != "" {
		pdfPage.AutoRow(func(row *gpdftemplate.RowBuilder) {
			row.Col(12, func(column *gpdftemplate.ColBuilder) {
				column.Text(page.Description, gpdftemplate.FontSize(12))
			})
		})
	}
	for _, question := range page.Questions {
		pdfPage.AutoRow(func(row *gpdftemplate.RowBuilder) {
			row.Col(11, func(column *gpdftemplate.ColBuilder) {
				column.Text(question, gpdftemplate.FontSize(14))
			})
		})

		pdfPage.AutoRow(func(row *gpdftemplate.RowBuilder) {
			row.Col(11, func(column *gpdftemplate.ColBuilder) {
				column.Text(page.Description, gpdftemplate.FontSize(12))
			})
			row.Col(1, func(column *gpdftemplate.ColBuilder) {
				for i := 1; i <= 5; i++ {
					column.Text(
						fmt.Sprintf("-%d", i),
						gpdftemplate.FontSize(7),
						gpdftemplate.TextColor(pdf.Gray(0.7)),
						gpdftemplate.AlignRight(),
					)
					column.Spacer(gpdfdocument.Mm(2))
				}
			})
		})
	}
	data, err := document.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate questions PDF: %w", err)
	}
	return data, nil
}
