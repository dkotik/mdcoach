---
title: "Markdown behavior"
weight: 50
---

MdCoach uses Markdown as the source format, with YAML front matter for document metadata. The command-line compiler and Go parser provide standard Markdown features along with extensions such as GitHub-style tables, task lists, strikethrough, footnotes, definition lists, and typographic substitutions.

## Front matter

Place metadata between matching `---` lines at the beginning of a Markdown file:

```yaml
---
title: Photosynthesis
description: An introduction to how plants use light.
author: A. Student
questions:
  - What inputs are needed for photosynthesis?
  - Where in a plant cell does it take place?
---
```

`title`, `description`, and `author` are useful presentation metadata. `questions` is a list of strings used by the `review` command to create a worksheet. When multiple files are combined with `compile`, metadata from the first file is retained.

## Slide boundaries

A level-one or level-two heading starts a new slide. The heading is part of the new slide. Use a thematic break such as `---` when you want to split content without adding a heading:

```markdown
# First slide

Content for the first slide.

---

Content for the next slide.
```

In the command-line presentation compiler, an asterisk-only thematic break such as `***` separates speaker notes from slide content. The notes are kept separately from the visible slide. A later slide boundary starts a new slide and ends the previous notes section.

## Images

Use standard Markdown image syntax:

```markdown
![Alt text](images/diagram.png "Optional title")
```

Relative image paths are resolved from the source file's directory. HTTP and HTTPS image URLs are also supported; those images are fetched while generating the presentation. Images are resized to configured limits and included in the generated output/cache.

## Center Emphasis

The text of any paragraph containing only emphasized text by itself is centered. In the example below, the text in second paragraph will be centered.

```markdown
paragraph1

**paragraph2**

paragraph3
```


## Review questions

The `review` command reads a case-insensitive `questions` key from front matter. Its value must be a YAML list of strings. Questions from all provided files are combined, shuffled, and optionally shortened with `--limit` before a PDF worksheet is written.

## Go parser transformations

The Goldmark v2 parser exposed by `mdcoach.NewParser` turns top-level image-only paragraphs into figures, parses `:name:` emote tokens, recognizes aside blocks, and groups document sections into slide nodes. The renderer supplied by `mdcoach.NewRenderer` converts those presentation nodes to HTML.
