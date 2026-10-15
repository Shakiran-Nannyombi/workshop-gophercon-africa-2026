package pages

import (
	"math"
	"strconv"

	"workshop/internal/data"
	"workshop/internal/views/components"
)

// Conversion is one converter result.
type Conversion struct {
	Amount   float64
	From, To string
	Result   float64
	Unit     float64 // 1 From in To
	AsOf     string
	Valid    bool
}

var ConverterCurrencies = []struct{ Code, Name string }{
	{"KES", "Kenyan shilling"}, {"USD", "US dollar"}, {"EUR", "Euro"}, {"GBP", "British pound"},
	{"UGX", "Ugandan shilling"}, {"TZS", "Tanzanian shilling"}, {"RWF", "Rwandan franc"}, {"ETB", "Ethiopian birr"},
}

// Convert converts amount between two currencies through the configured USD rates.
func Convert(s *data.Store, amount float64, from, to string) Conversion {
	r := s.Site.Rates.PerUSD
	c := Conversion{Amount: amount, From: from, To: to, AsOf: s.Site.Rates.AsOf}
	rf, okf := r[from]
	rt, okt := r[to]
	if !okf || !okt || rf == 0 || math.IsNaN(amount) || math.IsInf(amount, 0) || amount < 0 || amount > 1e12 {
		return c
	}
	c.Unit = rt / rf
	c.Result = amount * c.Unit
	c.Valid = true
	return c
}

func fmtAmount(v float64, code string) string { return components.FormatMoney(v, code) }

func amountValue(v float64) string {
	if v == math.Trunc(v) {
		return strconv.FormatFloat(v, 'f', 0, 64)
	}
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func pricesIn(ps []data.PriceExample, cat string) []data.PriceExample {
	var out []data.PriceExample
	for _, p := range ps {
		if p.Category == cat {
			out = append(out, p)
		}
	}
	return out
}

var priceCategories = []struct{ Key, Title, Icon string }{
	{"transport", "Getting around", "bus"},
	{"food", "Food & drink", "food"},
	{"stay", "A bed for the night", "city"},
	{"connectivity", "Staying connected", "phone"},
}

func flowBySlug(m data.Money, slug string) data.MpesaFlow {
	for _, f := range m.MpesaSteps {
		if f.Slug == slug {
			return f
		}
	}
	if len(m.MpesaSteps) > 0 {
		return m.MpesaSteps[0]
	}
	return data.MpesaFlow{}
}

// scaleWidth turns units-per-dollar into a bar width on a log scale (1 → ~8%, 10,000 → 100%).
func scaleWidth(perUSD float64) string {
	if perUSD <= 0 {
		return "width:0"
	}
	w := 8 + math.Log10(math.Max(perUSD, 1))/4*92
	return "width:" + strconv.FormatFloat(math.Min(w, 100), 'f', 1, 64) + "%"
}

func coffeeUSD(s *data.Store, c data.Country) string {
	r := s.Site.Rates.PerUSD[c.CurrencyCode]
	if r == 0 || c.CoffeeLocal == 0 {
		return "?"
	}
	return strconv.FormatFloat(float64(c.CoffeeLocal)/r, 'f', 2, 64)
}

func perUSD(s *data.Store, code string) string {
	return components.FormatMoney(s.Site.Rates.PerUSD[code], code)
}

func countryMoney(s *data.Store, code string) []data.MobileMoney {
	var out []data.MobileMoney
	for _, m := range s.Money.MobileMoney {
		if m.Country == code {
			out = append(out, m)
		}
	}
	return out
}

// fmtUnit prints a unit rate with enough precision to be useful (0.0078, 28.3, 3,650).
func fmtUnit(v float64) string {
	switch {
	case v < 0.1:
		return strconv.FormatFloat(v, 'f', 5, 64)
	case v < 100:
		return strconv.FormatFloat(v, 'f', 2, 64)
	}
	return components.Commas(int(math.Round(v)))
}
