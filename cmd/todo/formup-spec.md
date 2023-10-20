# Opinionated Markdown

. Heading 1
.. Heading 2
... Heading 3 @#tag
.... Heading 4 @#tag
..... Heading 5 @#tag

[Image Title]@http://www.google.com#location

[Linked text]@http://...

> Quote \

Horizontal Rule
* * * *

<!-- [golang]
Language-hint aside block.

[.cssClass]
Cass class hint aside block. -->

Text types: =italic= and ==bold==. Emoji are :first: class citizens.

    Comment block. No hint, means do not display.
    <!-- indented blocks fail after lists in markdown -->

Paragraph with a footnote@1. See below how footnotes are linked.@tag Inline foot note@[Like so].

    1: Footnote block \t\w\:
    tag: cont

<!-- - Reddit >!spoilers!< ? - no just make ``` spoiler ``` block code -->
