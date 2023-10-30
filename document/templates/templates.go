/*
Package templates embeds HTML, Javascript, and CSS presentation components.
*/
package templates

import (
	_ "embed"
)

//go:generate npm run build

//go:embed dist/index.css
var StyleSheet string

//go:embed dist/index.js
var Javascript string
