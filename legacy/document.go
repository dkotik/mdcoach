package mdcoach

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"path"
	"path/filepath"

	blackfriday "github.com/russross/blackfriday/v2"
)

// Extensions sets up proper markdown parsing options.
const Extensions = blackfriday.CommonExtensions | blackfriday.Footnotes | blackfriday.AutoHeadingIDs

// TODO: this should be part of sanetemplate/markdown sub package?

// Document models a markdown file.
type Document struct {
	Tree           *blackfriday.Node
	Footnotes      *blackfriday.Node
	totalFootnotes int
}

func (doc *Document) parseFile(file, relativePath string, l func(f string) ([]byte, error)) {
	var err error
	var b []byte
	if l == nil {
		err = fmt.Errorf(`could not load "%s", because access to the file system was not provided`, file)
	} else {
		b, err = l(filepath.Join(relativePath, file))
	}
	if err != nil {
		log.Printf(`Could not find file %s: %s.`, file, err.Error())
		doc.Tree.AppendChild(&blackfriday.Node{
			Type:    blackfriday.HTMLBlock,
			Literal: []byte(fmt.Sprintf(`<aside class="error">Could not find file %s: %s.</aside>`, file, err.Error()))})
	}
	doc.Parse(b, relativePath, l)
}

// Parse builds up markdown tree from provided bytes. Can be called multiple times.
func (doc *Document) Parse(b []byte, relativePath string, l func(f string) ([]byte, error)) {
	adjustPath := func(p []byte) []byte {
		if len(p) > 0 && bytes.Index(p, []byte(`://`)) == -1 && p[0] != '/' {
			if p[0] == '~' && p[1] == '/' { // this should be handled by the file system loader!!
				return []byte(path.Clean(path.Join(userHomePath, string(p[2:]))))
			}
			return []byte(path.Clean(path.Join(relativePath, string(p))))
		}
		return p
	}
	parser := blackfriday.New(blackfriday.WithExtensions(Extensions))
	tree := parser.Parse(b)
	if tree.LastChild != nil && tree.LastChild.IsFootnotesList {
		tree.LastChild.Unlink()
	}

	var current, next *blackfriday.Node
	current = tree.FirstChild
	for current != nil { // iterate through direct children in a way that allows moving them around
		next = current.Next
		doc.Tree.AppendChild(current)
		current.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
			if len(node.LinkData.Destination) > 0 && entering {
				if node.LinkData.NoteID > 0 { // adjust a footnote
					doc.totalFootnotes++
					node.LinkData.Footnote.RefLink = []byte(fmt.Sprintf(`%d-%s`, doc.totalFootnotes, node.LinkData.Footnote.RefLink))
					node.LinkData.Destination = node.LinkData.Footnote.RefLink
					doc.Footnotes.AppendChild(node.LinkData.Footnote)
					node.LinkData.NoteID = doc.totalFootnotes
					return blackfriday.GoToNext
				} else if node.Type == blackfriday.Image && node.Parent == current && bytes.HasSuffix(node.LinkData.Destination, []byte(`.md`)) {
					incl := string(node.LinkData.Destination)
					doc.parseFile(path.Base(incl), path.Clean(path.Join(relativePath, path.Dir(incl))), l)
					current.Unlink() // remove the actual *.md image node
					return blackfriday.SkipChildren
				}
				node.LinkData.Destination = adjustPath(node.LinkData.Destination)
			}
			return blackfriday.GoToNext
		})
		current = next
	}
}

// Load builds up markdown tree from provided sources. Can be called multiple times.
func (doc *Document) Load(loader func(file string) ([]byte, error), sources ...string) {
	for _, src := range sources {
		// // TODO: i can probably calculate relative path without having to pass it here?
		doc.parseFile(filepath.Base(src), filepath.Dir(src), loader)
	}
}

// Render generates HTML from the document Markdown tree.
func (doc *Document) Render(w io.Writer, r blackfriday.Renderer) {
	if doc.totalFootnotes > 0 && doc.Tree.LastChild != doc.Footnotes {
		doc.Tree.AppendChild(doc.Footnotes) // move footnotes to the end
	}
	r.RenderHeader(w, doc.Tree)
	doc.Tree.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return r.RenderNode(w, node, entering)
	})
	r.RenderFooter(w, doc.Tree)
}

// NewDocument builds a markdown tree with includes from provided file system.
func NewDocument() (doc *Document) {
	doc = new(Document)
	doc.Tree = &blackfriday.Node{Type: blackfriday.Document}
	doc.Footnotes = &blackfriday.Node{Type: blackfriday.List}
	doc.Footnotes.IsFootnotesList = true
	doc.Footnotes.ListData.ListFlags = blackfriday.ListTypeOrdered
	return doc
}
