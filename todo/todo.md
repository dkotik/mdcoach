## Confirm

- weasy print errors are not properly captured still!

## Primary

- Svelte dual window view using Javascript broadcasts
  - Svelte key block to replay transition: {#key ...}
    {#key expression}...{/key}
    Key blocks destroy and recreate their contents when the value of an expression changes.
- Toggleable progress line at the bottom of the screen!!
- img folder is not being included into Cache assets as of right now
- When file is not found, SaveMeta crashes:
  panic: interface conversion: interface {} is nil, not string
  mdcoach.(\*Environment).SaveMeta(0xc00992da40, 0x7fffe552474d, 0xb, 0x0, 0x0, 0xc009bde960)
  mdcoach@/environment.go:127 +0x530
  mdcoach.Paper(0xc00992da40, 0xbc9298, 0xb, 0xc001588cf0, 0x6, 0x9, 0x0, 0x0)
  mdcoach@/paper.go:128 +0x2b5
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
- scalp for features?
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
