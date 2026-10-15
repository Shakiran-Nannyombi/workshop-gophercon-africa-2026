package pages

import (
	"sort"

	"workshop/internal/data"
	"workshop/internal/views/components"
)

// ComparePair is the two cities on the compare page.
type ComparePair struct {
	A, B *data.City
}

// capitalPairs are quick capital-vs-capital comparisons.
var capitalPairs = [][2]string{
	{"nairobi", "kampala"}, {"nairobi", "kigali"}, {"kampala", "kigali"},
	{"nairobi", "addis-ababa"}, {"dodoma", "dar-es-salaam"}, {"gitega", "bujumbura"},
}

func countryName(s *data.Store, c *data.City) string {
	if co := countryOf(s, c.Country); co != nil {
		return co.Name
	}
	return c.Country
}

func dots(n int) []bool {
	out := make([]bool, 5)
	for i := range out {
		out[i] = i < n
	}
	return out
}

// currencyScale is how far KSh 1,000 goes in the city's currency.
func currencyScale(s *data.Store, c *data.City) string {
	co := countryOf(s, c.Country)
	if co == nil {
		return "—"
	}
	r, kes := s.Site.Rates.PerUSD[co.CurrencyCode], s.Site.Rates.PerUSD["KES"]
	if r == 0 || kes == 0 {
		return co.CurrencyCode
	}
	if co.CurrencyCode == "KES" {
		return "KSh 1,000 ≈ US$" + components.FormatMoney(1000/kes, "USD")
	}
	return "KSh 1,000 ≈ " + components.Sym(co.CurrencyCode) + " " + roughNum(1000*r/kes) + " (" + co.CurrencyCode + ")"
}

// gettingAroundCost lists the priced ways of moving about a capital, or the city's first transport note.
func gettingAroundCost(s *data.Store, c *data.City) []string {
	var out []string
	for _, m := range modesFor(s, c) {
		if m.Cost != "" {
			out = append(out, m.Name+": "+m.Cost)
		}
	}
	if len(out) == 0 && len(c.Transport) > 0 {
		out = append(out, c.Transport[0])
	}
	return out
}

// termsByLanguage groups glossary terms, alphabetically within each group.
func termsByLanguage(ts []data.Term) []termGroup {
	idx := map[string]int{}
	var out []termGroup
	for _, t := range ts {
		i, ok := idx[t.Language]
		if !ok {
			i = len(out)
			idx[t.Language] = i
			out = append(out, termGroup{Language: t.Language})
		}
		out[i].Terms = append(out[i].Terms, t)
	}
	for i := range out {
		sort.SliceStable(out[i].Terms, func(a, b int) bool { return data.Fold(out[i].Terms[a].Term) < data.Fold(out[i].Terms[b].Term) })
	}
	return out
}

type termGroup struct {
	Language string
	Terms    []data.Term
}

func slugify(v string) string {
	f := data.Fold(v)
	b := []rune(f)
	for i, r := range b {
		if r == ' ' {
			b[i] = '-'
		}
	}
	return string(b)
}
