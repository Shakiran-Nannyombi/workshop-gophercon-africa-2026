package pages

import (
	"math"
	"net/url"
	"strconv"
	"strings"

	"workshop/internal/data"
	"workshop/internal/views/components"
)

// CountyFilter is the state of the county grid's search and region pills.
type CountyFilter struct {
	Q, Region string
}

func (f CountyFilter) Active() bool { return strings.TrimSpace(f.Q) != "" || f.Region != "" }

// Href is the shareable URL for this filter, landing on the grid.
func (f CountyFilter) Href() string {
	v := url.Values{}
	if f.Q != "" {
		v.Set("q", f.Q)
	}
	if f.Region != "" {
		v.Set("region", f.Region)
	}
	if len(v) == 0 {
		return "/counties#all-counties"
	}
	return "/counties?" + v.Encode() + "#all-counties"
}

// FilterCounties returns the counties matching the region and every search word.
func FilterCounties(s *data.Store, f CountyFilter) []data.County {
	var out []data.County
	for _, c := range s.Counties {
		if f.Region != "" && c.Region != f.Region {
			continue
		}
		fields := append([]string{c.Name, c.Seat, c.Slug, countyRegionName(s, c.Region), strconv.Itoa(c.Code), components.CountyCode(c.Code)}, c.Towns...)
		fields = append(fields, c.KnownFor...)
		fields = append(fields, c.GoodFor...)
		if !data.Matches(f.Q, fields...) {
			continue
		}
		out = append(out, c)
	}
	return out
}

func countySlugs(cs []data.County) string { return strings.Join(slugList(cs), " ") }

func slugList(cs []data.County) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Slug
	}
	return out
}

func countyRegionName(s *data.Store, slug string) string {
	if r, ok := s.Region(slug); ok {
		return r.Name
	}
	return slug
}

// regionOrder is the order of the region pills and the legend.
var regionOrder = []string{"nairobi", "central", "coast", "rift-valley", "western", "nyanza", "eastern", "north-eastern"}

func regionCount(s *data.Store, slug string) int {
	n := 0
	for _, c := range s.Counties {
		if c.Region == slug {
			n++
		}
	}
	return n
}

func regionFilterHref(slug string) string { return CountyFilter{Region: slug}.Href() }

func countyShape(s *data.Store, slug string) data.Shape {
	sh, _ := components.ShapeFor(s.Geo.Kenya, slug)
	return sh
}

var regionMotifs = map[string]string{
	"coast": "coast", "nairobi": "skyline", "central": "highland", "western": "lake", "nyanza": "lake",
	"rift-valley": "savannah", "eastern": "hills", "north-eastern": "savannah",
}

func regionMotif(slug string) string {
	if m, ok := regionMotifs[slug]; ok {
		return m
	}
	return "savannah"
}

// countyDensity is people per km², rounded to a sensible precision.
func countyDensity(c data.County) string {
	if c.AreaKm2 <= 0 {
		return "?"
	}
	d := float64(c.Pop2019) / float64(c.AreaKm2)
	switch {
	case d >= 1000:
		d = math.Round(d/100) * 100
	case d >= 100:
		d = math.Round(d/10) * 10
	default:
		d = math.Round(d)
	}
	return components.Commas(int(d))
}

// countyNeighbours returns the previous and next county by official code, wrapping round.
func countyNeighbours(s *data.Store, c *data.County) (data.County, data.County) {
	n := len(s.Counties)
	for i := range s.Counties {
		if s.Counties[i].Slug == c.Slug {
			return s.Counties[(i+n-1)%n], s.Counties[(i+1)%n]
		}
	}
	return *c, *c
}

func levelTone(i int) string { return "tier--" + strconv.Itoa(min(i, 5)) }

func countyJoin(xs []string) string {
	switch len(xs) {
	case 0:
		return ""
	case 1:
		return xs[0]
	}
	return strings.Join(xs[:len(xs)-1], ", ") + " and " + xs[len(xs)-1]
}

// countyLede is the hero sentence for one county.
func countyLede(s *data.Store, c *data.County) string {
	l := "County " + components.CountyCode(c.Code) + " of 47, in the " + countyRegionName(s, c.Region) + " region. The county seat is " + c.Seat + "."
	if len(c.KnownFor) > 0 {
		l += " Known for " + countyJoin(countyFirst(c.KnownFor, 2)) + "."
	}
	return l
}

func countyFirst(xs []string, n int) []string { return xs[:min(len(xs), n)] }
