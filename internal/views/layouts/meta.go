package layouts

import "workshop/internal/data"

// Meta is what every page tells the layout: SEO fields and which nav item is current.
type Meta struct {
	Site        data.Site
	Path        string
	Section     string
	Title       string
	Description string
	NoIndex     bool
}

// FullTitle is the <title>, with the site name appended.
func (m Meta) FullTitle() string {
	if m.Title == "" {
		return m.Site.Name + " — " + m.Site.Tagline
	}
	return m.Title + " · " + m.Site.Name
}

func (m Meta) Canonical() string { return m.Site.BaseURL + m.Path }

type NavItem struct {
	Section, Href, Label, Short string
}

var Nav = []NavItem{
	{"money", "/money", "Money", "Money"},
	{"counties", "/counties", "Counties", "Counties"},
	{"nairobi", "/nairobi", "Nairobi", "Nairobi"},
	{"cities", "/cities", "Cities", "Cities"},
	{"compare", "/compare", "Compare", "Compare"},
	{"glossary", "/glossary", "Words", "Words"},
}

// Proverbs are printed like the "jina" on a kanga, one per page section.
var Proverbs = []struct{ Swahili, English string }{
	{"Haba na haba hujaza kibaba", "Little by little fills the measure"},
	{"Haraka haraka haina baraka", "Hurry, hurry has no blessing"},
	{"Asiyefunzwa na mamaye hufunzwa na ulimwengu", "Who isn't taught by their mother is taught by the world"},
	{"Mgeni siku mbili, siku ya tatu mpe jembe", "A guest for two days; on the third, hand them a hoe"},
	{"Kidole kimoja hakivunji chawa", "One finger cannot crush a louse"},
	{"Usipoziba ufa utajenga ukuta", "If you don't fill the crack, you'll build a wall"},
}

func (m Meta) Proverb() (string, string) {
	p := Proverbs[len(m.Path)%len(Proverbs)]
	return p.Swahili, p.English
}
