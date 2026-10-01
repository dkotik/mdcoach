package epub

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"slices"
	"testing"

	"github.com/dkotik/mdcoach/presentation"
)

func TestWriteEPUB(t *testing.T) {
	tests := []struct {
		name        string
		page        []byte
		author      string
		wantCreator string
	}{
		{
			name:        "author from parsed metadata",
			page:        []byte(`<!doctype html><html><head><title>Talk</title></head><body><h1>Hello</h1></body></html>`),
			author:      "Ada & Grace",
			wantCreator: `<dc:creator>Ada &amp; Grace</dc:creator>`,
		},
		{
			name: "does not infer author from rendered HTML",
			page: []byte(`<!doctype html><html><head><meta name="author" content="HTML author"></head><body></body></html>`),
		},
		{
			name: "empty HTML page and metadata",
			page: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var output bytes.Buffer
			if err := Write(&output, tt.page, presentation.Frontmatter{Author: tt.author}); err != nil {
				t.Fatal(err)
			}

			archive, err := zip.NewReader(bytes.NewReader(output.Bytes()), int64(output.Len()))
			if err != nil {
				t.Fatalf("open EPUB archive: %v", err)
			}

			wantEntries := []string{
				"mimetype",
				"META-INF/container.xml",
				"OEBPS/package.opf",
				"OEBPS/nav.xhtml",
				"OEBPS/presentation.xhtml",
			}
			gotEntries := make([]string, len(archive.File))
			for i, file := range archive.File {
				gotEntries[i] = file.Name
			}
			if !slices.Equal(gotEntries, wantEntries) {
				t.Fatalf("EPUB entries = %v, want %v", gotEntries, wantEntries)
			}

			mimetype := archive.File[0]
			if mimetype.Method != zip.Store {
				t.Errorf("mimetype compression method = %d, want zip.Store (%d)", mimetype.Method, zip.Store)
			}
			if mimetype.Flags&0x8 != 0 || len(mimetype.Extra) != 0 {
				t.Errorf("mimetype entry has data descriptor or extra fields: flags=%#x extra=%v", mimetype.Flags, mimetype.Extra)
			}
			if got := readEPUBEntry(t, mimetype); string(got) != epubMimetype {
				t.Errorf("mimetype = %q, want %q", got, epubMimetype)
			}

			container := readEPUBEntry(t, archive.File[1])
			if !bytes.Contains(container, []byte(`full-path="OEBPS/package.opf"`)) {
				t.Errorf("container does not point to package.opf: %s", container)
			}
			packageDocument := readEPUBEntry(t, archive.File[2])
			if tt.wantCreator != "" && !bytes.Contains(packageDocument, []byte(tt.wantCreator)) {
				t.Errorf("EPUB package does not contain creator %q: %s", tt.wantCreator, packageDocument)
			}
			if tt.wantCreator == "" && bytes.Contains(packageDocument, []byte("<dc:creator>")) {
				t.Errorf("EPUB package contains unexpected creator: %s", packageDocument)
			}
			page := readEPUBEntry(t, archive.File[4])
			if !bytes.Equal(page, tt.page) {
				t.Errorf("packaged page = %q, want %q", page, tt.page)
			}
		})
	}
}

func TestWriteEPUBStripsUnmarkedScriptAndStyleElements(t *testing.T) {
	tests := []struct {
		name string
		page string
		want string
	}{
		{
			name: "removes unmarked script and style elements",
			page: `<p>keep</p><script>drop()</script><style>.hidden { display: none }</style>`,
			want: `<p>keep</p>`,
		},
		{
			name: "preserves elements with role attributes",
			page: `<script role="application/json">{"key":"value"}</script><style role="presentation">p { color: red }</style>`,
			want: `<script role="application/json">{"key":"value"}</script><style role="presentation">p { color: red }</style>`,
		},
		{
			name: "does not mistake role text in another attribute for a role attribute",
			page: `<script data-label="role">drop()</script><style class="role">drop</style><script data-role="main">drop()</script>`,
			want: ``,
		},
		{
			name: "handles mixed case and multiline elements",
			page: `<ScRiPt
 TYPE="module"
 ROLE="application/javascript">
keep()
</sCrIpT>
<StYlE
media="screen">
drop
</STYLE>`,
			want: `<ScRiPt
 TYPE="module"
 ROLE="application/javascript">
keep()
</sCrIpT>
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var output bytes.Buffer
			if err := Write(&output, []byte(tt.page), presentation.Frontmatter{}); err != nil {
				t.Fatal(err)
			}

			archive, err := zip.NewReader(bytes.NewReader(output.Bytes()), int64(output.Len()))
			if err != nil {
				t.Fatalf("open EPUB archive: %v", err)
			}
			page := readEPUBEntry(t, archive.File[len(archive.File)-1])
			if string(page) != tt.want {
				t.Errorf("packaged page = %q, want %q", page, tt.want)
			}
		})
	}
}

func TestWriteEPUBReturnsWriterError(t *testing.T) {
	wantErr := errors.New("write failed")
	err := Write(errorWriter{err: wantErr}, []byte("<html></html>"), presentation.Frontmatter{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("Write() error = %v, want wrapped %v", err, wantErr)
	}
}

func readEPUBEntry(t *testing.T, file *zip.File) []byte {
	t.Helper()

	reader, err := file.Open()
	if err != nil {
		t.Fatalf("open EPUB entry %q: %v", file.Name, err)
	}
	content, readErr := io.ReadAll(reader)
	closeErr := reader.Close()
	if readErr != nil {
		t.Fatalf("read EPUB entry %q: %v", file.Name, readErr)
	}
	if closeErr != nil {
		t.Fatalf("close EPUB entry %q: %v", file.Name, closeErr)
	}
	return content
}

type errorWriter struct {
	err error
}

func (w errorWriter) Write([]byte) (int, error) {
	return 0, w.err
}
