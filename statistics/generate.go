package statistics

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/kyrremann/unparsd/v4/parsing"
	"gorm.io/gorm"
)

func GenerateAndSave(db *gorm.DB, path, allStyles string) error {
	missingStyles, err := MissingStyles(db, allStyles)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(missingStyles, filepath.Join(path, "missing_styles.json")); err != nil {
		return err
	}

	distinctStyles, err := DistinctStyles(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(distinctStyles, filepath.Join(path, "styles.json")); err != nil {
		return err
	}

	breweries, err := BreweryStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(breweries, filepath.Join(path, "breweries.json")); err != nil {
		return err
	}

	beers, err := BeerStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(beers, filepath.Join(path, "beers.json")); err != nil {
		return err
	}

	countries, err := CountryStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(countries, filepath.Join(path, "countries.json")); err != nil {
		return err
	}

	missingCountries, err := MissingCountries(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(missingCountries, filepath.Join(path, "missing_countries.json")); err != nil {
		return err
	}

	allMyStats, err := AllMyStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(allMyStats, filepath.Join(path, "allmy.json")); err != nil {
		return err
	}

	ratingDeltas, err := RatingDeltas(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(ratingDeltas, filepath.Join(path, "rating_deltas.json")); err != nil {
		return err
	}

	topVenues, err := TopVenues(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(topVenues, filepath.Join(path, "venues.json")); err != nil {
		return err
	}

	servingTypes, err := ServingTypeStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(servingTypes, filepath.Join(path, "serving_types.json")); err != nil {
		return err
	}

	flavors, err := FlavorProfileStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(flavors, filepath.Join(path, "flavors.json")); err != nil {
		return err
	}

	return nil
}

func GenerateMonthlyAndSave(db *gorm.DB, path string) error {
	err := os.MkdirAll(path, 0o750)
	if err != nil {
		return err
	}

	view := `---
layout: monthly
banner: In {{ .Year}} I started drinking {{ .StartDay }}th of {{ .StartMonth }} and I managed to drink {{ .Checkins }} beers, averaging {{ .BeersPerDay }} beers a day
year: "{{ .Year }}"
---
`
	tmpl, err := template.New("monthly").Parse(view)
	if err != nil {
		return err
	}

	monthlyData, err := GetMonthlyBannerSummary(db)
	if err != nil {
		return err
	}
	for _, y := range monthlyData {
		// #nosec G304 -- path is constructed from a trusted base directory and a year string from the DB
		output, err := os.Create(filepath.Join(path, fmt.Sprintf("%v.html", y.Year)))
		if err != nil {
			return err
		}

		err = tmpl.Execute(output, y)
		if err != nil {
			return err
		}

		if err := output.Close(); err != nil {
			return err
		}
	}

	return nil
}
