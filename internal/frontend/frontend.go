package frontend

import (
	_ "embed" // for assets
)

// CSS holds all the frontend CSS.
//go:embed public/build/bundle.css
var CSS []byte

// JS holds all the JavaScript required to run the frontend.
//go:embed public/build/bundle.js
var JS []byte
