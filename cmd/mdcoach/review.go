package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dkotik/mdcoach"
	"github.com/gpdf-dev/gpdf"
	gpdfdocument "github.com/gpdf-dev/gpdf/document"
	gpdftemplate "github.com/gpdf-dev/gpdf/template"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v3"
	"github.com/yuin/goldmark/v2/ast"
)

func reviewCmd() *cli.Command {
	return &cli.Command{
		Name:  "review",
		Usage: "generate a review sheet using questions from Markdown frontmatter",
		Flags: []cli.Flag{
			outputFlag,
			openFlag,
			overwriteFlag,
			silentFlag,
			titleFlag,
		},
		Action: func(_ context.Context, c *cli.Command) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot locate working directory: %w", err)
			}

			args := c.Args().Slice()
			if len(args) == 0 {
				return errors.New("review command requires a file path to a Markdown file")
			}

			questions := make([]string, 0)
			for _, filePath := range args {
				if filepath.IsLocal(filePath) {
					filePath = filepath.Join(cwd, filePath)
				}

				markdown, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("unable to read file %q: %w", filePath, err)
				}

				document, ok := mdcoach.NewParser().Parse(markdown).(*ast.Document)
				if !ok {
					return fmt.Errorf("parser returned a non-document AST for %q", filePath)
				}
				fileQuestions, err := questionsFromMetadata(document.Metadata())
				if err != nil {
					return fmt.Errorf("read questions from %q: %w", filePath, err)
				}
				questions = append(questions, fileQuestions...)
			}

			if len(questions) == 0 {
				return errors.New("no questions were found in Markdown frontmatter")
			}
			rand.New(rand.NewSource(time.Now().UnixNano())).Shuffle(
				len(questions),
				func(i, j int) { questions[i], questions[j] = questions[j], questions[i] },
			)

			output := c.String("output")
			if filepath.IsLocal(output) {
				output = filepath.Join(cwd, output)
			}
			switch ext := strings.ToLower(filepath.Ext(output)); ext {
			case ".pdf":
			case "":
				output = filepath.Join(
					output,
					"review"+time.Now().Format("2006-01-02")+".pdf",
				)
			default:
				return fmt.Errorf("output format %q is not supported; use a .pdf file or directory", ext)
			}

			pdf, err := renderQuestionsPDF(questions, c.String("title"))
			if err != nil {
				return err
			}
			if err := os.WriteFile(output, pdf, 0o644); err != nil {
				return fmt.Errorf("write PDF %q: %w", output, err)
			}

			fmt.Println(output)
			if c.IsSet("open") {
				return open.Run("file://" + output)
			}
			return nil
		},
	}
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

func renderQuestionsPDF(questions []string, title string) ([]byte, error) {
	document := gpdf.NewDocument(
		gpdf.WithPageSize(gpdf.A4),
		gpdf.WithMargins(gpdfdocument.UniformEdges(gpdfdocument.Mm(20))),
		gpdf.WithMetadata(gpdfdocument.DocumentMetadata{Title: title}),
	)
	page := document.AddPage()
	for _, question := range questions {
		question := question
		page.AutoRow(func(row *gpdftemplate.RowBuilder) {
			row.Col(12, func(column *gpdftemplate.ColBuilder) {
				column.Text(question, gpdftemplate.FontSize(14))
			})
		})
	}
	data, err := document.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate questions PDF: %w", err)
	}
	return data, nil
}
