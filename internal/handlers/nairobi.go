package handlers

import (
	"net/http"

	"workshop/internal/views/pages"

	"github.com/a-h/templ"
)

func (s *site) nairobi(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := pages.AreaFilter{Q: q.Get("q"), County: q.Get("county"), Cost: q.Get("cost"), Suits: q.Get("suits")}
	if f.County != "" {
		if _, ok := s.store.County(f.County); !ok {
			f.County = ""
		}
	}
	l := pages.WhichCounty(s.store, q.Get("where"))

	var partial templ.Component
	switch r.Header.Get("HX-Target") {
	case "which-result":
		partial = pages.WhichResult(l)
	case "area-results":
		partial = pages.AreaResults(s.store, f)
	}
	m := s.meta(r, "nairobi", "Nairobi explained: city, county, areas and getting around",
		"Nairobi is a city and a county. See which 'Nairobi' places are really in Kiambu, Kajiado or Machakos, browse 50+ areas, and learn matatus, SIMs, food, safety and weather.")
	m.Path = "/nairobi"
	s.render(w, r, http.StatusOK, pages.Nairobi(m, s.store, f, l), partial)
}
