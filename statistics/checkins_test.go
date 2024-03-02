package statistics

import (
	"testing"

	"github.com/kyrremann/unparsd/parsing"
	"github.com/stretchr/testify/assert"
)

func TestCheckins(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/untappd.json")
	assert.NoError(t, err)

	daily, err := CheckinsByDay(db)
	assert.NoError(t, err)
	assert.Len(t, daily, 72)

	checkins := 0
	for _, d := range daily {
		checkins += d.Checkins
	}
	assert.Equal(t, 126, checkins)
}
