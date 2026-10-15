package data

import "math"

// Geo holds the pre-simplified map shapes written by cmd/mapgen.
type Geo struct {
	Kenya      ShapeSet
	Nairobi    ShapeSet
	EastAfrica ShapeSet
}

type ShapeSet struct {
	ViewBox    string     `json:"viewBox"`
	Source     string     `json:"source"`
	License    string     `json:"license"`
	Projection Projection `json:"projection"`
	Counties   []Shape    `json:"counties"`
	Countries  []Shape    `json:"countries"`
	Lakes      []Shape    `json:"lakes"`
}

type Shape struct {
	Slug string  `json:"slug"`
	Code string  `json:"code"`
	Name string  `json:"name"`
	D    string  `json:"d"`
	CX   float64 `json:"cx"`
	CY   float64 `json:"cy"`
}

// Projection is the equirectangular box a ShapeSet was drawn in.
type Projection struct {
	Lon0   float64 `json:"lon0"`
	Lat0   float64 `json:"lat0"`
	Lon1   float64 `json:"lon1"`
	Lat1   float64 `json:"lat1"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Project maps a lat/lon to viewBox coordinates, matching cmd/mapgen.
func (p Projection) Project(lat, lon float64) (x, y float64) {
	if p.Width == 0 || p.Lon1 == p.Lon0 {
		return 0, 0
	}
	minLat, maxLat := math.Min(p.Lat0, p.Lat1), math.Max(p.Lat0, p.Lat1)
	x = (lon - p.Lon0) / (p.Lon1 - p.Lon0) * p.Width
	y = (maxLat - lat) / (maxLat - minLat) * p.Height
	return x, y
}
