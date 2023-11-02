package document

import (
	"io"
	"strings"
)

func WriteCascadingStyleSheet(w io.Writer, css string) (err error) {
	if _, err = io.WriteString(w, `<style type="text/css">`); err != nil {
		return err
	}
	// TODO: proper escaping, &gt; `>` breaks css.
	// if err = WriteEscapedHTML(w, css); err != nil {
	// 	return err
	// }
	if _, err = io.Copy(w, strings.NewReader(css)); err != nil {
		return err
	}
	if _, err = io.WriteString(w, `</style>`); err != nil {
		return err
	}
	return nil
}
