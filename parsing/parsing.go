package parsing

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/kyrremann/unparsd/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func LoadJsonIntoDatabase(file string) (*gorm.DB, error) {
	db, err := OpenInMemoryDatabase()
	if err != nil {
		return nil, err
	}

	checkins, err := ParseJsonToCheckins(file)
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
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Brewery{}, &models.Beer{}, &models.Venue{}, &models.Checkin{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ReadFile(file string) ([]byte, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
		}
	}(jsonFile)

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func ParseJsonFile(file string, v interface{}) error {
	bytes, err := ReadFile(file)
	if err != nil {
		return err
	}

	return UnmarshalJson(bytes, v)
}

func UnmarshalJson(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, v)
}

func ParseJsonToCheckins(file string) ([]models.JSONCheckin, error) {
	var checkins []models.JSONCheckin
	err := ParseJsonFile(file, &checkins)
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
	dbCheckin := models.Checkin{
		ID:             jsonCheckin.CheckinID,
		RatingScore:    jsonCheckin.RatingScore,
		Comment:        jsonCheckin.Comment,
		FlavorProfiles: jsonCheckin.FlavorProfiles,
		ServingTypes:   jsonCheckin.ServingTypes,
		TotalToasts:    jsonCheckin.TotalToasts,
		TaggedFriends:  jsonCheckin.TaggedFriends,
		TotalComments:  jsonCheckin.TotalComments,
		PhotoUrl:       jsonCheckin.PhotoUrl,
		PurchaseVenue:  jsonCheckin.PurchaseVenue,
		CheckinAt:      jsonCheckin.CheckinAt,
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
			Type:                      strings.TrimSpace(jsonCheckin.BeerType),
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

	res := db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"photo_url": dbCheckin.PhotoUrl}),
		}).
		Create(&dbCheckin)
	return res.Error
}

func SaveDataToJsonFile(v interface{}, fileName string) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
