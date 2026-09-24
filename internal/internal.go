/*
Package internal provides common utilities for the mdcoach and related packages.
*/
package internal

import "embed"

//go:embed assets/*
var Assets embed.FS
