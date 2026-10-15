// Package handlers serves Karibu Guide's pages. Every handler returns a full
// page normally, and only the changed fragment when htmx asks (HX-Request).
package handlers

import (
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	content "workshop/data"
	"workshop/internal/data"
	"workshop/internal/views/layouts"
	"workshop/internal/views/pages"
	"workshop/static"

	"github.com/a-h/templ"
)

// Options configure the site.
type Options struct {
	// Dev serves /static and /data from disk so CSS and JSON edits show on refresh.
	Dev bool
	// BaseURL overrides site.json's baseURL for canonical links and the sitemap.
	BaseURL string
	// RateKES overrides site.json's KES-per-USD rate.
	RateKES float64
}

type site struct {
	log   *slog.Logger
	store *data.Store
	opts  Options
}

// New loads the content and returns the site's handler.
func New(logger *slog.Logger, opts Options) (http.Handler, error) {
	var dataFS, staticFS fs.FS = content.FS, static.FS
	if opts.Dev {
		dataFS, staticFS = os.DirFS("data"), os.DirFS("static")
	}
	store, err := data.Load(dataFS)
	if err != nil {
		return nil, err
	}
	if opts.BaseURL != "" {
		store.Site.BaseURL = strings.TrimRight(opts.BaseURL, "/")
	}
	if opts.RateKES > 0 {
		store.Site.Rates.PerUSD["KES"] = opts.RateKES
	}
	s := &site{log: logger, store: store, opts: opts}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", s.static(staticFS))
	mux.HandleFunc("GET /{$}", s.home)
	mux.HandleFunc("GET /money", s.money)
	mux.HandleFunc("GET /nairobi", s.nairobi)
	mux.HandleFunc("GET /counties", s.counties)
	mux.HandleFunc("GET /counties/{slug}", s.county)
	mux.HandleFunc("GET /cities", s.cities)
	mux.HandleFunc("GET /cities/{slug}", s.city)
	mux.HandleFunc("GET /compare", s.compare)
	mux.HandleFunc("GET /glossary", s.glossary)
	mux.HandleFunc("GET /about", s.about)
	mux.HandleFunc("GET /design-system", s.designSystem)
	mux.HandleFunc("GET /sitemap.xml", s.sitemap)
	mux.HandleFunc("GET /robots.txt", s.robots)
	mux.HandleFunc("GET /", s.notFound)
	return mux, nil
}

func (s *site) static(fsys fs.FS) http.Handler {
	files := http.StripPrefix("/static/", http.FileServerFS(fsys))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			s.notFound(w, r)
			return
		}
		if !s.opts.Dev {
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		files.ServeHTTP(w, r)
	})
}

// meta builds the layout props shared by every page.
func (s *site) meta(r *http.Request, section, title, description string) layouts.Meta {
	return layouts.Meta{
		Site:        s.store.Site,
		Path:        r.URL.Path,
		Section:     section,
		Title:       title,
		Description: description,
	}
}

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true" && r.Header.Get("HX-Boosted") != "true"
}

// render writes the partial for htmx requests and the full page otherwise.
func (s *site) render(w http.ResponseWriter, r *http.Request, status int, full, partial templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Add("Vary", "HX-Request")
	c := full
	if partial != nil && isHTMX(r) {
		c = partial
	}
	w.WriteHeader(status)
	if err := c.Render(r.Context(), w); err != nil {
		s.log.ErrorContext(r.Context(), "handlers: rendering "+r.URL.Path, slog.Any("err", err))
	}
}

func (s *site) home(w http.ResponseWriter, r *http.Request) {
	m := s.meta(r, "home", "", "A visual, beginner-friendly guide to Kenyan money, the 47 counties, Nairobi's neighbourhoods and East Africa's main cities.")
	country := r.URL.Query().Get("country")
	if _, ok := s.store.Country(country); !ok {
		country = "KE"
	}
	s.render(w, r, http.StatusOK, pages.Home(m, s.store, country), pages.CityPicker(s.store, country))
}

func (s *site) designSystem(w http.ResponseWriter, r *http.Request) {
	m := s.meta(r, "", "Design system", "Colours, lettering and components used across Karibu Guide.")
	m.NoIndex = true
	s.render(w, r, http.StatusOK, pages.DesignSystem(m), nil)
}

func (s *site) notFound(w http.ResponseWriter, r *http.Request) {
	m := s.meta(r, "", "Page not found", "")
	m.NoIndex = true
	s.render(w, r, http.StatusNotFound, pages.NotFound(m), nil)
}
