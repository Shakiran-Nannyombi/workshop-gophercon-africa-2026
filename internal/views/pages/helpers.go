package pages

import (
	"fmt"
	"strconv"
	"strings"

	"workshop/internal/data"
)

// usdFor converts an amount of KES to a rounded USD string using the configured rate.
func usdFor(s *data.Store, kes float64) string {
	r := s.Site.Rates.PerUSD["KES"]
	if r == 0 {
		return "?"
	}
	return fmt.Sprintf("%.2f", kes/r)
}

func citiesIn(s *data.Store, country string) []data.City {
	var out []data.City
	for _, c := range s.Cities {
		if c.Country == country {
			out = append(out, c)
		}
	}
	return out
}

// Commas formats 1234567 as "1,234,567".
func Commas(n int) string {
	s := strconv.Itoa(n)
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	if neg {
		return "-" + b.String()
	}
	return b.String()
}
