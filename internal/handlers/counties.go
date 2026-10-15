package handlers

import "net/http"

// Placeholder: replaced when the counties page is built.
func (s *site) counties(w http.ResponseWriter, r *http.Request) { s.notFound(w, r) }

func (s *site) county(w http.ResponseWriter, r *http.Request) { s.notFound(w, r) }
