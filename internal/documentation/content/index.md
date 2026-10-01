---
title: "Presentation index"
weight: 90
---

The `presentation/index` package creates and updates an HTML page that lists generated presentations. It is useful when publishing several presentations in the same output directory: visitors can open one `index.html` and follow links to each presentation.

## Create and update an index

`index.New(indexPath)` creates a complete, empty index page at the given path. It does not overwrite an existing file; if the path already exists, the returned error wraps `os.ErrExist`.

`index.Add(indexPath, presentationPath, metadata)` adds a list item to an existing index. The example below creates the index and adds a compiled presentation:

```go
package main

import (
	"log"
	"os"
	"time"

	"github.com/dkotik/mdcoach/presentation"
	"github.com/dkotik/mdcoach/presentation/index"
)

func main() {
	if err := os.MkdirAll("public", 0o755); err != nil {
		log.Fatal(err)
	}
	indexPath := "public/index.html"
	presentationPath := "public/talk.html" // already generated

	if err := index.New(indexPath); err != nil {
		log.Fatal(err)
	}

	metadata := presentation.Frontmatter{
		Title:       "A short talk",
		Author:      "Ada Lovelace",
		Description: "An introduction to the topic.",
		Created:     time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC),
	}
	if err := index.Add(indexPath, presentationPath, metadata); err != nil {
		log.Fatal(err)
	}
}
```

The link and its `data-path` attribute use a slash-separated path relative to the index file's directory. For the example above, the entry points to `talk.html`. `Add` does not check that the presentation file exists; compile or copy the presentation before adding it.

## Entry content and updates

Each generated `<li>` carries a `data-path` attribute identifying its presentation. The link text comes from `metadata.Title`; if it is blank, the package uses the presentation file name without its extension. The entry can also show `Author`, `Description`, and `Created` from `presentation.Frontmatter`. Metadata and paths are HTML-escaped when written.

Adding a path already present in the list replaces that `<li>` in place, so metadata can be refreshed without duplicate links. A new path is appended to the page's `<ul id="presentations">`. The page created by `New` includes this list; `Add` returns an error if an existing index has no matching list. Other page content is preserved when an entry is added or replaced.

`New` generates the page from the package's `index.html` template and embeds the shared `dark-light-toggle` component. Dark mode is the default; the visitor can switch themes, and the preference is saved in browser local storage. To regenerate the embedded page after changing its template or component, run `go generate ./presentation/index`.

Calls to `Add` for the same file must be serialized because each call reads and rewrites the whole index. The CLI handles this when compiling multiple files concurrently.

## Use from the command line

Pass `--index` when compiling presentations to create or update `index.html` in the output directory:

```sh
mkdir -p ./public
mdcoach compile --index --output ./public ./talks/intro.md ./talks/results.md
```

With a directory output, each presentation is compiled separately and linked from the directory's index. The same option works when compiling a single output file; the index is placed alongside it. See [Command-line usage](/command.html) for other compilation options.
