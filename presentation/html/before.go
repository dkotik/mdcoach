//go:build ignore

package main

import (
	"fmt"
	"io"
	"os"
)

var styleSheets = []string{
	"layout.css",
	"section.css",
	"image.css",
}

func makeBefore(w io.Writer) error {
	if _, err := io.WriteString(w, "<style>\n"); err != nil {
		return fmt.Errorf("write style start tag: %w", err)
	}
	for _, name := range styleSheets {
		content, err := os.ReadFile("stylesheets/" + name)
		if err != nil {
			return fmt.Errorf("read stylesheet %q: %w", name, err)
		}
		if _, err := w.Write(content); err != nil {
			return fmt.Errorf("write stylesheet %q: %w", name, err)
		}
		if len(content) == 0 || content[len(content)-1] != '\n' {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return fmt.Errorf("write newline after stylesheet %q: %w", name, err)
			}
		}
	}
	if _, err := io.WriteString(w, "</style><main>\n"); err != nil {
		return fmt.Errorf("write style end tag: %w", err)
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
