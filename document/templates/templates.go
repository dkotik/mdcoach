/*
Package templates embeds HTML, Javascript, and CSS presentation components.
*/
package templates

import (
	_ "embed"
)

//go:embed dist/index.css
var StyleSheet string

//go:embed dist/index.js
var Javascript string
