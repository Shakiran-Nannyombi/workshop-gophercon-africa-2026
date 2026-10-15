package pages

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"workshop/internal/data"
	"workshop/internal/views/components"
)

type cityCountryGroup struct {
	Country data.Country
	Cities  []data.City
}

func citiesByCountry(s *data.Store) []cityCountryGroup {
	var out []cityCountryGroup
	for _, c := range s.Countries {
		if cs := citiesIn(s, c.Code); len(cs) > 0 {
			out = append(out, cityCountryGroup{c, cs})
		}
	}
	return out
}

func citiesCountryNames(s *data.Store) map[string]string {
	m := make(map[string]string, len(s.Countries))
	for _, c := range s.Countries {
		m[c.Code] = c.Name
	}
	return m
}

func citiesMap(s *data.Store, current string, compact bool) components.EastAfricaMapProps {
	return components.EastAfricaMapProps{Geo: s.Geo.EastAfrica, Cities: s.Cities, Countries: citiesCountryNames(s), Current: current, Compact: compact}
}

// cityPrevNext returns the previous and next city in data order, wrapping round.
func cityPrevNext(s *data.Store, slug string) (prev, next data.City) {
	n := len(s.Cities)
	for i, c := range s.Cities {
		if c.Slug == slug {
			return s.Cities[(i+n-1)%n], s.Cities[(i+1)%n]
		}
	}
	return
}

// cityCompareWith picks the nearest city in the same country plus the two nearest abroad.
func cityCompareWith(s *data.Store, c data.City) []data.City {
	others := make([]data.City, 0, len(s.Cities))
	for _, o := range s.Cities {
		if o.Slug != c.Slug {
			others = append(others, o)
		}
	}
	dist := func(o data.City) float64 { return math.Hypot(o.Lat-c.Lat, (o.Lon-c.Lon)*math.Cos(c.Lat*math.Pi/180)) }
	sort.SliceStable(others, func(i, j int) bool { return dist(others[i]) < dist(others[j]) })
	var out []data.City
	home, abroad := 0, 0
	for _, o := range others {
		if o.Country == c.Country && home < 1 {
			out, home = append(out, o), home+1
		} else if o.Country != c.Country && abroad < 2 {
			out, abroad = append(out, o), abroad+1
		}
	}
	for _, o := range others {
		if len(out) >= 3 {
			break
		}
		if !cityListHas(out, o.Slug) {
			out = append(out, o)
		}
	}
	return out
}

func cityListHas(cs []data.City, slug string) bool {
	for _, c := range cs {
		if c.Slug == slug {
			return true
		}
	}
	return false
}

func cityScoresEditorial(c data.City) bool {
	return strings.Contains(strings.ToLower(string(c.Verify)), "editorial estimate")
}

func cityCompareURL(a, b string) string { return "/compare?a=" + a + "&b=" + b }

func cityPageURL(slug string) string { return "/cities/" + slug }

func citiesPlural(n int) string {
	if n == 1 {
		return "city"
	}
	return "cities"
}

var cityTOC = []struct{ ID, Label string }{
	{"glance", "At a glance"},
	{"areas", "Areas"},
	{"money", "Money"},
	{"practical", "Transport, food, safety"},
	{"suits", "Who it suits"},
	{"locator", "Map & compare"},
}

func cityEyebrow(c data.City, country *data.Country) string {
	if country == nil {
		return "City guide"
	}
	return "City guide · " + country.Name
}

func hoodNum(i int) string { return fmt.Sprintf("%02d", i+1) }

func suitTone(i int) string {
	return []string{"chip--gold", "chip--green", "chip--red", "chip--surface"}[i%4]
}

func joinWords(xs []string) string { return strings.Join(xs, ", ") }
