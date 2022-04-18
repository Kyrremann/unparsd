package statistics

import (
	"testing"

	"github.com/kyrremann/unparsd/parsing"
	"github.com/stretchr/testify/assert"
)

func TestBreweries(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/untappd.json")
	assert.NoError(t, err)

	breweries, err := BreweryStats(db)
	assert.NoError(t, err)
	assert.Len(t, breweries, 80)

	brewery := breweries[0]
	assert.Equal(t, "United States", brewery.Country)
	assert.Equal(t, "US", brewery.ISO3166Alpha2)

	var checkins = 0
	for _, b := range breweries {
		checkins += b.Checkins
	}
	assert.Equal(t, 126, checkins)
}
