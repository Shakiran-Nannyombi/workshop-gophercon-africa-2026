package components

import "strings"

func costRank(c string) int {
	switch c {
	case "budget":
		return 1
	case "mid", "mixed":
		return 2
	case "upmarket":
		return 3
	}
	return 0
}

func costLabel(c string) string {
	switch c {
	case "budget":
		return "Budget"
	case "mid":
		return "Mid-range"
	case "upmarket":
		return "Upmarket"
	case "mixed":
		return "Mixed"
	}
	return "Unknown"
}

var suitLabels = map[string]string{
	"visitors": "Visitors", "expats": "Expats", "families": "Families",
	"students": "Students", "budget": "Tight budgets", "commuters": "Commuters",
}

func suitLabel(s string) string {
	if l, ok := suitLabels[s]; ok {
		return l
	}
	return s
}

func join(xs []string) string { return strings.Join(xs, " · ") }
