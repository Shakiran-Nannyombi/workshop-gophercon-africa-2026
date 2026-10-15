package data

// CountryProfile is the "at a glance" and practical detail on a country page.
type CountryProfile struct {
	Slug           string       `json:"slug"`
	Depth          string       `json:"depth"`          // "full" | "light"
	Capital        string       `json:"capital"`        // City.Slug of the official capital
	CommercialCity string       `json:"commercialCity"` // City.Slug, empty when it is the capital
	CapitalNote    string       `json:"capitalNote"`
	Population     string       `json:"population"` // "~56 million (2024 est.)"
	TimeZone       string       `json:"timeZone"`
	CallingCode    string       `json:"callingCode"`
	DrivingSide    string       `json:"drivingSide"` // "left" | "right"
	Plugs          []string     `json:"plugs"`       // "G", "C"...
	Airports       []Airport    `json:"airports"`
	Intro          string       `json:"intro"`
	Organisation   Organisation `json:"organisation"`
	GettingIn      GettingIn    `json:"gettingIn"`
	CapitalModes   []Mode       `json:"capitalModes"` // how people move inside the capital
	Advisory       *Advisory    `json:"advisory,omitempty"`
}

type Airport struct {
	IATA string `json:"iata"`
	Name string `json:"name"`
	City string `json:"city"`
	Main bool   `json:"main,omitempty"`
}

// Organisation explains how a country is divided, so a beginner can read an address.
type Organisation struct {
	Summary        string       `json:"summary"`
	Levels         []AdminLevel `json:"levels"` // largest first
	AddressExample string       `json:"addressExample"`
	KenyaCompare   string       `json:"kenyaCompare"`
	Verify         Verify       `json:"verify,omitempty"`
}

type AdminLevel struct {
	Name  string `json:"name"`            // "Region", "District"
	Local string `json:"local,omitempty"` // local-language name, if commonly seen
	Count string `json:"count"`           // "~31", "10"
	Note  string `json:"note"`
}

type GettingIn struct {
	AirportToCity []Mode   `json:"airportToCity"`
	RideHail      []string `json:"rideHail"` // apps known to operate in the capital
	Notes         []string `json:"notes"`
	Verify        Verify   `json:"verify,omitempty"`
}

// Mode is one way of getting about; Icon picks from the shared icon set.
type Mode struct {
	Icon   string `json:"icon"` // matatu, boda, bus, taxi, train, plane, ferry, walk, tuktuk, car
	Name   string `json:"name"`
	Cost   string `json:"cost,omitempty"` // "~KSh 100–300"
	Time   string `json:"time,omitempty"`
	Note   string `json:"note"`
	Verify Verify `json:"verify,omitempty"`
}

// Advisory flags a country with active conflict or government travel advisories.
type Advisory struct {
	Summary  string `json:"summary"`
	Reviewed string `json:"reviewed"` // YYYY-MM-DD
	Verify   Verify `json:"verify,omitempty"`
}

// Currency is one country's money in depth.
type Currency struct {
	Code         string         `json:"code"`
	Name         string         `json:"name"`
	Symbol       string         `json:"symbol"`
	Country      string         `json:"country"` // Country.Code
	Subunit      string         `json:"subunit,omitempty"`
	Issuer       string         `json:"issuer"`
	InPractice   string         `json:"inPractice,omitempty"` // e.g. "USD for most larger payments"
	Notes        []Denomination `json:"notes"`
	Coins        []Denomination `json:"coins"`
	PriceFormats []PriceFormat  `json:"priceFormats"`
	Cash         Cash           `json:"cash"`
	FXNotes      []DatedNote    `json:"fxNotes,omitempty"`
	Verify       Verify         `json:"verify,omitempty"`
}

type Cash struct {
	USD        string `json:"usd"` // how widely US dollars are accepted
	ATMs       string `json:"atms"`
	Cards      string `json:"cards"`
	Exchange   string `json:"exchange"`
	RefusedUSD string `json:"refusedUSD"` // which note years or series get refused
	Verify     Verify `json:"verify,omitempty"`
}

// DatedNote is a time-sensitive warning, such as a parallel exchange-rate gap.
type DatedNote struct {
	Date   string `json:"date"` // YYYY-MM-DD it was last checked
	Title  string `json:"title"`
	Body   string `json:"body"`
	Verify Verify `json:"verify,omitempty"`
}

// Route links two cities and is shown for either direction.
type Route struct {
	From    string        `json:"from"` // City.Slug
	To      string        `json:"to"`
	Summary string        `json:"summary"`
	Border  string        `json:"border,omitempty"` // BorderCrossing.Slug on the usual road route
	Options []RouteOption `json:"options"`
}

type RouteOption struct {
	Mode     string `json:"mode"` // flight, bus, sgr, train, shared-taxi, road, ferry
	Label    string `json:"label"`
	Duration string `json:"duration"` // "~1 h"
	CostMin  int    `json:"costMin"`
	CostMax  int    `json:"costMax"`
	Currency string `json:"currency"`
	Notes    string `json:"notes"`
	Verify   Verify `json:"verify,omitempty"`
}

// BorderCrossing is a land border post between two countries.
type BorderCrossing struct {
	Slug      string   `json:"slug"`
	Name      string   `json:"name"`
	Countries []string `json:"countries"` // two Country.Codes
	Towns     string   `json:"towns"`     // "Busia (KE) / Busia (UG)"
	Road      string   `json:"road"`      // which journey it is on
	OneStop   bool     `json:"oneStop"`   // One Stop Border Post
	Notes     string   `json:"notes"`
	Verify    Verify   `json:"verify,omitempty"`
}

// Borders is the cross-border guide.
type Borders struct {
	Intro     string           `json:"intro"`
	EAC       []string         `json:"eac"`
	EATV      []string         `json:"eatv"`
	Prepare   []string         `json:"prepare"`
	Crossings []BorderCrossing `json:"crossings"`
	Verify    Verify           `json:"verify,omitempty"`
}

// Connectivity is SIMs, data and roaming per country.
type Connectivity struct {
	Intro     string       `json:"intro"`
	Roaming   []string     `json:"roaming"`
	ESIM      []string     `json:"esim"`
	Countries []CountryNet `json:"countries"`
	Verify    Verify       `json:"verify,omitempty"`
}

type CountryNet struct {
	Country  string    `json:"country"`
	Carriers []Carrier `json:"carriers"`
	SIM      string    `json:"sim"`
	Data     string    `json:"data"`
	ESIM     string    `json:"esim"`
	Verify   Verify    `json:"verify,omitempty"`
}

type Carrier struct {
	Name  string `json:"name"`
	Note  string `json:"note"`
	Money string `json:"money,omitempty"` // mobile-money service it runs
}
