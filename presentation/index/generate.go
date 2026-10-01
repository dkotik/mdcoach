//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"os"
)

const scriptPlaceholder = "__DARK_LIGHT_TOGGLE_SCRIPT__"

func main() {
	page, err := os.ReadFile("index.html")
	if err != nil {
		panic(fmt.Errorf("read index HTML template: %w", err))
	}
	script, err := os.ReadFile("../javascript/components/dark-light-toggle.js")
	if err != nil {
		panic(fmt.Errorf("read dark-light-toggle component: %w", err))
	}
	if count := bytes.Count(page, []byte(scriptPlaceholder)); count != 1 {
		panic(fmt.Errorf("index HTML template has %d script placeholders, want 1", count))
	}

	generated := bytes.Replace(page, []byte(scriptPlaceholder), script, 1)
	if err := os.WriteFile("index.gen.html", generated, 0o644); err != nil {
		panic(fmt.Errorf("write generated index HTML: %w", err))
	}
}
