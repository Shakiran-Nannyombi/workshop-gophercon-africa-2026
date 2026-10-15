package pages

import (
	"sort"
	"strings"

	"workshop/internal/data"
)

func cityPtr(s *data.Store, slug string) *data.City {
	if c, ok := s.City(slug); ok {
		return c
	}
	return nil
}

func kenyaOrg(s *data.Store) data.Organisation {
	if k, ok := s.Country("KE"); ok {
		return k.Organisation
	}
	return data.Organisation{}
}

// citiesFor lists a country's cities ("" for all), capitals then commercial hubs first.
func citiesFor(s *data.Store, country string) []data.City {
	var out []data.City
	for _, c := range s.Cities {
		if country == "" || c.Country == country {
			out = append(out, c)
		}
	}
	rank := func(slug string) int {
		switch s.CityRole(slug) {
		case "capital":
			return 0
		case "commercial":
			return 1
		}
		return 2
	}
	sort.SliceStable(out, func(i, j int) bool { return rank(out[i].Slug) < rank(out[j].Slug) })
	return out
}

func crossingsFor(s *data.Store, code string) []data.BorderCrossing {
	var out []data.BorderCrossing
	for _, b := range s.Borders.Crossings {
		for _, c := range b.Countries {
			if c == code {
				out = append(out, b)
				break
			}
		}
	}
	return out
}

// routesFor lists routes with an end in the given country, or touching the given city.
func routesFor(s *data.Store, country, city string) []data.Route {
	var out []data.Route
	for _, r := range s.Routes {
		a, b := cityPtr(s, r.From), cityPtr(s, r.To)
		if a == nil || b == nil {
			continue
		}
		if (city != "" && (r.From == city || r.To == city)) || (country != "" && (a.Country == country || b.Country == country)) {
			out = append(out, r)
		}
	}
	return out
}

func cityName(s *data.Store, slug string) string {
	if c, ok := s.City(slug); ok {
		return c.Name
	}
	return slug
}

func countryOf(s *data.Store, code string) *data.Country {
	if c, ok := s.Country(code); ok {
		return c
	}
	return nil
}

func fullCountries(s *data.Store, depth string) []data.Country {
	var out []data.Country
	for _, c := range s.Countries {
		if c.Depth == depth || (depth == "full" && c.Depth == "") {
			out = append(out, c)
		}
	}
	return out
}

func roleLabel(role string) string {
	switch role {
	case "capital":
		return "Capital"
	case "commercial":
		return "Commercial hub"
	}
	return ""
}

func roleIcon(role string) string {
	if role == "capital" {
		return "star"
	}
	return "city"
}

func joinList(v []string) string { return strings.Join(v, ", ") }

func countryMotif(s *data.Store, c *data.Country) string {
	if cap := cityPtr(s, c.Capital); cap != nil && cap.Motif != "" {
		return cap.Motif
	}
	return "savannah"
}
