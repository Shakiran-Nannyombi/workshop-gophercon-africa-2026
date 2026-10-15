package handlers

import "net/http"

// Placeholder: replaced when the cities page is built.
func (s *site) cities(w http.ResponseWriter, r *http.Request) { s.notFound(w, r) }

func (s *site) city(w http.ResponseWriter, r *http.Request) { s.notFound(w, r) }
