// Command mapgen is a one-off generator that turns open boundary datasets
// into small, static SVG path data committed under data/geo/. The site then
// renders hand-styled inline SVG maps with no runtime map library and no
// remote requests.
//
// Run from the repository root:
//
//	go run ./cmd/mapgen
//
// Downloads are cached in data/geo/.cache/ (git-ignored); delete that
// directory to force a fresh download.
//
// # Projection
//
// Every output file carries a "projection" object {lon0, lat0, lon1, lat1,
// width, height}. The bounding box maps exactly onto the viewBox
// 0..width × 0..height with no padding:
//
//	x = (lon - lon0) / (lon1 - lon0) * width
//	y = (max(lat0, lat1) - lat) / |lat1 - lat0| * height
//
// The aspect ratio is equirectangular, scaled by the cosine of the centre
// latitude: height = width * |lat1-lat0| / ((lon1-lon0) * cos(centreLat)),
// rounded to one decimal (the rounded value is the one used for projecting).
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	outDir   = "data/geo"
	cacheDir = "data/geo/.cache"

	gbAPI      = "https://www.geoboundaries.org/api/current/gbOpen/KEN/ADM1/"
	gbFallback = "https://github.com/wmgeolab/geoBoundaries/raw/main/releaseData/gbOpen/KEN/ADM1/geoBoundaries-KEN-ADM1_simplified.geojson"
	neCountry  = "https://raw.githubusercontent.com/nvkelso/natural-earth-vector/master/geojson/ne_50m_admin_0_countries.geojson"
	neLakes    = "https://raw.githubusercontent.com/nvkelso/natural-earth-vector/master/geojson/ne_50m_lakes.geojson"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("mapgen: ")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cacheDir, ".gitignore"), []byte("*\n"), 0o644); err != nil {
		log.Fatal(err)
	}

	ken, meta, err := loadKenya()
	if err != nil {
		log.Fatal(err)
	}
	counties, err := matchCounties(ken)
	if err != nil {
		log.Fatal(err)
	}

	if err := writeKenya(counties, meta); err != nil {
		log.Fatal(err)
	}
	if err := writeNairobi(counties, meta); err != nil {
		log.Fatal(err)
	}
	if err := writeEastAfrica(); err != nil {
		log.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// Fetching

var client = &http.Client{Timeout: 10 * time.Minute}

// fetch returns the body of url, reading from / writing to cacheDir/name.
func fetch(url, name string) ([]byte, error) {
	path := filepath.Join(cacheDir, name)
	if b, err := os.ReadFile(path); err == nil && len(b) > 0 {
		return b, nil
	}
	log.Printf("downloading %s", url)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}
	tmp := path + ".part"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return nil, err
	}
	return b, os.Rename(tmp, path)
}

type gbMeta struct {
	BoundaryID      string `json:"boundaryID"`
	BoundaryYear    string `json:"boundaryYearRepresented"`
	BoundarySource  string `json:"boundarySource"`
	BoundaryLicense string `json:"boundaryLicense"`
	LicenseSource   string `json:"licenseSource"`
	Simplified      string `json:"simplifiedGeometryGeoJSON"`
	GJ              string `json:"gjDownloadURL"`
}

func (m gbMeta) source() string {
	s := "geoBoundaries gbOpen KEN ADM1"
	if m.BoundaryID != "" {
		s += " (" + m.BoundaryID
		if m.BoundaryYear != "" {
			s += ", " + m.BoundaryYear
		}
		s += ")"
	}
	if m.BoundarySource != "" {
		s += "; original source: " + m.BoundarySource
	}
	return s + "; https://www.geoboundaries.org"
}

func loadKenya() ([]feature, gbMeta, error) {
	var meta gbMeta
	var body []byte
	apiBody, err := fetch(gbAPI, "gb_ken_adm1_api.json")
	if err == nil {
		err = json.Unmarshal(apiBody, &meta)
	}
	if err == nil {
		u := meta.Simplified
		if u == "" {
			u = meta.GJ
		}
		if u == "" {
			err = errors.New("geoBoundaries API returned no GeoJSON URL")
		} else {
			body, err = fetch(u, "gb_ken_adm1.geojson")
		}
	}
	if err != nil {
		log.Printf("geoBoundaries API path failed (%v); trying fallback", err)
		meta = gbMeta{BoundaryLicense: "see https://www.geoboundaries.org (API unavailable)"}
		body, err = fetch(gbFallback, "gb_ken_adm1.geojson")
		if err != nil {
			return nil, meta, err
		}
	}
	fs, err := parseFeatures(body, func(p map[string]any) string { return str(p["shapeName"]) })
	return fs, meta, err
}

func str(v any) string {
	s, _ := v.(string)
	return s
}

// ---------------------------------------------------------------------------
// GeoJSON

type pt struct{ X, Y float64 }
type ring []pt
type polygon []ring // first ring is the outer ring, the rest are holes

type feature struct {
	Name  string
	Props map[string]any
	Polys []polygon // lon/lat until projected
}

func parseFeatures(b []byte, name func(map[string]any) string) ([]feature, error) {
	var fc struct {
		Features []struct {
			Properties map[string]any `json:"properties"`
			Geometry   *struct {
				Type        string          `json:"type"`
				Coordinates json.RawMessage `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}
	if err := json.Unmarshal(b, &fc); err != nil {
		return nil, fmt.Errorf("parse GeoJSON: %w", err)
	}
	var out []feature
	for _, f := range fc.Features {
		if f.Geometry == nil {
			continue
		}
		var polys [][][][]float64
		switch f.Geometry.Type {
		case "Polygon":
			var p [][][]float64
			if err := json.Unmarshal(f.Geometry.Coordinates, &p); err != nil {
				return nil, err
			}
			polys = [][][][]float64{p}
		case "MultiPolygon":
			if err := json.Unmarshal(f.Geometry.Coordinates, &polys); err != nil {
				return nil, err
			}
		default:
			continue
		}
		ft := feature{Name: name(f.Properties), Props: f.Properties}
		for _, p := range polys {
			var poly polygon
			for _, r := range p {
				var rr ring
				for _, c := range r {
					if len(c) >= 2 {
						rr = append(rr, pt{c[0], c[1]})
					}
				}
				// Drop the closing duplicate; rings are implicitly closed.
				if len(rr) > 1 && rr[0] == rr[len(rr)-1] {
					rr = rr[:len(rr)-1]
				}
				if len(rr) >= 3 {
					poly = append(poly, rr)
				}
			}
			if len(poly) > 0 {
				ft.Polys = append(ft.Polys, poly)
			}
		}
		out = append(out, ft)
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Counties

var countySlugs = []string{
	"mombasa", "kwale", "kilifi", "tana-river", "lamu", "taita-taveta", "garissa",
	"wajir", "mandera", "marsabit", "isiolo", "meru", "tharaka-nithi", "embu",
	"kitui", "machakos", "makueni", "nyandarua", "nyeri", "kirinyaga", "muranga",
	"kiambu", "turkana", "west-pokot", "samburu", "trans-nzoia", "uasin-gishu",
	"elgeyo-marakwet", "nandi", "baringo", "laikipia", "nakuru", "narok",
	"kajiado", "kericho", "bomet", "kakamega", "vihiga", "bungoma", "busia",
	"siaya", "kisumu", "homa-bay", "migori", "kisii", "nyamira", "nairobi",
}

var displayNames = map[string]string{
	"tana-river": "Tana River", "taita-taveta": "Taita-Taveta",
	"tharaka-nithi": "Tharaka-Nithi", "muranga": "Murang'a",
	"west-pokot": "West Pokot", "trans-nzoia": "Trans-Nzoia",
	"uasin-gishu": "Uasin Gishu", "elgeyo-marakwet": "Elgeyo-Marakwet",
	"homa-bay": "Homa Bay",
}

// slugAliases maps slugified source names to our canonical slugs.
var slugAliases = map[string]string{
	"nairobi-city":     "nairobi",
	"nairobi-county":   "nairobi",
	"tharaka":          "tharaka-nithi",
	"tharaka-nithi":    "tharaka-nithi",
	"keiyo-marakwet":   "elgeyo-marakwet",
	"elgeyo":           "elgeyo-marakwet",
	"elegeyo-marakwet": "elgeyo-marakwet",
	"murang-a":         "muranga",
	"taita":            "taita-taveta",
	"transnzoia":       "trans-nzoia",
	"homabay":          "homa-bay",
	"tanariver":        "tana-river",
	"westpokot":        "west-pokot",
	"uasingishu":       "uasin-gishu",
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.TrimSuffix(s, " county")
	var b strings.Builder
	dash := false
	for _, r := range s {
		switch {
		case r == '\'' || r == '’' || r == '`':
			// strip apostrophes
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
			dash = false
		default:
			if !dash && b.Len() > 0 {
				b.WriteByte('-')
				dash = true
			}
		}
	}
	out := strings.TrimSuffix(b.String(), "-")
	if a, ok := slugAliases[out]; ok {
		return a
	}
	return out
}

func displayName(slug string) string {
	if n, ok := displayNames[slug]; ok {
		return n
	}
	return strings.ToUpper(slug[:1]) + slug[1:]
}

func matchCounties(fs []feature) (map[string]feature, error) {
	want := map[string]bool{}
	for _, s := range countySlugs {
		want[s] = true
	}
	got := map[string]feature{}
	var problems []string
	for _, f := range fs {
		s := slugify(f.Name)
		if !want[s] {
			problems = append(problems, fmt.Sprintf("unmatched shape %q (slug %q)", f.Name, s))
			continue
		}
		if prev, dup := got[s]; dup {
			// Some datasets split a county into several features; merge.
			prev.Polys = append(prev.Polys, f.Polys...)
			got[s] = prev
			continue
		}
		f.Name = s
		got[s] = f
	}
	for _, s := range countySlugs {
		if _, ok := got[s]; !ok {
			problems = append(problems, "missing county "+s)
		}
	}
	if len(problems) > 0 {
		return nil, fmt.Errorf("county matching failed:\n  %s", strings.Join(problems, "\n  "))
	}
	return got, nil
}

// ---------------------------------------------------------------------------
// Projection

type projection struct {
	Lon0   float64 `json:"lon0"`
	Lat0   float64 `json:"lat0"`
	Lon1   float64 `json:"lon1"`
	Lat1   float64 `json:"lat1"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func newProjection(lon0, lat0, lon1, lat1, width float64) projection {
	latSpan := math.Abs(lat1 - lat0)
	k := math.Cos((lat0 + lat1) / 2 * math.Pi / 180)
	h := width * latSpan / ((lon1 - lon0) * k)
	return projection{lon0, lat0, lon1, lat1, width, math.Round(h*10) / 10}
}

func (p projection) project(lon, lat float64) pt {
	maxLat := math.Max(p.Lat0, p.Lat1)
	return pt{
		(lon - p.Lon0) / (p.Lon1 - p.Lon0) * p.Width,
		(maxLat - lat) / math.Abs(p.Lat1-p.Lat0) * p.Height,
	}
}

func (p projection) viewBox() string {
	return "0 0 " + ff(p.Width) + " " + ff(p.Height)
}

func projectFeatures(fs []feature, p projection, clip bool) []feature {
	out := make([]feature, 0, len(fs))
	for _, f := range fs {
		nf := feature{Name: f.Name, Props: f.Props}
		for _, poly := range f.Polys {
			var np polygon
			for _, r := range poly {
				rr := make(ring, len(r))
				for i, c := range r {
					rr[i] = p.project(c.X, c.Y)
				}
				if clip {
					rr = clipRect(rr, 0, 0, p.Width, p.Height)
				}
				if len(rr) >= 3 {
					np = append(np, rr)
				} else if len(np) == 0 {
					break // outer ring clipped away: drop polygon
				}
			}
			if len(np) > 0 {
				nf.Polys = append(nf.Polys, np)
			}
		}
		if len(nf.Polys) > 0 {
			out = append(out, nf)
		}
	}
	return out
}

// clipRect clips a ring to an axis-aligned rectangle (Sutherland–Hodgman).
func clipRect(r ring, x0, y0, x1, y1 float64) ring {
	type edge struct {
		inside func(pt) bool
		cross  func(a, b pt) pt
	}
	lerpX := func(a, b pt, x float64) pt { return pt{x, a.Y + (b.Y-a.Y)*(x-a.X)/(b.X-a.X)} }
	lerpY := func(a, b pt, y float64) pt { return pt{a.X + (b.X-a.X)*(y-a.Y)/(b.Y-a.Y), y} }
	edges := []edge{
		{func(p pt) bool { return p.X >= x0 }, func(a, b pt) pt { return lerpX(a, b, x0) }},
		{func(p pt) bool { return p.X <= x1 }, func(a, b pt) pt { return lerpX(a, b, x1) }},
		{func(p pt) bool { return p.Y >= y0 }, func(a, b pt) pt { return lerpY(a, b, y0) }},
		{func(p pt) bool { return p.Y <= y1 }, func(a, b pt) pt { return lerpY(a, b, y1) }},
	}
	out := r
	for _, e := range edges {
		if len(out) == 0 {
			break
		}
		in := out
		out = nil
		prev := in[len(in)-1]
		for _, cur := range in {
			switch ci, pi := e.inside(cur), e.inside(prev); {
			case ci && pi:
				out = append(out, cur)
			case ci && !pi:
				out = append(out, e.cross(prev, cur), cur)
			case !ci && pi:
				out = append(out, e.cross(prev, cur))
			}
			prev = cur
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Topology-aware simplification
//
// Neighbouring shapes share border vertices. To keep shared borders identical
// after simplification (no slivers or gaps), every ring is split at
// "junction" vertices - where the set of rings using a vertex changes - and
// each arc between junctions is simplified with Douglas–Peucker in a
// canonical direction, so both neighbours get exactly the same result.

type key [2]int64

func keyOf(p pt) key { return key{int64(math.Round(p.X * 1e4)), int64(math.Round(p.Y * 1e4))} }

type simplifyOpts struct {
	tol     func(featureIdx int) float64 // DP tolerance in viewBox units
	minArea float64                      // drop rings smaller than this (px²), except a feature's largest
}

func simplifyAll(fs []feature, o simplifyOpts) []feature {
	// owners[k] = set of feature indices touching vertex k (bitset over features).
	owners := map[key]map[int]bool{}
	for fi, f := range fs {
		for _, poly := range f.Polys {
			for _, r := range poly {
				for _, p := range r {
					k := keyOf(p)
					if owners[k] == nil {
						owners[k] = map[int]bool{}
					}
					owners[k][fi] = true
				}
			}
		}
	}
	sig := map[key]string{}
	sigOf := func(k key) string {
		if s, ok := sig[k]; ok {
			return s
		}
		ids := make([]int, 0, len(owners[k]))
		for id := range owners[k] {
			ids = append(ids, id)
		}
		sort.Ints(ids)
		s := fmt.Sprint(ids)
		sig[k] = s
		return s
	}

	out := make([]feature, len(fs))
	for fi, f := range fs {
		nf := feature{Name: f.Name, Props: f.Props}
		// Find the largest outer ring so we never drop it.
		bestP, bestA := -1, -1.0
		for pi, poly := range f.Polys {
			if a := math.Abs(area(poly[0])); a > bestA {
				bestP, bestA = pi, a
			}
		}
		for pi, poly := range f.Polys {
			var np polygon
			for ri, r := range poly {
				keep := pi == bestP && ri == 0
				if !keep && math.Abs(area(r)) < o.minArea {
					if ri == 0 {
						break
					}
					continue
				}
				sr := simplifyRing(r, owners, sigOf, o.tol)
				if len(sr) < 3 || (!keep && math.Abs(area(sr)) < o.minArea) {
					if keep {
						sr = r // never lose the main shape
					} else if ri == 0 {
						break
					} else {
						continue
					}
				}
				np = append(np, sr)
			}
			if len(np) > 0 {
				nf.Polys = append(nf.Polys, np)
			}
		}
		out[fi] = nf
	}
	return out
}

func simplifyRing(r ring, owners map[key]map[int]bool, sigOf func(key) string, tol func(int) float64) ring {
	n := len(r)
	keys := make([]key, n)
	for i, p := range r {
		keys[i] = keyOf(p)
	}
	fixed := make([]bool, n)
	nfixed := 0
	for i := range r {
		s := sigOf(keys[i])
		if s != sigOf(keys[(i+n-1)%n]) || s != sigOf(keys[(i+1)%n]) {
			fixed[i] = true
			nfixed++
		}
	}
	if nfixed < 2 {
		// Free-standing ring (island, or entirely inside one shared border):
		// anchor at the lexicographically smallest vertex and the farthest one.
		a := 0
		for i := range keys {
			if keys[i][0] < keys[a][0] || (keys[i][0] == keys[a][0] && keys[i][1] < keys[a][1]) {
				a = i
			}
		}
		b, bd := a, -1.0
		for i, p := range r {
			if d := dist2(p, r[a]); d > bd {
				b, bd = i, d
			}
		}
		fixed[a], fixed[b] = true, true
	}
	start := 0
	for !fixed[start] {
		start++
	}
	var out ring
	i := start
	for {
		// Collect the arc from fixed vertex i to the next fixed vertex j.
		arc := ring{r[i]}
		j := (i + 1) % n
		for !fixed[j] {
			arc = append(arc, r[j])
			j = (j + 1) % n
		}
		arc = append(arc, r[j])

		// Tolerance: the smallest of all features owning this arc.
		t := math.Inf(1)
		for _, p := range arc[1 : len(arc)-1] {
			for id := range owners[keyOf(p)] {
				t = math.Min(t, tol(id))
			}
		}
		if math.IsInf(t, 1) {
			t = 0
		}
		// Canonical direction so neighbours simplify identically.
		ka, kb := keyOf(arc[0]), keyOf(arc[len(arc)-1])
		rev := kb[0] < ka[0] || (kb[0] == ka[0] && kb[1] < ka[1])
		if rev {
			reverse(arc)
		}
		s := douglasPeucker(arc, t)
		if rev {
			reverse(s)
		}
		out = append(out, s[:len(s)-1]...)
		i = j
		if i == start {
			break
		}
	}
	return out
}

func reverse(r ring) {
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
}

func douglasPeucker(pts ring, tol float64) ring {
	if len(pts) <= 2 || tol <= 0 {
		return append(ring(nil), pts...)
	}
	keep := make([]bool, len(pts))
	keep[0], keep[len(pts)-1] = true, true
	var rec func(a, b int)
	rec = func(a, b int) {
		if b <= a+1 {
			return
		}
		idx, dmax := -1, -1.0
		for i := a + 1; i < b; i++ {
			if d := segDist(pts[i], pts[a], pts[b]); d > dmax {
				idx, dmax = i, d
			}
		}
		if dmax > tol {
			keep[idx] = true
			rec(a, idx)
			rec(idx, b)
		}
	}
	rec(0, len(pts)-1)
	var out ring
	for i, k := range keep {
		if k {
			out = append(out, pts[i])
		}
	}
	return out
}

func segDist(p, a, b pt) float64 {
	dx, dy := b.X-a.X, b.Y-a.Y
	l2 := dx*dx + dy*dy
	if l2 == 0 {
		return math.Sqrt(dist2(p, a))
	}
	t := math.Max(0, math.Min(1, ((p.X-a.X)*dx+(p.Y-a.Y)*dy)/l2))
	return math.Sqrt(dist2(p, pt{a.X + t*dx, a.Y + t*dy}))
}

func dist2(a, b pt) float64 { return (a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y) }

func area(r ring) float64 {
	s := 0.0
	for i := range r {
		j := (i + 1) % len(r)
		s += r[i].X*r[j].Y - r[j].X*r[i].Y
	}
	return s / 2
}

// ---------------------------------------------------------------------------
// Output helpers

// ff formats a float with at most one decimal.
func ff(v float64) string {
	v = math.Round(v*10) / 10
	if v == 0 {
		v = 0 // normalise -0
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

type f1 float64

func (v f1) MarshalJSON() ([]byte, error) { return []byte(ff(float64(v))), nil }

func pathData(f feature) string {
	var b strings.Builder
	for _, poly := range f.Polys {
		for _, r := range poly {
			var last string
			first := true
			for _, p := range r {
				c := ff(p.X) + "," + ff(p.Y)
				if c == last {
					continue
				}
				if first {
					b.WriteString("M")
					first = false
				} else {
					b.WriteString(" ")
				}
				b.WriteString(c)
				last = c
			}
			b.WriteString("Z")
		}
	}
	return b.String()
}

// labelPoint returns the centroid of the feature's largest ring, or - if
// that falls outside the shape - the midpoint of the widest horizontal
// interior span through the centroid's y.
func labelPoint(f feature) pt {
	var big ring
	bestA := -1.0
	for _, poly := range f.Polys {
		if a := math.Abs(area(poly[0])); a > bestA {
			big, bestA = poly[0], a
		}
	}
	if big == nil {
		return pt{}
	}
	c := centroid(big)
	if insideFeature(f, c) {
		return c
	}
	try := []float64{c.Y}
	for d := 1.0; d <= 20; d++ {
		try = append(try, c.Y+d, c.Y-d)
	}
	for _, y := range try {
		if p, ok := widestSpan(f, y); ok {
			return p
		}
	}
	return c
}

func centroid(r ring) pt {
	var cx, cy, a float64
	for i := range r {
		j := (i + 1) % len(r)
		cr := r[i].X*r[j].Y - r[j].X*r[i].Y
		a += cr
		cx += (r[i].X + r[j].X) * cr
		cy += (r[i].Y + r[j].Y) * cr
	}
	if a == 0 {
		return r[0]
	}
	return pt{cx / (3 * a), cy / (3 * a)}
}

func insideFeature(f feature, p pt) bool {
	in := false
	for _, poly := range f.Polys {
		for _, r := range poly {
			for i, j := 0, len(r)-1; i < len(r); j, i = i, i+1 {
				if (r[i].Y > p.Y) != (r[j].Y > p.Y) &&
					p.X < (r[j].X-r[i].X)*(p.Y-r[i].Y)/(r[j].Y-r[i].Y)+r[i].X {
					in = !in
				}
			}
		}
	}
	return in
}

func widestSpan(f feature, y float64) (pt, bool) {
	var xs []float64
	for _, poly := range f.Polys {
		for _, r := range poly {
			for i, j := 0, len(r)-1; i < len(r); j, i = i, i+1 {
				if (r[i].Y > y) != (r[j].Y > y) {
					xs = append(xs, (r[j].X-r[i].X)*(y-r[i].Y)/(r[j].Y-r[i].Y)+r[i].X)
				}
			}
		}
	}
	sort.Float64s(xs)
	best, bw := pt{}, 0.0
	for i := 0; i+1 < len(xs); i += 2 {
		if w := xs[i+1] - xs[i]; w > bw {
			best, bw = pt{(xs[i] + xs[i+1]) / 2, y}, w
		}
	}
	return best, bw > 0
}

func bbox(fs []feature) (minX, minY, maxX, maxY float64) {
	minX, minY = math.Inf(1), math.Inf(1)
	maxX, maxY = math.Inf(-1), math.Inf(-1)
	for _, f := range fs {
		for _, poly := range f.Polys {
			for _, p := range poly[0] {
				minX, maxX = math.Min(minX, p.X), math.Max(maxX, p.X)
				minY, maxY = math.Min(minY, p.Y), math.Max(maxY, p.Y)
			}
		}
	}
	return
}

type shapeOut struct {
	Slug string `json:"slug,omitempty"`
	Code string `json:"code,omitempty"`
	Name string `json:"name"`
	D    string `json:"d"`
	CX   f1     `json:"cx"`
	CY   f1     `json:"cy"`
}

type countyFile struct {
	ViewBox    string     `json:"viewBox"`
	Source     string     `json:"source"`
	License    string     `json:"license"`
	Projection projection `json:"projection"`
	Counties   []shapeOut `json:"counties"`
	Lakes      []shapeOut `json:"lakes,omitempty"`
}

type eaFile struct {
	ViewBox    string     `json:"viewBox"`
	Source     string     `json:"source"`
	License    string     `json:"license"`
	Projection projection `json:"projection"`
	Countries  []shapeOut `json:"countries"`
	Lakes      []shapeOut `json:"lakes"`
}

func writeJSON(name string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	path := filepath.Join(outDir, name)
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return err
	}
	log.Printf("wrote %s (%d bytes)", path, len(b))
	return nil
}

// ---------------------------------------------------------------------------
// Outputs

func orderedCounties(m map[string]feature, slugs []string) []feature {
	fs := make([]feature, 0, len(slugs))
	for _, s := range slugs {
		fs = append(fs, m[s])
	}
	return fs
}

func countyShapes(fs []feature) []shapeOut {
	out := make([]shapeOut, 0, len(fs))
	for _, f := range fs {
		c := labelPoint(f)
		out = append(out, shapeOut{Slug: f.Name, Name: displayName(f.Name), D: pathData(f), CX: f1(c.X), CY: f1(c.Y)})
	}
	return out
}

func writeKenya(m map[string]feature, meta gbMeta) error {
	fs := orderedCounties(m, countySlugs)
	minLon, minLat, maxLon, maxLat := bbox(fs)
	const pad = 0.1 // degrees of breathing room around Kenya
	r2 := func(v float64) float64 { return math.Round(v*100) / 100 }
	p := newProjection(r2(minLon-pad), r2(minLat-pad), r2(maxLon+pad), r2(maxLat+pad), 800)

	pf := projectFeatures(fs, p, false)
	// Small counties (Mombasa, Vihiga, Nairobi, ...) get a finer tolerance.
	tols := make([]float64, len(pf))
	for i, f := range pf {
		x0, y0, x1, y1 := bbox([]feature{f})
		if math.Max(x1-x0, y1-y0) < 40 {
			tols[i] = 0.15
		} else {
			tols[i] = 0.5
		}
	}
	sf := simplifyAll(pf, simplifyOpts{tol: func(i int) float64 { return tols[i] }, minArea: 1.5})
	// geoBoundaries county polygons extend over lake water (Victoria,
	// Turkana), so lakes are emitted separately to be drawn on top.
	src := meta.source()
	lakes, err := lakeShapes(p, 0.4)
	if err != nil {
		log.Printf("kenya lakes skipped: %v", err)
	} else {
		src += "; lakes: Natural Earth 1:50m lakes (public domain)"
	}
	return writeJSON("kenya_counties.json", countyFile{
		ViewBox: p.viewBox(), Source: src, License: meta.BoundaryLicense,
		Projection: p, Counties: countyShapes(sf), Lakes: lakes,
	})
}

func writeNairobi(m map[string]feature, meta gbMeta) error {
	slugs := []string{"nairobi", "kiambu", "kajiado", "machakos"}
	fs := orderedCounties(m, slugs)
	p := newProjection(36.55, -1.60, 37.20, -0.95, 1000)
	pf := projectFeatures(fs, p, true)
	if len(pf) != len(slugs) {
		return fmt.Errorf("nairobi region: expected %d counties inside bbox, got %d", len(slugs), len(pf))
	}
	sf := simplifyAll(pf, simplifyOpts{tol: func(int) float64 { return 0.6 }, minArea: 4})
	return writeJSON("nairobi_region.json", countyFile{
		ViewBox: p.viewBox(), Source: meta.source(), License: meta.BoundaryLicense,
		Projection: p, Counties: countyShapes(sf),
	})
}

// eaCountries lists Natural Earth ADM0_A3 codes to keep, with output codes.
// The first nine are the focus countries; the rest are drawn for context so
// the map has no holes inside the bbox. Somaliland (SOL) is merged into SO.
var eaCountries = []struct{ a3, code, name string }{
	{"KEN", "KE", "Kenya"}, {"UGA", "UG", "Uganda"}, {"TZA", "TZ", "Tanzania"},
	{"RWA", "RW", "Rwanda"}, {"BDI", "BI", "Burundi"}, {"ETH", "ET", "Ethiopia"},
	{"SOM", "SO", "Somalia"}, {"SDS", "SS", "South Sudan"},
	{"COD", "CD", "DR Congo"},
	{"SDN", "SD", "Sudan"}, {"ERI", "ER", "Eritrea"}, {"DJI", "DJ", "Djibouti"},
	{"YEM", "YE", "Yemen"}, {"CAF", "CF", "Central African Republic"},
	{"ZMB", "ZM", "Zambia"}, {"MWI", "MW", "Malawi"}, {"MOZ", "MZ", "Mozambique"},
}

var eaLakes = []string{
	"Lake Victoria", "Lake Turkana", "Lake Tanganyika", "Lake Kivu", "Lake Albert",
	"Lake Edward", "Lake Kyoga", "Lake Tana", "Lake Malawi", "Lake Rukwa",
}

func writeEastAfrica() error {
	cb, err := fetch(neCountry, "ne_50m_admin_0_countries.geojson")
	if err != nil {
		return err
	}
	all, err := parseFeatures(cb, func(p map[string]any) string { return str(p["ADM0_A3"]) })
	if err != nil {
		return err
	}
	byA3 := map[string]*feature{}
	for i := range all {
		byA3[all[i].Name] = &all[i]
	}
	if sol := byA3["SOL"]; sol != nil && byA3["SOM"] != nil {
		byA3["SOM"].Polys = append(byA3["SOM"].Polys, sol.Polys...)
	}
	var fs []feature
	for _, c := range eaCountries {
		f := byA3[c.a3]
		if f == nil {
			return fmt.Errorf("natural earth: country %s not found", c.a3)
		}
		fs = append(fs, feature{Name: c.code, Props: map[string]any{"name": c.name}, Polys: f.Polys})
	}

	p := newProjection(28.5, -12.0, 48.5, 15.5, 800)
	pf := projectFeatures(fs, p, true)
	sf := simplifyAll(pf, simplifyOpts{tol: func(int) float64 { return 0.4 }, minArea: 1})
	var countries []shapeOut
	for _, f := range sf {
		c := labelPoint(f)
		countries = append(countries, shapeOut{Code: f.Name, Name: str(f.Props["name"]), D: pathData(f), CX: f1(c.X), CY: f1(c.Y)})
	}
	for _, c := range eaCountries[:9] {
		found := false
		for _, s := range countries {
			found = found || s.Code == c.code
		}
		if !found {
			return fmt.Errorf("east africa: %s missing after clipping", c.code)
		}
	}

	lakes, err := lakeShapes(p, 0.4)
	lakeNote := "; lakes: Natural Earth 1:50m lakes"
	if err != nil {
		log.Printf("lakes skipped: %v", err)
		lakeNote = ""
	}

	return writeJSON("east_africa.json", eaFile{
		ViewBox:    p.viewBox(),
		Source:     "Natural Earth 1:50m Admin 0 Countries (naturalearthdata.com)" + lakeNote,
		License:    "Public Domain",
		Projection: p,
		Countries:  countries,
		Lakes:      lakes,
	})
}

// lakeShapes projects the named Natural Earth lakes into p, clipped to the
// viewBox. Lakes entirely outside the box are omitted.
func lakeShapes(p projection, tol float64) ([]shapeOut, error) {
	lb, err := fetch(neLakes, "ne_50m_lakes.geojson")
	if err != nil {
		return nil, err
	}
	lf, err := parseFeatures(lb, func(p map[string]any) string { return str(p["name"]) })
	if err != nil {
		return nil, err
	}
	var keep []feature
	for _, name := range eaLakes {
		for _, f := range lf {
			if strings.EqualFold(f.Name, name) {
				f.Name = name
				keep = append(keep, f)
			}
		}
	}
	pl := projectFeatures(keep, p, true)
	sl := simplifyAll(pl, simplifyOpts{tol: func(int) float64 { return tol }, minArea: 1})
	var lakes []shapeOut
	for _, f := range sl {
		c := labelPoint(f)
		lakes = append(lakes, shapeOut{Name: f.Name, D: pathData(f), CX: f1(c.X), CY: f1(c.Y)})
	}
	return lakes, nil
}
