package document

import "io"

// WriteJavascriptModuleES6 embeds an ES6 Javascript module into an [io.Writer]. ES6 modules cannot be loaded from local file system due to browser CORS restrictions. The `<script>` tag must be designated as a module.
func WriteJavascriptModuleES6(w io.Writer, js string) (err error) {
	if _, err = io.WriteString(w, `<script type="module">`); err != nil {
		return err
	}
	if err = writeEscapedHTML(w, js); err != nil {
		return err
	}
	if _, err = io.WriteString(w, `</script>`); err != nil {
		return err
	}
	return nil
}
