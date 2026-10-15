package data

// Nairobi is the long-form content of the Nairobi deep-dive page (data/nairobi.json).
type Nairobi struct {
	CityVsCounty []string     `json:"cityVsCounty"` // paragraphs
	MetroNote    string       `json:"metroNote"`
	Ruaka        Spotlight    `json:"ruaka"`
	Transport    []GuideBlock `json:"transport"`
	Matatu       MatatuGuide  `json:"matatu"`
	Traffic      GuideBlock   `json:"traffic"`
	SIM          GuideBlock   `json:"sim"`
	Food         GuideBlock   `json:"food"`
	Safety       GuideBlock   `json:"safety"`
	Etiquette    GuideBlock   `json:"etiquette"`
	Weather      Weather      `json:"weather"`
	Roads        []Road       `json:"roads"`
	Verify       Verify       `json:"verify,omitempty"`
}

type Spotlight struct {
	Title    string       `json:"title"`
	Intro    string       `json:"intro"`
	Sections []GuideBlock `json:"sections"`
	Verify   Verify       `json:"verify,omitempty"`
}

// GuideBlock is a titled chunk of practical advice.
type GuideBlock struct {
	Slug   string   `json:"slug"`
	Title  string   `json:"title"`
	Icon   string   `json:"icon"` // matatu, boda, ride, bus, train, walk, sim, food, shield, people, clock
	Body   string   `json:"body"`
	Points []string `json:"points"`
	Verify Verify   `json:"verify,omitempty"`
}

type MatatuGuide struct {
	Intro   string        `json:"intro"`
	Colours []ColourNote  `json:"colours"`
	Routes  []RouteNumber `json:"routes"`
	Tips    []string      `json:"tips"`
	Verify  Verify        `json:"verify,omitempty"`
}

type ColourNote struct {
	Swatch  string `json:"swatch"` // hex
	Name    string `json:"name"`
	Meaning string `json:"meaning"`
}

type RouteNumber struct {
	Number string `json:"number"`
	To     string `json:"to"`
	From   string `json:"from"`
	Verify Verify `json:"verify,omitempty"`
}

type Weather struct {
	Summary   string       `json:"summary"`
	Elevation string       `json:"elevation"`
	Months    []Month      `json:"months"` // 12, Jan..Dec
	Seasons   []GuideBlock `json:"seasons"`
	Source    string       `json:"source"`
	Verify    Verify       `json:"verify,omitempty"`
}

type Month struct {
	Name   string  `json:"name"` // "Jan"
	HighC  float64 `json:"highC"`
	LowC   float64 `json:"lowC"`
	RainMM float64 `json:"rainMM"`
}

// Road is a schematic road line through approximate lat/lon waypoints.
type Road struct {
	Name   string       `json:"name"`
	Kind   string       `json:"kind"`   // highway, bypass, road, expressway
	Points [][2]float64 `json:"points"` // [lat, lon]
	Verify Verify       `json:"verify,omitempty"`
}
