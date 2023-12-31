/*
Package document writes Markdown slides into a portable HTML bundle.
*/
package document

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"strings"
	"sync"

	"github.com/dkotik/mdcoach/document/templates"
)

var headMetaTemplate *template.Template

func WriteHeader(w io.Writer, meta *Metadata) (err error) {
	if meta == nil {
		return errors.New("cannot use a <nil> HTML metadata")
	}
	sync.OnceFunc(func() {
		headMetaTemplate = template.Must(template.New("header").Parse(`
      <title>{{ .Title }}</title>`))
	})()
	if _, err = io.Copy(w, strings.NewReader(`<!doctype html><html lang="en"><head><meta charset="UTF-8" /><meta name="viewport" content="width=device-width, initial-scale=1.0" />`)); err != nil {
		return err
	}

	if err = headMetaTemplate.Execute(w, meta); err != nil {
		return fmt.Errorf("header template execution error: %w", err)
	}
	if err = WriteCascadingStyleSheet(w, templates.StyleSheet); err != nil {
		return err
	}

	// TODO: add favicon.
	// <link rel="icon" type="image/svg+xml" href="/favicon.png" />

	if _, err = io.Copy(w, strings.NewReader(`</head><body class="dark"><div id="app"></div>`)); err != nil {
		return err
	}

	if err = WriteJavascriptModuleES6(w, templates.Javascript); err != nil {
		return err
	}

	_, err = io.Copy(w, strings.NewReader(`<div style="display: none;"><div id="slideData">[`))
	return err
}

func WriteFooter(w io.Writer) (err error) {
	_, err = io.Copy(w, strings.NewReader(`"&lt;p&gt;&amp;hellip;&lt;/p&gt;","",""]</div></div></body></html>`))
	return
}
