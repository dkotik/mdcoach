---
title: "Command-line usage"
weight: 40
---

The executable is named `mdcoach`. Run `mdcoach --help` or `mdcoach <command> --help` to see the flags available in your build.

## Compile Markdown slides

```sh
mdcoach compile [flags] file.md [more-files.md ...]
```

For one combined presentation, provide an HTML output path:

```sh
mdcoach compile --output ./talk.html ./opening.md ./conclusion.md
```

When combining files, front matter is taken from the first file; the remaining files' front matter is removed before their Markdown content is appended.

If `--output` names an existing directory, each input is compiled to a separate HTML file in that directory, using the input file's base name. Create the directory first:

```sh
mkdir -p ./public
mdcoach compile --output ./public ./talks/intro.md ./talks/results.md
```

`--output` also has the short form `-o`. Its default is the system temporary directory. Use `--open` (or `-p`) to open the generated HTML in the default browser. Relative input and output paths are resolved from the current working directory.

## Generate a PDF review sheet

Add a YAML list of questions to a Markdown file's front matter:

```yaml
---
title: Cell biology
questions:
  - What is the role of the cell membrane?
  - How does diffusion differ from osmosis?
---
```

Then run:

```sh
mdcoach review --output ./cell-biology-review.pdf --title "Cell Biology Review" --limit 12 ./lesson.md
```

Questions from multiple source files can be combined. They are shuffled before the optional `--limit` is applied. If `--output` is a `.pdf` path, that file is written directly; if it is a directory (or omitted), a dated `reviewYYYY-MM-DD.pdf` file is created there. `--open` / `-p` opens the result in the default PDF viewer.

## Other commands

- `mdcoach demo` prints a built-in Markdown presentation example to standard output.
- `mdcoach snippets pulsar` prints editor snippets for Pulsar/Atom.

## Common flags

| Flag | Commands | Purpose |
| --- | --- | --- |
| `--output`, `-o` | `compile`, `review` | Set the output file or directory. |
| `--open`, `-p` | `compile`, `review` | Open generated output after the command succeeds. |
| `--title`, `-t` | `review` | Set the title shown in the generated PDF. |
| `--limit` | `review` | Maximum number of questions to include; zero leaves the list untrimmed. |
