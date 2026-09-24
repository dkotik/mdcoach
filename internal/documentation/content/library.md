---
title: "Using MdCoach as a Go library"
weight: 30
---

The `presentation` package parses Markdown sources, processes their images, and writes a complete HTML presentation to an `io.Writer`.

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/dkotik/mdcoach/presentation"
)

func main() {
    output, err := os.Create("presentation.html")
    if err != nil {
        log.Fatal(err)
    }
    defer output.Close()

    err = presentation.New(
        context.Background(),
        output,
        []string{"slides.md"},
        presentation.WithImageSizeLimit(1600, 1200),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

`presentation.New` accepts a context, destination writer, a slice of Markdown file paths, and optional settings. The default parser and renderer are used unless you pass `presentation.WithParser` or `presentation.WithRenderer`. Image handling can be configured with `presentation.WithImageSizeLimit(width, height)` and `presentation.WithImageQuality(quality)`.

The output includes the presentation shell and the styles, scripts, and image data needed by the generated slides. Local image paths are resolved relative to each Markdown source file. Remote images are downloaded during generation; the context can be canceled to stop work.

## Lower-level packages

The root `mdcoach` package exposes `NewParser` and `NewRenderer` for applications that need to work with Goldmark ASTs directly. `NewImageLoader` can process image nodes and populate an `ImageCache` before rendering. For most applications, `presentation.New` is the simplest entry point because it coordinates parsing, image loading, rendering, and writing.

The `review` package provides `LoadQuestions`, `New`, and `Write` for generating PDF worksheets from front-matter question lists. See [Command-line usage](/command.html) for a complete Markdown example.
