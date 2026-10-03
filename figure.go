package mdcoach

import (
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

// KindFigure is the node kind for image-only figures.
var KindFigure = ast.NewNodeKind("Figure")

// Figure is a figure containing a single image.
type Figure struct {
	ast.BaseBlock
}

var _ ast.Node = (*Figure)(nil)
var _ html.NodeRenderer = (*figureRenderer)(nil)
var _ html.NodeRenderer = (*videoFigureRenderer)(nil)
var _ parser.ASTTransformer = (*figureTransformer)(nil)

// Kind returns FigureKind.
func (*Figure) Kind() ast.NodeKind {
	return KindFigure
}

// Dump dumps the figure and its children.
func (f *Figure) Dump(_ []byte) *ast.NodeDump {
	return ast.NewNodeDump(f, nil)
}

type figureRenderer struct{}

// NewFigureRenderer returns a Goldmark v2 HTML renderer for Figure nodes.
func NewFigureRenderer() html.NodeRenderer {
	return &videoFigureRenderer{}
}

type videoFigureRenderer struct {
	figureRenderer
}

func (r *videoFigureRenderer) Render(
	writer io.Writer,
	source []byte,
	node ast.Node,
	entering bool,
	rc renderer.Context,
) (ast.WalkStatus, error) {
	figure, ok := node.(*Figure)
	if !ok {
		return r.figureRenderer.Render(writer, source, node, entering, rc)
	}
	image, ok := figure.FirstChild().(*ast.Image)
	if !ok {
		return r.figureRenderer.Render(writer, source, node, entering, rc)
	}
	player, ok := videoPlayerForDestination(image.Destination.Value(source))
	if !ok {
		return r.figureRenderer.Render(writer, source, node, entering, rc)
	}
	if !entering {
		return r.figureRenderer.Render(writer, source, node, entering, rc)
	}

	w, ok := writer.(util.BufWriter)
	if !ok {
		w = util.NewErrorBufWriter(writer)
	}
	if _, err := io.WriteString(w, "<figure>"); err != nil {
		return ast.WalkStop, err
	}
	if player.embedURL != "" {
		if _, err := io.WriteString(w, `<iframe src="`); err != nil {
			return ast.WalkStop, err
		}
		if _, err := html.ContextLinkURLWriter(rc).WriteString(player.embedURL); err != nil {
			return ast.WalkStop, err
		}
		if _, err := io.WriteString(w, `" title="`+player.title+`" allowfullscreen="allowfullscreen"></iframe>`); err != nil {
			return ast.WalkStop, err
		}
	} else {
		if _, err := io.WriteString(w, `<video controls="controls" preload="metadata"><source src="`); err != nil {
			return ast.WalkStop, err
		}
		if _, err := html.ContextLinkURLWriter(rc).WriteString(player.sourceURL); err != nil {
			return ast.WalkStop, err
		}
		if _, err := io.WriteString(w, `" type="`+player.mediaType+`" /></video>`); err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkSkipChildren, nil
}

type videoPlayer struct {
	embedURL  string
	title     string
	sourceURL string
	mediaType string
}

func videoPlayerForDestination(destination string) (videoPlayer, bool) {
	if html.IsDangerousURL(destination) {
		return videoPlayer{}, false
	}
	parsed, err := url.Parse(destination)
	if err != nil {
		return videoPlayer{}, false
	}

	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
	switch host {
	case "youtube.com", "m.youtube.com", "youtube-nocookie.com":
		videoID := ""
		switch parsed.Path {
		case "/watch":
			videoID = parsed.Query().Get("v")
		default:
			segments := strings.Split(strings.Trim(parsed.Path, "/"), "/")
			if len(segments) == 2 && (segments[0] == "embed" || segments[0] == "shorts" || segments[0] == "live") {
				videoID = segments[1]
			}
		}
		if validVideoID(videoID) {
			return videoPlayer{
				embedURL: "https://www.youtube-nocookie.com/embed/" + videoID,
				title:    "YouTube video player",
			}, true
		}
	case "youtu.be":
		videoID := strings.Trim(parsed.Path, "/")
		if validVideoID(videoID) {
			return videoPlayer{
				embedURL: "https://www.youtube-nocookie.com/embed/" + videoID,
				title:    "YouTube video player",
			}, true
		}
	case "vimeo.com", "player.vimeo.com":
		segments := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		videoID := ""
		if host == "player.vimeo.com" && len(segments) == 2 && segments[0] == "video" {
			videoID = segments[1]
		} else if len(segments) > 0 {
			videoID = segments[0]
		}
		if validNumericVideoID(videoID) {
			embedURL := "https://player.vimeo.com/video/" + videoID
			if privacyHash := parsed.Query().Get("h"); privacyHash != "" {
				embedURL += "?h=" + url.QueryEscape(privacyHash)
			}
			return videoPlayer{embedURL: embedURL, title: "Vimeo video player"}, true
		}
	case "dailymotion.com", "dai.ly":
		segments := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		videoID := ""
		if host == "dai.ly" && len(segments) > 0 {
			videoID = segments[0]
		} else if len(segments) == 2 && segments[0] == "video" {
			videoID = segments[1]
		}
		if validVideoID(videoID) {
			return videoPlayer{
				embedURL: "https://www.dailymotion.com/embed/video/" + videoID,
				title:    "Dailymotion video player",
			}, true
		}
	}

	if parsed.Scheme != "" && parsed.Scheme != "http" && parsed.Scheme != "https" {
		return videoPlayer{}, false
	}
	mediaTypes := map[string]string{
		".mp4":  "video/mp4",
		".m4v":  "video/mp4",
		".webm": "video/webm",
		".ogv":  "video/ogg",
		".ogg":  "video/ogg",
	}
	mediaType, ok := mediaTypes[strings.ToLower(path.Ext(parsed.Path))]
	if !ok {
		return videoPlayer{}, false
	}
	return videoPlayer{sourceURL: destination, mediaType: mediaType}, true
}

func validVideoID(id string) bool {
	if len(id) < 5 || len(id) > 64 {
		return false
	}
	for _, char := range id {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '_' || char == '-' {
			continue
		}
		return false
	}
	return true
}

func validNumericVideoID(id string) bool {
	if id == "" {
		return false
	}
	for _, char := range id {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func (*figureRenderer) Render(
	writer io.Writer,
	source []byte,
	node ast.Node,
	entering bool,
	rc renderer.Context,
) (ast.WalkStatus, error) {
	w := writer.(util.BufWriter)
	if entering {
		_, _ = w.WriteString("<figure>")
		return ast.WalkContinue, nil
	}

	image, ok := node.FirstChild().(*ast.Image)
	if ok && !image.Title.IsEmpty() {
		_, _ = w.WriteString("<figcaption>")
		if _, err := image.Title.WriteTo(html.ContextTextWriter(rc), source); err != nil {
			return ast.WalkStop, err
		}
		_, _ = w.WriteString("</figcaption>")
	}
	_, _ = w.WriteString("</figure>")
	return ast.WalkContinue, nil
}

// figureTransformer replaces top-level paragraphs containing only one image
// with Figure nodes.
type figureTransformer struct{}

// NewFigureTransformer returns an AST transformer that turns image-only
// document children into Figure nodes.
func NewFigureTransformer() parser.ASTTransformer {
	return &figureTransformer{}
}

// Transform replaces only direct Document children that are paragraphs with
// exactly one image child.
func (*figureTransformer) Transform(document *ast.Document, _ text.Reader, _ parser.Context) {
	for child := document.FirstChild(); child != nil; {
		next := child.NextSibling()
		paragraph, ok := child.(*ast.Paragraph)
		if !ok || paragraph.ChildCount() != 1 || paragraph.FirstChild().Kind() != ast.KindImage {
			child = next
			continue
		}

		image := paragraph.FirstChild()
		paragraph.RemoveChild(image)
		figure := &Figure{}
		figure.AppendChild(image)
		document.ReplaceChild(paragraph, figure)
		child = next
	}
}
