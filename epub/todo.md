# EPUB AST rendering plan

## Goal and API

Add `epub.New` to render an already-parsed `ast.Node` into a complete EPUB archive written to an `io.Writer`. A practical signature is:

```go
func New(ctx context.Context, w io.Writer, source []byte, node ast.Node, options ...Option) error
```

The source bytes are needed to resolve Goldmark text segments, link destinations, image paths, and attributes; the context is needed for remote image loading. Options should provide publication metadata and image-loading configuration. Keep the current `Write` function as the lower-level archive writer, evolving it to accept the rendered XHTML, stylesheet, and asset set.

Expect callers to pass the tree produced by `mdcoach.NewParser()`. Reuse its AST transformations rather than duplicating parsing: slide cutting, image-only figure conversion, blockquote citation extraction, and per-slide footnote placement/indexing happen before rendering.

## Renderer parity and EPUB markup

Inventory every renderer registered by `mdcoach.NewRenderer` and implement an EPUB renderer/dispatch layer for those node kinds, while retaining the standard Goldmark renderers for ordinary Markdown nodes:

- Preserve normal Markdown structures (headings, paragraphs, emphasis, links, lists, code, blockquotes, thematic breaks) and the registered GFM table, strikethrough, and task-list behavior.
- Render `Slide` as a simple `<section>` containing its heading and content. Omit the browser presentation grid, layout-dependent `data-*` attributes, and custom slide controls.
- Render `Figure` as `<figure>`, an image, and optional `<figcaption>` using the image's accessible text; avoid the current nested image `<div>`.
- Render images as `<img src="images/<hash>.png" alt="...">`, preserving title and dimensions where useful. Never depend on CSS background images, custom data attributes, or data URLs.
- Render known emotes with the same PNG asset pipeline and an `<img>`/alt label; retain readable text for unknown emotes.
- Render `SlideNotes` using a common element such as `<aside>` with a visible “Speaker notes” label, not the `<slide-notes>` custom element.
- Render footnote references and definitions with ordinary `<sup>`, `<a href="#...">`, and `<aside>`/`<section>` markup. Preserve per-slide numbering and stable fragment targets while replacing browser-specific attributes.
- Keep blockquote citations as escaped `<footer>` content. Continue escaping text and attributes, and do not pass arbitrary raw HTML through unchecked; only allow the generated citation footer or render unsupported raw HTML as text.

Use a conservative subset of commonly supported EPUB HTML tags (for example `a`, `aside`, `blockquote`, `br`, `code`, `em`, `figure`, `figcaption`, `footer`, `h1`–`h6`, `img`, `li`, `ol`, `p`, `pre`, `section`, `span`, `strong`, `sup`, `table`, `tbody`, `td`, `th`, `thead`, `tr`, and `ul`). Prefer text markers over form controls for task lists if checkbox support is inconsistent.

## Image assets and EPUB assembly

1. Use the EPUB `ImageLoader` and `ImageCache` to load local/remote image destinations before serialization, resize to configured limits, and keep encoded PNG bytes raw in the cache.
2. Resolve image nodes to cached hashes and write each unique image once under `OEBPS/images/<hash>.png`; render relative `src` references and preserve alt text. Include emote assets through the same cache path.
3. Build the EPUB ZIP with `mimetype` first and uncompressed, followed by `META-INF/container.xml`, package metadata, navigation, XHTML, stylesheet, and images. List every document, stylesheet, and image in the OPF manifest and reference the reading document in the spine. Keep navigation targets valid.
4. Propagate image, renderer, ZIP, and writer errors with context; do not close the caller-owned writer.

## Implementation and verification steps

- [x] Add an EPUB-specific renderer for built-in Goldmark nodes and custom node kinds; keep rendering independent of web-presentation JavaScript and CSS.
- [x] Add options for title/author and media filesystem, dimensions, and HTTP client; define behavior for missing metadata, empty documents, and unsupported node kinds.
- [x] Add table-driven rendering coverage for custom nodes and representative standard Markdown constructs. Assert output uses EPUB-friendly syntax and escapes text and attributes.
- [x] Add end-to-end tests that build an AST with `mdcoach.NewParser()`, call `epub.New`, open the ZIP, verify required entry order and OPF/nav links, parse the XHTML, and confirm PNG entries are valid, deduplicated, and referenced with relative paths rather than base64 URLs.
- [x] Retain tests for writer failures and image-loading errors.
- [x] Verify with `go test ./...`, `go test -race ./epub`, and `git diff --check`.
