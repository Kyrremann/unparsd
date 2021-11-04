package statistics

import (
	"github.com/kyrremann/unparsd/parsing"
	"gorm.io/gorm"
)

func GenerateAndSave(db *gorm.DB, path string) error {
	missingStyles, err := MissingStyles(db)
	if err != nil {
		return err
	}

	err = parsing.SaveDataToJsonFile(missingStyles, path+"/missing_styles.json")
	if err != nil {
		return err
	}

	distinctStyles, err := DistinctStyles(db)
	if err != nil {
		return err
	}

	err = parsing.SaveDataToJsonFile(distinctStyles, path+"/styles.json")
	if err != nil {
		return err
	}

	breweries, err := BreweryStats(db)
	if err != nil {
		return err
	}

	err = parsing.SaveDataToJsonFile(breweries, path+"/breweries.json")
	if err != nil {
		return err
	}

	beers, err := BeerStats(db)
	if err != nil {
		return err
	}

	err = parsing.SaveDataToJsonFile(beers, path+"/beers.json")
	if err != nil {
		return err
	}

	countries, err := CountryStats(db)
	if err != nil {
		return err
	}

	err = parsing.SaveDataToJsonFile(countries, path+"/countries.json")
	if err != nil {
		return err
	}

	allMyStats, err := AllMyStats(db)
	if err != nil {
		return err
	}

	err = parsing.SaveDataToJsonFile(allMyStats, path+"/allmy.json")
	return err
}
