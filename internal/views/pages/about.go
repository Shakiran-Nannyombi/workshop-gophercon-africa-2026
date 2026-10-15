package pages

import (
	"net/url"
	"sort"
	"strings"

	"workshop/internal/data"

	"github.com/a-h/templ"
)

// sourceLink returns u as a safe link target, but only for absolute http(s) URLs.
func sourceLink(u string) (templ.SafeURL, bool) {
	p, err := url.Parse(strings.TrimSpace(u))
	if err != nil || p.Host == "" || (p.Scheme != "http" && p.Scheme != "https") {
		return "", false
	}
	return templ.SafeURL(p.String()), true
}

func displayURL(u string) string {
	u = strings.TrimPrefix(strings.TrimPrefix(strings.TrimSpace(u), "https://"), "http://")
	return strings.TrimSuffix(u, "/")
}

type rateRow struct {
	Code string
	Rate float64
}

// rateRows lists the configured per-USD rates, KES first, then alphabetically.
func rateRows(s *data.Store) []rateRow {
	var out []rateRow
	for c, r := range s.Site.Rates.PerUSD {
		if c != "USD" {
			out = append(out, rateRow{c, r})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if (out[i].Code == "KES") != (out[j].Code == "KES") {
			return out[i].Code == "KES"
		}
		return out[i].Code < out[j].Code
	})
	return out
}

type mapCredit struct {
	Title, Source, License string
}

func mapCredits(s *data.Store) []mapCredit {
	var out []mapCredit
	for _, c := range []mapCredit{
		{"Kenya county map", s.Geo.Kenya.Source, s.Geo.Kenya.License},
		{"East Africa map", s.Geo.EastAfrica.Source, s.Geo.EastAfrica.License},
	} {
		if c.Source != "" || c.License != "" {
			out = append(out, c)
		}
	}
	return out
}

var dataFiles = []struct{ Name, What string }{
	{"site.json", "Site name, “last reviewed” date and exchange rates"},
	{"money.json", "Notes, coins, M-Pesa, prices and scams"},
	{"counties.json", "The 47 counties"},
	{"regions.json", "The eight former provinces"},
	{"nairobi.json", "Nairobi deep dive: transport, SIM, food, weather"},
	{"nairobi_areas.json", "Nairobi area guide"},
	{"cities.json", "East African city pages and compare scores"},
	{"countries.json", "Countries, currencies and mobile money"},
	{"glossary.json", "Glossary words"},
	{"phrases.json", "Swahili phrasebook"},
	{"faq.json", "Frequently asked questions"},
	{"sources.json", "The sources listed on this page"},
}

var aboutTOC = []struct{ ID, Label string }{
	{"what", "What this is"},
	{"edit", "Edit the content"},
	{"figures", "Approximate figures"},
	{"rates", "Exchange rates"},
	{"sources", "Sources"},
	{"maps", "Map data"},
	{"fair-use", "Fair use"},
}
