package mdcoach

import (
	"bytes"
	"io"
	"testing"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/text"
)

func TestImageRendererRendersOnlyImageNodes(t *testing.T) {
	imageNode := ast.NewImage(text.NewSingleLineValueFromString(
		"media/cat_1.jpg",
		text.IdentityDecoder,
	))
	imageNode.AppendChild(ast.NewText(text.NewSingleLineValueFromString(
		"cat",
		text.IdentityDecoder,
	)))

	tree := ast.NewDocument()
	paragraph := ast.NewParagraph()
	paragraph.AppendChild(imageNode)
	tree.AppendChild(paragraph)
	tree.AppendChild(ast.NewParagraph())

	var rendered bytes.Buffer
	if err := renderImagesOnly(&rendered, tree); err != nil {
		t.Fatal(err)
	}

	const want = `<img src="media/cat_1.jpg" alt="cat">`
	if rendered.String() != want {
		t.Fatalf("rendered image = %q, want %q", rendered.String(), want)
	}
}

func renderImagesOnly(writer io.Writer, tree ast.Node) error {
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

		if err := renderer.Render(writer, nil, node); err != nil {
			renderErr = err
			return ast.WalkStop, err
		}
		return ast.WalkSkipChildren, nil
	})
	return renderErr
}
