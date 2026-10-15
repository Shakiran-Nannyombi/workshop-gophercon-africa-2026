package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"workshop/internal/views/pages"

	"github.com/a-h/templ"
)

func (s *site) money(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	conv := s.conversion(q.Get("amount"), q.Get("from"), q.Get("to"))
	flow := q.Get("flow")

	var partial templ.Component
	switch r.Header.Get("HX-Target") {
	case "convert-result":
		partial = pages.ConvertResult(conv)
	case "mpesa-flow":
		partial = pages.MpesaFlow(s.store, flow)
	}
	m := s.meta(r, "money", "Kenyan money explained: shillings, M-Pesa and prices",
		"Kenyan shilling notes and coins, M-Pesa from zero, a currency converter, real prices in KES, scam tips, and how Kenya compares with Uganda, Tanzania, Rwanda and Ethiopia.")
	s.render(w, r, http.StatusOK, pages.Money(m, s.store, flow, conv), partial)
}

// conversion parses converter input, defaulting to 1,000 KES in USD.
func (s *site) conversion(amount, from, to string) pages.Conversion {
	a := 1000.0
	if amount = strings.TrimSpace(strings.ReplaceAll(amount, ",", "")); amount != "" {
		v, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			v = -1
		}
		a = v
	}
	if _, ok := s.store.Site.Rates.PerUSD[from]; !ok {
		from = "KES"
	}
	if _, ok := s.store.Site.Rates.PerUSD[to]; !ok {
		to = "USD"
	}
	return pages.Convert(s.store, a, from, to)
}
