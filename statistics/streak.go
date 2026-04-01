package statistics

import (
	"sort"
	"time"

	"gorm.io/gorm"
)

// StreakStats contains the longest and current consecutive-day check-in streaks.
type StreakStats struct {
	Longest      int    `json:"longest"`
	LongestStart string `json:"longest_start"`
	LongestEnd   string `json:"longest_end"`
	Current      int    `json:"current"`
	CurrentStart string `json:"current_start"`
}

// CheckinStreak computes the longest and current check-in streaks (consecutive calendar days).
// Pass a non-empty year to restrict the calculation to that year only.
func CheckinStreak(db *gorm.DB, year string) (StreakStats, error) {
	// Fetch all distinct check-in dates in ascending order.
	var dates []string
	tx := db.
		Table("checkins").
		Distinct("strftime('%Y-%m-%d', checkin_at) as d").
		Order("d ASC")
	if year != "" {
		tx = tx.Where("strftime('%Y', checkin_at) = ?", year)
	}
	res := tx.Pluck("d", &dates)
	if res.Error != nil {
		return StreakStats{}, res.Error
	}
	if len(dates) == 0 {
		return StreakStats{}, nil
	}

	sort.Strings(dates)

	parse := func(s string) time.Time {
		t, _ := time.Parse("2006-01-02", s)
		return t
	}

	// Walk through dates and find the longest streak.
	streakStart := dates[0]
	streakLen := 1

	bestLen := 1
	bestStart := dates[0]
	bestEnd := dates[0]

	for i := 1; i < len(dates); i++ {
		prev := parse(dates[i-1])
		curr := parse(dates[i])
		if curr.Sub(prev) == 24*time.Hour {
			streakLen++
		} else {
			if streakLen > bestLen {
				bestLen = streakLen
				bestStart = streakStart
				bestEnd = dates[i-1]
			}
			streakStart = dates[i]
			streakLen = 1
		}
	}
	// Final segment.
	if streakLen > bestLen {
		bestLen = streakLen
		bestStart = streakStart
		bestEnd = dates[len(dates)-1]
	}

	// Current streak: last check-in must be today or yesterday.
	today := time.Now().Truncate(24 * time.Hour)
	currentLen := 0
	currentStart := ""

	lastDate := parse(dates[len(dates)-1])
	daysFromToday := int(today.Sub(lastDate).Hours() / 24)
	if daysFromToday <= 1 {
		currentLen = 1
		currentStart = dates[len(dates)-1]
		for i := len(dates) - 2; i >= 0; i-- {
			prev := parse(dates[i+1])
			curr := parse(dates[i])
			if prev.Sub(curr) == 24*time.Hour {
				currentLen++
				currentStart = dates[i]
			} else {
				break
			}
		}
	}

	return StreakStats{
		Longest:      bestLen,
		LongestStart: bestStart,
		LongestEnd:   bestEnd,
		Current:      currentLen,
		CurrentStart: currentStart,
	}, nil
}
