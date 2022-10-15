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

	var brewery Brewery
	for _, b := range breweries {
		if b.Name == "Ting Tar Tid" {
			brewery = b
			break
		}
	}

	assert.Equal(t, "Ting Tar Tid", brewery.Name)
	assert.Equal(t, "Norway", brewery.Country)
	assert.Equal(t, "NO", brewery.ISO3166Alpha2)
	sortedBeers := "En Sjakkmester Weissbier\nEplegløgg\nKaffekværna - 100g\nKaffekværna - 250g\nKaffekværna - Base\nKværna\nMer De Montagne\nRustikk Koriander\nSingle-hop Lucky Jack"
	assert.Equal(t, sortedBeers, brewery.ListOfBeers)

	var checkins = 0
	for _, b := range breweries {
		checkins += b.Checkins
	}
	assert.Equal(t, 126, checkins)
}
