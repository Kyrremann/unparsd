package statistics

import (
	"gorm.io/gorm"
)

// DayOfWeekStat holds the check-in count for a single day of the week.
type DayOfWeekStat struct {
	// Day is the full English name of the weekday (Monday … Sunday).
	Day   string `json:"day"`
	Count int    `json:"count"`
}

// DayOfWeekStats returns check-in counts grouped by day of the week,
// ordered Monday through Sunday.
// Pass a non-empty year to restrict the calculation to that year only.
func DayOfWeekStats(db *gorm.DB, year string) ([]DayOfWeekStat, error) {
	// SQLite strftime('%w', …) returns 0=Sunday … 6=Saturday.
	type row struct {
		DayNum int `gorm:"column:day_num"`
		Count  int `gorm:"column:count"`
	}
	var rows []row
	tx := db.
		Table("checkins").
		Select("CAST(strftime('%w', checkin_at) AS INTEGER) as day_num, count(*) as count").
		Group("day_num").
		Order("day_num ASC")
	if year != "" {
		tx = tx.Where("strftime('%Y', checkin_at) = ?", year)
	}
	res := tx.Find(&rows)
	if res.Error != nil {
		return nil, res.Error
	}

	// Build a map so we can fill in missing days with 0.
	countByNum := make(map[int]int, 7)
	for _, r := range rows {
		countByNum[r.DayNum] = r.Count
	}

	// Emit Monday (1) … Saturday (6) … Sunday (0), i.e. ISO week order.
	isoOrder := []struct {
		num  int
		name string
	}{
		{1, "Monday"},
		{2, "Tuesday"},
		{3, "Wednesday"},
		{4, "Thursday"},
		{5, "Friday"},
		{6, "Saturday"},
		{0, "Sunday"},
	}

	stats := make([]DayOfWeekStat, 0, 7)
	for _, d := range isoOrder {
		stats = append(stats, DayOfWeekStat{Day: d.name, Count: countByNum[d.num]})
	}
	return stats, nil
}
