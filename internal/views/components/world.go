package components

import (
	"fmt"
	"math"
	"strings"

	"workshop/internal/data"
)

// MapMarker is a city drawn on the East Africa map.
type MapMarker struct {
	Name, Href, Role string // Role: "capital" | "commercial"
	X, Y             float64
}

// MapLine is one route drawn between two cities.
type MapLine struct {
	From, To      MapMarker
	Href          string
	Title         string
	Detail        string
	CrossesBorder bool
}

// Markers places every country's capital and commercial city on the map.
func Markers(s *data.Store) []MapMarker {
	p := s.Geo.EastAfrica.Projection
	var out []MapMarker
	for _, c := range s.Countries {
		for _, m := range []struct{ slug, role string }{{c.Capital, "capital"}, {c.CommercialCity, "commercial"}} {
			if city, ok := s.City(m.slug); ok {
				x, y := p.Project(city.Lat, city.Lon)
				out = append(out, MapMarker{Name: city.Name, Href: "/cities/" + city.Slug, Role: m.role, X: x, Y: y})
			}
		}
	}
	return out
}

// CityMarker places one city on the map.
func CityMarker(s *data.Store, slug string) (MapMarker, bool) {
	city, ok := s.City(slug)
	if !ok {
		return MapMarker{}, false
	}
	x, y := s.Geo.EastAfrica.Projection.Project(city.Lat, city.Lon)
	return MapMarker{Name: city.Name, Href: "/cities/" + city.Slug, Role: s.CityRole(slug), X: x, Y: y}, true
}

func f1(v float64) string { return fmt.Sprintf("%.1f", v) }

// starPath is a five-pointed star centred on x,y.
func starPath(x, y, r float64) string {
	var b strings.Builder
	for i := 0; i < 10; i++ {
		rr := r
		if i%2 == 1 {
			rr = r * 0.45
		}
		a := -math.Pi/2 + float64(i)*math.Pi/5
		cmd := "L"
		if i == 0 {
			cmd = "M"
		}
		fmt.Fprintf(&b, "%s%.1f %.1f", cmd, x+rr*math.Cos(a), y+rr*math.Sin(a))
	}
	return b.String() + "Z"
}

// curve bends a route line gently so overlapping routes stay apart.
func curve(a, b MapMarker) string {
	mx, my := (a.X+b.X)/2, (a.Y+b.Y)/2
	dx, dy := b.X-a.X, b.Y-a.Y
	k := 0.18
	return fmt.Sprintf("M%.1f %.1f Q%.1f %.1f %.1f %.1f", a.X, a.Y, mx-dy*k, my+dx*k, b.X, b.Y)
}

// countryHref links a map shape to its country page, or "" for countries the guide doesn't cover.
func countryHref(s *data.Store, code string) string {
	if c, ok := s.Country(code); ok {
		return "/countries/" + c.Slug
	}
	return ""
}

// ModeIcon maps a transport mode to an icon in the shared set.
func ModeIcon(mode string) string {
	switch mode {
	case "matatu", "daladala", "minibus", "shared-taxi":
		return "matatu"
	case "boda", "moto":
		return "boda"
	case "taxi", "ride-hail":
		return "taxi"
	case "train", "sgr", "rail", "light-rail":
		return "train"
	case "plane", "flight":
		return "plane"
	case "ferry", "boat":
		return "ferry"
	case "tuktuk", "bajaj":
		return "tuktuk"
	case "walk":
		return "walk"
	case "car", "road", "drive":
		return "car"
	}
	return "bus"
}

func drivingLabel(side string) string {
	switch side {
	case "left":
		return "On the left"
	case "right":
		return "On the right"
	}
	return "—"
}

func joinOr(v []string, empty string) string {
	if len(v) == 0 {
		return empty
	}
	return strings.Join(v, ", ")
}

func airportCodes(as []data.Airport) string {
	var codes []string
	for _, a := range as {
		codes = append(codes, a.IATA)
	}
	return joinOr(codes, "—")
}

func markerLabel(m MapMarker) string {
	if m.Role == "capital" {
		return m.Name + " (capital)"
	}
	return m.Name + " (commercial city)"
}

// routeEnds lists each city at the end of a route once.
func routeEnds(lines []MapLine) []MapMarker {
	seen := map[string]bool{}
	var out []MapMarker
	for _, l := range lines {
		for _, m := range []MapMarker{l.From, l.To} {
			if !seen[m.Href] {
				seen[m.Href] = true
				out = append(out, m)
			}
		}
	}
	return out
}

// ladderWidth narrows each step of the organisation ladder.
func ladderWidth(i, n int) string {
	if n <= 1 {
		return "--w:100%"
	}
	return fmt.Sprintf("--w:%d%%", 100-i*40/(n-1))
}
