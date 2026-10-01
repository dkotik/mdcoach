package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dkotik/mdcoach/presentation"
	presentationindex "github.com/dkotik/mdcoach/presentation/index"
)

func TestCompileAddsOutputToIndex(t *testing.T) {
	tests := []struct {
		name        string
		extension   string
		index       bool
		prepopulate bool
		wantIndex   bool
	}{
		{
			name:      "HTML creates and populates an index",
			extension: ".html",
			index:     true,
			wantIndex: true,
		},
		{
			name:      "EPUB creates and populates an index",
			extension: ".epub",
			index:     true,
			wantIndex: true,
		},
		{
			name:        "HTML replaces existing entry while preserving other entries",
			extension:   ".html",
			index:       true,
			prepopulate: true,
			wantIndex:   true,
		},
		{
			name:      "disabled flag does not create an index",
			extension: ".html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			directory := t.TempDir()
			source := filepath.Join(directory, "slides.md")
			markdown := `---
title: Updated title
author: Ada Lovelace
description: A short introduction.
---
# Opening

Presentation body.
`
			if err := os.WriteFile(source, []byte(markdown), 0o600); err != nil {
				t.Fatal(err)
			}
			output := filepath.Join(directory, "slides"+tt.extension)
			indexPath := filepath.Join(directory, "index.html")
			if tt.prepopulate {
				if err := presentationindex.New(indexPath); err != nil {
					t.Fatal(err)
				}
				if err := presentationindex.Add(indexPath, output, presentation.Frontmatter{Title: "Old title"}); err != nil {
					t.Fatal(err)
				}
				if err := presentationindex.Add(indexPath, filepath.Join(directory, "other.html"), presentation.Frontmatter{Title: "Other entry"}); err != nil {
					t.Fatal(err)
				}
			}

			var err error
			if tt.extension == ".epub" {
				err = compileMarkdownToEPUB(context.Background(), output, []string{source}, true, tt.index)
			} else {
				err = compileMarkdownToHTML(context.Background(), output, []string{source}, true, tt.index)
			}
			if err != nil {
				t.Fatal(err)
			}

			_, err = os.Stat(indexPath)
			if !tt.wantIndex {
				if !os.IsNotExist(err) {
					t.Fatalf("index file stat error = %v, want not-exist", err)
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
			indexContent := string(content)
			want := []string{
				`<li data-path="slides` + tt.extension + `">`,
				`<a href="slides` + tt.extension + `">Updated title</a>`,
				`<small>Ada Lovelace</small>`,
				`<p>A short introduction.</p>`,
			}
			for _, expected := range want {
				if !strings.Contains(indexContent, expected) {
					t.Errorf("index does not contain %q:\n%s", expected, indexContent)
				}
			}
			if got := strings.Count(indexContent, `data-path="slides`+tt.extension+`"`); got != 1 {
				t.Errorf("index contains %d entries for output, want 1:\n%s", got, indexContent)
			}
			if tt.prepopulate {
				if strings.Contains(indexContent, "Old title") {
					t.Errorf("index still contains replaced entry:\n%s", indexContent)
				}
				if !strings.Contains(indexContent, `<a href="other.html">Other entry</a>`) {
					t.Errorf("index lost unrelated entry:\n%s", indexContent)
				}
			}
		})
	}
}
