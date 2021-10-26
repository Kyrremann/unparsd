package parsing

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/kyrremann/unparsd/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func LoadJsonIntoDatabase(file string) (*gorm.DB, error) {
	db, err := OpenInMemoryDatabase()
	if err != nil {
		return nil, err
	}

	checkins, err := ParseJSONToCheckins(file)
	if err != nil {
		return nil, err
	}

	err = insertAllIntoDatabase(checkins, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func OpenInMemoryDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Brewery{}, &models.Beer{}, &models.Venue{}, &models.Checkin{})
	return db, nil
}

func ParseJSON(file string, v interface{}) error {
	jsonFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, v)
	if err != nil {
		return err
	}

	return nil
}

func ParseJSONToCheckins(file string) ([]models.JSONCheckin, error) {
	var checkins []models.JSONCheckin
	err := ParseJSON(file, &checkins)
	if err != nil {
		return nil, err
	}

	return checkins, nil
}

func insertAllIntoDatabase(checkins []models.JSONCheckin, db *gorm.DB) error {
	for _, checkin := range checkins {
		err := InsertIntoDatabase(checkin, db)
		if err != nil {
			return err
		}
	}

	return nil
}

func InsertIntoDatabase(jsonCheckin models.JSONCheckin, db *gorm.DB) error {
	var ratingScore float32
	if len(jsonCheckin.RatingScore) > 0 {
		ratingScore64, err := strconv.ParseFloat(jsonCheckin.RatingScore, 32)
		if err != nil {
			return err
		}
		ratingScore = float32(ratingScore64)
	}

	dbCheckin := models.Checkin{
		ID:             jsonCheckin.CheckinID,
		RatingScore:    ratingScore,
		Comment:        jsonCheckin.Comment,
		FlavorProfiles: jsonCheckin.FlavorProfiles,
		ServingTypes:   jsonCheckin.ServingTypes,
		TotalToasts:    jsonCheckin.TotalToasts,
		TaggedFriends:  jsonCheckin.TaggedFriends,
		TotalComments:  jsonCheckin.TotalComments,
		PhotoUrl:       jsonCheckin.PhotoUrl,
		PurchaseVenue:  jsonCheckin.PurchaseVenue,
		CreatedAt:      jsonCheckin.CreatedAt,
		Venue: models.Venue{
			Name:    jsonCheckin.VenueName,
			City:    jsonCheckin.VenueCity,
			State:   jsonCheckin.VenueCity,
			Country: jsonCheckin.VenueCountry,
			Lat:     jsonCheckin.VenueLat,
			Lng:     jsonCheckin.VenueLng,
		},
		Beer: models.Beer{
			ID:                        jsonCheckin.BID,
			Name:                      jsonCheckin.BeerName,
			Type:                      jsonCheckin.BeerType,
			Abv:                       jsonCheckin.BeerAbv,
			Ibu:                       jsonCheckin.BeerIbu,
			GlobalWeightedRatingScore: jsonCheckin.GlobalWeightedRatingScore,
			GlobalRatingScore:         jsonCheckin.GlobalRatingScore,
			Brewery: models.Brewery{
				ID:      jsonCheckin.BreweryID,
				Name:    jsonCheckin.BreweryName,
				City:    jsonCheckin.BreweryCity,
				State:   jsonCheckin.BreweryState,
				Country: jsonCheckin.BreweryCountry,
			},
		},
	}

	res := db.Create(&dbCheckin)
	return res.Error
}
