package handlers

import (
	"net/http"

	"workshop/internal/views/pages"

	"github.com/a-h/templ"
)

func (s *site) glossary(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := pages.GlossaryFilter{Q: q.Get("q"), Lang: q.Get("lang")}
	if f.Lang != "" {
		known := false
		for _, l := range pages.GlossaryLangs(s.store) {
			known = known || pages.LangSlug(l) == f.Lang
		}
		if !known {
			f.Lang = ""
		}
	}

	var partial templ.Component
	switch r.Header.Get("HX-Target") {
	case "glossary-results":
		partial = pages.GlossaryResults(s.store, f)
	}
	m := s.meta(r, "glossary", "Kenyan words, Sheng and a Swahili phrasebook",
		"Matatu, boda boda, Paybill, poa: a searchable glossary of Swahili, Sheng and Kenyan English, a pocket Swahili phrasebook and answers to first-timers' questions.")
	m.Path = "/glossary"
	s.render(w, r, http.StatusOK, pages.Glossary(m, s.store, f), partial)
}
