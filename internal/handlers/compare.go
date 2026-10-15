package handlers

import (
	"net/http"

	"workshop/internal/views/pages"

	"github.com/a-h/templ"
)

func (s *site) compare(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	a, b := q.Get("a"), q.Get("b")
	if a == "" && b == "" {
		a, b = "nairobi", "kampala"
	}
	swap := q.Get("swap") != ""
	if swap {
		a, b = b, a
	}
	v := pages.NewCompare(s.store, a, b)
	canonical := pages.CompareURL(v.A.Slug, v.B.Slug)
	if swap && !isHTMX(r) {
		http.Redirect(w, r, canonical+"#compare-result", http.StatusSeeOther)
		return
	}

	var partial templ.Component
	switch r.Header.Get("HX-Target") {
	case "compare-result":
		v.OOB = swap || r.Header.Get("HX-Trigger") != "compare-form" || v.A.Slug != a || v.B.Slug != b
		w.Header().Set("HX-Push-Url", canonical)
		partial = pages.CompareResult(s.store, v)
	}
	m := s.meta(r, "compare", v.A.Name+" vs "+v.B.Name+": compare East African cities",
		"Compare "+v.A.Name+" and "+v.B.Name+" side by side: cost of living, getting around, climate, safety perception, nightlife, currency and the price of a coffee.")
	s.render(w, r, http.StatusOK, pages.Compare(m, s.store, v), partial)
}
