// Package review generates PDF review sheets from Markdown frontmatter questions.
package review

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dkotik/mdcoach"
	"github.com/yuin/goldmark/v2/ast"
)

// LoadQuestions parses sources and returns all questions from their Markdown
// frontmatter in randomized order.
func LoadQuestions(ctx context.Context, sources ...string) ([]string, error) {
	questions := make([]string, 0)
	for _, sourcePath := range sources {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		markdown, err := os.ReadFile(sourcePath)
		if err != nil {
			return nil, fmt.Errorf("read Markdown file %q: %w", sourcePath, err)
		}
		document, ok := mdcoach.NewParser().Parse(markdown).(*ast.Document)
		if !ok {
			return nil, fmt.Errorf("parser returned a non-document AST for %q", sourcePath)
		}
		fileQuestions, err := questionsFromMetadata(document.Metadata())
		if err != nil {
			return nil, fmt.Errorf("read questions from %q: %w", sourcePath, err)
		}
		questions = append(questions, fileQuestions...)
	}
	if len(questions) == 0 {
		return nil, fmt.Errorf("no questions were found in Markdown frontmatter")
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})
	return questions, nil
}

// Write parses sources, collects and shuffles their questions, and writes a
// generated PDF to w.
func Write(ctx context.Context, w io.Writer, title string, sources ...string) error {
	questions, err := LoadQuestions(ctx, sources...)
	if err != nil {
		return err
	}
	data, err := New(questions, title)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write review PDF: %w", err)
	}
	return nil
}

// WriteFile generates a review PDF and writes it to path.
func WriteFile(ctx context.Context, path, title string, sources ...string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create review PDF %q: %w", path, err)
	}
	if err := Write(ctx, file, title, sources...); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close review PDF %q: %w", path, err)
	}
	return nil
}

func questionsFromMetadata(metadata map[string]any) ([]string, error) {
	var value any
	for key, candidate := range metadata {
		if strings.EqualFold(key, "questions") {
			value = candidate
			break
		}
	}
	if value == nil {
		return nil, nil
	}

	items, ok := value.([]any)
	if !ok {
		if stringsValue, ok := value.([]string); ok {
			return stringsValue, nil
		}
		return nil, fmt.Errorf("questions must be a list, got %T", value)
	}

	questions := make([]string, 0, len(items))
	for i, item := range items {
		question, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("question %d must be a string, got %T", i+1, item)
		}
		questions = append(questions, question)
	}
	return questions, nil
}
