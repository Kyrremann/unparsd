package models

import "gorm.io/gorm"

type Checkin struct {
	gorm.Model

	ID             int
	RatingScore    float32
	Comment        string
	CheckinAt      string
	FlavorProfiles string
	ServingTypes   string
	TaggedFriends  string
	TotalToasts    int
	TotalComments  int
	PhotoUrl       string
	PurchaseVenue  string
	BeerID         int
	Beer           Beer
	VenueName      string
	Venue          Venue `gorm:"references:Name"`
}

type JSONCheckin struct {
	BID                       int     `json:"bid,string"`
	BeerAbv                   float32 `json:"beer_abv,string"`
	BeerIbu                   int     `json:"beer_ibu,string"`
	BeerName                  string  `json:"beer_name"`
	BeerType                  string  `json:"beer_type"`
	BreweryCity               string  `json:"brewery_city"`
	BreweryCountry            string  `json:"brewery_country"`
	BreweryID                 int     `json:"brewery_id,string"`
	BreweryName               string  `json:"brewery_name"`
	BreweryState              string  `json:"brewery_state"`
	CheckinAt                 string  `json:"created_at"`
	CheckinID                 int     `json:"checkin_id,string"`
	Comment                   string  `json:"comment"`
	FlavorProfiles            string  `json:"flavor_profiles"`
	GlobalRatingScore         float32 `json:"global_rating_score"`
	GlobalWeightedRatingScore float32 `json:"global_weighted_rating_score"`
	PhotoUrl                  string  `json:"photo_url"`
	PurchaseVenue             string  `json:"purchase_venue"`
	RatingScore               string  `json:"rating_score"`
	ServingTypes              string  `json:"serving_types"`
	TaggedFriends             string  `json:"tagged_friends"`
	TotalComments             int     `json:"total_comments,string"`
	TotalToasts               int     `json:"total_toasts,string"`
	VenueCity                 string  `json:"venue_city"`
	VenueCountry              string  `json:"venue_country"`
	VenueLat                  string  `json:"venue_lat"`
	VenueLng                  string  `json:"venue_lng"`
	VenueName                 string  `json:"venue_name"`
	VenueState                string  `json:"venue_state"`
}
