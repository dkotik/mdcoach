package mdcoach

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
	"time"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"

	"github.com/microcosm-cc/bluemonday"
)

var (
	// Minifier compresses data for smaller output.
	Minifier *minify.M
	// Sanitizer cleans all HTML output of the program.
	Sanitizer *bluemonday.Policy

	userHomePath  = `~`
	slideBoundary = []byte("\n---\n")
)

// Exec runs an external command with a timeout and tracks errors.
func Exec(cmd string, args ...string) (err error) {
	// TODO: set time out in program arguments? or at least raise a proper time out error? right now just quits without doing anything
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Minute)
	defer cancel() // interrupts the process running
	p := exec.CommandContext(ctx, cmd, args...)
	// log.Printf(`Running %s %v`, cmd, args) // @TODO: manage this under verbose flag?

	var errorMessage bytes.Buffer
	p.Stdout = os.Stdout
	p.Stderr = &errorMessage
	err = p.Run()

	if err == nil && errorMessage.Len() > 0 {
		// log.Fatal(errorMessage.String())
		err = fmt.Errorf(`could not run %s command, reason: %s`, cmd, strings.TrimSpace(errorMessage.String()))
	}

	if err != nil {
		// log.Printf(`external command %s error: %s`, cmd, err.Error())
		// https://www.amazon.com/gp/feature.html?ie=UTF8&docId=1000765211 - kindlegen
		if _, err2 := exec.LookPath(`weasyprint`); err2 != nil {
			log.Printf(`Weasyprint PDF generator has not been found. Run "sudo python3 -m pip install weasyprint".`)
		}
		if _, err2 := exec.LookPath(`youtube-dl`); err2 != nil {
			// sudo pip install -U youtube-dl
			log.Printf(`Youtude-dl video downloader has not been found. Run "sudo python3 -m pip install youtube-dl".`)
		}
	}
	return err
}

// TODO: check if this is even needed - I am cleaning inside nodes now.
// IOcleanHTML : encode, sanitize, and minize data to clean HTML string.
func IOcleanHTML(in []byte) []byte {
	// return in
	b := bytes.NewBuffer(nil)
	Minifier.Minify(`text/html`, b, bytes.NewReader(Sanitizer.SanitizeBytes(in)))
	return b.Bytes()
}

func init() {
	if usr, err := user.Current(); err == nil {
		userHomePath = usr.HomeDir
	}
	Sanitizer = bluemonday.UGCPolicy()
	Sanitizer.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code", "blockquote", "figure")
	Sanitizer.AllowAttrs("target").Matching(regexp.MustCompile("^\\_blank$")).OnElements("a")
	Sanitizer.AllowAttrs("class").Matching(regexp.MustCompile(`^\w[\w\-]+\w$`)).OnElements("a")
	Sanitizer.AllowAttrs("class").Matching(regexp.MustCompile("^[a-zA-Z0-9]+$")).OnElements("aside", `ul`, `li`, `section`)
	Sanitizer.AllowAttrs("class").Matching(regexp.MustCompile("^break$")).OnElements("hr")
	Sanitizer.AllowAttrs("class").Matching(regexp.MustCompile(`^emote emote\-\w[\w\-]+\w$`)).OnElements("span")

	// Video processing
	Sanitizer.AllowElements(`video`, `source`)
	Sanitizer.AllowAttrs(`controls`).OnElements(`video`)
	Sanitizer.AllowAttrs(`src`, `type`).OnElements(`source`)

	Minifier = minify.New()
	Minifier.AddFunc(`text/html`, html.Minify)
}
