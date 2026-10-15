package pages

import (
	"net/url"
	"sort"
	"strings"

	"workshop/internal/data"
)

// GlossaryFilter is the state of the glossary search box and language pills.
type GlossaryFilter struct {
	Q, Lang string // Lang is a language slug, e.g. "sheng"
}

func (f GlossaryFilter) active() bool { return strings.TrimSpace(f.Q) != "" || f.Lang != "" }

// LangSlug turns "English (Kenya)" into "english" for URLs and CSS classes.
func LangSlug(lang string) string {
	if w := strings.Fields(data.Fold(lang)); len(w) > 0 {
		return w[0]
	}
	return ""
}

func langLabel(lang string) string {
	if lang == "English (Kenya)" {
		return "Kenyan English"
	}
	return lang
}

// langAttr is the HTML lang for a term; Sheng has no code of its own.
func langAttr(lang string) string {
	switch LangSlug(lang) {
	case "swahili":
		return "sw"
	case "english":
		return "en-KE"
	}
	return ""
}

// GlossaryLangs lists the languages present in the glossary, in a fixed friendly order.
func GlossaryLangs(s *data.Store) []string {
	rank := map[string]int{"swahili": 0, "sheng": 1, "english": 2}
	seen := map[string]bool{}
	var out []string
	for _, t := range s.Glossary {
		if t.Language != "" && !seen[t.Language] {
			seen[t.Language] = true
			out = append(out, t.Language)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		ri, oki := rank[LangSlug(out[i])]
		rj, okj := rank[LangSlug(out[j])]
		if oki != okj {
			return oki
		}
		return ri < rj
	})
	return out
}

func langOptions(s *data.Store) []struct{ Value, Label string } {
	opts := []struct{ Value, Label string }{{"", "All"}}
	for _, l := range GlossaryLangs(s) {
		opts = append(opts, struct{ Value, Label string }{LangSlug(l), langLabel(l)})
	}
	return opts
}

type letterGroup struct {
	Letter string
	Terms  []data.Term
}

// FilterTerms returns the matching terms grouped A–Z ("#" for anything else).
func FilterTerms(s *data.Store, f GlossaryFilter) ([]letterGroup, int) {
	var terms []data.Term
	for _, t := range s.Glossary {
		if f.Lang != "" && LangSlug(t.Language) != f.Lang {
			continue
		}
		if !data.Matches(f.Q, t.Term, t.Meaning, t.Example, t.Language, langLabel(t.Language)) {
			continue
		}
		terms = append(terms, t)
	}
	sort.SliceStable(terms, func(i, j int) bool { return data.Fold(terms[i].Term) < data.Fold(terms[j].Term) })
	var groups []letterGroup
	for _, t := range terms {
		l := initial(t.Term)
		if n := len(groups); n == 0 || groups[n-1].Letter != l {
			groups = append(groups, letterGroup{Letter: l})
		}
		groups[len(groups)-1].Terms = append(groups[len(groups)-1].Terms, t)
	}
	return groups, len(terms)
}

func initial(term string) string {
	f := data.Fold(term)
	if f != "" && f[0] >= 'a' && f[0] <= 'z' {
		return strings.ToUpper(f[:1])
	}
	return "#"
}

var alphabet = strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")

func hasLetter(groups []letterGroup, l string) bool {
	for _, g := range groups {
		if g.Letter == l {
			return true
		}
	}
	return false
}

func glossaryQuery(f GlossaryFilter) string {
	var parts []string
	if f.Q != "" {
		parts = append(parts, "q="+url.QueryEscape(f.Q))
	}
	if f.Lang != "" {
		parts = append(parts, "lang="+f.Lang)
	}
	if len(parts) == 0 {
		return ""
	}
	return "?" + strings.Join(parts, "&")
}

type phraseGroup struct {
	Slug, English, Swahili string
	Proverb, ProverbEn     string
	Phrases                []data.Phrase
}

// phraseCats are the phrasebook chapters, each with the Swahili name and a kanga-style proverb.
var phraseCats = []phraseGroup{
	{Slug: "greetings", English: "Greetings", Swahili: "Salamu", Proverb: "Mgeni njoo, mwenyeji apone", ProverbEn: "Let the guest come, so the host may prosper"},
	{Slug: "polite", English: "Being polite", Swahili: "Adabu", Proverb: "Maneno matamu humtoa nyoka pangoni", ProverbEn: "Sweet words draw the snake out of its hole"},
	{Slug: "money", English: "Money and shopping", Swahili: "Pesa", Proverb: "Haba na haba hujaza kibaba", ProverbEn: "Little by little fills the measure"},
	{Slug: "getting-around", English: "Getting around", Swahili: "Usafiri", Proverb: "Pole pole ndio mwendo", ProverbEn: "Slowly, slowly is the way to go"},
	{Slug: "food", English: "Food and drink", Swahili: "Chakula", Proverb: "Chovya chovya humaliza buyu la asali", ProverbEn: "Dip by dip empties the gourd of honey"},
	{Slug: "emergencies", English: "Emergencies", Swahili: "Dharura", Proverb: "Akufaaye kwa dhiki ndiye rafiki", ProverbEn: "A friend in need is a friend indeed"},
}

// PhraseGroups returns the phrasebook chapters in order; unknown categories go last.
func PhraseGroups(s *data.Store) []phraseGroup {
	out := make([]phraseGroup, 0, len(phraseCats))
	idx := map[string]int{}
	for _, c := range phraseCats {
		idx[c.Slug] = len(out)
		out = append(out, c)
	}
	for _, p := range s.Phrases {
		i, ok := idx[p.Category]
		if !ok {
			idx[p.Category] = len(out)
			out = append(out, phraseGroup{Slug: p.Category, English: strings.ReplaceAll(p.Category, "-", " ")})
			i = idx[p.Category]
		}
		out[i].Phrases = append(out[i].Phrases, p)
	}
	kept := out[:0]
	for _, g := range out {
		if len(g.Phrases) > 0 {
			kept = append(kept, g)
		}
	}
	return kept
}

func twoDigit(n int) string {
	if n < 10 {
		return "0" + itoa(n)
	}
	return itoa(n)
}

var glossaryTOC = []struct{ ID, Label string }{
	{"words", "Words A–Z"},
	{"phrasebook", "Swahili phrasebook"},
	{"faq", "Questions"},
}

func termID(term string) string { return "term-" + strings.ReplaceAll(data.Fold(term), " ", "-") }
