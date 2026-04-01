package statistics

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/kyrremann/unparsd/parsing"
	"gorm.io/gorm"
)

func GenerateAndSave(db *gorm.DB, path, allStyles string) error {
	dataPath := filepath.Join(path, "_data")
	if err := os.MkdirAll(dataPath, 0o755); err != nil {
		return err
	}

	missingStyles, err := MissingStyles(db, allStyles)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(missingStyles, filepath.Join(dataPath, "missing_styles.json")); err != nil {
		return err
	}

	distinctStyles, err := DistinctStyles(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(distinctStyles, filepath.Join(dataPath, "styles.json")); err != nil {
		return err
	}

	breweries, err := BreweryStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(breweries, filepath.Join(dataPath, "breweries.json")); err != nil {
		return err
	}

	beers, err := BeerStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(beers, filepath.Join(dataPath, "beers.json")); err != nil {
		return err
	}

	countries, err := CountryStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(countries, filepath.Join(dataPath, "countries.json")); err != nil {
		return err
	}

	missingCountries, err := MissingCountries(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(missingCountries, filepath.Join(dataPath, "missing_countries.json")); err != nil {
		return err
	}

	allMyStats, err := AllMyStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(allMyStats, filepath.Join(dataPath, "allmy.json")); err != nil {
		return err
	}

	ratingDeltas, err := RatingDeltas(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(ratingDeltas, filepath.Join(dataPath, "rating_deltas.json")); err != nil {
		return err
	}

	topVenues, err := TopVenues(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(topVenues, filepath.Join(dataPath, "venues.json")); err != nil {
		return err
	}

	servingTypes, err := ServingTypeStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(servingTypes, filepath.Join(dataPath, "serving_types.json")); err != nil {
		return err
	}

	flavors, err := FlavorProfileStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(flavors, filepath.Join(dataPath, "flavors.json")); err != nil {
		return err
	}

	return nil
}

func GenerateMonthlyAndSave(db *gorm.DB, path string) error {
	monthlyPath := filepath.Join(path, "_monthly")
	err := os.MkdirAll(monthlyPath, 0o755)
	if err != nil {
		return err
	}

	view := `---
layout: monthly
banner: In {{ .Year}} I started drinking {{ .StartDay }}th of {{ .StartMonth }} and I managed to drink {{ .Checkins }} beers, averaging {{ .BeersPerDay }} beers a day
---

{% for value in site.data.allmy.years['{{ .Year }}'].months %}
  {% cycle 'add row' : '<div class="boxes-tables pure-g">', '', '' %}
  {% include infoboxes.html data=value %}
  {% cycle 'end row' : '', '', '</div>' %}
{% endfor %}
{% cycle 'end row' : '', '</div>', '</div>' %}
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
		output, err := os.Create(filepath.Join(monthlyPath, fmt.Sprintf("%v.html", y.Year)))
		if err != nil {
			return err
		}
		err = tmpl.Execute(output, y)
		output.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
