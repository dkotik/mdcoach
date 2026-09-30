package main

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestWriteEPUB(t *testing.T) {
	tests := []struct {
		name string
		page []byte
	}{
		{
			name: "complete HTML page",
			page: []byte(`<!doctype html><html><head><title>Talk</title></head><body><h1>Hello</h1></body></html>`),
		},
		{
			name: "empty HTML page",
			page: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var output bytes.Buffer
			if err := writeEPUB(&output, tt.page); err != nil {
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
			page := readEPUBEntry(t, archive.File[4])
			if !bytes.Equal(page, tt.page) {
				t.Errorf("packaged page = %q, want %q", page, tt.page)
			}
		})
	}
}

func TestCompileMarkdownToEPUB(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		wantText string
	}{
		{
			name:     "single slide",
			markdown: "# Opening\n\nWelcome to the presentation.\n",
			wantText: "Welcome to the presentation.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			directory := t.TempDir()
			source := filepath.Join(directory, "slides.md")
			if err := os.WriteFile(source, []byte(tt.markdown), 0o600); err != nil {
				t.Fatal(err)
			}
			output := filepath.Join(directory, "slides.epub")
			if err := compileMarkdownToEPUB(context.Background(), output, []string{source}, true); err != nil {
				t.Fatal(err)
			}

			data, err := os.ReadFile(output)
			if err != nil {
				t.Fatal(err)
			}
			archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
			if err != nil {
				t.Fatalf("open compiled EPUB: %v", err)
			}
			page := readEPUBEntry(t, archive.File[len(archive.File)-1])
			if !bytes.Contains(page, []byte(tt.wantText)) {
				t.Errorf("EPUB page = %q, want it to contain %q", page, tt.wantText)
			}
		})
	}
}

func TestWriteEPUBReturnsWriterError(t *testing.T) {
	wantErr := errors.New("write failed")
	err := writeEPUB(errorWriter{err: wantErr}, []byte("<html></html>"))
	if !errors.Is(err, wantErr) {
		t.Fatalf("writeEPUB() error = %v, want wrapped %v", err, wantErr)
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
