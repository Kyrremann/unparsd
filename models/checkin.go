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
