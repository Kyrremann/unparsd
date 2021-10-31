package statistics

import (
	"testing"

	"github.com/kyrremann/unparsd/parsing"
	"github.com/stretchr/testify/assert"
)

func TestAllMy(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/untappd.json")
	assert.NoError(t, err)

	globalStats, err := AllMyStats(db)
	assert.NoError(t, err)
	assert.Equal(t, 125, globalStats.Checkins)
	assert.Equal(t, 112, globalStats.UniqueBeers)
	assert.Equal(t, "2016-03-01", globalStats.StartDate)
	assert.Len(t, globalStats.Periodes, 5)
	assert.Len(t, globalStats.Periodes[0].Months, 3)
}
