package components

import (
	"fmt"
	"strings"

	"workshop/internal/data"
)

type EastAfricaMapProps struct {
	Geo       data.ShapeSet
	Cities    []data.City
	Countries map[string]string // code -> name
	Current   string            // slug of the city to highlight; "" for none
	Compact   bool
}

type mapCity struct {
	Slug, Name, Country string
	x, y                float64
	X, Y, LX, LY        string
	Anchor              string
	Capital, Current    bool
}

type mapLabel struct {
	dx, dy float64
	anchor string
}

// cityLabels keeps neighbouring labels apart (Kampala/Entebbe/Jinja, Nairobi/Nakuru, Dar/Zanzibar).
var cityLabels = map[string]mapLabel{
	"nairobi":       {15, 6, "start"},
	"nakuru":        {11, 2, "start"},
	"eldoret":       {0, -13, "middle"},
	"kisumu":        {0, 24, "middle"},
	"mombasa":       {12, 6, "start"},
	"kampala":       {-15, 0, "end"},
	"entebbe":       {-11, 18, "end"},
	"jinja":         {2, -14, "middle"},
	"dar-es-salaam": {12, 10, "start"},
	"zanzibar-city": {12, -4, "start"},
	"arusha":        {-12, 6, "end"},
	"kigali":        {15, 6, "start"},
	"addis-ababa":   {15, 6, "start"},
}

var capitals = map[string]bool{"nairobi": true, "kampala": true, "kigali": true, "addis-ababa": true}

var featuredCountries = map[string]bool{"KE": true, "UG": true, "TZ": true, "RW": true, "ET": true}

// countryLabels are hand-placed so names sit clear of the city dots.
var countryLabels = map[string][2]float64{
	"KE": {425, 555}, "UG": {118, 520}, "TZ": {255, 820}, "RW": {48, 668}, "ET": {470, 345},
	"SO": {560, 540}, "SS": {135, 390}, "SD": {130, 200}, "CD": {40, 470},
}

const eaViewBox = "0 150 620 800"

func mapCities(p EastAfricaMapProps) []mapCity {
	out := make([]mapCity, 0, len(p.Cities))
	for _, c := range p.Cities {
		x, y := p.Geo.Projection.Project(c.Lat, c.Lon)
		l, ok := cityLabels[c.Slug]
		if !ok {
			l = mapLabel{0, -12, "middle"}
		}
		if p.Compact {
			l = mapLabel{20, 9, "start"}
			if x > 360 {
				l = mapLabel{-20, 9, "end"}
			}
		}
		out = append(out, mapCity{
			Slug: c.Slug, Name: c.Name, Country: c.Country, x: x, y: y,
			X: f1(x), Y: f1(y), LX: f1(x + l.dx), LY: f1(y + l.dy), Anchor: l.anchor,
			Capital: capitals[c.Slug], Current: c.Slug == p.Current,
		})
	}
	return out
}

func eaCountryClass(code string) string {
	if featuredCountries[code] {
		return "ea-country ea-country--featured ea-country--" + strings.ToLower(code)
	}
	return "ea-country"
}

func eaCountryLabel(code string) (string, string, bool) {
	p, ok := countryLabels[code]
	return f1(p[0]), f1(p[1]), ok
}

func eaCityClass(c mapCity) string {
	cls := "ea-city"
	if c.Capital {
		cls += " ea-city--capital"
	}
	if c.Current {
		cls += " ea-city--current"
	}
	return cls
}

// starPath draws a five-point star of radius r centred on the city.
func starPath(c mapCity, r float64) string {
	pts := [10][2]float64{{0, -1}, {.2245, -.309}, {.951, -.309}, {.363, .118}, {.588, .809}, {0, .382}, {-.588, .809}, {-.363, .118}, {-.951, -.309}, {-.2245, -.309}}
	var b strings.Builder
	for i, p := range pts {
		cmd := "L"
		if i == 0 {
			cmd = "M"
		}
		fmt.Fprintf(&b, "%s%.1f %.1f ", cmd, c.x+p[0]*r, c.y+p[1]*r)
	}
	b.WriteString("Z")
	return b.String()
}

func cityAria(c mapCity, country string) string {
	s := c.Name + ", " + country
	if c.Capital {
		s += ", national capital"
	}
	if c.Current {
		s += " (this page)"
	}
	return s
}

func dodomaXY(p data.Projection) (string, string) {
	x, y := p.Project(-6.163, 35.7516)
	return f1(x), f1(y)
}

func showLabel(p EastAfricaMapProps, c mapCity) bool { return !p.Compact || c.Current }

func eaTitleID(p EastAfricaMapProps) string {
	if p.Compact {
		return "ea-map-title-" + p.Current
	}
	return "ea-map-title"
}

func cityCountry(p EastAfricaMapProps) string {
	for _, c := range p.Cities {
		if c.Slug == p.Current {
			return c.Country
		}
	}
	return ""
}

func cityName(p EastAfricaMapProps) string {
	for _, c := range p.Cities {
		if c.Slug == p.Current {
			return c.Name
		}
	}
	return ""
}
