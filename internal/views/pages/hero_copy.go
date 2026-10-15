package pages

import (
	"fmt"
	"workshop/internal/data"
	"workshop/internal/views/components"
)

// heroCopyID is the id of the home hero's copy block, swapped when a country is picked.
const heroCopyID = "hero-copy"

// countryHeroes is the home hero's wording for each country in the city picker.
var countryHeroes = map[string]components.HeroProps{
	"KE": {
		Eyebrow: "Karibu · welcome",
		Title:   "Kenya, explained from zero.",
		Lede:    "Confused by the money, the map, or why your Nairobi address says Kiambu? This guide starts where you are: never been, slightly lost, very curious.",
	},
	"UG": {
		Eyebrow: "Tukusanyukidde · welcome",
		Title:   "Uganda, explained from zero.",
		Lede:    "Shillings with lots of zeros, MTN MoMo at every kiosk and boda bodas on every hill. Start in Kampala, then follow the lake shore to Entebbe and the Nile at Jinja.",
	},
	"TZ": {
		Eyebrow: "Karibu · welcome",
		Title:   "Tanzania, explained from zero.",
		Lede:    "From Dar es Salaam’s harbour to Arusha under Kilimanjaro and Zanzibar’s Stone Town: shillings, M-Pesa and Swahili will carry you a long way.",
	},
	"RW": {
		Eyebrow: "Murakaza neza · welcome",
		Title:   "Rwanda, explained from zero.",
		Lede:    "A thousand hills, francs instead of shillings, and Kigali: clean, calm and easy to find your feet in. Here’s what to know before you land.",
	},
	"ET": {
		Eyebrow: "Enkwan dehna metu · welcome",
		Title:   "Ethiopia, explained from zero.",
		Lede:    "Birr, telebirr, its own calendar and the best macchiato in the region. Addis Ababa sits high in the highlands, and it rewards the curious.",
	},
}

// homeHero returns the home hero for a country, falling back to Kenya.
func homeHero(country string) components.HeroProps {
	p, ok := countryHeroes[country]
	if !ok {
		p = countryHeroes["KE"]
	}
	p.Motif = "savannah"
	p.Tall = true
	p.CopyID = heroCopyID
	return p
}

type homeFact struct {
	Eyebrow string
	Big     string
	Body    string
	USD     string
}

func homeFacts(s *data.Store, country string) []homeFact {
	money := moneyFact(s, country)
	switch country {
	case "UG":
		return []homeFact{
			money,
			{Eyebrow: "Map", Big: "Districts, not counties", Body: "Kampala is the capital. Uganda is organised into districts."},
			{Eyebrow: "Pay", Big: "MTN MoMo", Body: "Most kiosks take mobile money before they take a card."},
		}
	case "TZ":
		return []homeFact{
			money,
			{Eyebrow: "Map", Big: "Mainland + Zanzibar", Body: "Dar es Salaam is the biggest city; Dodoma is the official capital."},
			{Eyebrow: "Pay", Big: "M-Pesa", Body: "Vodacom M-Pesa, Mixx and Airtel Money are how most people pay."},
		}
	case "RW":
		return []homeFact{
			money,
			{Eyebrow: "Map", Big: "30 districts", Body: "Kigali is the capital, plus four provinces around the hills."},
			{Eyebrow: "Pay", Big: "MTN MoMo", Body: "Cards work in Kigali hotels; phones work at almost every kiosk."},
		}
	case "ET":
		return []homeFact{
			money,
			{Eyebrow: "Map", Big: "Addis Ababa", Body: "The capital, and home of the African Union."},
			{Eyebrow: "Pay", Big: "telebirr", Body: "Cards are less common than in Nairobi; mobile money is growing fast."},
		}
	default:
		return []homeFact{
			money,
			{Eyebrow: "Map", Big: "47 counties", Body: "Nairobi is one of them, and it’s also a city."},
			{Eyebrow: "Pay", Big: "M-Pesa", Body: "Most kiosks take it before they take a card."},
		}
	}
}

func moneyFact(s *data.Store, country string) homeFact {
	co, ok := s.Country(country)
	if !ok {
		co, _ = s.Country("KE")
	}
	amount := 1000
	switch country {
	case "UG":
		amount = 10000
	case "TZ":
		amount = 10000
	case "RW":
		amount = 5000
	case "ET":
		amount = 500
	}
	return homeFact{
		Eyebrow: "Money",
		Big:     co.Symbol + " " + comma(amount),
		Body:    "is about ",
		USD:     usdForCode(s, float64(amount), co.CurrencyCode),
	}
}

func usdForCode(s *data.Store, amount float64, code string) string {
	r := s.Site.Rates.PerUSD[code]
	if r == 0 {
		return "?"
	}
	return fmt.Sprintf("%.2f", amount/r)
}

func comma(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) < 4 {
		return s
	}
	out := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out += ","
		}
		out += string(c)
	}
	return out
}
