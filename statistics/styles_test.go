package statistics

import (
	"testing"

	"github.com/kyrremann/unparsd/v4/parsing"
	"github.com/stretchr/testify/assert"
)

func TestDistinctStyles(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/checkins")
	assert.NoError(t, err)

	styles, err := DistinctStyles(db)
	assert.NoError(t, err)
	assert.Len(t, styles, 57)
}
