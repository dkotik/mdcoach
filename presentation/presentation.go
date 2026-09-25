/*
Package presentation makes
HTML presentations from Markdown files.
*/
package presentation

import (
	"bytes"
	"context"
	_ "embed" // for html/before.gen.html and html/after.gen.html
	"encoding/base64"
	"fmt"
	"html/template"
	stdImage "image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/dkotik/mdcoach"
	"github.com/nfnt/resize"
	"github.com/yuin/goldmark/v2/ast"
)

//go:generate go run ./html/after.go
//go:embed html/after.gen.html
var afterMain []byte

//go:generate go run ./html/before.go
//go:embed html/before.gen.html
var beforeMain []byte

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
		withDefaultTemplates,
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

		metadata.Favicon, err = faviconFromFigure(firstTree, firstSource, sources[0])
		if err != nil {
			return fmt.Errorf("failed to create favicon from %s: %w", sources[0], err)
		}
	}
	if err = o.HeaderTemplate.Execute(w, metadata); err != nil {
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

	if err = o.FooterTemplate.Execute(w, metadata); err != nil {
		return fmt.Errorf("failed to render footer: %w", err)
	}

	return nil
}

func faviconFromFigure(tree ast.Node, source []byte, sourcePath string) (template.HTML, error) {
	for child := tree.FirstChild(); child != nil; child = child.NextSibling() {
		if child.Kind() != mdcoach.KindFigure {
			continue
		}

		imageNode, ok := child.FirstChild().(*ast.Image)
		if !ok {
			return "", fmt.Errorf("figure's first child must be *ast.Image, got %T", child.FirstChild())
		}

		destination := imageNode.Destination.Value(source)
		imagePath := filepath.FromSlash(destination)
		if !filepath.IsAbs(imagePath) {
			imagePath = filepath.Join(filepath.Dir(sourcePath), imagePath)
		}
		data, err := os.ReadFile(imagePath)
		if err != nil {
			return "", fmt.Errorf("read figure image %q: %w", imagePath, err)
		}

		decoded, _, err := stdImage.Decode(bytes.NewReader(data))
		if err != nil {
			return "", fmt.Errorf("decode figure image %q: %w", imagePath, err)
		}

		resized := resize.Thumbnail(64, 64, decoded, resize.Lanczos3)
		var encoded bytes.Buffer
		if err := png.Encode(&encoded, resized); err != nil {
			return "", fmt.Errorf("encode favicon PNG from %q: %w", imagePath, err)
		}

		favicon := "<link rel=\"icon\" type=\"image/png\" href=\"data:image/png;base64," +
			base64.StdEncoding.EncodeToString(encoded.Bytes()) + `">`
		return template.HTML(favicon), nil
	}
	return "", nil
}
