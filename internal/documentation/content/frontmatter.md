---
title: "Front matter"
weight: 40
---

Presentation metadata is read from YAML front matter at the start of a Markdown file. Keys are matched without regard to letter case. The fields below correspond to the `presentation.Frontmatter` values used when rendering a presentation.

## ID

Set `id` to give the presentation a document identifier. If it is omitted, MdCoach will generate a deterministic identifier based on the contents abstract structure tree.

## Title

Set `title` to provide the presentation title. It is displayed as the HTML document title.

## Description

Set `description` to provide a short description of the presentation. It is included in the generated page's description metadata.

## Keywords

Set `keywords` to provide keywords for the generated page's metadata.

## Author

Set `author` to identify the presentation's author. The Russian key `автор` is also accepted when `author` is empty.

## Created

Set `created` to specify the presentation's creation date. If it is omitted, `date` is used instead. Accepted values include RFC 3339 timestamps, dates in `YYYY-MM-DD` format, and month names such as `January 2, 2006` or `January 2006`.

## Duration

Set `duration` to specify a duration using Go's `time.ParseDuration` syntax, for example `90s`, `5m`, or `1h30m`.

## Favicon

The favicon is generated from the first figure image in the Markdown document. It is not set directly with a front-matter string.

## Stylesheet

Set `stylesheet` to the path of a CSS file to include in the generated page. Relative paths are resolved from the directory containing the Markdown source, and the file contents are embedded in the HTML header.
