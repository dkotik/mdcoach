package review

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFiles(t *testing.T) {
	temp := t.TempDir()
	sources := []string{
		"presentation-1.md",
		"presentation-2.md",
		"presentation-3.md",
		"presentation-4.md",
	}

	for _, sourceName := range sources {
		source := filepath.Join("..", "testdata", sourceName)
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
	sourceNames := []string{
		"presentation-1.md",
		"presentation-2.md",
		"presentation-3.md",
		"presentation-4.md",
	}

	var questions []string
	for _, sourceName := range sourceNames {
		source := filepath.Join("..", "testdata", sourceName)
		fileQuestions, err := LoadQuestions(source)
		if err != nil {
			t.Fatal(err)
		}
		questions = append(questions, fileQuestions...)
	}
	if len(questions) < 5 {
		t.Fatalf("got %d questions, want at least 5", len(questions))
	}
}
