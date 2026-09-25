package presentation

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/yuin/goldmark/v2/ast"
)

type Metadata struct {
	ID          string
	Title       string
	Description string
	Keywords    string
	Author      string
	Created     time.Time
	Favicon     template.HTML
}

func metadataFromTree(tree ast.Node) (Metadata, error) {
	document, ok := tree.(*ast.Document)
	if !ok {
		return Metadata{}, fmt.Errorf("expected Goldmark document, got %T", tree)
	}

	frontmatter := document.Metadata()
	metadata := Metadata{
		ID:          frontmatterString(frontmatter, "id"),
		Title:       frontmatterString(frontmatter, "title"),
		Description: frontmatterString(frontmatter, "description"),
		Keywords:    frontmatterString(frontmatter, "keywords"),
		Author:      frontmatterString(frontmatter, "author"),
	}
	if metadata.Author == "" {
		metadata.Author = frontmatterString(frontmatter, "автор")
	}
	if metadata.ID == "" {
		if id, exists := document.Attribute("id"); exists {
			metadata.ID = string(id.Value(nil))
		}
	}

	created, exists := frontmatterValue(frontmatter, "created")
	if !exists {
		created, exists = frontmatterValue(frontmatter, "date")
	}
	if exists {
		var err error
		metadata.Created, err = parseCreated(created)
		if err != nil {
			return Metadata{}, fmt.Errorf("parse frontmatter created date: %w", err)
		}
	}
	return metadata, nil
}

func frontmatterValue(frontmatter map[string]interface{}, name string) (interface{}, bool) {
	for key, value := range frontmatter {
		if strings.EqualFold(key, name) {
			return value, true
		}
	}
	return nil, false
}

func frontmatterString(frontmatter map[string]interface{}, name string) string {
	value, exists := frontmatterValue(frontmatter, name)
	if !exists || value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func parseCreated(value interface{}) (time.Time, error) {
	if created, ok := value.(time.Time); ok {
		return created, nil
	}

	date := fmt.Sprint(value)
	for _, layout := range []string{
		time.RFC3339Nano,
		"2006-01-02",
		"January 2, 2006",
		"January 2006",
	} {
		if created, err := time.Parse(layout, date); err == nil {
			return created, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported date value %q", date)
}
