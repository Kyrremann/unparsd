package statistics

import (
	"github.com/kyrremann/unparsd/models"
	"github.com/pariz/gountries"
	"gorm.io/gorm"
)

type Country struct {
	ID        string   `json:"id"`
	Breweries int      `json:"breweries"`
	Checkins  int      `json:"checkins"`
	Name      string   `json:"name"`
	State     string   `json:"state"`
	Settings  Settings `json:"settings" gorm:"-"`
}

type Settings struct {
	Fill        string `json:"fill"`
	TooltipText string `json:"tooltipText"`
	ToggleKey   string `json:"toggleKey"`
	Interactive string `json:"interactive"`
}

func getDefaultSettings() Settings {
	return Settings{
		Fill:        "#c96e12",
		TooltipText: "{name}\nBreweries: {breweries}\nCheckins: {checkins}",
		ToggleKey:   "active",
		Interactive: "true",
	}
}

type State struct {
	Checkins  int    `json:"checkins"`
	Breweries int    `json:"breweries"`
	Country   string `json:"country"`
	Name      string `json:"name"`
}

type ISO3166Alpha2 struct {
	Query *gountries.Query
}

func (iso ISO3166Alpha2) getISO3166Alpha2(country, state string) (string, error) {
	switch country {
	case "China / People's Republic of China":
		return "CN", nil
	case "Palestinian Territories":
		return "PS", nil
	case "Principality of Monaco":
		return "MC", nil
	case "Wales", "England", "Scotland", "Northern Ireland":
		return "GB", nil
	case "Surinam":
		return "SR", nil
	case "North Macedonia":
		return "MK", nil
	case "France":
		switch state {
		case "Guyane":
			return "GF", nil
		default:
			break
		}
	}

	gountry, err := iso.Query.FindCountryByName(country)
	if err != nil {
		return "", err
	}
	return gountry.Alpha2, nil
}

func CountryStats(db *gorm.DB) ([]Country, error) {
	var dbCountries []Country
	res := db.
		Model(models.Checkin{}).
		Select("count(checkins.id) as checkins," +
			"count(DISTINCT(breweries.id)) as breweries," +
			"breweries.country as name," +
			"breweries.state as state").
		Joins("LEFT JOIN beers on checkins.beer_id == beers.id").
		Joins("LEFT JOIN breweries on beers.brewery_id == breweries.id").
		Group("breweries.country").
		Find(&dbCountries)
	if res.Error != nil {
		return nil, res.Error
	}

	iso := ISO3166Alpha2{
		Query: gountries.New(),
	}
	var countries []Country
	for _, c := range dbCountries {
		ISO3166Alpha2, err := iso.getISO3166Alpha2(c.Name, c.State)
		if err != nil {
			return nil, err
		}

		c.ID = ISO3166Alpha2
		c.Settings = getDefaultSettings()
		countries = append(countries, c)
	}

	return countries, nil
}

func CountryStateStats(db *gorm.DB) (map[string]State, error) {
	var dbStates []State
	res := db.
		Model(models.Checkin{}).
		Select("count(checkins.id) as checkins," +
			"count(DISTINCT(breweries.id)) as breweries," +
			"breweries.country as country," +
			"breweries.state as name").
		Joins("LEFT JOIN beers on checkins.beer_id == beers.id").
		Joins("LEFT JOIN breweries on beers.brewery_id == breweries.id").
		Group("breweries.state").
		Find(&dbStates)
	if res.Error != nil {
		return nil, res.Error
	}

	iso := ISO3166Alpha2{
		Query: gountries.New(),
	}
	states := make(map[string]State)
	for _, c := range dbStates {
		ISO3166Alpha2, err := iso.getISO3166Alpha2(c.Country, c.Name)
		if err != nil {
			return nil, err
		}

		states[ISO3166Alpha2] = c
	}

	return states, nil
}
