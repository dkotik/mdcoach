/*
Package presentation makes
HTML presentations from Markdown files.
*/
package presentation

import (
	"bytes"
	"context"
	_ "embed" // for generated HTML and CSS assets
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/dkotik/mdcoach"
	"github.com/nfnt/resize"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
)

//go:generate go run ./stylesheets/style.go
//go:embed stylesheets/style.gen.css
var defaultStylesheet []byte

//go:generate go run ./html/after.go
//go:embed html/after.gen.html
var afterMain []byte

//go:generate go run ./html/before.go
//go:embed html/before.gen.html
var beforeMain []byte

// Source contains a Markdown source file and its parsed AST.
type Source struct {
	Path         string
	Source       []byte
	Presentation ast.Node
}

// Parse reads and parses Markdown sources, extracting frontmatter from the first source.
func Parse(sources []string) ([]Source, Frontmatter, error) {
	return parseSources(sources, mdcoach.NewParser())
}

func parseSources(sources []string, markdownParser parser.Parser) ([]Source, Frontmatter, error) {
	parsedSources := make([]Source, 0, len(sources))
	metadata := Frontmatter{}
	for i, sourcePath := range sources {
		source, err := os.ReadFile(sourcePath)
		if err != nil {
			return nil, Frontmatter{}, fmt.Errorf("failed to read source file %s: %w", sourcePath, err)
		}
		tree := markdownParser.Parse(source)
		if i == 0 {
			metadata, err = frontmatterFromTree(tree, sourcePath)
			if err != nil {
				return nil, Frontmatter{}, fmt.Errorf("failed to read frontmatter from %s: %w", sourcePath, err)
			}
			metadata.Favicon, err = faviconFromFigure(tree, source, sourcePath)
			if err != nil {
				return nil, Frontmatter{}, fmt.Errorf("failed to create favicon from %s: %w", sourcePath, err)
			}
		}
		parsedSources = append(parsedSources, Source{
			Path:         sourcePath,
			Source:       source,
			Presentation: tree,
		})
	}
	return parsedSources, metadata, nil
}

// Render writes parsed sources and their metadata as an HTML document.
func Render(
	ctx context.Context,
	w io.Writer,
	sources []Source,
	metadata Frontmatter,
	withOptions ...Option,
) error {
	o, err := newOptions(withOptions)
	if err != nil {
		return err
	}
	return renderSources(ctx, w, sources, metadata, o)
}

func renderSources(
	ctx context.Context,
	w io.Writer,
	sources []Source,
	metadata Frontmatter,
	o *options,
) error {
	metadata.Stylesheet = template.CSS(
		string(defaultStylesheet) + string(metadata.Stylesheet),
	)
	if err := o.HeaderTemplate.Execute(w, metadata); err != nil {
		return fmt.Errorf("failed to render header: %w", err)
	}
	if _, err := w.Write(beforeMain); err != nil {
		return fmt.Errorf("failed to write before main: %w", err)
	}

	mediaOptions := mdcoach.MediaOptions{
		Path:        ".",
		WidthLimit:  o.ImageWidthLimit,
		HeightLimit: o.ImageHeightLimit,
		Quality:     o.ImageQuality,
	}
	for _, source := range sources {
		mediaOptions.Path = filepath.Dir(source.Path)
		if err := mdcoach.NewImageLoader(o.ImageCache, mediaOptions).LoadImages(
			ctx,
			source.Source,
			source.Presentation,
		); err != nil {
			return fmt.Errorf("failed to load images from %s: %w", source.Path, err)
		}
		if err := o.Renderer.Render(w, source.Source, source.Presentation); err != nil {
			return fmt.Errorf("failed to render file %s: %w", source.Path, err)
		}
	}

	if _, err := w.Write(afterMain); err != nil {
		return fmt.Errorf("failed to write after main: %w", err)
	}
	if err := o.ImageCache.WriteImageDataCSS(w); err != nil {
		return fmt.Errorf("failed to write image data CSS: %w", err)
	}
	if err := o.FooterTemplate.Execute(w, metadata); err != nil {
		return fmt.Errorf("failed to render footer: %w", err)
	}
	return nil
}

func newOptions(withOptions []Option) (*options, error) {
	o := &options{}
	for _, opt := range withOptions {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	if o.ImageCache == nil {
		o.ImageCache = mdcoach.NewImageCache()
	}
	if o.Renderer == nil {
		o.Renderer = mdcoach.NewRenderer(o.ImageCache)
	}
	if err := withDefaultTemplates(o); err != nil {
		return nil, err
	}
	return o, nil
}

// New parses Markdown sources and renders them as an HTML presentation.
func New(
	ctx context.Context,
	w io.Writer,
	sources []string,
	withOptions ...Option,
) error {
	o, err := newOptions(withOptions)
	if err != nil {
		return err
	}
	if o.Parser == nil {
		o.Parser = mdcoach.NewParser()
	}
	parsedSources, metadata, err := parseSources(sources, o.Parser)
	if err != nil {
		return err
	}
	return renderSources(ctx, w, parsedSources, metadata, o)
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
				return "", fmt.Errorf("unable to read figure image %q: %w", imagePath, err)
			}

			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				return "", fmt.Errorf("unable to decode figure image %q: %w", imagePath, err)
			}

			resized := resizeCover64(decoded)
			resized = roundImageBorders(resized)

			var encoded bytes.Buffer
			if err := png.Encode(&encoded, resized); err != nil {
				return "", fmt.Errorf("unable to encode favicon PNG from %q: %w", imagePath, err)
			}

			favicon := "<link rel=\"icon\" type=\"image/png\" href=\"data:image/png;base64," +
				base64.StdEncoding.EncodeToString(encoded.Bytes()) + `">`
			return template.HTML(favicon), nil
		}
	}
	return "", nil
}

func resizeCover64(img image.Image) image.Image {
	bounds := img.Bounds()
	if bounds.Dx() >= bounds.Dy() {
		img = resize.Resize(0, 64, img, resize.Lanczos3)
	} else {
		img = resize.Resize(64, 0, img, resize.Lanczos3)
	}

	bounds = img.Bounds()
	cropMin := image.Pt(
		bounds.Min.X+(bounds.Dx()-64)/2,
		bounds.Min.Y+(bounds.Dy()-64)/2,
	)
	cropped := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	draw.Draw(cropped, cropped.Bounds(), img, cropMin, draw.Src)
	return cropped
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
