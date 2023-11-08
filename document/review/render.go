package review

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"text/template"
	"time"
)

//go:embed templates/render.html
var tmpl string

func (r *Review) RenderToFile(p string) (err error) {
	w, err := os.Create(p)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, w.Close())
	}()
	return r.Render(w)
}

func (r *Review) Render(w io.Writer) error {
	t, err := template.New("").Funcs(map[string]any{
		"plusOne": func(v any) int {
			i, ok := v.(int)
			if !ok {
				return 1
			}
			return i + 1
		},
	}).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("cannot prepare review template: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	missed := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	grades := make([]int, 0)
	max := len(r.questions) * 3 // max 3 points per question
	for i, v := range missed {
		grade := (max - v) * 100 / max
		if grade < 60 {
			missed = missed[:i]
			break
		}
		grades = append(grades, grade)
	}
	reverseIntSlice(missed)
	reverseIntSlice(grades)

	return t.Execute(w, struct {
		Title       string
		Description string
		Author      string
		Keywords    string
		Copyright   string
		Header      string
		Footer      string
		Created     time.Time
		Updated     time.Time
		Questions   []string
		Missed      []int
		Grades      []int
	}{
		Questions: r.questions,
		Missed:    missed,
		Grades:    grades,
	})
}

func reverseIntSlice(slice []int) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i] //reverse the slice
	}
}
