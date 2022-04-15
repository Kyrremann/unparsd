package statistics

import (
	"testing"

	"github.com/kyrremann/unparsd/parsing"
	"github.com/stretchr/testify/assert"
)

func TestBeers(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/untappd.json")
	assert.NoError(t, err)

	beers, err := BeerStats(db)
	assert.NoError(t, err)
	assert.Len(t, beers, 112)

	var checkins = 0
	for _, b := range beers {
		checkins += b.Checkins
	}
	assert.Equal(t, 126, checkins)
}
