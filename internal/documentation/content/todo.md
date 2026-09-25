---
title: "Roadmap to v1.0.0"
weight: 60
---

- **Что общего между дарами? Золотой цвет.** in paragraphy by itself - center text?

```css
/*
  Stack children vertically on each other and center
  them into the middle of the grid:

  https://www.youtube.com/watch?v=CU8Plk-53RU
*/
.pile {
  display: grid;
  grid-template-areas: "pile";
  place-items: "center";

  > * {
    grid-area: "pile";
  }
}

.featured-card-content {
  align-self: start end; /* hug top right corner */
  align-self: start; /* hug top */
  align-self: end; /* hug bottom */
  background-image: linear-gradient(to top, rgba(0,0,0,0.8), rgba(0,0,0,0.0));
}
```

- add "g11" key stroke combination jumper to any slide number
- <blockquote><footer> instead of <cite>? https://andybrewer.github.io/mvp/
- [ ] https://github.com/gpdf-dev/gpdf for PDF generation // there is already a goldmark-pdf extension - try that first?
- [ ] implement aside.go
- [ ] Support EPUB notes output: https://willcrichton.net/notes/portable-epubs/#epub-content%2FEPUB%2Findex.xhtml$ - it is just a ZIP bundle of HTML files with CSS and images
- [ ] add header insertion for review
- parsers and renderers should be paired
  - create new node kind for SlideCut instead of NotesBreak
- [ ] External PDF to Markdown coverter: https://github.com/VikParuchuri/marker (also nougat)
- Toggleable progress line at the bottom of the screen
- support figure with a footnote! syntax
- include documentation into .cache dir for all cache distributions?
- Utilize Presentation API: https://developer.mozilla.org/en-US/docs/Web/API/Presentation_API
- double HR as early end of the presentation?
- https://keleshev.com/my-book-writing-setup/ - pandoc to pdf and to epub
  > to solve above issue, use webpack to bundle in fonts | or packer JS?
  > https://survivejs.com/webpack/loading/fonts/

## Considerations

- Markdown javascript mind map
- <https://voussoir.net/writing/css_for_printing>
- scalp for features?
  - <https://github.com/quail-ink/goldmark-enclave> - more embeds
  - https://www.deckset.com/features/
  - https://godoc.org/golang.org/x/tools/present
  - https://casual-effects.com/markdeep/
  - https://github.com/maaslalani/slides - another presenter
- use `...` for front-matter termination?
- release templating engine as open source sanetemplate: emoji, markdown, templating
- http://criticmarkup.com/spec.php - add criticmark support? including comments?
- document compressor: https://github.com/mzucker/noteshrink/blob/master/README.md
- https://github.com/alecthomas/chroma
- http://gravizo.com/
- gif does resize correctly; but i should probably support animated GIFs? or terrible idea?
- md to video with a synced audio track: https://www.videopuppet.com/docs/script/

// TODO: allow stylesheet override?
// if \_, ok := meta[`stylesheet`]; !ok {
// meta[`stylesheet`] = styleSheet
// }

- Tufte renderer? https://edwardtufte.github.io/tufte-css/
- A well-designed presentation rendered: http://bencane.com/stories/2020/07/06/how-i-structure-go-packages/#/eof-bio
- https://markodenic.com/html-tips/

MdCoach is under active development. The items below are areas for future exploration, not release promises or scheduled milestones.

## Ideas in the project notes

- Improve speaker-note and presenter-view workflows.
- Add more output formats, including EPUB.
- Make slide themes and styles easier to override.
- Add optional progress indicators and presentation controls.
- Extend review sheets with additional layout and question-selection options.
- Continue refining image, animation, and figure handling.

## How to contribute

Bug reports, focused feature proposals, documentation improvements, and pull requests are welcome in the [GitHub repository](https://github.com/dkotik/mdcoach). For changes that affect Markdown syntax or the generated HTML, include a small example and tests where practical.
