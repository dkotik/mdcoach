package main

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestCompileMarkdownToEPUB(t *testing.T) {
	tests := []struct {
		name        string
		markdown    string
		wantText    string
		wantAuthor  string
		wantCreator string
	}{
		{
			name:        "single slide with frontmatter author",
			markdown:    "---\nauthor: Ada & Lovelace\n---\n# Opening\n\nWelcome to the presentation.\n",
			wantText:    "Welcome to the presentation.",
			wantAuthor:  `name="author" content="Ada &amp; Lovelace"`,
			wantCreator: `<dc:creator>Ada &amp; Lovelace</dc:creator>`,
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
			if err := compileMarkdownToEPUB(context.Background(), output, []string{source}, true, false); err != nil {
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
			page := readCompiledEPUBEntry(t, archive.File[4])
			if !bytes.Contains(page, []byte(tt.wantText)) {
				t.Errorf("EPUB page = %q, want it to contain %q", page, tt.wantText)
			}
			if !bytes.Contains(page, []byte(tt.wantAuthor)) {
				t.Errorf("EPUB page does not contain author metadata %q", tt.wantAuthor)
			}
			packageDocument := readCompiledEPUBEntry(t, archive.File[2])
			if !bytes.Contains(packageDocument, []byte(tt.wantCreator)) {
				t.Errorf("EPUB package does not contain creator %q: %s", tt.wantCreator, packageDocument)
			}
		})
	}
}

func TestCompileMarkdownToEPUBCombinesSources(t *testing.T) {
	directory := t.TempDir()
	firstSource := filepath.Join(directory, "first.md")
	secondSource := filepath.Join(directory, "second.md")
	for source, content := range map[string]string{
		firstSource:  "# First slide\n\nFirst source.\n",
		secondSource: "# Second slide\n\nSecond source.\n",
	} {
		if err := os.WriteFile(source, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	output := filepath.Join(directory, "slides.epub")
	if err := compileMarkdownToEPUB(context.Background(), output, []string{firstSource, secondSource}, true, false); err != nil {
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
	page := readCompiledEPUBEntry(t, archive.File[4])
	for _, expected := range []string{"First slide", "First source.", "Second slide", "Second source."} {
		if !bytes.Contains(page, []byte(expected)) {
			t.Errorf("EPUB page does not contain %q: %s", expected, page)
		}
	}
}

func readCompiledEPUBEntry(t *testing.T, file *zip.File) []byte {
	t.Helper()

	reader, err := file.Open()
	if err != nil {
		t.Fatalf("open EPUB entry %q: %v", file.Name, err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read EPUB entry %q: %v", file.Name, err)
	}
	return content
}
