package mdcoach

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif" // GIF format, has to be imported so image reader can read it
	"image/jpeg"
	_ "image/png" // PNG format, has to be imported so image reader can read it
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	blackfriday "github.com/russross/blackfriday/v2"
)

const httpTimeout = time.Second * 15

var (
	httpError error // block all HTTP requests after one error
	reYoutube = regexp.MustCompile(`^https\:\/\/(youtu.be\/|www\.youtube\.com\/watch\?v\=)(?P<query>[^\"]+)$`)
)

// ImageRenderer replaces images with cached versions.
type ImageRenderer struct {
	blackfriday.Renderer
	e *Environment
}

// RenderNode allows image and youtube caching and rendering.
func (r *ImageRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	if node.Type == blackfriday.Image {
		if m := reYoutube.FindSubmatch(node.LinkData.Destination); m != nil {
			token, created := r.e.EnsureToken(string(node.LinkData.Destination))
			if created {
				Exec(`youtube-dl`, string(m[2]),
					`-o`, filepath.Join(r.e.CachePath, token+`.mp4`),
					`--format`, `mp4`,
					`--cache-dir`, r.e.CachePath)
			}
			// TODO: download and cache the video using youtube-dl? insert iframe if youtube-dl makes error
			// fmt.Fprintf(w, `<iframe type="text/html" src="https://www.youtube.com/embed/%s" frameborder="0" allow="autoplay; encrypted-media" allowfullscreen></iframe>`, m[2])
			fmt.Fprintf(w, `<video controls><source src="%s.mp4" type="video/mp4"></video>`, filepath.Join(r.e.CachePath, token))
		} else {
			token, err := r.e.ImageToken(`path`, string(node.LinkData.Destination))
			if err == nil {
				fmt.Fprintf(w, `<img src="%s.thumb.jpg" alt="`, token) // this will not modify underlying dest.
			} else {
				fmt.Printf("Could not locate image <%s>. Reason: %s.\n", node.LinkData.Destination, err.Error())
				fmt.Fprintf(w, `<img src="%s" alt="`, filepath.Join(r.e.CachePath, `img`, `notfound.jpg`))
			}
			// TODO: this is some kind of hack, may be glitch together with my other processors.
			r.Renderer.RenderNode(bytes.NewBuffer(nil), node, true) // render <img into a void
			r.Renderer.RenderNode(w, node.FirstChild, true)         // render alt
			r.Renderer.RenderNode(w, node, false)                   // render title; close tag
			// return blackfriday.GoToNext
		}
		return blackfriday.SkipChildren
	}
	return r.Renderer.RenderNode(w, node, entering)
}

// ImageToken ensures a correct pointer to a cached image.
func (e *Environment) ImageToken(source, uri string) (token string, err error) {
	var created bool
	remote := strings.HasPrefix(uri, `https://`) || strings.HasPrefix(uri, `http://`)
	if remote {
		token, created = e.EnsureToken(uri)
	} else {
		token, created = e.EnsureToken(source + uri)
	}
	if created {
		err := GenerateCachedImages(filepath.Join(e.CachePath, token), uri, remote)
		if err != nil { // roll back tokens
			e.DeleteToken(token)
			return token, err
		}
	}
	return token, nil
}

// DecodeImage reads data from uri into an image object.
func DecodeImage(uri string, local bool) (m image.Image, err error) {
	var handle io.ReadCloser
	if local {
		handle, err = os.Open(uri)
		if err != nil {
			return m, err
		}
	} else if httpError != nil {
		return nil, fmt.Errorf(`network access blocked by previous failure`)
	} else {
		resp, err2 := (&http.Client{Timeout: httpTimeout}).Get(uri)
		if err2 != nil {
			httpError = err2
			return m, httpError
		}
		handle = resp.Body
	}
	defer handle.Close()
	m, _, err = image.Decode(handle)
	return m, err
}

// GenerateCachedImages creates a cache of images from a given URI.
func GenerateCachedImages(tokenpath, source string, remote bool) error {
	m, err := DecodeImage(source, !remote)
	if err != nil {
		return err
	}
	if err = WriteImage(tokenpath+`.webp`, &m, 1080, 720, 40); err != nil {
		return err
	}
	return WriteImage(tokenpath+`.thumb.jpg`, &m, 360, 240, 40)
}

// WriteImage saves provided image to disk.
func WriteImage(path string, m *image.Image, width, height, quality uint) error {
	handle, err := os.Create(path)
	if err != nil {
		return err
	}
	defer handle.Close()
	t := resize.Thumbnail(width, height, *m, resize.Lanczos3)
	switch filepath.Ext(path) {
	case `.webp`:
		return webp.Encode(handle, *m, &webp.Options{Lossless: false, Quality: float32(quality)})
	default:
		return jpeg.Encode(handle, t, &jpeg.Options{Quality: int(quality)})
	}
}
