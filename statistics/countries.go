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
	case "Australia":
		switch state {
		case "The Territory of Christmas Island":
			return "CX", nil
		case "The Territory of Cocos (Keeling) Islands":
			return "CC", nil
		case "The Territory of Heard Island and McDonald Islands":
			return "HM", nil
		case "The Territory of Norfolk Island":
			return "NF", nil
		}
	case "China":
		switch state {
		case "The Hong Kong Special Administrative Region of China":
			return "HK", nil
		case "The Macao Special Administrative Region of China":
			return "MO", nil
		}
	case "Denmark":
		switch state {
		case "The Faroe Islands":
			return "FO", nil
		case "Kalaallit Nunaat":
			return "GL", nil
		}
	case "Finland":
		switch state {
		case "Åland":
			return "AX", nil
		}
	case "France":
		switch state {
		case "Guyane":
			return "GF", nil
		case "French Polynesia":
			return "PF", nil
		case "The French Southern and Antarctic Lands":
			return "TF", nil
		case "Guadeloupe":
			return "GP", nil
		case "Martinique":
			return "MQ", nil
		case "The Department of Mayotte":
			return "YT", nil
		case "New Caledonia":
			return "NC", nil
		case "Réunion":
			return "RE", nil
		case "The Collectivity of Saint-Barthélemy":
			return "BL", nil
		case "The Collectivity of Saint-Martin":
			return "MF", nil
		case " 	The Overseas Collectivity of Saint-Pierre and Miquelon":
			return "PM", nil
		case "The Territory of the Wallis and Futuna Islands":
			return "WF", nil
		}
	case "Netherlands":
		switch state {
		case "Aruba":
			return "AW", nil
		case "Bonaire":
			return "BQ", nil
		case "Curaçao":
			return "CW", nil
		case "Saba":
			return "BQ", nil
		case "Sint Eustatius":
			return "BQ", nil
		case "Sint Maarten":
			return "SQ", nil
		}
	case "New Zealand":
		switch state {
		case "The Cook Islands":
			return "CK", nil
		case "Niue":
			return "NU", nil
		case "Tokelau":
			return "TK", nil
		}
	case "Norway":
		switch state {
		case "Bouvet Island":
			return "BV", nil
		case "Svalbard and Jan Mayen":
			return "SJ", nil
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
