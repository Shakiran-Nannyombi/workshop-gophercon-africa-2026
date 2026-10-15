package components

import (
	"fmt"
	"strconv"
	"strings"

	"workshop/internal/data"
)

// CompareDim is one scored dimension of the compare chart and table.
type CompareDim struct {
	Label, Five string
	Pick        func(data.Scores) data.Score
}

// CompareDims are the scored dimensions, in display order.
var CompareDims = []CompareDim{
	{"Cost of living", "5 = most expensive", func(s data.Scores) data.Score { return s.Cost }},
	{"Ease of transport", "5 = easiest to get around", func(s data.Scores) data.Score { return s.Transport }},
	{"Climate", "5 = mildest weather", func(s data.Scores) data.Score { return s.Climate }},
	{"Safety perception", "5 = most relaxed feel", func(s data.Scores) data.Score { return s.Safety }},
	{"Nightlife", "5 = liveliest", func(s data.Scores) data.Score { return s.Nightlife }},
}

// cmpMax is the share of the plot width that a score of 5 fills, leaving room for labels.
const cmpMax = 80.0

// cmpPct maps a 0-5 score to an x position as a percentage of the plot width.
func cmpPct(v int) string {
	return strconv.FormatFloat(float64(min(max(v, 0), 5))*cmpMax/5, 'f', 1, 64) + "%"
}

// CityCode is a short three-letter tag that labels a city's bars.
func CityCode(c data.City) string {
	n := strings.ToUpper(strings.ReplaceAll(c.Name, " ", ""))
	if len(n) > 3 {
		return n[:3]
	}
	return n
}

func cmpTitle(c data.City, d CompareDim, sc data.Score) string {
	return fmt.Sprintf("%s, %s: %d of 5. %s", c.Name, strings.ToLower(d.Label), sc.Value, sc.Reason)
}

func cmpRowLabel(a, b data.City, d CompareDim) string {
	return fmt.Sprintf("%s (%s): %s %d, %s %d out of 5", d.Label, d.Five, a.Name, d.Pick(a.Scores).Value, b.Name, d.Pick(b.Scores).Value)
}

func cmpTicks() []int { return []int{0, 1, 2, 3, 4, 5} }

func cmpTickAnchor(t int) string {
	if t == 0 {
		return "start"
	}
	return "middle"
}

func cmpList(xs []string) string {
	if len(xs) == 0 {
		return "Not listed"
	}
	return strings.Join(xs, ", ")
}

func cmpCurrency(c *data.Country) string {
	if c == nil {
		return "Not listed"
	}
	return c.Currency + " (" + c.CurrencyCode + ")"
}

func cmpMobileMoney(c *data.Country) []string {
	if c == nil {
		return nil
	}
	return c.MobileMoney
}
