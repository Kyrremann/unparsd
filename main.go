package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kyrremann/unparsd/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}

	db.AutoMigrate(&models.Brewery{}, &models.Beer{}, &models.Venue{}, &models.Checkin{})

	// Open our jsonFile
	jsonFile, err := os.Open("untappd.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened untappd.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var checkins []models.Checkin

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'checkins' which we defined above
	json.Unmarshal(byteValue, &checkins)

	for i, c := range checkins {
		_, err = fmt.Printf("%d: %s", i, c.Beer.Name)
		if err != nil {
			fmt.Println(err)
		}
	}
}
