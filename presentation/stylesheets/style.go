//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
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
	var source bytes.Buffer
	for _, name := range styleSheets {
		content, err := os.ReadFile("stylesheets/" + name)
		if err != nil {
			return fmt.Errorf("read stylesheet %q: %w", name, err)
		}
		if _, err := source.Write(content); err != nil {
			return fmt.Errorf("combine stylesheet %q: %w", name, err)
		}
		if len(content) == 0 || content[len(content)-1] != '\n' {
			if err := source.WriteByte('\n'); err != nil {
				return fmt.Errorf("add newline after stylesheet %q: %w", name, err)
			}
		}
	}

	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)
	if err := minifier.Minify("text/css", w, &source); err != nil {
		return fmt.Errorf("minify stylesheet: %w", err)
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
