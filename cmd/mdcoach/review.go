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

	"github.com/dkotik/mdcoach/review"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v3"
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
			&cli.IntFlag{
				Name:  "limit",
				Value: 0,
				Usage: "maximum number of regular questions to include",
			},
			&cli.Uint8Flag{
				Name:  "bonus",
				Value: 2,
				Usage: "number of bonus questions to include",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot locate working directory: %w", err)
			}

			args := c.Args().Slice()
			if len(args) == 0 {
				return errors.New("review command requires a file path to a Markdown file")
			}
			for i, filePath := range args {
				if filepath.IsLocal(filePath) {
					args[i] = filepath.Join(cwd, filePath)
				}
			}

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

			questions, err := review.LoadQuestions(args...)
			if err != nil {
				return err
			}
			rand.Shuffle(len(questions), func(i, j int) {
				questions[i], questions[j] = questions[j], questions[i]
			})
			bonusCount := min(int(c.Uint8("bonus")), len(questions))
			bonusQuestions := append([]string(nil), questions[:bonusCount]...)
			questions = questions[bonusCount:]
			if limit := c.Int("limit"); limit > 0 && limit < len(questions) {
				questions = questions[:limit]
			}

			if err := confirmOverwrite(output, c.Bool("force")); err != nil {
				if errors.Is(err, errSkip) {
					return nil
				}
				return err
			}

			w, err := os.Create(output)
			if err != nil {
				return fmt.Errorf("create output file: %w", err)
			}
			defer w.Close()

			if err := review.Write(w, review.Page{
				Title:          c.String("title"),
				Description:    c.String("description"),
				Questions:      questions,
				BonusQuestions: bonusQuestions,
			}); err != nil {
				return err
			}

			fmt.Println("review:", output)
			if c.IsSet("open") {
				return open.Run("file://" + output)
			}
			return nil
		},
	}
}
