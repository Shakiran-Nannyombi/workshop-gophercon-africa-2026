package components

import "workshop/internal/domain/sighting"

// WithPhotos keeps the first n sightings that have a still photograph.
func WithPhotos(sightings []sighting.Sighting, n int) []sighting.Sighting {
	out := make([]sighting.Sighting, 0, n)
	for _, s := range sightings {
		if _, ok := stillImage(s); !ok {
			continue
		}
		out = append(out, s)
		if len(out) == n {
			break
		}
	}
	return out
}

func stillImage(s sighting.Sighting) (sighting.Media, bool) {
	for _, m := range s.Media {
		if m.Type == "StillImage" && m.URL != "" {
			return m, true
		}
	}
	return sighting.Media{}, false
}
