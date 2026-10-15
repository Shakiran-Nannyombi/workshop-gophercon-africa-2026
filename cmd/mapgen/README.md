# mapgen

One-off generator for the static SVG map data in `data/geo/`. The site renders
these paths as inline SVG, so there's no runtime map library and no remote
requests. The output JSON is committed. Rerun this tool only when you want to
change the sources, the bounding boxes or the simplification.

```sh
go run ./cmd/mapgen          # from the repository root
```

It uses only the Go standard library. Downloads are cached in
`data/geo/.cache/`, which is git-ignored. Delete that directory to download
everything again. The run fails if any of the 47 county slugs is missing, or if
a source shape doesn't match a slug.

## Outputs

| File | Contents | viewBox width |
|------|----------|---------------|
| `kenya_counties.json` | 47 counties (`slug`, `name`, `d`, `cx`, `cy`), plus Lake Victoria and Lake Turkana in `lakes` | 800 |
| `nairobi_region.json` | nairobi, kiambu, kajiado and machakos, clipped to lon 36.55–37.20, lat -1.60 – -0.95 | 1000 |
| `east_africa.json` | countries (`code`, `name`, `d`, `cx`, `cy`) and lakes (`name`, `d`), clipped to lon 28.5–48.5, lat -12.0 – 15.5 | 800 |

`cx`/`cy` is a label point. It is the centroid of the largest ring. If that
centroid falls outside the shape, the midpoint of the widest interior
horizontal span through it is used instead. The geoBoundaries county polygons
extend over lake water, so draw `lakes` on top of the counties.

## Sources and licences

- **Kenya counties**: geoBoundaries gbOpen KEN ADM1 (`KEN-ADM1-32016919`, 2020),
  simplified geometry. The API reports `boundarySource` as "RCMRD GeoPortal,
  Africa GeoPortal" and `boundaryLicense` as "Public Domain".
  <https://www.geoboundaries.org/api/current/gbOpen/KEN/ADM1/>
- **Countries and lakes**: Natural Earth 1:50m Admin 0 Countries and 1:50m
  Lakes. Public domain. <https://www.naturalearthdata.com>. Somaliland is merged
  into Somalia (`SO`). The first nine countries (KE UG TZ RW BI ET SO SS CD) are
  the focus countries. SD ER DJ YE CF ZM MW MZ are included only as context, so
  the box has no gaps.

## Projection

Each file has `"projection": {lon0, lat0, lon1, lat1, width, height}`. The
bounding box maps exactly onto `0..width × 0..height` (no padding):

```
x = (lon - lon0) / (lon1 - lon0) * width
y = (max(lat0, lat1) - lat) / |lat1 - lat0| * height
```

`height = width * |lat1 - lat0| / ((lon1 - lon0) * cos(centreLat))`, rounded to
one decimal. That rounded value is the one used for projecting. This is
equirectangular projection scaled by the cosine of the centre latitude.
`internal/data/geo.go`'s `Projection.Project` implements the same formula.

## Processing

1. Project lon/lat into viewBox units. For the region and East Africa maps,
   clip each ring to the viewBox with Sutherland–Hodgman.
2. Simplify with topology awareness. Rings are split where the set of shapes
   sharing a vertex changes. Each shared arc is then run through
   Douglas–Peucker in a canonical direction, so neighbouring borders stay
   identical.
3. Drop tiny islets (Kenya: under 1.5 px²). A shape's largest ring is always kept.
4. Round coordinates to 1 decimal.

Tolerances: Kenya uses 0.5 px, and 0.15 px for counties smaller than 40 px
(Mombasa, Vihiga, Nairobi and so on). The Nairobi region uses 0.6 px and East
Africa uses 0.4 px.
