package pages

import (
	"fmt"
	"strconv"

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

func itoa(n int) string { return strconv.Itoa(n) }
