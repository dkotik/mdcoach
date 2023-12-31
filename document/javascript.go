package document

import (
	"io"
	"strings"
)

// WriteJavascriptModuleES6 embeds an ES6 Javascript module into an [io.Writer]. ES6 modules cannot be loaded from local file system due to browser CORS restrictions. The `<script>` tag must be designated as a module.
func WriteJavascriptModuleES6(w io.Writer, js string) (err error) {
	if _, err = io.WriteString(w, `<script type="module">`); err != nil {
		return err
	}
	// TODO: proper escaping!
	// if err = WriteEscapedHTML(w, js); err != nil {
	// 	return err
	// }
	if _, err = io.Copy(w, strings.NewReader(js)); err != nil {
		return err
	}
	if _, err = io.WriteString(w, `</script>`); err != nil {
		return err
	}
	return nil
}
