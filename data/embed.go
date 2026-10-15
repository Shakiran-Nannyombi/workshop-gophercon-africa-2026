// Package content embeds the JSON content files that drive Karibu Guide.
// Edit the .json files here; no template changes are needed.
package content

import "embed"

//go:embed *.json geo/*.json
var FS embed.FS
