package handlers

import (
	"net/http"

	"workshop/internal/views/pages"
)

func (s *site) cities(w http.ResponseWriter, r *http.Request) {
	m := s.meta(r, "cities", "East African cities explained: Nairobi, Kampala, Dar, Kigali, Addis and more",
		"A visual guide to 13 cities in Kenya, Uganda, Tanzania, Rwanda and Ethiopia: a map, what each city is like, its currency and mobile money, and links to compare any two.")
	s.render(w, r, http.StatusOK, pages.Cities(m, s.store), nil)
}

func (s *site) city(w http.ResponseWriter, r *http.Request) {
	c, ok := s.store.City(r.PathValue("slug"))
	if !ok {
		s.notFound(w, r)
		return
	}
	country, _ := s.store.Country(c.Country)
	where := c.Name
	if country != nil {
		where += ", " + country.Name
	}
	m := s.meta(r, "cities", c.Name+" for first-time visitors: areas, money, transport and food",
		c.Tagline+". "+where+" for first-time visitors: neighbourhoods, currency and mobile money, getting around, food, safety and who it suits.")
	m.Path = "/cities/" + c.Slug
	s.render(w, r, http.StatusOK, pages.City(m, s.store, *c, country), nil)
}
