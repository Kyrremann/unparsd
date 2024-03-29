package statistics

import (
	"testing"

	"github.com/kyrremann/unparsd/parsing"
	"github.com/stretchr/testify/assert"
)

func TestMonthly(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/untappd.json")
	assert.NoError(t, err)

	monthly, err := GetMonthlyBannerSummary(db)
	assert.NoError(t, err)
	assert.Len(t, monthly, 5)
	assert.Equal(t, 2016, monthly[0].Year)
}
