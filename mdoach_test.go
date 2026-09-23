package mdcoach

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/sebdah/goldie/v2"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/text"
)

func TestPresentationAST(t *testing.T) {
	source, err := os.ReadFile("testdata/presentation-1.md")
	if err != nil {
		t.Fatal(err)
	}

	tree := NewParser().Parse(source)
	var dump bytes.Buffer
	dumpAST(tree, source, &dump)

	goldie.New(t).Assert(t, "presentation_ast", dump.Bytes())
}

func dumpAST(tree ast.Node, source []byte, dump *bytes.Buffer) {
	level := 0
	_ = ast.Walk(tree, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			content := nodePreview(node, source)
			if content == "" {
				fmt.Fprintf(dump, "%*s%s", level*2, "", node.Kind())
			} else {
				fmt.Fprintf(dump, "%*s%s: %s", level*2, "", node.Kind(), content)
			}
			if node.HasChildren() {
				dump.WriteString(" {")
			}
			dump.WriteByte('\n')
			level++
			for _, attribute := range node.Attributes() {
				fmt.Fprintf(
					dump,
					"%*s[%s=%q]\n",
					level*2,
					"",
					attribute.Name,
					attribute.Value.Value(source),
				)
			}
			// properties := node.Dump(source).Properties
			// propertyNames := make([]string, 0, len(properties))
			// for name := range properties {
			// 	if name != "attributes" && name != "children" {
			// 		propertyNames = append(propertyNames, name)
			// 	}
			// }
			// sort.Strings(propertyNames)
			// for _, name := range propertyNames {
			// 	fmt.Fprintf(
			// 		dump,
			// 		"%*s[%s=%q]\n",
			// 		level*2,
			// 		"",
			// 		name,
			// 		propertyValue(properties[name], source),
			// 	)
			// }
		} else {
			level--
			if node.HasChildren() {
				fmt.Fprintf(dump, "%*s}\n", level*2, "")
			}
		}
		return ast.WalkContinue, nil
	})
}

func propertyValue(value any, source []byte) string {
	switch value := value.(type) {
	case text.SingleLineValue:
		return value.Value(source)
	case text.MultiLineValue:
		return value.Value(source)
	default:
		return fmt.Sprint(value)
	}
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
