//go:build ignore

package main

import (
	"fmt"
	"io"
	"os"
)

var scripts = []string{
	"debounce.js",
	"components/all.js",
	"components/keystroke-combo.js",
	"components/presentation-menu.js",
	"components/presentation-side-button.js",
	"components/dark-light-toggle.js",
	"components/fullscreen-toggle.js",
	"components/presentation-clock.js",
	"components/presentation-timer.js",
	"components/presentation-progress.js",
}

func makeBefore(w io.Writer) error {
	if _, err := io.WriteString(w, "<script>\n"); err != nil {
		return fmt.Errorf("write script start tag: %w", err)
	}
	for _, entry := range scripts {
		content, err := os.ReadFile("javascript/" + entry)
		if err != nil {
			return fmt.Errorf("read javascript file %q: %w", entry, err)
		}
		if _, err := w.Write(content); err != nil {
			return fmt.Errorf("write javascript file %q: %w", entry, err)
		}
		if len(content) == 0 || content[len(content)-1] != '\n' {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return fmt.Errorf("write newline after javascript file %q: %w", entry, err)
			}
		}
	}
	if _, err := io.WriteString(w, "</script><main>\n"); err != nil {
		return fmt.Errorf("write script end and main start: %w", err)
	}
	return nil
}

func main() {
	file, err := os.Create("html/before.gen.html")
	if err != nil {
		panic(fmt.Errorf("create before.gen.html: %w", err))
	}
	if err := makeBefore(file); err != nil {
		_ = file.Close()
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(fmt.Errorf("close before.gen.html: %w", err))
	}
}
