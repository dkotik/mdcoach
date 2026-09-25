/*
Package presentation makes
HTML presentations from Markdown files.
*/
package presentation

import (
	"context"
	_ "embed" // for html/before.gen.html and html/after.gen.html
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/dkotik/mdcoach"
	"github.com/yuin/goldmark/v2/ast"
)

//go:generate go run ./html/after.go
//go:embed html/after.gen.html
var afterMain []byte

//go:generate go run ./html/before.go
//go:embed html/before.gen.html
var beforeMain []byte

//go:embed html/header.html
var header []byte

var headerTemplate = template.Must(template.New("header").Parse(string(header)))

//go:embed html/footer.html
var footer []byte

func New(
	ctx context.Context,
	w io.Writer,
	sources []string,
	withOptions ...Option,
) (err error) {
	o := &options{}
	for _, opt := range append(withOptions,
		func(o *options) error {
			if o.Parser == nil {
				o.Parser = mdcoach.NewParser()
			}
			if o.ImageCache == nil {
				o.ImageCache = mdcoach.NewImageCache()
			}
			if o.Renderer == nil {
				o.Renderer = mdcoach.NewRenderer(o.ImageCache)
			}
			return nil
		},
	) {
		if err := opt(o); err != nil {
			return err
		}
	}

	mo := mdcoach.MediaOptions{
		Path:        ".",
		WidthLimit:  o.ImageWidthLimit,
		HeightLimit: o.ImageHeightLimit,
		Quality:     o.ImageQuality,
	}

	var firstSource []byte
	var firstTree ast.Node
	metadata := Metadata{}
	if len(sources) > 0 {
		firstSource, err = os.ReadFile(sources[0])
		if err != nil {
			return fmt.Errorf("failed to read source file %s: %w", sources[0], err)
		}
		firstTree = o.Parser.Parse(firstSource)
		metadata, err = metadataFromTree(firstTree)
		if err != nil {
			return fmt.Errorf("failed to read metadata from %s: %w", sources[0], err)
		}
	}
	if err = headerTemplate.Execute(w, metadata); err != nil {
		return fmt.Errorf("failed to render header: %w", err)
	}
	if _, err = w.Write(beforeMain); err != nil {
		return fmt.Errorf("failed to write before main: %w", err)
	}

	for i, sourcePath := range sources {
		source, tree := firstSource, firstTree
		if i > 0 {
			source, err = os.ReadFile(sourcePath)
			if err != nil {
				return fmt.Errorf("failed to read source file %s: %w", sourcePath, err)
			}
			tree = o.Parser.Parse(source)
		}
		mo.Path = filepath.Dir(sourcePath)
		if err = mdcoach.NewImageLoader(o.ImageCache, mo).LoadImages(
			ctx,
			source,
			tree,
		); err != nil {
			return fmt.Errorf("failed to load images from %s: %w", sourcePath, err)
		}

		err = o.Renderer.Render(w, source, tree)
		if err != nil {
			return fmt.Errorf("failed to render file %s: %w", sourcePath, err)
		}
	}

	_, err = w.Write(afterMain)
	if err != nil {
		return fmt.Errorf("failed to write after main: %w", err)
	}

	if err = o.ImageCache.WriteImageDataCSS(w); err != nil {
		return fmt.Errorf("failed to write image data CSS: %w", err)
	}

	_, err = w.Write(footer)
	if err != nil {
		return fmt.Errorf("failed to write footer: %w", err)
	}

	return nil
}
