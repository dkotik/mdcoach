package presentation

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dkotik/mdcoach"
	"github.com/sebdah/goldie/v2"
)

func TestFrontmatterDuration(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		want      time.Duration
		wantError bool
	}{
		{name: "valid", value: "1h30m", want: 90 * time.Minute},
		{name: "invalid", value: "several hours", wantError: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			source := []byte("---\nduration: " + testCase.value + "\n---\n# Slide\n")
			tree := mdcoach.NewParser().Parse(source)
			frontmatter, err := frontmatterFromTree(tree, "")
			if testCase.wantError {
				if err == nil {
					t.Fatal("frontmatterFromTree() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if frontmatter.Duration != testCase.want {
				t.Errorf("frontmatter duration = %s, want %s", frontmatter.Duration, testCase.want)
			}
		})
	}
}

func TestNewPresentationLoadsFrontmatterStylesheet(t *testing.T) {
	directory := t.TempDir()
	stylesheet := []byte("body { color: rebeccapurple; }\n")
	if err := os.WriteFile(filepath.Join(directory, "custom.css"), stylesheet, 0o600); err != nil {
		t.Fatal(err)
	}

	source := []byte("---\nstylesheet: custom.css\n---\n# Slide\n")
	sourcePath := filepath.Join(directory, "presentation.md")
	if err := os.WriteFile(sourcePath, source, 0o600); err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	if err := New(context.Background(), &output, []string{sourcePath}); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(output.Bytes(), stylesheet) {
		t.Fatalf("rendered presentation does not contain stylesheet %q", stylesheet)
	}
}

func TestNewPresentation(t *testing.T) {
	var output bytes.Buffer
	if err := New(
		context.Background(),
		&output,
		[]string{"../testdata/presentation-1.md"},
	); err != nil {
		t.Fatal(err)
	}

	goldie.New(t, goldie.WithNameSuffix(".html")).Assert(t, "presentation", output.Bytes())
}
