package statistics

import (
	"testing"

	"github.com/kyrremann/unparsd/models"
	"github.com/kyrremann/unparsd/parsing"
	"github.com/pariz/gountries"
	"github.com/stretchr/testify/assert"
)

func TestCountries(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/untappd.json")
	assert.NoError(t, err)

	var brewery models.Brewery
	res := db.First(&brewery)
	assert.NoError(t, res.Error)

	assert.Equal(t, "United States", brewery.Country)

	query := gountries.New()
	gCountry, err := query.FindCountryByName(brewery.Country)
	assert.NoError(t, err)
	assert.Equal(t, "US", gCountry.Alpha2)

	countries, err := CountryStats(db)
	assert.NoError(t, err)
	country := countries[0]
	assert.Equal(t, "Australia", country.Name)
	assert.Equal(t, 1, country.Checkins)
	assert.Equal(t, 1, country.Breweries)
}
