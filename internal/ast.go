package internal

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/yuin/goldmark/v2/ast"
)

// WriteAST writes a deterministic, indented representation of tree to writer.
func WriteAST(writer io.Writer, tree ast.Node, source []byte) {
	level := 0
	_ = ast.Walk(tree, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			content := nodePreview(node, source)
			if content == "" {
				_, _ = fmt.Fprintf(writer, "%*s%s", level*2, "", node.Kind())
			} else {
				_, _ = fmt.Fprintf(writer, "%*s%s: %s", level*2, "", node.Kind(), content)
			}
			if node.HasChildren() {
				_, _ = fmt.Fprint(writer, " {")
			}
			_, _ = fmt.Fprintln(writer)
			level++
			for _, attribute := range node.Attributes() {
				_, _ = fmt.Fprintf(
					writer,
					"%*s[%s=%q]\n",
					level*2,
					"",
					attribute.Name,
					attribute.Value.Value(source),
				)
			}
		} else {
			level--
			if node.HasChildren() {
				_, _ = fmt.Fprintf(writer, "%*s}\n", level*2, "")
			}
		}
		return ast.WalkContinue, nil
	})
}

func nodePreview(node ast.Node, source []byte) string {
	var content string

	switch node := node.(type) {
	case *ast.Text:
		content = node.Value.Value(source)
	case ast.BlockNode:
		segments := node.Source()
		if len(segments) > 0 {
			content = string(segments[0].Bytes(source))
		}
	}

	if content == "" {
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			if content = nodePreview(child, source); content != "" {
				break
			}
		}
	}

	content = strings.Join(strings.Fields(content), " ")
	if utf8.RuneCountInString(content) > 60 {
		content = string([]rune(content)[:60])
	}
	return content
}
