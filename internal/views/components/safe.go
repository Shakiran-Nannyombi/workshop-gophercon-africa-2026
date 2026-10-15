package components

import "regexp"

var hexColour = regexp.MustCompile(`^#[0-9a-fA-F]{3,8}$`)

// safeColour lets only hex colours from the data files into inline styles.
func safeColour(c string) string {
	if hexColour.MatchString(c) {
		return c
	}
	return "#999999"
}

// SafeColour is safeColour for other packages.
func SafeColour(c string) string { return safeColour(c) }
