package statistics

import (
	"fmt"
	"os"
	"text/template"

	"github.com/kyrremann/unparsd/parsing"
	"gorm.io/gorm"
)

func GenerateAndSave(db *gorm.DB, path, allStyles string) error {
	path = path + "/_data"
	if err := os.MkdirAll(path, 0o755); err != nil {
		return err
	}

	missingStyles, err := MissingStyles(db, allStyles)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(missingStyles, path+"/missing_styles.json"); err != nil {
		return err
	}

	distinctStyles, err := DistinctStyles(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(distinctStyles, path+"/styles.json"); err != nil {
		return err
	}

	breweries, err := BreweryStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(breweries, path+"/breweries.json"); err != nil {
		return err
	}

	beers, err := BeerStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(beers, path+"/beers.json"); err != nil {
		return err
	}

	countries, err := CountryStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(countries, path+"/countries.json"); err != nil {
		return err
	}

	missingCountries, err := MissingCountries(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(missingCountries, path+"/missing_countries.json"); err != nil {
		return err
	}

	allMyStats, err := AllMyStats(db)
	if err != nil {
		return err
	}

	if err := parsing.SaveDataToJsonFile(allMyStats, path+"/allmy.json"); err != nil {
		return err
	}

	checkinsByDay, err := CheckinsByDay(db)
	if err != nil {
		return err
	}

	if err = parsing.SaveDataToJsonFile(checkinsByDay, path+"/checkins_by_day.json"); err != nil {
		return err
	}

	return nil
}

func GenerateMonthlyAndSave(db *gorm.DB, path string) error {
	path = path + "/_monthly"
	err := os.MkdirAll(path, 0o755)
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
		output, err := os.Create(path + fmt.Sprintf("/%v.html", y.Year))
		if err != nil {
			return err
		}
		err = tmpl.Execute(output, y)
		if err != nil {
			return err
		}
	}

	return nil
}
