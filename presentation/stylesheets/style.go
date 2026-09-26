//go:build ignore

package main

import (
	"fmt"
	"io"
	"os"
)

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

func makeStylesheet(w io.Writer) error {
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
	return nil
}

func main() {
	file, err := os.Create("stylesheets/style.gen.css")
	if err != nil {
		panic(fmt.Errorf("create style.gen.css: %w", err))
	}
	if err := makeStylesheet(file); err != nil {
		_ = file.Close()
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(fmt.Errorf("close style.gen.css: %w", err))
	}
}
