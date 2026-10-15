package components

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"workshop/internal/data"
)

// Commas formats 1234567 as "1,234,567".
func Commas(n int) string {
	s := strconv.Itoa(int(math.Abs(float64(n))))
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	if n < 0 {
		return "-" + b.String()
	}
	return b.String()
}

// FormatMoney prints an amount with grouping and sensible decimals for the currency.
func FormatMoney(v float64, code string) string {
	switch code {
	case "USD", "EUR", "GBP", "ETB":
		whole := math.Floor(math.Abs(v))
		frac := int(math.Round((math.Abs(v) - whole) * 100))
		if frac == 100 {
			whole, frac = whole+1, 0
		}
		sign := ""
		if v < 0 {
			sign = "-"
		}
		return fmt.Sprintf("%s%s.%02d", sign, Commas(int(whole)), frac)
	}
	if math.Abs(v) < 10 {
		return strconv.FormatFloat(v, 'f', 1, 64)
	}
	return Commas(int(math.Round(v)))
}

// Symbols are how each currency is usually written on a price tag.
var Symbols = map[string]string{
	"KES": "KSh", "UGX": "USh", "TZS": "TSh", "RWF": "FRw", "ETB": "Br",
	"USD": "US$", "EUR": "€", "GBP": "£",
}

func Sym(code string) string {
	if s, ok := Symbols[code]; ok {
		return s
	}
	return code
}

// KESRange prints "KSh 30–150", or a single figure when min == max.
func KESRange(p data.PriceExample) string {
	if p.MaxKES <= p.MinKES {
		return "KSh " + Commas(p.MinKES)
	}
	return "KSh " + Commas(p.MinKES) + "–" + Commas(p.MaxKES)
}

// USDRange converts a KES range with the configured rate, rounded for display.
func USDRange(p data.PriceExample, perUSD float64) string {
	if perUSD == 0 {
		return ""
	}
	lo, hi := float64(p.MinKES)/perUSD, float64(p.MaxKES)/perUSD
	f := func(v float64) string {
		if v < 10 {
			return strconv.FormatFloat(v, 'f', 2, 64)
		}
		return Commas(int(math.Round(v)))
	}
	if p.MaxKES <= p.MinKES {
		return "≈ US$" + f(lo)
	}
	return "≈ US$" + f(lo) + "–" + f(hi)
}

// guilloche returns rosette path data for a banknote's security pattern.
func guilloche(cx, cy, r float64, petals int) string {
	var b strings.Builder
	for i := 0; i <= 360; i += 3 {
		t := float64(i) * math.Pi / 180
		rr := r * (0.72 + 0.28*math.Cos(float64(petals)*t))
		x, y := cx+rr*math.Cos(t), cy+rr*math.Sin(t)
		if i == 0 {
			fmt.Fprintf(&b, "M%.1f %.1f", x, y)
		} else {
			fmt.Fprintf(&b, "L%.1f %.1f", x, y)
		}
	}
	return b.String() + "Z"
}

func noteTilt(i int) string {
	return fmt.Sprintf("--i:%d;--tilt:%ddeg", i, (i-2)*5)
}

func itoa(n int) string { return strconv.Itoa(n) }
