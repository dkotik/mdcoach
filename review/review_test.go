package review

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFiles(t *testing.T) {
	temp := t.TempDir()
	sources := []string{
		filepath.Join("..", "testdata", "presentation-1.md"),
		filepath.Join("..", "testdata", "presentation-2.md"),
		filepath.Join("..", "testdata", "presentation-3.md"),
		filepath.Join("..", "testdata", "presentation-4.md"),
	}

	for _, source := range sources {
		questions, err := LoadQuestions(source)
		if err != nil {
			t.Fatalf("load questions from %q: %v", source, err)
		}

		name := filepath.Base(source) + ".pdf"
		output := filepath.Join(temp, name)
		file, err := os.Create(output)
		if err != nil {
			t.Fatalf("create review for %q: %v", source, err)
		}
		if err := Write(file, Page{Title: "Review", Questions: questions}); err != nil {
			_ = file.Close()
			t.Fatalf("write review for %q: %v", source, err)
		}
		if err := file.Close(); err != nil {
			t.Fatalf("close review for %q: %v", source, err)
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

	questions, err := LoadQuestions(sources...)
	if err != nil {
		t.Fatal(err)
	}
	if len(questions) < 5 {
		t.Fatalf("got %d questions, want at least 5", len(questions))
	}
}
