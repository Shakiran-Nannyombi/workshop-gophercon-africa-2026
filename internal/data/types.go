// Package data loads Karibu Guide's content from the JSON files in /data.
package data

// Verify marks a value that still needs checking; collected into VERIFY.md.
type Verify string

// Site holds site-wide settings.
type Site struct {
	Name         string `json:"name"`
	Tagline      string `json:"tagline"`
	BaseURL      string `json:"baseURL"`
	LastReviewed string `json:"lastReviewed"` // YYYY-MM-DD
	Rates        Rates  `json:"rates"`
}

// Rates are units of each currency per 1 USD, used by the converter.
type Rates struct {
	AsOf       string             `json:"asOf"`
	Source     string             `json:"source"`
	Disclaimer string             `json:"disclaimer"`
	PerUSD     map[string]float64 `json:"perUSD"` // "KES": 129.2
	Verify     Verify             `json:"verify,omitempty"`
}

// Money is the content of the money page.
type Money struct {
	Notes          []Denomination `json:"notes"`
	Coins          []Denomination `json:"coins"`
	PriceFormats   []PriceFormat  `json:"priceFormats"`
	MobileMoney    []MobileMoney  `json:"mobileMoney"`
	MpesaSteps     []MpesaFlow    `json:"mpesaSteps"`
	MpesaFees      []FeeBand      `json:"mpesaFees"`
	MpesaFeesNote  string         `json:"mpesaFeesNote"`
	WhyMpesa       []string       `json:"whyMpesa"`
	Prices         []PriceExample `json:"prices"`
	Scams          []Tip          `json:"scams"`
	Bargaining     []string       `json:"bargaining"`
	BigNumbersNote string         `json:"bigNumbersNote"`
}

// Denomination is one banknote or coin.
type Denomination struct {
	Value   int    `json:"value"`
	Colour  string `json:"colour"` // CSS colour for the illustration
	Feature string `json:"feature"`
	Verify  Verify `json:"verify,omitempty"`
}

type PriceFormat struct {
	Example string `json:"example"`
	Meaning string `json:"meaning"`
}

type MobileMoney struct {
	Name     string   `json:"name"`
	Provider string   `json:"provider"`
	Country  string   `json:"country"` // ISO code
	Summary  string   `json:"summary"`
	Verify   Verify   `json:"verify,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Expect   []string `json:"expect,omitempty"` // "what to expect" notes for a first-timer
}

// MpesaFlow explains one way of paying with M-Pesa.
type MpesaFlow struct {
	Slug    string   `json:"slug"`
	Title   string   `json:"title"`
	When    string   `json:"when"`
	Steps   []string `json:"steps"`
	Example string   `json:"example"`
}

type FeeBand struct {
	Range string `json:"range"`
	Fee   string `json:"fee"`
}

// PriceExample is an approximate price in KES.
type PriceExample struct {
	Item     string `json:"item"`
	Category string `json:"category"`
	MinKES   int    `json:"minKES"`
	MaxKES   int    `json:"maxKES"`
	Note     string `json:"note"`
	Verify   Verify `json:"verify,omitempty"`
}

type Tip struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Country is an East African country with its money.
type Country struct {
	Code         string   `json:"code"` // KE, UG...
	Name         string   `json:"name"`
	Flag         []string `json:"flag"` // stripe colours for the mini flag
	Currency     string   `json:"currency"`
	CurrencyCode string   `json:"currencyCode"`
	Symbol       string   `json:"symbol"`
	Languages    []string `json:"languages"`
	MobileMoney  []string `json:"mobileMoney"`
	Scale        string   `json:"scale"` // "a coffee is ~KSh 250"
	CoffeeLocal  int      `json:"coffeeLocal"`
	Verify       Verify   `json:"verify,omitempty"`
	CountryProfile
}

// Region is one of Kenya's eight former provinces, still used informally.
type Region struct {
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Blurb string `json:"blurb"`
}

// County is one of Kenya's 47 counties.
type County struct {
	Code       int      `json:"code"` // official 1-47
	Slug       string   `json:"slug"`
	Name       string   `json:"name"`
	Region     string   `json:"region"` // Region.Slug
	Seat       string   `json:"seat"`
	AreaKm2    int      `json:"areaKm2"`
	Pop2019    int      `json:"pop2019"` // 2019 census
	KnownFor   []string `json:"knownFor"`
	Climate    string   `json:"climate"`
	Food       string   `json:"food"`
	Transport  string   `json:"transport"`
	Vibe       string   `json:"vibe"`
	Towns      []string `json:"towns"`
	GoodFor    []string `json:"goodFor"`
	Confidence string   `json:"confidence"` // "good" | "limited"
	Verify     Verify   `json:"verify,omitempty"`
}

// Area is a Nairobi neighbourhood, or a "Nairobi" place in a neighbouring county.
type Area struct {
	Slug       string   `json:"slug"`
	Name       string   `json:"name"`
	County     string   `json:"county"` // County.Slug
	Zone       string   `json:"zone"`   // "CBD", "Westlands side", "Eastlands"...
	Character  string   `json:"character"`
	Cost       string   `json:"cost"` // "budget" | "mid" | "upmarket" | "mixed"
	SafetyNote string   `json:"safetyNote"`
	Suits      []string `json:"suits"` // visitors, expats, families, students
	Landmarks  []string `json:"landmarks"`
	Transport  string   `json:"transport"`
	AlsoCalled []string `json:"alsoCalled,omitempty"`
	Lat        float64  `json:"lat"` // approximate centre, for the schematic map
	Lon        float64  `json:"lon"`
	Confidence string   `json:"confidence"`
	Verify     Verify   `json:"verify,omitempty"`
}

// City is a featured East African city; every city shares one page layout.
type City struct {
	Slug           string          `json:"slug"`
	Name           string          `json:"name"`
	Country        string          `json:"country"` // Country.Code
	Tagline        string          `json:"tagline"`
	Intro          string          `json:"intro"`
	Population     string          `json:"population"` // "~4.4 million (2019 census)"
	Elevation      string          `json:"elevation"`
	Climate        string          `json:"climate"`
	Airport        string          `json:"airport"`
	Languages      []string        `json:"languages"`
	Lat            float64         `json:"lat"`
	Lon            float64         `json:"lon"`
	Motif          string          `json:"motif"` // illustration key: skyline, coast, lake, mountain...
	Neighbourhoods []Neighbourhood `json:"neighbourhoods"`
	MoneyTips      []string        `json:"moneyTips"`
	Transport      []string        `json:"transport"`
	Food           []string        `json:"food"`
	Safety         []string        `json:"safety"`
	Suits          []string        `json:"suits"`
	Scores         Scores          `json:"scores"`
	Confidence     string          `json:"confidence"`
	Verify         Verify          `json:"verify,omitempty"`
}

type Neighbourhood struct {
	Name string `json:"name"`
	Note string `json:"note"`
}

// Scores are rough 1-5 ratings for the compare tool, each with a reason.
type Scores struct {
	Cost      Score  `json:"cost"`      // 5 = most expensive
	Transport Score  `json:"transport"` // 5 = easiest to get around
	Climate   Score  `json:"climate"`   // 5 = mildest
	Safety    Score  `json:"safety"`    // 5 = most relaxed perception
	Nightlife Score  `json:"nightlife"` // 5 = liveliest
	Vibe      string `json:"vibe"`
}

type Score struct {
	Value  int    `json:"value"`
	Reason string `json:"reason"`
}

type Term struct {
	Term     string `json:"term"`
	Language string `json:"language"` // Swahili, Sheng, English (Kenya)
	Country  string `json:"country,omitempty"`
	Meaning  string `json:"meaning"`
	Example  string `json:"example,omitempty"`
}

type Phrase struct {
	Swahili  string `json:"swahili"`
	English  string `json:"english"`
	Say      string `json:"say"` // pronunciation hint
	Category string `json:"category"`
}

type FAQ struct {
	Q string `json:"q"`
	A string `json:"a"`
}

type Source struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Usage   string `json:"usage"`
	Kind    string `json:"kind,omitempty"`    // "" (content), "advisory", "entry"
	Country string `json:"country,omitempty"` // Country.Code, for entry/visa sites
}
