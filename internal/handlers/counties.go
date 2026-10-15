package handlers

import (
	"net/http"
	"strings"

	"workshop/internal/views/components"
	"workshop/internal/views/pages"

	"github.com/a-h/templ"
)

func (s *site) counties(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := pages.CountyFilter{Q: strings.TrimSpace(q.Get("q")), Region: q.Get("region")}
	if _, ok := s.store.Region(f.Region); !ok {
		f.Region = ""
	}

	var partial templ.Component
	switch r.Header.Get("HX-Target") {
	case "county-results":
		partial = pages.CountyResults(s.store, f)
	}
	m := s.meta(r, "counties", "Kenya's 47 counties: map, regions and what each is like",
		"An interactive map of Kenya's 47 counties by region, with each county's seat, population, area and character, plus a plain guide to counties, constituencies, wards and estates.")
	s.render(w, r, http.StatusOK, pages.Counties(m, s.store, f), partial)
}

func (s *site) county(w http.ResponseWriter, r *http.Request) {
	c, ok := s.store.County(r.PathValue("slug"))
	if !ok {
		s.notFound(w, r)
		return
	}
	region := c.Region
	if rg, ok := s.store.Region(c.Region); ok {
		region = rg.Name
	}
	desc := components.CountyTitle(*c) + " (county " + components.CountyCode(c.Code) + ", " + region + " region): seat " + c.Seat +
		", ~" + components.Commas(c.Pop2019) + " people (2019 census), ~" + components.Commas(c.AreaKm2) + " km². Climate, food, transport, main towns and what it's good for."
	m := s.meta(r, "counties", components.CountyTitle(*c)+", Kenya: map, facts and what to expect", desc)
	s.render(w, r, http.StatusOK, pages.CountyPage(m, s.store, c), nil)
}
