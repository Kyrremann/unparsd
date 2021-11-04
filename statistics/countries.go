package statistics

import (
	"github.com/kyrremann/unparsd/models"
	"github.com/pariz/gountries"
	"gorm.io/gorm"
)

type Country struct {
	FillKey   string `json:"fillKey"`
	Breweries int    `json:"breweries"`
	Checkins  int    `json:"checkins"`
	Country   string `json:"-"`
}

type ISO3166Alpha3 struct {
	Query *gountries.Query
}

func (iso ISO3166Alpha3) getISO3166Alpha3(country string) (string, error) {
	switch country {
	case "China / People's Republic of China":
		return "CHN", nil
	case "Palestinian Territories":
		return "PSE", nil
	case "Principality of Monaco":
		return "MCO", nil
	case "Wales", "England", "Scotland", "Northern Ireland":
		return "GBR", nil
	}

	gountry, err := iso.Query.FindCountryByName(country)
	if err != nil {
		return "", err
	}
	return gountry.Alpha3, nil
}

func CountryStats(db *gorm.DB) (map[string]Country, error) {
	var dbCountries []Country
	res := db.
		Model(models.Checkin{}).
		Select("count(checkins.id) as checkins," +
			"count(DISTINCT(breweries.id)) as breweries," +
			"breweries.country," +
			"'brewery' as fill_key").
		Joins("LEFT JOIN beers on checkins.beer_id == beers.id").
		Joins("LEFT JOIN breweries on beers.brewery_id == breweries.id").
		Group("breweries.country").
		Find(&dbCountries)
	if res.Error != nil {
		return nil, res.Error
	}

	iso := ISO3166Alpha3{
		Query: gountries.New(),
	}
	countries := make(map[string]Country)
	for _, c := range dbCountries {
		ISO3166Alpha3, err := iso.getISO3166Alpha3(c.Country)
		if err != nil {
			return nil, err
		}

		countries[ISO3166Alpha3] = c
	}

	return countries, nil
}
