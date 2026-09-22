package mdcoach

import (
	"bytes"
	"io"
	"testing"

	"github.com/yuin/goldmark/v2/ast"
)

func TestImageRendererRendersOnlyImageNodes(t *testing.T) {
	source := []byte(`![cat](media/cat_1.jpg)`)
	tree, ok := NewParser().Parse(source).(*ast.Document)
	if !ok {
		t.Fatal("parser returned a non-document AST")
	}

	var rendered bytes.Buffer
	if err := renderImagesOnly(&rendered, source, tree); err != nil {
		t.Fatal(err)
	}

	const want = `<div data-src="media/cat_1.jpg" data-alt="cat"></div>`
	if rendered.String() != want {
		t.Fatalf("rendered image = %q, want %q", rendered.String(), want)
	}
}

func renderImagesOnly(writer io.Writer, source []byte, tree ast.Node) error {
	var (
		renderer  = NewRenderer()
		renderErr error
	)
	_ = ast.Walk(tree, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if node.Kind() != ast.KindImage {
			return ast.WalkContinue, nil
		}

		if err := renderer.Render(writer, source, node); err != nil {
			renderErr = err
			return ast.WalkStop, err
		}
		return ast.WalkSkipChildren, nil
	})
	return renderErr
}
