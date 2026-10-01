package index

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dkotik/mdcoach/presentation"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(t *testing.T, directory string) string
		wantErr error
	}{
		{
			name: "creates an empty index",
			prepare: func(_ *testing.T, directory string) string {
				return filepath.Join(directory, "index.html")
			},
		},
		{
			name: "refuses to overwrite an existing file",
			prepare: func(t *testing.T, directory string) string {
				indexPath := filepath.Join(directory, "index.html")
				if err := os.WriteFile(indexPath, []byte("keep"), 0o600); err != nil {
					t.Fatal(err)
				}
				return indexPath
			},
			wantErr: fs.ErrExist,
		},
		{
			name: "fails when the directory does not exist",
			prepare: func(_ *testing.T, directory string) string {
				return filepath.Join(directory, "missing", "index.html")
			},
			wantErr: fs.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			indexPath := tt.prepare(t, t.TempDir())
			err := New(indexPath)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("New() error = %v, want %v", err, tt.wantErr)
				}
				if errors.Is(tt.wantErr, fs.ErrExist) {
					content, readErr := os.ReadFile(indexPath)
					if readErr != nil {
						t.Fatal(readErr)
					}
					if string(content) != "keep" {
						t.Errorf("existing file content = %q, want it untouched", content)
					}
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			content, err := os.ReadFile(indexPath)
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{
				"<!DOCTYPE html>",
				`<meta charset="utf-8">`,
				"<title>Presentations</title>",
				`<ul id="` + ListID + `">` + "\n    </ul>",
			} {
				if !strings.Contains(string(content), want) {
					t.Errorf("index does not contain %q:\n%s", want, content)
				}
			}
			if strings.Contains(string(content), "<li") {
				t.Errorf("new index already contains a list item:\n%s", content)
			}
		})
	}
}

type addition struct {
	path     string // "{dir}" stands for the temporary directory
	metadata presentation.Frontmatter
}

func TestAdd(t *testing.T) {
	created := time.Date(2024, time.March, 5, 14, 30, 0, 0, time.UTC)
	tests := []struct {
		name         string
		index        string // initial index content; empty uses New
		missingIndex bool
		additions    []addition
		want         []string
		wantNot      []string
		wantOrder    []string       // substrings expected in this order
		wantCount    map[string]int // exact substring occurrences
		wantErr      error
		wantErrText  string
	}{
		{
			name: "appends a complete entry inside the list",
			additions: []addition{{
				path: "{dir}/talk.html",
				metadata: presentation.Frontmatter{
					Title:       "Opening Talk",
					Author:      "Ada Lovelace",
					Created:     created,
					Description: "Where it all begins.",
				},
			}},
			want: []string{`    <ul id="presentations">
      <li data-path="talk.html">
        <a href="talk.html">Opening Talk</a>
        <small>Ada Lovelace · <time datetime="2024-03-05">March 5, 2024</time></small>
        <p>Where it all begins.</p>
      </li>
    </ul>
  </body>`},
		},
		{
			name:      "falls back to the file name as the title",
			additions: []addition{{path: "{dir}/talks/intro.html"}},
			want: []string{`<li data-path="talks/intro.html">
        <a href="talks/intro.html">intro</a>
      </li>`},
			wantNot: []string{"<small>", "<p>", "<time"},
		},
		{
			name: "omits the author separator when only the date is known",
			additions: []addition{{
				path:     "{dir}/talk.html",
				metadata: presentation.Frontmatter{Created: created},
			}},
			want:    []string{`<small><time datetime="2024-03-05">March 5, 2024</time></small>`},
			wantNot: []string{" · "},
		},
		{
			name: "replaces the entry that points to the same path",
			additions: []addition{
				{path: "{dir}/talk.html", metadata: presentation.Frontmatter{Title: "First"}},
				{path: "{dir}/talk.html", metadata: presentation.Frontmatter{Title: "Second"}},
			},
			want:      []string{`<a href="talk.html">Second</a>`},
			wantNot:   []string{"First"},
			wantCount: map[string]int{`data-path="talk.html"`: 1},
		},
		{
			name: "keeps the position of a replaced entry among its neighbours",
			additions: []addition{
				{path: "{dir}/a.html", metadata: presentation.Frontmatter{Title: "A"}},
				{path: "{dir}/b.html", metadata: presentation.Frontmatter{Title: "B"}},
				{path: "{dir}/a.html", metadata: presentation.Frontmatter{Title: "A again"}},
			},
			wantOrder: []string{`<a href="a.html">A again</a>`, `<a href="b.html">B</a>`, "</ul>"},
			wantCount: map[string]int{`data-path="a.html"`: 1, `data-path="b.html"`: 1},
		},
		{
			name: "normalizes paths relative to the index directory",
			additions: []addition{
				{path: "{dir}/./talk.html", metadata: presentation.Frontmatter{Title: "Dotted"}},
				{path: "{dir}/talks/../talk.html", metadata: presentation.Frontmatter{Title: "Cleaned"}},
				{path: "{dir}/../outside.html"},
			},
			want: []string{
				`<a href="talk.html">Cleaned</a>`,
				`<li data-path="../outside.html">`,
				`<a href="../outside.html">outside</a>`,
			},
			wantNot:   []string{"Dotted"},
			wantCount: map[string]int{`data-path="talk.html"`: 1},
		},
		{
			name: "escapes metadata and matches escaped paths",
			additions: []addition{
				{
					path: "{dir}/a&b.html",
					metadata: presentation.Frontmatter{
						Title:       "Tom & <Jerry>",
						Author:      `"Quoted" <author>`,
						Description: "Less < more",
					},
				},
				{path: "{dir}/a&b.html", metadata: presentation.Frontmatter{Title: "Replaced"}},
			},
			want: []string{
				`<li data-path="a&amp;b.html">`,
				`<a href="a&amp;b.html">Replaced</a>`,
			},
			wantNot:   []string{"Tom", "<Jerry>", "<author>", "Less < more"},
			wantCount: map[string]int{`data-path="a&amp;b.html"`: 1},
		},
		{
			name: "replaces an item that has no closing tag",
			index: `<ul id="presentations">
<li data-path="talk.html">Old
<li data-path="other.html">Other</li>
</ul>`,
			additions: []addition{{path: "{dir}/talk.html", metadata: presentation.Frontmatter{Title: "New"}}},
			want: []string{`<ul id="presentations">
<li data-path="talk.html">
  <a href="talk.html">New</a>
</li>
<li data-path="other.html">Other</li>
</ul>`},
			wantNot: []string{"Old"},
		},
		{
			name:      "replaces an item that contains a nested list",
			index:     `<ul id="presentations"><li data-path="talk.html">Old<ul><li>nested</li></ul></li><li data-path="other.html">Other</li></ul>`,
			additions: []addition{{path: "{dir}/talk.html", metadata: presentation.Frontmatter{Title: "New"}}},
			want: []string{`<ul id="presentations"><li data-path="talk.html">
  <a href="talk.html">New</a>
</li><li data-path="other.html">Other</li></ul>`},
			wantNot: []string{"Old", "nested"},
		},
		{
			name: "replaces the last item of a list that omits closing tags",
			index: `<ul id="presentations">
  <li data-path="other.html">Other
  <li data-path="talk.html">Old
</ul>`,
			additions: []addition{{path: "{dir}/talk.html", metadata: presentation.Frontmatter{Title: "New"}}},
			want: []string{`<ul id="presentations">
  <li data-path="other.html">Other
  <li data-path="talk.html">
    <a href="talk.html">New</a>
  </li>
</ul>`},
			wantNot: []string{"Old"},
		},
		{
			name:      "appends to an inline list and preserves the rest of the document",
			index:     `<!DOCTYPE html><HTML><BODY><h1>My Talks</h1><!-- keep --><UL ID="presentations"></UL></BODY></HTML>`,
			additions: []addition{{path: "{dir}/talk.html"}},
			want: []string{`<!DOCTYPE html><HTML><BODY><h1>My Talks</h1><!-- keep --><UL ID="presentations">
  <li data-path="talk.html">
    <a href="talk.html">talk</a>
  </li>
</UL></BODY></HTML>`},
		},
		{
			name:      "appends to the list with the presentations id rather than the first list",
			index:     "<ul>\n  <li>navigation</li>\n</ul>\n<ol id=\"presentations\">\n</ol>\n",
			additions: []addition{{path: "{dir}/talk.html"}},
			want:      []string{"<ul>\n  <li>navigation</li>\n</ul>\n<ol id=\"presentations\">\n  <li data-path=\"talk.html\">\n    <a href=\"talk.html\">talk</a>\n  </li>\n</ol>\n"},
		},
		{
			name: "ignores markup inside scripts and comments",
			index: `<script>document.write('<li data-path="talk.html">scripted</li>')</script>
<!-- <li data-path="talk.html">commented</li> -->
<ul id="presentations">
</ul>`,
			additions: []addition{{path: "{dir}/talk.html", metadata: presentation.Frontmatter{Title: "Real"}}},
			want: []string{
				`<script>document.write('<li data-path="talk.html">scripted</li>')</script>`,
				`<!-- <li data-path="talk.html">commented</li> -->`,
				"<ul id=\"presentations\">\n  <li data-path=\"talk.html\">\n    <a href=\"talk.html\">Real</a>\n  </li>\n</ul>",
			},
		},
		{
			name:      "keeps CRLF line endings",
			index:     "<ul id=\"presentations\">\r\n</ul>\r\n",
			additions: []addition{{path: "{dir}/talk.html"}},
			want:      []string{"<ul id=\"presentations\">\r\n  <li data-path=\"talk.html\">\r\n    <a href=\"talk.html\">talk</a>\r\n  </li>\r\n</ul>\r\n"},
		},
		{
			name:        "fails when the index has no presentations list",
			index:       "<ul>\n  <li>navigation</li>\n</ul>\n",
			additions:   []addition{{path: "{dir}/talk.html"}},
			wantErrText: `no <ul id="presentations"> list`,
		},
		{
			name:         "fails when the index does not exist",
			missingIndex: true,
			additions:    []addition{{path: "{dir}/talk.html"}},
			wantErr:      fs.ErrNotExist,
		},
		{
			name:        "fails for an empty presentation path",
			additions:   []addition{{path: ""}},
			wantErrText: "empty presentation path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			directory := t.TempDir()
			indexPath := filepath.Join(directory, "index.html")
			switch {
			case tt.missingIndex:
			case tt.index != "":
				if err := os.WriteFile(indexPath, []byte(tt.index), 0o600); err != nil {
					t.Fatal(err)
				}
			default:
				if err := New(indexPath); err != nil {
					t.Fatal(err)
				}
			}

			var err error
			for _, a := range tt.additions {
				err = Add(indexPath, strings.ReplaceAll(a.path, "{dir}", directory), a.metadata)
				if err != nil {
					break
				}
			}
			if tt.wantErr != nil || tt.wantErrText != "" {
				if err == nil {
					t.Fatal("Add() error = nil, want error")
				}
				if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
					t.Fatalf("Add() error = %v, want %v", err, tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErrText) {
					t.Fatalf("Add() error = %q, want it to contain %q", err, tt.wantErrText)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			content, err := os.ReadFile(indexPath)
			if err != nil {
				t.Fatal(err)
			}
			index := string(content)
			for _, want := range tt.want {
				if !strings.Contains(index, want) {
					t.Errorf("index does not contain %q:\n%s", want, index)
				}
			}
			for _, unwanted := range tt.wantNot {
				if strings.Contains(index, unwanted) {
					t.Errorf("index unexpectedly contains %q:\n%s", unwanted, index)
				}
			}
			position := 0
			for _, want := range tt.wantOrder {
				offset := strings.Index(index[position:], want)
				if offset < 0 {
					t.Errorf("index does not contain %q after position %d:\n%s", want, position, index)
					break
				}
				position += offset + len(want)
			}
			for substring, want := range tt.wantCount {
				if got := strings.Count(index, substring); got != want {
					t.Errorf("index contains %q %d times, want %d:\n%s", substring, got, want, index)
				}
			}
		})
	}
}

func TestAddResolvesRelativePathsAgainstWorkingDirectory(t *testing.T) {
	directory := t.TempDir()
	t.Chdir(directory)
	if err := os.Mkdir("site", 0o750); err != nil {
		t.Fatal(err)
	}
	indexPath := filepath.Join("site", "index.html")
	if err := New(indexPath); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		path     string
		wantHref string
	}{
		{name: "inside the index directory", path: filepath.Join("site", "talks", "intro.html"), wantHref: "talks/intro.html"},
		{name: "outside the index directory", path: "shared.html", wantHref: "../shared.html"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Add(indexPath, tt.path, presentation.Frontmatter{}); err != nil {
				t.Fatal(err)
			}
			content, err := os.ReadFile(indexPath)
			if err != nil {
				t.Fatal(err)
			}
			want := `<li data-path="` + tt.wantHref + `">` + "\n        " + `<a href="` + tt.wantHref + `">`
			if !strings.Contains(string(content), want) {
				t.Errorf("index does not contain %q:\n%s", want, content)
			}
		})
	}
}
