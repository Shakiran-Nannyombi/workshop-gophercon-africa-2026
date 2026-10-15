package components

import (
	"fmt"
	"strings"

	"workshop/internal/data"
)

// roadPath projects a road's lat/lon waypoints into an SVG path.
func roadPath(p data.Projection, r data.Road) string {
	var b strings.Builder
	for i, pt := range r.Points {
		x, y := p.Project(pt[0], pt[1])
		if i == 0 {
			fmt.Fprintf(&b, "M%.1f %.1f", x, y)
		} else {
			fmt.Fprintf(&b, " L%.1f %.1f", x, y)
		}
	}
	return b.String()
}

// roadLabelPoint is where a road's name sits: part-way along it.
func roadLabelPoint(p data.Projection, r data.Road) (string, string) {
	if len(r.Points) == 0 {
		return "0", "0"
	}
	pt := r.Points[len(r.Points)*2/3]
	if len(r.Points) == 2 {
		pt = [2]float64{(r.Points[0][0] + r.Points[1][0]) / 2, (r.Points[0][1] + r.Points[1][1]) / 2}
	}
	x, y := p.Project(pt[0], pt[1])
	return f1(x), f1(y)
}

func areaXY(p data.Projection, a data.Area) (string, string) {
	x, y := p.Project(a.Lat, a.Lon)
	return f1(x), f1(y)
}

func f1(v float64) string { return fmt.Sprintf("%.1f", v) }

func countyClass(slug string) string { return "map-county map-county--" + slug }

func dotClass(a data.Area) string { return "map-dot map-dot--" + a.County }


// nairobiLabels places county names where they're readable in the default Nairobi view.
var nairobiLabels = map[string][2]float64{
	"nairobi": {745, 480}, "kiambu": {420, 240}, "kajiado": {270, 790}, "machakos": {760, 730},
}

func countyLabelXY(c data.Shape) (string, string) {
	if p, ok := nairobiLabels[c.Slug]; ok {
		return f1(p[0]), f1(p[1])
	}
	return f1(c.CX), f1(c.CY)
}
