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
	"image"
	"image/color"
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
	for slide := range tree.Children() {
		for child := range slide.Children() {
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

			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				return "", fmt.Errorf("decode figure image %q: %w", imagePath, err)
			}

			resized := resize.Thumbnail(64, 64, decoded, resize.Lanczos3)
			resized = roundImageBorders(resized)

			var encoded bytes.Buffer
			if err := png.Encode(&encoded, resized); err != nil {
				return "", fmt.Errorf("encode favicon PNG from %q: %w", imagePath, err)
			}

			favicon := "<link rel=\"icon\" type=\"image/png\" href=\"data:image/png;base64," +
				base64.StdEncoding.EncodeToString(encoded.Bytes()) + `">`
			return template.HTML(favicon), nil
		}
	}
	return "", nil
}

func roundImageBorders(img image.Image) image.Image {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	rounded := image.NewNRGBA(bounds)
	if width == 0 || height == 0 {
		return rounded
	}

	radius := float64(min(width, height)) / 6
	const samplesPerAxis = 4
	const totalSamples = samplesPerAxis * samplesPerAxis

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			inside := 0
			for sampleY := 0; sampleY < samplesPerAxis; sampleY++ {
				py := float64(y-bounds.Min.Y) + (float64(sampleY)+0.5)/samplesPerAxis
				for sampleX := 0; sampleX < samplesPerAxis; sampleX++ {
					px := float64(x-bounds.Min.X) + (float64(sampleX)+0.5)/samplesPerAxis
					if insideRoundedRectangle(px, py, float64(width), float64(height), radius) {
						inside++
					}
				}
			}
			if inside == 0 {
				continue
			}

			pixel := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			pixel.A = uint8((uint16(pixel.A)*uint16(inside) + totalSamples/2) / totalSamples)
			rounded.SetNRGBA(x, y, pixel)
		}
	}
	return rounded
}

func insideRoundedRectangle(x, y, width, height, radius float64) bool {
	centerX, centerY := x, y
	if centerX < radius {
		centerX = radius
	} else if centerX > width-radius {
		centerX = width - radius
	}
	if centerY < radius {
		centerY = radius
	} else if centerY > height-radius {
		centerY = height - radius
	}

	dx, dy := x-centerX, y-centerY
	return dx*dx+dy*dy <= radius*radius
}
