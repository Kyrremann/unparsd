package main

import (
	_ "embed"
	"log"
	"path/filepath"

	"github.com/kyrremann/unparsd/parsing"
	"github.com/kyrremann/unparsd/statistics"
)

//go:embed fixture/all_styles.json
var allStylesJson []byte

func main() {
	statistics.AllStylesJson = allStylesJson
	db, err := parsing.LoadJsonIntoDatabase("untappd.json")
	if err != nil {
		log.Fatal(err)
	}

	base, err := filepath.Abs("output/")
	if err != nil {
		log.Fatal(err)
	}

	err = statistics.GenerateAndSave(db, base)
	if err != nil {
		log.Fatal(err)
	}
}
