---
title: "MdCoach"
weight: 10
---

MdCoach turns Markdown documents into browser-based slide presentations. It is a Go command-line application and a set of Go packages for parsing, rendering, and publishing presentations.

## What it provides

- **Write slides in Markdown.** Top-level level-one and level-two headings start slides; horizontal rules provide manual breaks.
- **Keep speaker notes with the source.** The command-line compiler separates note sections from slide content.
- **Use familiar Markdown.** The parser includes tables, task lists, strikethrough, footnotes, definition lists, and typographic substitutions.
- **Handle presentation media.** Local and remote images can be loaded, resized, cached, and embedded in generated HTML.
- **Create review sheets.** A `questions` list in YAML front matter can be turned into a shuffled PDF worksheet.
- **Use it as a Go library.** The `presentation` package writes a complete HTML presentation to any `io.Writer`.

## Get started

Install the command-line tool with Go 1.25 or later:

```sh
go install github.com/dkotik/mdcoach/cmd/mdcoach@latest
```

Create `talk.md`:

```markdown
---
title: My presentation
author: Your name
---

# Opening

A slide written in Markdown.

## A second slide

Add an image with standard Markdown syntax:

![A description](images/example.png)
```

Compile it to HTML:

```sh
mdcoach compile --output talk.html talk.md
```

See [Installation](/intallation.html), [Command line](/command.html), [Go library](/library.html), and [Markdown behavior](/specification.html) for details.

## License

MdCoach is distributed under the [MIT License](https://github.com/dkotik/mdcoach/blob/main/LICENSE).
