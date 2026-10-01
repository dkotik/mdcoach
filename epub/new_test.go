package epub

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"image/png"
	"io"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/dkotik/mdcoach"
)

func TestNewASTRendering(t *testing.T) {
	imageData := encodePNG(t, loaderTestImage(3, 2))
	tests := []struct {
		name       string
		source     string
		withAssets bool
		want       []string
		wantNot    []string
	}{
		{
			name:   "empty document uses default title",
			source: "",
			want:   []string{`<title>Presentation</title>`, "<body>"},
			wantNot: []string{
				"<section",
			},
		},
		{
			name: "renders Markdown and custom nodes as portable XHTML",
			source: `# EPUB & Slides

![A < cat](media/cat.png)

![Duplicate](media/cat-copy.png)

> Quote text (Ada & Bob)

- [x] Done
- [ ] Pending

~~struck~~

| Name | Value |
| --- | --- |
| A & B | cell |

See footnote[^n].

:smile:

<script>alert("unsafe")</script>

***

Speaker notes.

[^n]: Footnote body.
`,
			withAssets: true,
			want: []string{
				`<title>A &amp; &lt;Book&gt;</title>`,
				`<section id="slide-1">`,
				`<figure>`,
				`<figcaption>A &lt; cat</figcaption>`,
				`<img src="images/`,
				`alt=":smile:"`,
				`<footer>Ada &amp; Bob</footer>`,
				`[x]`,
				`[ ]`,
				`<del>struck</del>`,
				`<table>`,
				`<aside class="speaker-notes"><h2>Speaker notes</h2>`,
				`<aside class="footnote"`,
				`<a href="#fnref-fn-bg-0">↩</a>`,
				`<a href="#fn-`,
				`&lt;script&gt;alert(`,
				`&lt;/script&gt;`,
			},
			wantNot: []string{
				"data-",
				"base64,",
				"slide-notes",
				`class="grid"`,
				"<script>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			source := []byte(tt.source)
			node := mdcoach.NewParser().Parse(source)
			var options []Option
			if tt.withAssets {
				options = append(options, WithTitle("A & <Book>"), WithAuthor("Ada & Grace"))
				media := fstest.MapFS{
					"media/cat.png":      &fstest.MapFile{Data: imageData},
					"media/cat-copy.png": &fstest.MapFile{Data: imageData},
				}
				emotes := fstest.MapFS{
					"smile.png": &fstest.MapFile{Data: imageData},
				}
				options = append(options,
					WithMediaOptions(MediaOptions{FS: media}),
					WithEmoteFS(emotes),
				)
			}

			var output bytes.Buffer
			if err := New(context.Background(), &output, source, node, options...); err != nil {
				t.Fatal(err)
			}
			archive, err := zip.NewReader(bytes.NewReader(output.Bytes()), int64(output.Len()))
			if err != nil {
				t.Fatalf("open EPUB archive: %v", err)
			}

			page := readEPUBEntry(t, archive.File[4])
			assertWellFormedXML(t, page)
			for _, expected := range tt.want {
				if !bytes.Contains(page, []byte(expected)) {
					t.Errorf("XHTML does not contain %q:\n%s", expected, page)
				}
			}
			for _, unexpected := range tt.wantNot {
				if bytes.Contains(page, []byte(unexpected)) {
					t.Errorf("XHTML unexpectedly contains %q:\n%s", unexpected, page)
				}
			}

			if !tt.withAssets {
				if len(archive.File) != 6 {
					t.Fatalf("EPUB entry count = %d, want 6 (including stylesheet)", len(archive.File))
				}
				return
			}

			wantEntries := []string{
				"mimetype",
				"META-INF/container.xml",
				"OEBPS/package.opf",
				"OEBPS/nav.xhtml",
				"OEBPS/presentation.xhtml",
				"OEBPS/styles.css",
			}
			if len(archive.File) != len(wantEntries)+1 {
				t.Fatalf("EPUB entry count = %d, want %d with deduplicated image", len(archive.File), len(wantEntries)+1)
			}
			for i, name := range wantEntries {
				if archive.File[i].Name != name {
					t.Errorf("EPUB entry %d = %q, want %q", i, archive.File[i].Name, name)
				}
			}
			imageEntry := archive.File[len(archive.File)-1]
			if !strings.HasPrefix(imageEntry.Name, "OEBPS/images/") || !strings.HasSuffix(imageEntry.Name, ".png") {
				t.Fatalf("image entry = %q, want OEBPS/images/<hash>.png", imageEntry.Name)
			}
			imageFile, err := imageEntry.Open()
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := png.Decode(imageFile)
			closeErr := imageFile.Close()
			if err != nil {
				t.Fatalf("decode archived PNG: %v", err)
			}
			if closeErr != nil {
				t.Fatalf("close archived PNG: %v", closeErr)
			}
			if got := decoded.Bounds().Size(); got.X != 3 || got.Y != 2 {
				t.Errorf("archived PNG dimensions = %v, want (3, 2)", got)
			}

			packageDocument := readEPUBEntry(t, archive.File[2])
			for _, expected := range []string{
				`<dc:title>A &amp; &lt;Book&gt;</dc:title>`,
				`<dc:creator>Ada &amp; Grace</dc:creator>`,
				`href="styles.css" media-type="text/css"`,
				`href="images/` + strings.TrimSuffix(strings.TrimPrefix(imageEntry.Name, "OEBPS/images/"), ".png") + `.png" media-type="image/png"`,
			} {
				if !bytes.Contains(packageDocument, []byte(expected)) {
					t.Errorf("package document does not contain %q:\n%s", expected, packageDocument)
				}
			}
			navigation := readEPUBEntry(t, archive.File[3])
			if !bytes.Contains(navigation, []byte(`href="presentation.xhtml"`)) {
				t.Errorf("navigation does not link to the presentation: %s", navigation)
			}
		})
	}
}

func TestNewUsesFrontmatterMetadata(t *testing.T) {
	source := []byte("---\ntitle: Story title\nauthor: Ada\n---\n# Chapter\n")
	node := mdcoach.NewParser().Parse(source)
	var output bytes.Buffer
	if err := New(context.Background(), &output, source, node); err != nil {
		t.Fatal(err)
	}
	archive, err := zip.NewReader(bytes.NewReader(output.Bytes()), int64(output.Len()))
	if err != nil {
		t.Fatal(err)
	}
	packageDocument := readEPUBEntry(t, archive.File[2])
	if !bytes.Contains(packageDocument, []byte("<dc:title>Story title</dc:title>")) {
		t.Errorf("package does not use frontmatter title: %s", packageDocument)
	}
	if !bytes.Contains(packageDocument, []byte("<dc:creator>Ada</dc:creator>")) {
		t.Errorf("package does not use frontmatter author: %s", packageDocument)
	}
}

func TestNewErrors(t *testing.T) {
	tests := []struct {
		name string
		call func(*bytes.Buffer) error
	}{
		{
			name: "nil AST node",
			call: func(output *bytes.Buffer) error {
				return New(context.Background(), output, nil, nil)
			},
		},
		{
			name: "nil writer",
			call: func(_ *bytes.Buffer) error {
				return New(context.Background(), nil, nil, mdcoach.NewParser().Parse(nil))
			},
		},
		{
			name: "writer failure",
			call: func(_ *bytes.Buffer) error {
				return New(context.Background(), errorWriter{err: errors.New("write failed")}, nil, mdcoach.NewParser().Parse(nil))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.call(&bytes.Buffer{}); err == nil {
				t.Fatal("New() error = nil, want error")
			}
		})
	}
}

func assertWellFormedXML(t *testing.T, data []byte) {
	t.Helper()

	decoder := xml.NewDecoder(bytes.NewReader(data))
	for {
		if _, err := decoder.Token(); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			t.Fatalf("parse generated XHTML: %v", err)
		}
	}
}
