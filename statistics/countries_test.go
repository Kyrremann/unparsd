package statistics

import (
	"fmt"
	"testing"

	"github.com/kyrremann/unparsd/models"
	"github.com/kyrremann/unparsd/parsing"

	"github.com/stretchr/testify/assert"
)

func TestCountries(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/untappd.json")
	assert.NoError(t, err)

	var brewery models.Brewery
	res := db.First(&brewery)
	assert.NoError(t, res.Error)

	assert.Equal(t, "United States", brewery.Country)

	countries, err := CountryStats(db)
	assert.NoError(t, err)
	country := countries[0]
	assert.Equal(t, "Australia", country.Name)
	assert.Equal(t, 1, country.Checkins)
	assert.Equal(t, 1, country.Breweries)
	assert.Equal(t, "AU", country.ID)
}

func TestMissingCountries(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/untappd.json")
	assert.NoError(t, err)

	missingCountries, err := MissingCountries(db)
	assert.NoError(t, err)
	fmt.Println(missingCountries)
	assert.Equal(t, 236, len(missingCountries))
}
