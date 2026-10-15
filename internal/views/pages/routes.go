package pages

import (
	"strings"

	"workshop/internal/data"
	"workshop/internal/views/components"
)

// Plan is the route planner's state: the chosen cities and the route found, if any.
type Plan struct {
	From, To string
	Route    *data.Route
}

// costRange prints "~KSh 1,500–3,000", or "" when the option has no fare.
func costRange(o data.RouteOption) string {
	if o.Currency == "" || o.CostMax <= 0 {
		return ""
	}
	sym := components.Sym(o.Currency)
	if o.CostMin <= 0 || o.CostMin == o.CostMax {
		return "~" + sym + " " + components.Commas(o.CostMax)
	}
	return "~" + sym + " " + components.Commas(o.CostMin) + "–" + components.Commas(o.CostMax)
}

func modeName(mode string) string {
	switch mode {
	case "sgr":
		return "SGR train"
	case "shared-taxi":
		return "Shared taxi"
	case "road":
		return "By road"
	}
	if mode == "" {
		return ""
	}
	return strings.ToUpper(mode[:1]) + mode[1:]
}

// routeDetail is a one-line summary of a route's options, for the map.
func routeDetail(r data.Route) string {
	var parts []string
	for _, o := range r.Options {
		p := modeName(o.Mode)
		if o.Duration != "" {
			p += " " + o.Duration
		}
		parts = append(parts, p)
	}
	return strings.Join(parts, " · ")
}

// routeLines turns every route into a line on the routes map.
func routeLines(s *data.Store) []components.MapLine {
	var out []components.MapLine
	for _, r := range s.Routes {
		a, okA := components.CityMarker(s, r.From)
		b, okB := components.CityMarker(s, r.To)
		if !okA || !okB {
			continue
		}
		out = append(out, components.MapLine{
			From: a, To: b,
			Href:          "/routes?from=" + r.From + "&to=" + r.To + "#planner",
			Title:         a.Name + " → " + b.Name,
			Detail:        routeDetail(r),
			CrossesBorder: r.Border != "",
		})
	}
	return out
}

type cityGroup struct {
	Country string
	Cities  []data.City
}

// cityGroups lists the cities by country for the planner's dropdowns.
func cityGroups(s *data.Store) []cityGroup {
	var out []cityGroup
	for _, c := range s.Countries {
		if cs := citiesFor(s, c.Code); len(cs) > 0 {
			out = append(out, cityGroup{c.Name, cs})
		}
	}
	return out
}

// routesFrom lists routes starting or ending at a city, for "try instead" links.
func routesFrom(s *data.Store, city string) []data.Route { return routesFor(s, "", city) }

func borderName(s *data.Store, slug string) string {
	if b, ok := s.Border(slug); ok {
		return b.Name
	}
	return slug
}

func netFor(s *data.Store, code string) *data.CountryNet {
	for i := range s.Connectivity.Countries {
		if s.Connectivity.Countries[i].Country == code {
			return &s.Connectivity.Countries[i]
		}
	}
	return nil
}

func routesVia(s *data.Store, border string) []data.Route {
	var out []data.Route
	for _, r := range s.Routes {
		if r.Border == border {
			out = append(out, r)
		}
	}
	return out
}
