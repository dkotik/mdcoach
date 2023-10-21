package renderer

import "github.com/yuin/goldmark"

var testMarkdown = goldmark.New(goldmark.WithRenderer(Must()))
