package parsing

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kyrremann/unparsd/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func LoadJsonIntoDatabase(path string) (*gorm.DB, error) {
	db, err := OpenInMemoryDatabase()
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var checkins []models.JSONCheckin
	if info.IsDir() {
		checkins, err = ParseJsonDirToCheckins(path)
	} else {
		checkins, err = ParseJsonToCheckins(path)
	}
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
	defer func() {
		if cerr := jsonFile.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

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

// ParseJsonDirToCheckins reads all *.json files from dir and returns the
// merged slice of check-ins.
func ParseJsonDirToCheckins(dir string) ([]models.JSONCheckin, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var all []models.JSONCheckin
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		checkins, err := ParseJsonToCheckins(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", entry.Name(), err)
		}
		all = append(all, checkins...)
	}
	return all, nil
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
		RatingScore:    convertRatingScore(jsonCheckin.RatingScore),
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
			State:   jsonCheckin.VenueState,
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
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"photo_url":      dbCheckin.PhotoUrl,
				"rating_score":   dbCheckin.RatingScore,
				"comment":        dbCheckin.Comment,
				"total_toasts":   dbCheckin.TotalToasts,
				"total_comments": dbCheckin.TotalComments,
			}),
		}).
		Create(&dbCheckin)
	return res.Error
}

func convertRatingScore(score any) float32 {
	switch v := score.(type) {
	case float64:
		return float32(v)
	case string:
		return 0
	default:
		log.Printf("unexpected type %T", v)
		return 0
	}
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
	defer file.Close()

	_, err = file.Write(bytes)
	return err
}
