package statistics

import (
	"testing"

	"github.com/kyrremann/unparsd/parsing"
	"github.com/stretchr/testify/assert"
)

func TestMostCheckinsPerDay(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/untappd.json")
	assert.NoError(t, err)

	mostCheckinsPerDay, err := MostCheckinsPerDay(db, "2020", "")
	assert.NoError(t, err)
	assert.Equal(t, 4, mostCheckinsPerDay.Count)
	assert.Equal(t, "2020-01-23", mostCheckinsPerDay.Date)

	mostCheckinsPerDay, err = MostCheckinsPerDay(db, "2016", "05")
	assert.NoError(t, err)
	assert.Equal(t, 1, mostCheckinsPerDay.Count)
	assert.Equal(t, "2016-05-15", mostCheckinsPerDay.Date)

	mostUniqueBeersPerDay, err := MostUniqueBeersPerDay(db, "2020", "")
	assert.NoError(t, err)
	assert.Equal(t, 3, mostUniqueBeersPerDay.Count)
	assert.Equal(t, "2020-01-04", mostUniqueBeersPerDay.Date)

	mostUniqueBeersPerDay, err = MostUniqueBeersPerDay(db, "2016", "05")
	assert.NoError(t, err)
	assert.Equal(t, 1, mostUniqueBeersPerDay.Count)
	assert.Equal(t, "2016-05-15", mostUniqueBeersPerDay.Date)
}
