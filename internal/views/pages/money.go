package pages

import (
	"math"
	"sort"
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

type currencyOption struct{ Code, Name string }

// converterCurrencies is every East African currency with a rate, then USD, EUR and GBP.
func converterCurrencies(s *data.Store) []currencyOption {
	var out []currencyOption
	for _, c := range s.Currencies {
		if s.Site.Rates.PerUSD[c.Code] > 0 {
			out = append(out, currencyOption{c.Code, c.Name})
		}
	}
	return append(out, currencyOption{"USD", "US dollar"}, currencyOption{"EUR", "Euro"}, currencyOption{"GBP", "British pound"})
}

// ValidCurrency reports whether code can be converted.
func ValidCurrency(s *data.Store, code string) bool { return s.Site.Rates.PerUSD[code] > 0 }

// MentalTip is a "quick mental maths" rule of thumb for one currency.
type MentalTip struct {
	Code, Symbol  string
	VsKES, VsUSD  string // "~28 UGX per KSh 1", "~3,650 per US$1"
	Example, Rule string
	Valid         bool
}

// mentalTip builds the rule of thumb from the configured rates.
func mentalTip(s *data.Store, code string) MentalTip {
	r := s.Site.Rates.PerUSD
	x, kes := r[code], r["KES"]
	t := MentalTip{Code: code, Symbol: components.Sym(code)}
	if x == 0 || kes == 0 {
		return t
	}
	t.Valid = true
	t.VsUSD = "~" + roughNum(x) + " " + code + " to US$1"
	if code == "KES" {
		t.VsKES = "KSh 1,000 is about US$" + strconv.FormatFloat(1000/kes, 'f', 2, 64)
		t.Rule = "Roughly: drop the last two zeros and knock off a quarter to get US dollars (KSh 1,000 → 10 → ~US$7.50)."
		return t
	}
	perKES := x / kes
	if perKES >= 1 {
		t.VsKES = "~" + roughNum(perKES) + " " + code + " per KSh 1"
		t.Rule = "Divide a " + code + " price by ~" + roughNum(perKES) + " to get Kenyan shillings."
	} else {
		t.VsKES = "~KSh " + roughNum(1/perKES) + " per 1 " + code
		t.Rule = "Multiply a " + code + " price by ~" + roughNum(1/perKES) + " to get Kenyan shillings."
	}
	sample := roundNice(x * 10)
	t.Example = t.Symbol + " " + components.Commas(int(sample)) + " ≈ US$" + components.FormatMoney(sample/x, "USD") + " ≈ KSh " + components.Commas(int(math.Round(sample/x*kes)))
	return t
}

// roughNum rounds to two significant figures for mental maths.
func roughNum(v float64) string {
	if v <= 0 {
		return "0"
	}
	if v < 10 {
		return strconv.FormatFloat(math.Round(v*10)/10, 'f', -1, 64)
	}
	p := math.Pow(10, math.Floor(math.Log10(v))-1)
	return components.Commas(int(math.Round(v/p) * p))
}

// roundNice rounds up to one significant figure, e.g. 36,500 → 40,000.
func roundNice(v float64) float64 {
	if v <= 0 {
		return 0
	}
	p := math.Pow(10, math.Floor(math.Log10(v)))
	return math.Ceil(v/p) * p
}

// ScaleRow is one bar in "why are the numbers so different?".
type ScaleRow struct {
	Code, Name, Label string
	Width             string
}

// scaleRows prices the same ~US$5 item in every currency, biggest number last.
func scaleRows(s *data.Store, usd float64) []ScaleRow {
	var out []ScaleRow
	for _, c := range append([]currencyOption{{"USD", "US dollar"}}, converterCurrencies(s)...) {
		if c.Code == "EUR" || c.Code == "GBP" || (c.Code == "USD" && len(out) > 0) {
			continue
		}
		r := s.Site.Rates.PerUSD[c.Code]
		if r == 0 {
			continue
		}
		out = append(out, ScaleRow{c.Code, c.Name, components.Sym(c.Code) + " " + components.Commas(int(math.Round(usd*r))), scaleWidth(r)})
	}
	sort.SliceStable(out, func(i, j int) bool {
		return s.Site.Rates.PerUSD[out[i].Code] < s.Site.Rates.PerUSD[out[j].Code]
	})
	return out
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
