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

func TestDaysInMonth(t *testing.T) {
	days, err := daysInMonth("2016", "02")
	assert.NoError(t, err)
	assert.Equal(t, 29, days)

	days, err = daysInMonth("2017", "02")
	assert.NoError(t, err)
	assert.Equal(t, 28, days)
}

func TestDaysInYear(t *testing.T) {
	days, err := daysInYear("2016")
	assert.NoError(t, err)
	assert.Equal(t, 366, days)

	days, err = daysInYear("2017")
	assert.NoError(t, err)
	assert.Equal(t, 365, days)
}
