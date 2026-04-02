package statistics

import (
	"testing"

	"github.com/kyrremann/unparsd/parsing"
	"github.com/stretchr/testify/assert"
)

func TestDayOfWeekStats(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/checkins")
	assert.NoError(t, err)

	stats, err := DayOfWeekStats(db, "")
	assert.NoError(t, err)
	assert.Len(t, stats, 7)

	// Verify the order is Monday … Sunday.
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for i, s := range stats {
		assert.Equal(t, days[i], s.Day)
	}

	// Total check-ins across all days must equal the fixture total.
	total := 0
	for _, s := range stats {
		total += s.Count
	}
	assert.Equal(t, 126, total)
}

func TestCheckinStreak(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/checkins")
	assert.NoError(t, err)

	streak, err := CheckinStreak(db, "")
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, streak.Longest, 1)
	assert.NotEmpty(t, streak.LongestStart)
	assert.NotEmpty(t, streak.LongestEnd)
	// LongestStart <= LongestEnd (lexicographic date comparison is safe for YYYY-MM-DD).
	assert.LessOrEqual(t, streak.LongestStart, streak.LongestEnd)
}

func TestABVDistribution(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/checkins")
	assert.NoError(t, err)

	buckets, err := ABVDistribution(db, "")
	assert.NoError(t, err)
	assert.Len(t, buckets, 5)

	total := 0
	for _, b := range buckets {
		total += b.Count
	}
	assert.Equal(t, 126, total)
}

func TestRatingDeltas(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/checkins")
	assert.NoError(t, err)

	deltas, err := RatingDeltas(db)
	assert.NoError(t, err)

	// All returned beers must have at least one rated check-in.
	for _, d := range deltas {
		assert.Greater(t, d.AvgPersonal, 0.0)
	}

	// Result must be sorted by |delta| descending.
	for i := 1; i < len(deltas); i++ {
		prev := deltas[i-1].Delta
		if prev < 0 {
			prev = -prev
		}
		curr := deltas[i].Delta
		if curr < 0 {
			curr = -curr
		}
		assert.GreaterOrEqual(t, prev, curr)
	}
}

func TestTopVenues(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/checkins")
	assert.NoError(t, err)

	venues, err := TopVenues(db)
	assert.NoError(t, err)

	if len(venues) == 0 {
		// Fixture may have no venue data — that is acceptable.
		return
	}

	// Must be sorted by check-in count descending.
	for i := 1; i < len(venues); i++ {
		assert.GreaterOrEqual(t, venues[i-1].Checkins, venues[i].Checkins)
	}

	for _, v := range venues {
		assert.NotEmpty(t, v.Name)
		assert.Greater(t, v.Checkins, 0)
	}
}

func TestServingTypeStats(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/checkins")
	assert.NoError(t, err)

	stats, err := ServingTypeStats(db)
	assert.NoError(t, err)

	for _, s := range stats {
		assert.NotEmpty(t, s.ServingType)
		assert.Greater(t, s.Count, 0)
	}
}

func TestFlavorProfileStats(t *testing.T) {
	db, err := parsing.LoadJsonIntoDatabase("../fixture/checkins")
	assert.NoError(t, err)

	stats, err := FlavorProfileStats(db)
	assert.NoError(t, err)

	for _, s := range stats {
		assert.NotEmpty(t, s.Flavor)
		assert.Greater(t, s.Count, 0)
	}

	// Must be sorted by count descending.
	for i := 1; i < len(stats); i++ {
		assert.GreaterOrEqual(t, stats[i-1].Count, stats[i].Count)
	}
}
