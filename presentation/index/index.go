// Package index maintains an HTML page that links to compiled presentations.
//
// New creates an empty index page. Add inserts a list entry for a presentation
// into an existing index and replaces the entry that already points to the
// same presentation, so an index can be rebuilt incrementally. Entries are
// <li> elements inside the <ul id="presentations"> list; each one carries a
// data-path attribute holding the presentation path relative to the index
// file, which identifies the entry on later updates. The rest of the file is
// preserved byte for byte, so the page may be customized by hand.
package index

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"html"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dkotik/mdcoach/presentation"
	"github.com/tdewolff/parse/v2"
	parsehtml "github.com/tdewolff/parse/v2/html"
)

// ListID is the id attribute of the list element that holds presentation
// entries. The page written by New refers to it literally in its markup and
// stylesheet.
const ListID = "presentations"

//go:embed dark-light-toggle.js
var darkLightToggleScript string

const emptyIndex = `<!DOCTYPE html>
<html lang="en" data-theme="dark">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Presentations</title>
    <style>
      :root {
        color-scheme: dark;
        font-family: system-ui, sans-serif;
        line-height: 1.5;
        --index-background: #111;
        --index-text: #eee;
        --index-muted: #aaa;
        --index-link: #8ab4f8;
      }
      :root[data-theme="light"] {
        color-scheme: light;
        --index-background: #fff;
        --index-text: #222;
        --index-muted: #666;
        --index-link: #1558b0;
      }
      body {
        background: var(--index-background);
        color: var(--index-text);
        margin: 2rem auto;
        max-width: 40rem;
        padding: 0 1rem;
      }
      .index-header {
        align-items: center;
        display: flex;
        justify-content: space-between;
      }
      .index-header h1 {
        margin: 0;
      }
      dark-light-toggle {
        color: inherit;
        font-size: 1.5rem;
      }
      dark-light-toggle:focus-visible {
        outline: 2px solid currentColor;
        outline-offset: 3px;
      }
      dark-light-toggle svg {
        height: 1em;
        width: 1em;
      }
      #presentations {
        list-style: none;
        padding: 0;
      }
      #presentations > li {
        margin: 1.5rem 0;
      }
      #presentations a {
        color: var(--index-link);
        font-size: 1.25rem;
        font-weight: 600;
      }
      #presentations small {
        color: var(--index-muted);
        display: block;
      }
      #presentations p {
        margin: 0.25rem 0 0;
      }
    </style>
  </head>
  <body>
    <header class="index-header">
      <h1>Presentations</h1>
      <dark-light-toggle aria-label="Toggle color theme" title="Toggle color theme">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" aria-hidden="true" focusable="false">
          <circle cx="12" cy="12" r="4" fill="none" stroke="currentColor" stroke-width="2" />
          <path d="M12 2v2m0 16v2M4.93 4.93l1.42 1.42m11.3 11.3 1.42 1.42M2 12h2m16 0h2M4.93 19.07l1.42-1.42m11.3-11.3 1.42-1.42" fill="none" stroke="currentColor" stroke-linecap="round" stroke-width="2" />
        </svg>
      </dark-light-toggle>
    </header>
    <ul id="presentations">
    </ul>
    <script>
{{darkLightToggleScript}}
    </script>
  </body>
</html>
`

var entryTemplate = template.Must(template.New("entry").Parse(`<li data-path="{{.Path}}">
  <a href="{{.Path}}">{{.Title}}</a>
{{- if or .Author .Created}}
  <small>{{.Author}}{{if and .Author .Created}} · {{end}}{{if .Created}}<time datetime="{{.CreatedMachine}}">{{.Created}}</time>{{end}}</small>
{{- end}}
{{- if .Description}}
  <p>{{.Description}}</p>
{{- end}}
</li>`))

type entry struct {
	Path           string
	Title          string
	Author         string
	Created        string
	CreatedMachine string
	Description    string
}

// New creates an empty index page at indexPath. It fails with an error that
// wraps [os.ErrExist] when the file already exists, so a populated index is
// never discarded; remove the file first to start over.
func New(indexPath string) error {
	file, err := os.OpenFile(indexPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("create presentation index: %w", err)
	}
	page := strings.Replace(emptyIndex, "{{darkLightToggleScript}}", darkLightToggleScript, 1)
	_, writeErr := io.WriteString(file, page)
	closeErr := file.Close()
	if writeErr != nil {
		return fmt.Errorf("write presentation index %q: %w", indexPath, writeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close presentation index %q: %w", indexPath, closeErr)
	}
	return nil
}

// Add links the presentation file at presentationPath from the index page at
// indexPath, which must already exist (see New); a missing index yields an
// error that wraps [os.ErrNotExist]. Relative paths are resolved against the
// current directory. The entry links to the presentation relative to the
// index file and records that relative path in its data-path attribute. An
// existing entry with the same data-path is replaced in place; otherwise the
// entry is appended to the <ul id="presentations"> list. Metadata supplies
// the link text (falling back to the file name) along with the author,
// creation date, and description shown in the entry.
//
// Add rewrites the whole file, so concurrent calls for the same index must be
// serialized by the caller.
func Add(indexPath, presentationPath string, metadata presentation.Frontmatter) error {
	if presentationPath == "" {
		return errors.New("add presentation to index: empty presentation path")
	}
	document, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("read presentation index: %w", err)
	}
	relativePath, err := relativePresentationPath(indexPath, presentationPath)
	if err != nil {
		return err
	}
	rendered, err := renderEntry(relativePath, metadata)
	if err != nil {
		return err
	}
	updated, err := insertEntry(document, relativePath, rendered)
	if err != nil {
		return fmt.Errorf("update presentation index %q: %w", indexPath, err)
	}
	if err := os.WriteFile(indexPath, updated, 0o644); err != nil {
		return fmt.Errorf("write presentation index: %w", err)
	}
	return nil
}

// relativePresentationPath returns the slash-separated path of the
// presentation relative to the directory that holds the index.
func relativePresentationPath(indexPath, presentationPath string) (string, error) {
	indexDirectory, err := filepath.Abs(filepath.Dir(indexPath))
	if err != nil {
		return "", fmt.Errorf("resolve index directory of %q: %w", indexPath, err)
	}
	absolutePresentationPath, err := filepath.Abs(presentationPath)
	if err != nil {
		return "", fmt.Errorf("resolve presentation path %q: %w", presentationPath, err)
	}
	relativePath, err := filepath.Rel(indexDirectory, absolutePresentationPath)
	if err != nil {
		return "", fmt.Errorf("relate presentation %q to index %q: %w", presentationPath, indexPath, err)
	}
	return filepath.ToSlash(relativePath), nil
}

func renderEntry(relativePath string, metadata presentation.Frontmatter) ([]byte, error) {
	title := strings.TrimSpace(metadata.Title)
	if title == "" {
		base := path.Base(relativePath)
		title = strings.TrimSuffix(base, path.Ext(base))
		if title == "" {
			title = base
		}
	}
	values := entry{
		Path:        relativePath,
		Title:       title,
		Author:      strings.TrimSpace(metadata.Author),
		Description: strings.TrimSpace(metadata.Description),
	}
	if !metadata.Created.IsZero() {
		values.CreatedMachine = metadata.Created.Format("2006-01-02")
		values.Created = metadata.Created.Format("January 2, 2006")
	}

	var rendered bytes.Buffer
	if err := entryTemplate.Execute(&rendered, values); err != nil {
		return nil, fmt.Errorf("render index entry for %q: %w", relativePath, err)
	}
	return rendered.Bytes(), nil
}

// insertEntry replaces the existing entry for relativePath in document or,
// when there is none, appends rendered to the presentations list. Only the
// affected byte range changes; the inserted lines adopt the indentation and
// line endings found around the insertion point.
func insertEntry(document []byte, relativePath string, rendered []byte) ([]byte, error) {
	existing, listEnd, err := locate(document, relativePath)
	if err != nil {
		return nil, err
	}
	newline := "\n"
	if bytes.Contains(document, []byte("\r\n")) {
		newline = "\r\n"
	}

	var start, end int
	var replacement strings.Builder
	if existing != nil {
		start, end = existing.start, existing.end
		indent, _ := lineIndent(document, start)
		replacement.Write(indentLines(rendered, newline, indent))
	} else {
		if listEnd < 0 {
			return nil, fmt.Errorf("no <ul id=%q> list to append to", ListID)
		}
		start, end = listEnd, listEnd
		indent, atLineStart := lineIndent(document, start)
		if !atLineStart {
			replacement.WriteString(newline + indent)
		}
		replacement.WriteString("  ")
		replacement.Write(indentLines(rendered, newline, indent+"  "))
		replacement.WriteString(newline + indent)
	}

	updated := make([]byte, 0, len(document)-(end-start)+replacement.Len())
	updated = append(updated, document[:start]...)
	updated = append(updated, replacement.String()...)
	return append(updated, document[end:]...), nil
}

func indentLines(rendered []byte, newline, indent string) []byte {
	return bytes.ReplaceAll(rendered, []byte("\n"), []byte(newline+indent))
}

// lineIndent returns the leading whitespace of the line that contains offset
// and reports whether only that whitespace precedes offset on the line.
func lineIndent(document []byte, offset int) (indent string, atLineStart bool) {
	lineStart := bytes.LastIndexByte(document[:offset], '\n') + 1
	indentEnd := lineStart
	for indentEnd < offset && (document[indentEnd] == ' ' || document[indentEnd] == '\t') {
		indentEnd++
	}
	return string(document[lineStart:indentEnd]), indentEnd == offset
}

// span is the half-open byte range [start, end) of an existing entry.
type span struct {
	start int
	end   int
}

// locate scans document for the <li> whose data-path names relativePath and
// for the end tag of the presentations list. It returns a nil span when no
// entry matches and -1 when the list or its end tag is missing. An entry
// without a closing tag ends where the next item begins, where its list
// closes, or at the end of the document.
func locate(document []byte, relativePath string) (*span, int, error) {
	// The lexer lowercases tag and attribute names in place, so lex a copy
	// and apply the offsets to the original document.
	input := parse.NewInputBytes(bytes.Clone(document))
	lexer := parsehtml.NewLexer(input)

	var (
		entry      *span
		listEnd    = -1
		inList     bool // inside the presentations list
		listDepth  int  // lists nested inside the presentations list
		inEntry    bool // inside the matched entry
		entryLists int  // lists nested inside the matched entry
		tagName    string
		tagStart   int
		tagMatches bool
	)
	for {
		tokenType, data := lexer.Next()
		end := input.Offset()
		start := end - len(data)
		switch tokenType {
		case parsehtml.ErrorToken:
			if err := lexer.Err(); !errors.Is(err, io.EOF) {
				return nil, -1, fmt.Errorf("parse index HTML: %w", err)
			}
			if inEntry {
				entry.end = len(document)
			}
			return entry, listEnd, nil
		case parsehtml.StartTagToken:
			tagName, tagStart, tagMatches = string(lexer.Text()), start, false
			switch tagName {
			case "li":
				if inEntry && entryLists == 0 {
					entry.end, inEntry = implicitEntryEnd(document, tagStart), false
				}
			case "ul", "ol":
				if inEntry {
					entryLists++
				}
				if inList {
					listDepth++
				}
			}
		case parsehtml.AttributeToken:
			key, value := string(lexer.AttrKey()), attributeValue(lexer.AttrVal())
			switch tagName {
			case "li":
				tagMatches = tagMatches ||
					key == "data-path" && entry == nil && path.Clean(value) == relativePath
			case "ul", "ol":
				tagMatches = tagMatches ||
					key == "id" && value == ListID && listEnd < 0 && !inList
			}
		case parsehtml.StartTagCloseToken, parsehtml.StartTagVoidToken:
			if !tagMatches {
				continue
			}
			switch tagName {
			case "li":
				entry, inEntry, entryLists = &span{start: tagStart, end: -1}, true, 0
			case "ul", "ol":
				inList, listDepth = true, 0
			}
		case parsehtml.EndTagToken:
			switch string(lexer.Text()) {
			case "li":
				if inEntry && entryLists == 0 {
					entry.end, inEntry = end, false
				}
			case "ul", "ol":
				if inEntry {
					if entryLists == 0 {
						entry.end, inEntry = implicitEntryEnd(document, start), false
					} else {
						entryLists--
					}
				}
				if inList {
					if listDepth == 0 {
						listEnd, inList = start, false
					} else {
						listDepth--
					}
				}
			}
		}
	}
}

// implicitEntryEnd preserves the line break before a following list item or
// list end tag when an HTML list item omits its explicit closing tag.
func implicitEntryEnd(document []byte, tagStart int) int {
	_, atLineStart := lineIndent(document, tagStart)
	if !atLineStart {
		return tagStart
	}
	lineBreak := bytes.LastIndexByte(document[:tagStart], '\n')
	if lineBreak < 0 {
		return tagStart
	}
	if lineBreak > 0 && document[lineBreak-1] == '\r' {
		return lineBreak - 1
	}
	return lineBreak
}

// attributeValue unquotes and unescapes a raw attribute value from the lexer.
func attributeValue(raw []byte) string {
	if len(raw) >= 2 && (raw[0] == '"' || raw[0] == '\'') && raw[len(raw)-1] == raw[0] {
		raw = raw[1 : len(raw)-1]
	}
	return html.UnescapeString(string(raw))
}
