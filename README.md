Markdown Coach
==============
CommonExtensions Extensions = NoIntraEmphasis | Tables | FencedCode |
    Autolink | Strikethrough | SpaceHeadings | HeadingIDs |
    BackslashLineBreak | DefinitionLists
AutoHeadingIDs // Create the heading ID from the text
Footnotes
[included] HeadingIDs // specify heading IDs  with {#id} ?
[included] BackslashLineBreak // Translate trailing backslashes into line breaks
[included] DefinitionLists // Render definition lists

## Slide Rendering

Headings of the first and second level and horizontal rules cut the content into slides. Headings of the first level also change the slide background to action color.

## Synchronization

Browser windows containing presentations with the same document title will all synchronize their current slide position using `window.localStorage` events. Open multiple windows to present the same presentation on different monitors.

Use different viewing modes to give yourself convenient access to notes. Put two windows side by side in order to see notes next to the replica of the slides shown to the audience. Press `C` button to clone the current presentation window.

## Markdown Behavior

Mdcoach departs from normal Markdown grammar in the following cases.

### Included Files
An image element in a root-level paragraph that points to a Markdown file causes the entire paragraph to be replaced with the rendered content of that file. Relative links and images are all updated according to the relative file path of the included file.

### Aside Elements
``` markdown
> > Nested double block-quote is rendered as an `<aside>` element instead of a block-quote, which then transforms the slide in the following ways.
```
1. If found at the end of slide content, they float to the right. First image of the `<aside>` block fills the right side.
2. If found after the title of the slide, they float to the left. First image of the `<aside>` block fills the left side.
3. If found after the title of the slide, but there is nothing else in the slide, they become a splash element.  First image of the `<aside>` block fills the entire screen as a background element.
4. An empty aside element `> >` confines all the following slide elements to notes.

### Double Horizontal Rule
``` markdown
- - -
- - -
```
Two horizontal rule elements are rendered as one horizontal rule `<hr class="double" />`. It will force a page break in books or handout notes.

### PDF Utilities
Running `mdcoach <path-to-pdf.pdf>` will print the extracted text from each page of the chosen PDF file. Giving multiple PDF files as arguments will result in a single merged PDF file.

### Installation
license='Custom Lincense', or use GPL license - project can't be rolled into commercial software. https://drewdevault.com/2019/06/13/My-journey-from-MIT-to-GPL.html ? or should I try to sell this?
    keywords='markdown coach', 'images/backgrounds/*'],

producing HTML presentations, notes, and tests from markdown files.

LICENSE
=======

Copyright (c) 2019 Dmitry Kotik

Property of Dmitry Kotik, which can be used only with written and digitally signed permission from Dmitry Kotik @keybase.io/dkotik.
