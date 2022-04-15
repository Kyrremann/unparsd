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

	var checkins = 0
	for _, b := range breweries {
		checkins += b.Checkins
	}
	assert.Equal(t, 126, checkins)
}
