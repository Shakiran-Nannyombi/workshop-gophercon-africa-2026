package pages

import (
	"net/url"
	"strconv"

	"workshop/internal/data"
	"workshop/internal/views/components"
)

// CompareView is one resolved pair of cities for the compare page.
type CompareView struct {
	A, B   data.City
	CA, CB *data.Country
	Notice string
	OOB    bool // also refresh the picker out of band (htmx)
}

// ComparePairs are the quick-pick chips, as city slugs.
var ComparePairs = [][2]string{
	{"nairobi", "kampala"}, {"nairobi", "kigali"}, {"mombasa", "zanzibar-city"},
	{"dar-es-salaam", "nairobi"}, {"addis-ababa", "nairobi"},
}

// NewCompare resolves two city slugs, falling back to a sensible pair with a friendly note.
func NewCompare(s *data.Store, a, b string) CompareView {
	var v CompareView
	ca, okA := s.City(a)
	if !okA {
		ca, _ = s.City("nairobi")
	}
	cb, okB := s.City(b)
	if !okB {
		cb, _ = s.City("kampala")
	}
	switch {
	case ca == nil || cb == nil:
		return v
	case (a != "" && !okA) || (b != "" && !okB):
		v.Notice = "We don’t have one of those cities in the guide yet, so we’ve filled the gap with a popular pick."
	}
	if ca.Slug == cb.Slug {
		alt := "kampala"
		if ca.Slug == alt {
			alt = "nairobi"
		}
		if c, ok := s.City(alt); ok {
			cb = c
		}
		v.Notice = "You picked " + ca.Name + " twice, so we’ve paired it with " + cb.Name + ". Choose two different cities to compare."
	}
	v.A, v.B = *ca, *cb
	v.CA, _ = s.Country(ca.Country)
	v.CB, _ = s.Country(cb.Country)
	return v
}

// CompareURL is the canonical address of a comparison.
func CompareURL(a, b string) string {
	return "/compare?" + url.Values{"a": {a}, "b": {b}}.Encode()
}

func compareCountries(s *data.Store) []data.Country {
	var out []data.Country
	for _, c := range s.Countries {
		if len(citiesIn(s, c.Code)) > 0 {
			out = append(out, c)
		}
	}
	return out
}

func compareCityName(s *data.Store, slug string) string {
	if c, ok := s.City(slug); ok {
		return c.Name
	}
	return slug
}

func compareCoffee(c *data.Country) string { return c.Symbol + " " + components.Commas(c.CoffeeLocal) }

func compareCoffeeUSD(s *data.Store, c *data.Country) string {
	r := s.Site.Rates.PerUSD[c.CurrencyCode]
	if r == 0 {
		return "?"
	}
	return strconv.FormatFloat(float64(c.CoffeeLocal)/r, 'f', 2, 64)
}

func seriesLetter(series string) string { return map[string]string{"a": "A", "b": "B"}[series] }
