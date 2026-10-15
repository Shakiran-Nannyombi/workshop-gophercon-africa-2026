package pages

import (
	"strconv"

	"workshop/internal/data"
)

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
