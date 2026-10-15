package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// pagePaths lists every indexable URL, for the sitemap.
func (s *site) pagePaths() []string {
	return []string{"/"}
}

func (s *site) sitemap(w http.ResponseWriter, r *http.Request) {
	type url struct {
		Loc     string `xml:"loc"`
		LastMod string `xml:"lastmod"`
	}
	type urlset struct {
		XMLName xml.Name `xml:"urlset"`
		NS      string   `xml:"xmlns,attr"`
		URLs    []url    `xml:"url"`
	}
	set := urlset{NS: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	for _, p := range s.pagePaths() {
		set.URLs = append(set.URLs, url{Loc: s.store.Site.BaseURL + p, LastMod: s.store.Site.LastReviewed})
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	fmt.Fprint(w, xml.Header)
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(set); err != nil {
		s.log.ErrorContext(r.Context(), "handlers: writing sitemap: "+err.Error())
	}
}

func (s *site) robots(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "User-agent: *\nAllow: /\nDisallow: /design-system\n\nSitemap: %s/sitemap.xml\n", s.store.Site.BaseURL)
}
