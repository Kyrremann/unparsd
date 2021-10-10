package main

import (
	"testing"

	"github.com/kyrremann/unparsd/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	//db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&models.Brewery{}, &models.Beer{}, &models.Venue{}, &models.Checkin{})

	checkin := models.Checkin{
		ID:             283107883,
		RatingScore:    3,
		Comment:        "Fin begynnerøl.",
		FlavorProfiles: "",
		ServingTypes:   "",
		TotalToasts:    1,
		TaggedFriends:  "",
		TotalComments:  0,
		PhotoUrl:       "https://untappd.s3.amazonaws.com/photo/2016_03_01/b507b21dc14f5a0d5fbc48b260f89721_raw.jpg",
		PurchaseVenue:  "",
		CreatedAt:      "2016-03-01 19:06:42",
		Venue: models.Venue{
			Name:    "Mad Fork",
			City:    "",
			State:   "",
			Country: "Norge",
			Lat:     "59.929",
			Lng:     "10.7183",
		},
		Beer: models.Beer{
			ID:                        5939,
			Name:                      "1664",
			Type:                      "Lager - Euro Pale",
			Abv:                       5.5,
			Ibu:                       20,
			GlobalWeightedRatingScore: 3.11,
			GlobalRatingScore:         3.11,
			Brewery: models.Brewery{
				ID:      203,
				Name:    "Kronenbourg Brewery",
				City:    "Obernai",
				State:   "Grand-Est",
				Country: "France",
			},
		},
	}
	res := db.Create(&checkin)
	assert.NoError(t, res.Error)

	var c models.Checkin
	res = db.Preload("Venue").Preload("Beer.Brewery").First(&c, 283107883)
	assert.NoError(t, res.Error)
	assert.Equal(t, 283107883, c.ID)
	assert.Equal(t, "1664", c.Beer.Name)
	assert.Equal(t, "Kronenbourg Brewery", c.Beer.Brewery.Name)
	assert.Equal(t, "Mad Fork", c.Venue.Name)

	var b models.Beer
	res = db.Preload("Brewery").First(&b, 5939)
	assert.NoError(t, res.Error)
	assert.Equal(t, 5939, b.ID)
	assert.Equal(t, 203, b.BreweryID)
	assert.Equal(t, "Kronenbourg Brewery", b.Brewery.Name)

	var brw models.Brewery
	res = db.Preload("Beers").Find(&brw, 203)
	assert.NoError(t, res.Error)
	assert.Equal(t, "Kronenbourg Brewery", brw.Name)
	assert.Len(t, brw.Beers, 1)

	var v models.Venue
	res = db.First(&v)
	assert.NoError(t, res.Error)
	assert.Equal(t, "Mad Fork", v.Name)
}
