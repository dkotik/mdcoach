package internal

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

func NewSanitizer() *bluemonday.Policy {
	sanitizer := bluemonday.UGCPolicy()
	sanitizer.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code", "blockquote", "figure")
	sanitizer.AllowAttrs("target").Matching(regexp.MustCompile("^_blank$")).OnElements("a")
	sanitizer.AllowAttrs("class").Matching(regexp.MustCompile(`^\w[\w\-]+\w$`)).OnElements("a")
	sanitizer.AllowAttrs("class").Matching(regexp.MustCompile("^[a-zA-Z0-9]+$")).OnElements("aside", "ul", "li", "section")
	sanitizer.AllowAttrs("class").Matching(regexp.MustCompile("^break$")).OnElements("hr")
	sanitizer.AllowAttrs("class").Matching(regexp.MustCompile(`^emote emote\-\w[\w\-]+\w$`)).OnElements("span")

	// Video processing.
	sanitizer.AllowElements("video", "source")
	sanitizer.AllowAttrs("controls").OnElements("video")
	sanitizer.AllowAttrs("src", "type").OnElements("source")
	return sanitizer
}
