package handlers

import (
	"net/http"

	"workshop/internal/views/pages"
)

func (s *site) about(w http.ResponseWriter, r *http.Request) {
	m := s.meta(r, "", "About, sources and how to edit",
		"What Karibu Guide is, who it's for, how to edit its JSON content, how approximate figures and exchange rates work, and the sources and map data behind it.")
	m.Path = "/about"
	s.render(w, r, http.StatusOK, pages.About(m, s.store), nil)
}
