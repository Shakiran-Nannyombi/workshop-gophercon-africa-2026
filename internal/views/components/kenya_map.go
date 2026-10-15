package components

import (
	"math"
	"strconv"

	"workshop/internal/data"
)

// KenyaMapProps configure the clickable county map; Current set means locator mode.
type KenyaMapProps struct {
	Store   *data.Store
	Current string
	Slugs   []string
}

func (p KenyaMapProps) locator() bool { return p.Current != "" }

func (p KenyaMapProps) county(slug string) data.County {
	if c, ok := p.Store.County(slug); ok {
		return *c
	}
	return data.County{Slug: slug, Name: slug}
}

func (p KenyaMapProps) regionName(slug string) string {
	if r, ok := p.Store.Region(slug); ok {
		return r.Name
	}
	return slug
}

func (p KenyaMapProps) dimmed(slug string) bool {
	if p.Slugs == nil {
		return false
	}
	for _, s := range p.Slugs {
		if s == slug {
			return false
		}
	}
	return true
}

func (p KenyaMapProps) current() (data.Shape, bool) {
	for _, s := range p.Store.Geo.Kenya.Counties {
		if s.Slug == p.Current {
			return s, true
		}
	}
	return data.Shape{}, false
}

func (p KenyaMapProps) countyLabel(c data.County) string {
	return CountyTitle(c) + ", " + p.regionName(c.Region) + " region. Seat: " + c.Seat
}

func (p KenyaMapProps) mapClass(c data.County) string {
	cls := "kmap__county region-fill--" + c.Region
	if c.Slug == p.Current {
		cls += " is-current"
	}
	if p.dimmed(c.Slug) {
		cls += " is-dim"
	}
	return cls
}

// CountyTitle is the county's name with "County" appended.
func CountyTitle(c data.County) string { return c.Name + " County" }

// kmapCallout is a leader-line label for a county too small to tap reliably.
type kmapCallout struct {
	X1, Y1, X2, Y2, TX, TY float64
	Anchor                 string
}

var kmapCallouts = map[string]kmapCallout{
	"mombasa": {575, 929, 650, 975, 656, 983, "start"},
	"nairobi": {296, 664, 232, 618, 226, 612, "end"},
}

func kmapCalloutFor(slug string) (kmapCallout, bool) { c, ok := kmapCallouts[slug]; return c, ok }

func kmapNum(v float64) string { return strconv.FormatFloat(math.Round(v*10)/10, 'f', -1, 64) }

// kmapRing is the radius of the locator ring, sized to the county's own box.
func kmapRing(s data.Shape) string {
	x0, y0, x1, y1, ok := countyPathBounds(s.D)
	if !ok {
		return "20"
	}
	r := max(x1-x0, y1-y0)/2 + 14
	return kmapNum(max(r, 26))
}
