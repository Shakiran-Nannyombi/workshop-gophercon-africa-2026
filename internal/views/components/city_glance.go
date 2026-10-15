package components

import (
	"regexp"
	"strings"
)

var iataCode = regexp.MustCompile(`\(([A-Z]{3})\)`)

// airportCode is the first IATA code in an airport description, e.g. "NBO".
func airportCode(s string) string {
	if m := iataCode.FindStringSubmatch(s); m != nil {
		return m[1]
	}
	return "---"
}

func noAirport(s string) bool { return strings.HasPrefix(strings.ToLower(s), "no major airport") }

func isApprox(s string) bool { return strings.HasPrefix(strings.TrimSpace(s), "~") }

func trimApprox(s string) string { return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "~")) }

func joinList(xs []string) string { return strings.Join(xs, " · ") }
