package pages

import (
	"sort"
	"strconv"
	"strings"

	"workshop/internal/data"
)

// AreaFilter is the state of the Nairobi area guide's filters.
type AreaFilter struct {
	Q, County, Cost, Suits string
}

func (f AreaFilter) Query() string {
	var parts []string
	add := func(k, v string) {
		if v != "" {
			parts = append(parts, k+"="+v)
		}
	}
	add("q", f.Q)
	add("county", f.County)
	add("cost", f.Cost)
	add("suits", f.Suits)
	return strings.Join(parts, "&")
}

// FilterAreas returns the areas matching every set filter.
func FilterAreas(s *data.Store, f AreaFilter) []data.Area {
	var out []data.Area
	for _, a := range s.Areas {
		if f.County != "" && a.County != f.County {
			continue
		}
		if f.Cost != "" && a.Cost != f.Cost {
			continue
		}
		if f.Suits != "" && !contains(a.Suits, f.Suits) {
			continue
		}
		if !data.Matches(f.Q, append([]string{a.Name, a.Zone, a.Character, a.County, s.CountyName(a.County)}, append(a.AlsoCalled, a.Landmarks...)...)...) {
			continue
		}
		out = append(out, a)
	}
	return out
}

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

// Lookup is the answer to "Which county am I really in?".
type Lookup struct {
	Query  string
	Best   *data.Area
	Others []data.Area
}

// WhichCounty finds the area best matching q: exact names first, then partial matches.
func WhichCounty(s *data.Store, q string) Lookup {
	l := Lookup{Query: strings.TrimSpace(q)}
	fq := data.Fold(q)
	if fq == "" {
		return l
	}
	type scored struct {
		a     data.Area
		score int
	}
	var hits []scored
	for _, a := range s.Areas {
		names := append([]string{a.Name, a.Slug}, a.AlsoCalled...)
		best := 0
		for _, n := range names {
			fn := data.Fold(n)
			switch {
			case fn == fq:
				best = max(best, 3)
			case strings.HasPrefix(fn, fq) || strings.HasPrefix(fq, fn):
				best = max(best, 2)
			case strings.Contains(fn, fq) || strings.Contains(fq, fn):
				best = max(best, 1)
			}
		}
		if best == 0 && len(fq) >= 3 && data.Matches(q, a.Landmarks...) {
			best = 1
		}
		if best > 0 {
			hits = append(hits, scored{a, best})
		}
	}
	sort.SliceStable(hits, func(i, j int) bool { return hits[i].score > hits[j].score })
	for i, h := range hits {
		if i == 0 {
			a := h.a
			l.Best = &a
			continue
		}
		if len(l.Others) < 4 {
			l.Others = append(l.Others, h.a)
		}
	}
	return l
}

// outsideByCounty groups the "Nairobi" places that are really in other counties.
func outsideByCounty(s *data.Store) []struct {
	Slug, Name string
	Areas      []data.Area
} {
	order := []string{"kiambu", "kajiado", "machakos"}
	var out []struct {
		Slug, Name string
		Areas      []data.Area
	}
	for _, c := range order {
		g := struct {
			Slug, Name string
			Areas      []data.Area
		}{Slug: c, Name: s.CountyName(c)}
		for _, a := range s.Areas {
			if a.County == c {
				g.Areas = append(g.Areas, a)
			}
		}
		out = append(out, g)
	}
	return out
}

func limitedCount(as []data.Area) int {
	n := 0
	for _, a := range as {
		if a.Confidence == "limited" {
			n++
		}
	}
	return n
}

var areaCounties = []struct{ Value, Label string }{{"", "All counties"}, {"nairobi", "Nairobi"}, {"kiambu", "Kiambu"}, {"kajiado", "Kajiado"}, {"machakos", "Machakos"}}
var areaCosts = []struct{ Value, Label string }{{"", "Any cost"}, {"budget", "Budget"}, {"mid", "Mid-range"}, {"upmarket", "Upmarket"}, {"mixed", "Mixed"}}
var areaSuits = []struct{ Value, Label string }{{"", "Anyone"}, {"visitors", "Visitors"}, {"expats", "Expats"}, {"families", "Families"}, {"students", "Students"}}

var whichExamples = []string{"Ruaka", "Rongai", "Syokimau", "Kitengela", "Karen", "Ruiru"}

var countyTitles = map[string]string{"nairobi": "Nairobi City", "kiambu": "Kiambu", "kajiado": "Kajiado", "machakos": "Machakos"}

func countyTitle(slug string) string {
	if t, ok := countyTitles[slug]; ok {
		return t
	}
	return slug
}

func fmtNum(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }

var foods = []struct{ Art, Name, Note string }{
	{"ugali", "Ugali", "Firm maize-meal staple, eaten by hand"},
	{"nyama", "Nyama choma", "Roast meat, ordered by the kilo"},
	{"chapati", "Chapati", "Soft, flaky flatbread"},
	{"mandazi", "Mandazi", "Lightly sweet fried dough"},
	{"chai", "Chai", "Milky spiced tea, all day"},
	{"maize", "Mahindi choma", "Roast maize from the roadside"},
	{"swahili", "Pilau", "Spiced Swahili rice"},
}
