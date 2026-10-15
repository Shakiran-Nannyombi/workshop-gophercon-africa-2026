package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

// Store is all site content, loaded once at start-up and read-only afterwards.
type Store struct {
	Site      Site
	Money     Money
	Countries []Country
	Regions   []Region
	Counties  []County
	Areas     []Area
	Cities    []City
	Glossary  []Term
	Phrases   []Phrase
	FAQ       []FAQ
	Sources   []Source
	Geo       Geo

	Currencies   []Currency
	Routes       []Route
	Borders      Borders
	Connectivity Connectivity

	countyBySlug   map[string]*County
	currencyByCode map[string]*Currency
	countryBySlug  map[string]*Country
	borderBySlug   map[string]*BorderCrossing
	cityBySlug     map[string]*City
	countryByCode  map[string]*Country
	regionBySlug   map[string]*Region
	areaBySlug     map[string]*Area
}

// Load reads every content file from fsys and checks it hangs together.
func Load(fsys fs.FS) (*Store, error) {
	s := &Store{}
	files := []struct {
		name     string
		into     any
		optional bool
	}{
		{"site.json", &s.Site, false},
		{"rates.json", &s.Site.Rates, true},
		{"currencies.json", &s.Currencies, true},
		{"routes.json", &s.Routes, true},
		{"borders.json", &s.Borders, true},
		{"connectivity.json", &s.Connectivity, true},
		{"money.json", &s.Money, true},
		{"countries.json", &s.Countries, true},
		{"regions.json", &s.Regions, true},
		{"counties.json", &s.Counties, true},
		{"nairobi_areas.json", &s.Areas, true},
		{"cities.json", &s.Cities, true},
		{"glossary.json", &s.Glossary, true},
		{"phrases.json", &s.Phrases, true},
		{"faq.json", &s.FAQ, true},
		{"sources.json", &s.Sources, true},
		{"geo/kenya_counties.json", &s.Geo.Kenya, true},
		{"geo/nairobi_region.json", &s.Geo.Nairobi, true},
		{"geo/east_africa.json", &s.Geo.EastAfrica, true},
	}
	for _, f := range files {
		b, err := fs.ReadFile(fsys, f.name)
		if errors.Is(err, fs.ErrNotExist) && f.optional {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("data: reading %s: %w", f.name, err)
		}
		if err := json.Unmarshal(b, f.into); err != nil {
			return nil, fmt.Errorf("data: parsing %s: %w", f.name, err)
		}
	}
	if err := s.index(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) index() error {
	s.countyBySlug = map[string]*County{}
	s.cityBySlug = map[string]*City{}
	s.countryByCode = map[string]*Country{}
	s.regionBySlug = map[string]*Region{}
	s.areaBySlug = map[string]*Area{}

	for i := range s.Regions {
		s.regionBySlug[s.Regions[i].Slug] = &s.Regions[i]
	}
	for i := range s.Countries {
		s.countryByCode[s.Countries[i].Code] = &s.Countries[i]
	}
	sort.Slice(s.Counties, func(i, j int) bool { return s.Counties[i].Code < s.Counties[j].Code })
	for i := range s.Counties {
		c := &s.Counties[i]
		if _, dup := s.countyBySlug[c.Slug]; dup {
			return fmt.Errorf("data: duplicate county slug %q", c.Slug)
		}
		if len(s.Regions) > 0 && s.regionBySlug[c.Region] == nil {
			return fmt.Errorf("data: county %q has unknown region %q", c.Slug, c.Region)
		}
		s.countyBySlug[c.Slug] = c
	}
	if n := len(s.Counties); n != 0 && n != 47 {
		return fmt.Errorf("data: expected 47 counties, found %d", n)
	}
	for i := range s.Cities {
		c := &s.Cities[i]
		if _, dup := s.cityBySlug[c.Slug]; dup {
			return fmt.Errorf("data: duplicate city slug %q", c.Slug)
		}
		if len(s.Countries) > 0 && s.countryByCode[c.Country] == nil {
			return fmt.Errorf("data: city %q has unknown country %q", c.Slug, c.Country)
		}
		s.cityBySlug[c.Slug] = c
	}
	for i := range s.Areas {
		a := &s.Areas[i]
		if _, dup := s.areaBySlug[a.Slug]; dup {
			return fmt.Errorf("data: duplicate area slug %q", a.Slug)
		}
		if len(s.Counties) > 0 && s.countyBySlug[a.County] == nil {
			return fmt.Errorf("data: area %q has unknown county %q", a.Slug, a.County)
		}
		s.areaBySlug[a.Slug] = a
	}
	return s.indexWorld()
}

// indexWorld indexes and cross-checks countries, currencies, routes and borders.
func (s *Store) indexWorld() error {
	s.currencyByCode = map[string]*Currency{}
	s.countryBySlug = map[string]*Country{}
	s.borderBySlug = map[string]*BorderCrossing{}

	for i := range s.Countries {
		c := &s.Countries[i]
		if c.Slug == "" {
			c.Slug = strings.ToLower(strings.ReplaceAll(c.Name, " ", "-"))
		}
		s.countryBySlug[c.Slug] = c
		for _, city := range []string{c.Capital, c.CommercialCity} {
			if city != "" && s.cityBySlug[city] == nil {
				return fmt.Errorf("data: country %q names unknown city %q", c.Code, city)
			}
		}
	}
	for i := range s.Currencies {
		c := &s.Currencies[i]
		if s.countryByCode[c.Country] == nil {
			return fmt.Errorf("data: currency %q has unknown country %q", c.Code, c.Country)
		}
		s.currencyByCode[c.Code] = c
	}
	for i := range s.Borders.Crossings {
		b := &s.Borders.Crossings[i]
		for _, code := range b.Countries {
			if s.countryByCode[code] == nil {
				return fmt.Errorf("data: border %q has unknown country %q", b.Slug, code)
			}
		}
		s.borderBySlug[b.Slug] = b
	}
	for _, r := range s.Routes {
		if s.cityBySlug[r.From] == nil || s.cityBySlug[r.To] == nil {
			return fmt.Errorf("data: route %s-%s names an unknown city", r.From, r.To)
		}
		if r.Border != "" && s.borderBySlug[r.Border] == nil {
			return fmt.Errorf("data: route %s-%s has unknown border %q", r.From, r.To, r.Border)
		}
	}
	return nil
}

func (s *Store) Currency(code string) (*Currency, bool) {
	c, ok := s.currencyByCode[code]
	return c, ok
}
func (s *Store) CountryBySlug(slug string) (*Country, bool) {
	c, ok := s.countryBySlug[slug]
	return c, ok
}
func (s *Store) Border(slug string) (*BorderCrossing, bool) {
	b, ok := s.borderBySlug[slug]
	return b, ok
}

// RouteBetween finds the route linking two cities, in either direction.
func (s *Store) RouteBetween(a, b string) (*Route, bool) {
	for i := range s.Routes {
		r := &s.Routes[i]
		if (r.From == a && r.To == b) || (r.From == b && r.To == a) {
			return r, true
		}
	}
	return nil, false
}

// IsCapital reports whether a city is its country's official capital.
func (s *Store) IsCapital(city string) bool {
	for _, c := range s.Countries {
		if c.Capital == city {
			return true
		}
	}
	return false
}

// CityRole is "capital", "commercial" or "" for a city slug.
func (s *Store) CityRole(city string) string {
	for _, c := range s.Countries {
		switch city {
		case c.Capital:
			return "capital"
		case c.CommercialCity:
			return "commercial"
		}
	}
	return ""
}

// SourcesOf returns the sources of one kind ("" for content sources).
func (s *Store) SourcesOf(kind string) []Source {
	var out []Source
	for _, src := range s.Sources {
		if src.Kind == kind {
			out = append(out, src)
		}
	}
	return out
}

func (s *Store) County(slug string) (*County, bool)   { c, ok := s.countyBySlug[slug]; return c, ok }
func (s *Store) City(slug string) (*City, bool)       { c, ok := s.cityBySlug[slug]; return c, ok }
func (s *Store) Country(code string) (*Country, bool) { c, ok := s.countryByCode[code]; return c, ok }
func (s *Store) Region(slug string) (*Region, bool)   { r, ok := s.regionBySlug[slug]; return r, ok }
func (s *Store) Area(slug string) (*Area, bool)       { a, ok := s.areaBySlug[slug]; return a, ok }

// CountyName returns a county's display name, or the slug if unknown.
func (s *Store) CountyName(slug string) string {
	if c, ok := s.countyBySlug[slug]; ok {
		return c.Name
	}
	return slug
}

// Fold lowercases and strips punctuation so "Murang'a" matches "muranga".
func Fold(v string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(v) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == ' ':
			b.WriteRune(r)
		case r == '-' || r == '/':
			b.WriteRune(' ')
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

// Matches reports whether every word of q appears in one of the fields.
func Matches(q string, fields ...string) bool {
	q = Fold(q)
	if q == "" {
		return true
	}
	hay := Fold(strings.Join(fields, " "))
	for _, w := range strings.Fields(q) {
		if !strings.Contains(hay, w) {
			return false
		}
	}
	return true
}
