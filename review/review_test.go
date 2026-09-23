package review

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFiles(t *testing.T) {
	sources := []string{
		filepath.Join("..", "testdata", "presentation-1.md"),
		filepath.Join("..", "testdata", "presentation-2.md"),
		filepath.Join("..", "testdata", "presentation-3.md"),
		filepath.Join("..", "testdata", "presentation-4.md"),
	}

	for _, source := range sources {
		name := filepath.Base(source) + ".pdf"
		output := filepath.Join(t.TempDir(), name)
		if err := WriteFile(context.Background(), output, "Review", source); err != nil {
			t.Fatalf("write review for %q: %v", source, err)
		}

		data, err := os.ReadFile(output)
		if err != nil {
			t.Fatalf("read generated review for %q: %v", source, err)
		}
		if len(data) < 5 || string(data[:5]) != "%PDF-" {
			t.Fatalf("generated review for %q is not a PDF", source)
		}
	}
}

func TestQuestionsFromAllSources(t *testing.T) {
	sources := []string{
		filepath.Join("..", "testdata", "presentation-1.md"),
		filepath.Join("..", "testdata", "presentation-2.md"),
		filepath.Join("..", "testdata", "presentation-3.md"),
		filepath.Join("..", "testdata", "presentation-4.md"),
	}

	questions, err := LoadQuestions(context.Background(), sources...)
	if err != nil {
		t.Fatal(err)
	}
	if len(questions) < 5 {
		t.Fatalf("got %d questions, want at least 5", len(questions))
	}
}
