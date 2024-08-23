package statistics

import (
	"fmt"
	"slices"
	"strings"

	"github.com/biter777/countries"
	"github.com/kyrremann/unparsd/models"
	"gorm.io/gorm"
)

type Country struct {
	ID        string    `json:"id"`
	Breweries int       `json:"breweries,omitempty"`
	Checkins  int       `json:"checkins,omitempty"`
	Name      string    `json:"name"`
	State     string    `json:"-"`
	Settings  *Settings `json:"settings,omitempty" gorm:"-"`
}

type Settings struct {
	Fill        string `json:"fill"`
	TooltipText string `json:"tooltipText"`
	ToggleKey   string `json:"toggleKey"`
	Interactive string `json:"interactive"`
}

func getDefaultSettings() *Settings {
	return &Settings{
		Fill:        "#c96e12",
		TooltipText: "{name}\nBreweries: {breweries}\nCheckins: {checkins}",
		ToggleKey:   "active",
		Interactive: "true",
	}
}

func getISO3166Alpha2(country, state string) (string, string, error) {
	var alpha2 string

	switch country {
	case "Australia":
		switch state {
		case "The Territory of Christmas Island":
			alpha2 = "CX"
		case "The Territory of Cocos (Keeling) Islands":
			alpha2 = "CC"
		case "The Territory of Heard Island and McDonald Islands":
			alpha2 = "HM"
		case "The Territory of Norfolk Island":
			alpha2 = "NF"
		}
	case "China":
		switch state {
		case "The Hong Kong Special Administrative Region of China":
			alpha2 = "HK"
		case "The Macao Special Administrative Region of China":
			alpha2 = "MO"
		}
	case "China / People's Republic of China":
		alpha2 = "CN"
	case "Curaçao", "Curacao":
		alpha2 = "CW"
	case "Denmark":
		switch state {
		case "The Faroe Islands":
			alpha2 = "FO"
		case "Kalaallit Nunaat":
			alpha2 = "GL"
		}
	case "England", "Northern Ireland", "Scotland", "Wales":
		alpha2 = "GB"
	case "Finland":
		switch state {
		case "Åland":
			alpha2 = "AX"
		}
	case "France":
		switch state {
		case "Guyane":
			alpha2 = "GF"
			country = "Guyane"
		case "French Polynesia":
			alpha2 = "PF"
		case "The French Southern and Antarctic Lands":
			alpha2 = "TF"
		case "Guadeloupe":
			alpha2 = "GP"
		case "Martinique":
			alpha2 = "MQ"
		case "The Department of Mayotte":
			alpha2 = "YT"
		case "New Caledonia":
			alpha2 = "NC"
		case "Réunion", "La Réunion":
			alpha2 = "RE"
			country = "Réunion"
		case "The Collectivity of Saint-Barthélemy":
			alpha2 = "BL"
		case "The Collectivity of Saint-Martin":
			alpha2 = "MF"
		case "The Overseas Collectivity of Saint-Pierre and Miquelon":
			alpha2 = "PM"
		case "The Territory of the Wallis and Futuna Islands":
			alpha2 = "WF"
		}
	case "Guyane":
		alpha2 = "GF"
	case "Netherlands":
		switch state {
		case "Aruba":
			alpha2 = "AW"
		case "Bonaire":
			alpha2 = "BQ"
		case "Curaçao", "Curacao":
			alpha2 = "CW"
		case "Saba":
			alpha2 = "BQ"
		case "Sint Eustatius":
			alpha2 = "BQ"
		case "Sint Maarten":
			alpha2 = "SQ"
		}
	case "New Zealand":
		switch state {
		case "The Cook Islands":
			alpha2 = "CK"
		case "Niue":
			alpha2 = "NU"
		case "Tokelau":
			alpha2 = "TK"
		}
	case "North Macedonia":
		alpha2 = "MK"
	case "Norway":
		switch state {
		case "Bouvet Island":
			alpha2 = "BV"
		case "Svalbard and Jan Mayen":
			alpha2 = "SJ"
		}
	case "Palestinian Territories":
		alpha2 = "PS"
	case "Principality of Monaco":
		alpha2 = "MC"
	case "Spain":
		switch state {
		case "Canarias", "Canary Islands":
			alpha2 = "IC"
			country = "Canarias"
		}
	case "St. Lucia":
		alpha2 = "LC"
	case "Surinam":
		alpha2 = "SR"
	case "United States Virgin Islands":
		alpha2 = "VI"
	}

	if alpha2 != "" {
		return country, alpha2, nil
	}

	gountry := countries.ByName(country)
	if gountry == countries.Unknown {
		return "", "", fmt.Errorf("unknown country: %s", country)
	}

	return country, gountry.Alpha2(), nil
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
		Group("breweries.country, breweries.state").
		Find(&dbCountries)
	if res.Error != nil {
		return nil, res.Error
	}

	countriesMap := make(map[string]Country, len(dbCountries))
	for _, c := range dbCountries {
		country, ISO3166Alpha2, err := getISO3166Alpha2(c.Name, c.State)
		if err != nil {
			return nil, err
		}

		c.Name = country
		c.ID = ISO3166Alpha2

		if countriesMap[c.ID].Name == "" {
			c.Settings = getDefaultSettings()
			countriesMap[c.ID] = c
		} else {
			tmp := countriesMap[c.ID]
			tmp.Breweries += c.Breweries
			tmp.Checkins += c.Checkins
			countriesMap[c.ID] = tmp
		}
	}

	countriesSlice := make([]Country, 0, len(countriesMap))
	for _, c := range countriesMap {
		countriesSlice = append(countriesSlice, c)
	}

	slices.SortFunc(countriesSlice, func(a, b Country) int {
		return strings.Compare(a.Name, b.Name)
	})

	return countriesSlice, nil
}

func MissingCountries(db *gorm.DB) ([]Country, error) {
	stats, err := CountryStats(db)
	if err != nil {
		return nil, err
	}

	var checkedInCountries []string
	for _, country := range stats {
		checkedInCountries = append(checkedInCountries, country.ID)
	}

	var allCountries []string
	for _, country := range countries.All() {
		allCountries = append(allCountries, country.Alpha2())
	}

	missingCountriesAsString := intersection(allCountries, checkedInCountries)
	missingCountries := make([]Country, 0, len(missingCountriesAsString))
	for _, alpha2 := range missingCountriesAsString {
		country := countries.ByName(alpha2)
		name := country.String()

		if name == "Taiwan (Province of China)" {
			name = "Taiwan" // Taiwan is not a province of China
		}

		missingCountries = append(missingCountries, Country{
			ID:   alpha2,
			Name: name,
		})
	}

	slices.SortFunc(missingCountries, func(a, b Country) int {
		return strings.Compare(a.Name, b.Name)
	})

	return missingCountries, nil
}
