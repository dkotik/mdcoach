---
title: "Development"
weight: 60
---

- [ ] updateMainMargin in [@slide-notes.js] updates the wrong slide.
- [ ] Safari sizing is off.
- [ ] https://github.com/gpdf-dev/gpdf for PDF generation // there is already a goldmark-pdf extension - try that first?
- [ ] External PDF to Markdown coverter: https://github.com/VikParuchuri/marker (also nougat)
- [ ] Complete EPUB rendering by stripping scripts and styles | Support EPUB notes output: https://willcrichton.net/notes/portable-epubs/#epub-content%2FEPUB%2Findex.xhtml$ - it is just a ZIP bundle of HTML files with CSS and images

## Considerations

- https://keleshev.com/my-book-writing-setup/ - pandoc to pdf and to epub
- [ ] Utilize Presentation API: https://developer.mozilla.org/en-US/docs/Web/API/Presentation_API
  - [ ] html[data-presentation-receiver] to hide elements
- <https://revealjs.com/>
- <https://voussoir.net/writing/css_for_printing>
- <https://github.com/quail-ink/goldmark-enclave> - more embeds
- https://www.deckset.com/features/
- https://godoc.org/golang.org/x/tools/present
- https://casual-effects.com/markdeep/
- https://github.com/maaslalani/slides - another presenter
- release templating engine as open source sanetemplate: emoji, markdown, templating
- http://criticmarkup.com/spec.php - add criticmark support? including comments?
- document compressor: https://github.com/mzucker/noteshrink/blob/master/README.md
- https://github.com/alecthomas/chroma
- http://gravizo.com/
- gif does resize correctly; but i should probably support animated GIFs? or terrible idea?
- md to video with a synced audio track: https://www.videopuppet.com/docs/script/
- Tufte renderer? https://edwardtufte.github.io/tufte-css/
- A well-designed presentation rendered: http://bencane.com/stories/2020/07/06/how-i-structure-go-packages/#/eof-bio
- https://markodenic.com/html-tips/
- https://andybrewer.github.io/mvp/

## Ideas in the project notes

- Improve speaker-note and presenter-view workflows.
- Add more output formats, including EPUB.
- Add optional progress indicators and presentation controls.
- Extend review sheets with additional layout and question-selection options.
- Wrap presentations into executable Wails app.

## How to contribute

Bug reports, focused feature proposals, documentation improvements, and pull requests are welcome in the [GitHub repository](https://github.com/dkotik/mdcoach). For changes that affect Markdown syntax or the generated HTML, include a small example and tests where practical.
