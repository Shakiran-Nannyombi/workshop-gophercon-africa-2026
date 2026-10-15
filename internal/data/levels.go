package data

// Level is one layer of Kenya's geography, from county down to estate (data/levels.json).
type Level struct {
	Name    string `json:"name"`
	Count   string `json:"count"` // "47", "~290", "informal"
	Who     string `json:"who"`   // who runs it / what it's for
	Example string `json:"example"`
	Note    string `json:"note"`
	Verify  Verify `json:"verify,omitempty"`
}
