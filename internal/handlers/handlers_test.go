package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestSite(t *testing.T) http.Handler {
	t.Helper()
	h, err := New(slog.New(slog.NewTextHandler(io.Discard, nil)), Options{})
	if err != nil {
		t.Fatalf("loading site: %v", err)
	}
	return h
}

func get(t *testing.T, h http.Handler, url string, headers map[string]string) (int, string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func TestPagesRenderFullHTML(t *testing.T) {
	h := newTestSite(t)
	for _, url := range []string{"/", "/money", "/nairobi"} {
		code, body := get(t, h, url, nil)
		if code != http.StatusOK {
			t.Errorf("%s: status %d", url, code)
		}
		if !strings.Contains(body, "<!doctype html>") || !strings.Contains(body, `og:title`) {
			t.Errorf("%s: expected a full page with Open Graph tags", url)
		}
	}
}

func TestHTMXReturnsPartials(t *testing.T) {
	h := newTestSite(t)
	tests := []struct{ url, target, want string }{
		{"/money?amount=1000&from=KES&to=USD", "convert-result", `id="convert-result"`},
		{"/money?flow=paybill", "mpesa-flow", `id="mpesa-flow"`},
		{"/nairobi?where=ruaka", "which-result", "Kiambu"},
		{"/nairobi?county=kajiado", "area-results", `id="area-results"`},
		{"/?country=UG", "", `id="city-picker"`},
	}
	for _, tt := range tests {
		code, body := get(t, h, tt.url, map[string]string{"HX-Request": "true", "HX-Target": tt.target})
		if code != http.StatusOK {
			t.Errorf("%s: status %d", tt.url, code)
		}
		if strings.Contains(body, "<!doctype html>") {
			t.Errorf("%s: got a full page, want a partial", tt.url)
		}
		if !strings.Contains(body, tt.want) {
			t.Errorf("%s: partial missing %q", tt.url, tt.want)
		}
	}
}

func TestNotFound(t *testing.T) {
	h := newTestSite(t)
	code, body := get(t, h, "/no-such-page", nil)
	if code != http.StatusNotFound || !strings.Contains(body, "wrong matatu") {
		t.Fatalf("got %d, want the 404 page", code)
	}
}

func TestConversionRejectsBadInput(t *testing.T) {
	h := newTestSite(t)
	_, body := get(t, h, "/money?amount=-5&from=KES&to=USD", map[string]string{"HX-Request": "true", "HX-Target": "convert-result"})
	if !strings.Contains(body, "Enter an amount") {
		t.Fatal("negative amount should ask for a valid amount")
	}
}

func TestSitemapListsCounties(t *testing.T) {
	h := newTestSite(t)
	_, body := get(t, h, "/sitemap.xml", nil)
	if !strings.Contains(body, "/counties/kiambu") || !strings.Contains(body, "/cities/kampala") {
		t.Fatal("sitemap should list county and city pages")
	}
}
