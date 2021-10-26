package statistics

import (
	"testing"

	"github.com/kyrremann/unparsd/parsing"
	"github.com/stretchr/testify/assert"
)

func TestMissingStyles(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/untappd.json")
	assert.NoError(t, err)

	styles, err := MissingStyles(db)
	assert.NoError(t, err)
	assert.Len(t, styles, 174)
}

func TestUniqueStyles(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/untappd.json")
	assert.NoError(t, err)

	styles, err := DistinctStyles(db)
	assert.NoError(t, err)
	assert.Len(t, styles, 40)
}
