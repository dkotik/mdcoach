## Primary

- [ ] implement aside.go
- [ ] Support EPUB notes output: https://willcrichton.net/notes/portable-epubs/#epub-content%2FEPUB%2Findex.xhtml$ - it is just a ZIP bundle of HTML files with CSS and images
- [ ] review should allow a percentage or count of bonus questions as a flag
- [ ] add header insertion for review
- [ ] native PDF creation (might have a problem with tables):
  - [ ] https://github.com/ebuckley/write/blob/main/write/lib/pdf.go
- parsers and renderers should be paired
  - create new node kind for SlideCut instead of NotesBreak
- [ ] External PDF to Markdown coverter: https://github.com/VikParuchuri/marker (also nougat)
- Toggleable progress line at the bottom of the screen
- support figure with a footnote! syntax
- include documentation into .cache dir for all cache distributions?
- log errors to .cache directory - check last connection error on repeat program runs
- **Что общего между дарами? Золотой цвет.** in paragraphy by itself - center text?
- Utilize Presentation API: https://developer.mozilla.org/en-US/docs/Web/API/Presentation_API
- double HR as early end of the presentation?
- <blockquote><footer> instead of <cite>? https://andybrewer.github.io/mvp/
- Reddit >!spoilers!< ? - no just make `spoiler` block code
  --format command something like CLI.PersistentFlags().StringArrayP(`format`) ?

      is there a way to check this by Javascript?
      What you need to do is type about:config Into the address bar, search for security.fileuri.strict_origin_policy and double click it / disable it. (set it to false)

- https://keleshev.com/my-book-writing-setup/ - pandoc to pdf and to epub
  > to solve above issue, use webpack to bundle in fonts | or packer JS?
  > https://survivejs.com/webpack/loading/fonts/
- replace webp with MozJPEG: "I think MozJPEG is the clear winner here with consistently about 10% better compression than libjpeg." webp only worked better for images under 500px in size https://siipo.la/blog/is-webp-really-better-than-jpeg

## Considerations

- Markdown javascript mind map
- <https://voussoir.net/writing/css_for_printing>
- scalp for features?
  - <https://github.com/quail-ink/goldmark-enclave> - more embeds
  - https://www.deckset.com/features/
  - https://godoc.org/golang.org/x/tools/present
  - https://casual-effects.com/markdeep/
  - https://github.com/maaslalani/slides - another presenter
- os.UserConfigDir()
- <script type="module" src="/src/app.js"></script> // ESM module syntax supported in all but IE11
- use `...` for front-matter termination?
- release templating engine as open source sanetemplate: emoji, markdown, templating
- http://criticmarkup.com/spec.php - add criticmark support? including comments?
- document compressor: https://github.com/mzucker/noteshrink/blob/master/README.md
- https://github.com/alecthomas/chroma
- http://gravizo.com/
- gif does resize correctly; but i should probably support animated GIFs? or terrible idea?
- instead of "cache" directory embed files as data urls? like here: https://github.com/Y2Z/monolith
- md to video with a synced audio track: https://www.videopuppet.com/docs/script/

// TODO: allow stylesheet override?
// if \_, ok := meta[`stylesheet`]; !ok {
// meta[`stylesheet`] = styleSheet
// }

- Tufte renderer? https://edwardtufte.github.io/tufte-css/
- A well-designed presentation rendered: http://bencane.com/stories/2020/07/06/how-i-structure-go-packages/#/eof-bio
- https://markodenic.com/html-tips/
