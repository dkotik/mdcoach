//go:build ignore

package main

import (
	"fmt"
	"io"
	"os"
)

var scripts = []string{
	"input.js",
	"slides.js",
	"scale.js",
	"synchronize.js",
}

func makeAfter(w io.Writer) error {
	if _, err := io.WriteString(w, "</main><script>\n"); err != nil {
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
	if _, err := io.WriteString(w, "</script>\n"); err != nil {
		return fmt.Errorf("write script end tag: %w", err)
	}
	return nil
}

func main() {
	file, err := os.Create("html/after.gen.html")
	if err != nil {
		panic(fmt.Errorf("create after.gen.html: %w", err))
	}
	if err := makeAfter(file); err != nil {
		_ = file.Close()
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(fmt.Errorf("close after.gen.html: %w", err))
	}
}
