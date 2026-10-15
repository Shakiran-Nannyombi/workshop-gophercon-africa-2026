package handlers

import (
	"bytes"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"workshop/internal/data"
	"workshop/internal/views/pages"
)

func (s *site) countries(w http.ResponseWriter, r *http.Request) {
	m := s.meta(r, "countries", "Countries", "Ten East African countries: capitals, currencies, how addresses work and how to get in.")
	s.render(w, r, http.StatusOK, pages.CountriesHub(m, s.store), nil)
}

func (s *site) country(w http.ResponseWriter, r *http.Request) {
	c, ok := s.store.CountryBySlug(r.PathValue("slug"))
	if !ok {
		s.notFound(w, r)
		return
	}
	m := s.meta(r, "countries", c.Name, c.Name+" at a glance: capital, currency, how the country is organised, getting in and getting around.")
	s.render(w, r, http.StatusOK, pages.CountryPage(m, s.store, c), nil)
}

func (s *site) cities(w http.ResponseWriter, r *http.Request) {
	m := s.meta(r, "cities", "City guides", "Guides to East Africa's capitals and major cities, all in one shared layout.")
	country := s.pickerCountry(r, "")
	s.render(w, r, http.StatusOK, pages.CitiesHub(m, s.store, country), pages.CityPicker(s.store, country, "/cities"))
}

func (s *site) city(w http.ResponseWriter, r *http.Request) {
	c, ok := s.store.City(r.PathValue("slug"))
	if !ok {
		s.notFound(w, r)
		return
	}
	m := s.meta(r, "cities", c.Name, c.Name+": "+c.Tagline)
	s.render(w, r, http.StatusOK, pages.CityPage(m, s.store, c), nil)
}

// pickerCountry reads ?country= for the city picker: "all" (or unknown) means every country.
func (s *site) pickerCountry(r *http.Request, fallback string) string {
	q := r.URL.Query().Get("country")
	if q == "all" {
		return ""
	}
	if _, ok := s.store.Country(q); ok {
		return q
	}
	return fallback
}

// conversion reads amount/from/to from the query, falling back to the given pair.
func (s *site) conversion(r *http.Request, from, to string) pages.Conversion {
	q := r.URL.Query()
	amount := 1000.0
	if v, err := strconv.ParseFloat(strings.ReplaceAll(q.Get("amount"), ",", ""), 64); err == nil {
		amount = v
	}
	if f := strings.ToUpper(q.Get("from")); pages.ValidCurrency(s.store, f) {
		from = f
	}
	if t := strings.ToUpper(q.Get("to")); pages.ValidCurrency(s.store, t) {
		to = t
	}
	return pages.Convert(s.store, amount, from, to)
}

func (s *site) currencies(w http.ResponseWriter, r *http.Request) {
	m := s.meta(r, "currencies", "Currencies", "Every East African currency from zero: notes, coins, mobile money, cash and a converter.")
	s.render(w, r, http.StatusOK, pages.CurrenciesHub(m, s.store, s.conversion(r, "KES", "USD")), nil)
}

func (s *site) convert(w http.ResponseWriter, r *http.Request) {
	c := s.conversion(r, "KES", "USD")
	if !isHTMX(r) {
		http.Redirect(w, r, "/currencies?"+r.URL.RawQuery+"#converter", http.StatusSeeOther)
		return
	}
	s.render(w, r, http.StatusOK, pages.ConverterResult(s.store, c), nil)
}

func (s *site) plan(r *http.Request) pages.Plan {
	p := pages.Plan{}
	if _, ok := s.store.City(r.URL.Query().Get("from")); ok {
		p.From = r.URL.Query().Get("from")
	}
	if _, ok := s.store.City(r.URL.Query().Get("to")); ok {
		p.To = r.URL.Query().Get("to")
	}
	if p.From != "" && p.To != "" && p.From != p.To {
		p.Route, _ = s.store.RouteBetween(p.From, p.To)
	}
	return p
}

func (s *site) routes(w http.ResponseWriter, r *http.Request) {
	m := s.meta(r, "routes", "Routes between cities", "Flights, buses, the SGR train and ferries between East Africa's cities, with rough times and fares.")
	s.render(w, r, http.StatusOK, pages.RoutesPage(m, s.store, s.plan(r)), nil)
}

func (s *site) planFragment(w http.ResponseWriter, r *http.Request) {
	if !isHTMX(r) {
		http.Redirect(w, r, "/routes?"+r.URL.RawQuery+"#planner", http.StatusSeeOther)
		return
	}
	s.render(w, r, http.StatusOK, pages.PlanResult(s.store, s.plan(r)), nil)
}

func (s *site) borders(w http.ResponseWriter, r *http.Request) {
	m := s.meta(r, "borders", "Crossing borders", "East Africa's main land border crossings, the East African Community and the East Africa Tourist Visa.")
	s.render(w, r, http.StatusOK, pages.BordersPage(m, s.store), nil)
}

func (s *site) connectivity(w http.ResponseWriter, r *http.Request) {
	m := s.meta(r, "connectivity", "SIM cards & data", "SIM cards, mobile data, eSIMs and roaming across East Africa, country by country.")
	s.render(w, r, http.StatusOK, pages.ConnectivityPage(m, s.store), nil)
}

func (s *site) compare(w http.ResponseWriter, r *http.Request) {
	pick := func(key, fallback string) *data.City {
		if c, ok := s.store.City(r.URL.Query().Get(key)); ok {
			return c
		}
		c, _ := s.store.City(fallback)
		return c
	}
	p := pages.ComparePair{A: pick("a", "nairobi"), B: pick("b", "kampala")}
	m := s.meta(r, "compare", "Compare cities", "Compare two East African cities side by side: cost, getting around, climate and currency.")
	s.render(w, r, http.StatusOK, pages.ComparePage(m, s.store, p), nil)
}

func (s *site) glossary(w http.ResponseWriter, r *http.Request) {
	m := s.meta(r, "glossary", "Words", "Everyday words for getting around and paying in East Africa: matatu, daladala, boda boda, paybill, MoMo and more.")
	s.render(w, r, http.StatusOK, pages.GlossaryPage(m, s.store), nil)
}

func (s *site) about(w http.ResponseWriter, r *http.Request) {
	m := s.meta(r, "about", "About & sources", "Where Karibu Guide's facts come from, official travel advisories and official visa sites.")
	s.render(w, r, http.StatusOK, pages.AboutPage(m, s.store, s.verifyCount), nil)
}

// countVerify counts the "verify" notes in the content files, for the About page.
func countVerify(fsys fs.FS) int {
	n := 0
	files, _ := fs.Glob(fsys, "*.json")
	for _, f := range files {
		if b, err := fs.ReadFile(fsys, f); err == nil {
			n += bytes.Count(b, []byte(`"verify":`))
		}
	}
	return n
}

func (s *site) currency(w http.ResponseWriter, r *http.Request) {
	code := strings.ToUpper(r.PathValue("code"))
	c, ok := s.store.Currency(code)
	if !ok {
		s.notFound(w, r)
		return
	}
	if r.PathValue("code") != code {
		http.Redirect(w, r, "/currencies/"+code, http.StatusMovedPermanently)
		return
	}
	other := "KES"
	if code == "KES" {
		other = "USD"
	}
	m := s.meta(r, "currencies", c.Name+" ("+c.Code+")", "The "+c.Name+" from zero: notes and coins, how prices are written, mobile money and cash.")
	s.render(w, r, http.StatusOK, pages.CurrencyPage(m, s.store, c, s.conversion(r, code, other)), nil)
}
