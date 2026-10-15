package components

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"workshop/internal/data"
)

// countyPathBounds reads the bounding box from an absolute "M x,y x,y ... Z" path.
func countyPathBounds(d string) (x0, y0, x1, y1 float64, ok bool) {
	fields := strings.FieldsFunc(d, func(r rune) bool {
		return r == ' ' || r == ',' || r == '\n' || r == '\t' || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
	})
	x0, y0, x1, y1 = math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for i := 0; i+1 < len(fields); i += 2 {
		x, err1 := strconv.ParseFloat(fields[i], 64)
		y, err2 := strconv.ParseFloat(fields[i+1], 64)
		if err1 != nil || err2 != nil {
			continue
		}
		x0, y0, x1, y1 = math.Min(x0, x), math.Min(y0, y), math.Max(x1, x), math.Max(y1, y)
		ok = true
	}
	return
}

// ShapeViewBox is a tight, padded viewBox around one county's outline.
func ShapeViewBox(d string) string {
	x0, y0, x1, y1, ok := countyPathBounds(d)
	if !ok {
		return "0 0 100 100"
	}
	w, h := x1-x0, y1-y0
	pad := math.Max(w, h) * 0.06
	return fmt.Sprintf("%.1f %.1f %.1f %.1f", x0-pad, y0-pad, w+2*pad, h+2*pad)
}

// CountyCode pads an official county code to two digits.
func CountyCode(n int) string { return fmt.Sprintf("%02d", n) }

var goodForLabels = map[string]string{"remote-travel": "Remote travel", "tech": "Tech"}

// GoodForLabel turns a goodFor tag like "remote-travel" into display text.
func GoodForLabel(t string) string {
	if l, ok := goodForLabels[t]; ok {
		return l
	}
	t = strings.ReplaceAll(t, "-", " ")
	if t == "" {
		return t
	}
	return strings.ToUpper(t[:1]) + t[1:]
}

func firstN(xs []string, n int) []string { return xs[:min(len(xs), n)] }

// ShapeFor finds a county's outline in a shape set.
func ShapeFor(set data.ShapeSet, slug string) (data.Shape, bool) {
	for _, s := range set.Counties {
		if s.Slug == slug {
			return s, true
		}
	}
	return data.Shape{}, false
}
