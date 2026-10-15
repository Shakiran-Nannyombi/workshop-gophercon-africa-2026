package handlers

import (
	"net/http"
	"strconv"
	"strings"

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
