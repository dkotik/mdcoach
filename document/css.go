package document

import "io"

func WriteCascadingStyleSheet(w io.Writer, css string) (err error) {
	if _, err = io.WriteString(w, `<style type="text/css">`); err != nil {
		return err
	}
	if err = WriteEscapedHTML(w, css); err != nil {
		return err
	}
	if _, err = io.WriteString(w, `</style>`); err != nil {
		return err
	}
	return nil
}
