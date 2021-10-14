package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/kyrremann/unparsd/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {}

func OpenInMemoryDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Brewery{}, &models.Beer{}, &models.Venue{}, &models.Checkin{})
	return db, nil
}

func ParseJSON(file string) ([]models.JSONCheckin, error) {
	jsonFile, err := os.Open("untappd.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var checkins []models.JSONCheckin
	err = json.Unmarshal(byteValue, &checkins)
	if err != nil {
		return nil, err
	}

	return checkins, nil
}
