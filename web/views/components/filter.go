package components

import (
	"strconv"
	"strings"

	"github.com/a-h/templ"

	"workshop/internal/domain/sighting"
)

// Animal names the animal seen, as the card does.
func Animal(s sighting.Sighting) string { return displayName(s.Species) }

// Country is where it was seen.
func Country(s sighting.Sighting) string { return s.Location.Country }

// Year is when, to the year.
func Year(s sighting.Sighting) string { return strconv.Itoa(s.HappenedAt.Year()) }

// FilterOptions counts how often each key appears, keeping the n most common; 0 keeps all.
func FilterOptions(sightings []sighting.Sighting, n int, key func(sighting.Sighting) string) []sighting.Tally {
	return sighting.Top(sighting.CountBy(sightings, key), n)
}

// FilterAttrs marks a list item with what the filter matches on.
func FilterAttrs(s sighting.Sighting) templ.Attributes {
	return templ.Attributes{
		"data-animal":  Animal(s),
		"data-country": Country(s),
		"data-year":    Year(s),
		"data-search": strings.ToLower(strings.Join([]string{
			Animal(s), s.Species.ScientificName, place(s.Location), s.RecordedBy,
		}, " ")),
	}
}
