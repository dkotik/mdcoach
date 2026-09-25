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
}

var styleSheets = []string{
	"theme.css",
	"layout.css",
	"section.css",
	"typography.css",
	"heading.css",
	"blockquote.css",
	"code.css",
	"list_item.css",
	"definition.css",
	"table.css",
	"image.css",
	"links.css",
	"pygments.css",
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
	if _, err := io.WriteString(w, "</style><script>\n"); err != nil {
		return fmt.Errorf("write style end tag and script start: %w", err)
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
