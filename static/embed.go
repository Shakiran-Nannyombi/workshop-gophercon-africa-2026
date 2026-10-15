// Package static embeds the CSS, JS, fonts and SVG served under /static/.
package static

import "embed"

//go:embed css js fonts svg
var FS embed.FS
