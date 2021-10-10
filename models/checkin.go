package models

import "gorm.io/gorm"

type Checkin struct {
	gorm.Model

	ID             int
	RatingScore    float32
	Comment        string
	CreatedAt      string
	FlavorProfiles string
	ServingTypes   string
	TaggedFriends  string
	TotalToasts    int
	TotalComments  int
	PhotoUrl       string
	PurchaseVenue  string
	BeerID         int
	Beer           Beer
	VenueID        string
	Venue          Venue `gorm:"references:Name"`
}

type JSONCheckin struct {
	CheckinID                 int     `json:"checkin_id,string"`
	RatingScore               float32 `json:"rating_score,string"`
	Comment                   string
	CreatedAt                 string  `json:"created_at"`
	FlavorProfiles            string  `json:"flavor_profiles"`
	ServingTypes              string  `json:"serving_types"`
	TaggedFriends             string  `json:"tagged_friends"`
	TotalToasts               int     `json:"total_toasts,string"`
	TotalComments             int     `json:"total_comments,string"`
	PhotoUrl                  string  `json:"photo_url"`
	PurchaseVenue             string  `json:"purchase_venue"`
	VenueName                 string  `json:"venue_name"`
	VenueCity                 string  `json:"venue_city"`
	VenueState                string  `json:"venue_state"`
	VenueCountry              string  `json:"venue_country"`
	VenueLat                  string  `json:"venue_lat"`
	VenueLng                  string  `json:"venue_lng"`
	BreweryID                 int     `json:"brewery_id,string"`
	BreweryName               string  `json:"brewery_name"`
	BreweryCity               string  `json:"brewery_city"`
	BreweryState              string  `json:"brewery_state"`
	BreweryCountry            string  `json:"brewery_country"`
	BID                       int     `json:"bid,string"`
	BeerName                  string  `json:"beer_name"`
	BeerType                  string  `json:"beer_type"`
	BeerAbv                   float32 `json:"beer_abv,string"`
	BeerIbu                   int     `json:"beer_ibu,string"`
	GlobalWeightedRatingScore float32 `json:"global_weighted_rating_score,string"`
	GlobalRatingScore         float32 `json:"global_rating_score,string"`
}
