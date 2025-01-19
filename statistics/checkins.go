package statistics

import (
	"github.com/kyrremann/unparsd/models"
	"gorm.io/gorm"
)

func MostCheckinsPerDay(db *gorm.DB, year, month string) (MostPerDay, error) {
	return mostPerDay(db, "count(id) as count,", year, month)
}

func MostUniqueBeersPerDay(db *gorm.DB, year, month string) (MostPerDay, error) {
	return mostPerDay(db, "count(DISTINCT(beer_id)) as count,", year, month)
}

func mostPerDay(db *gorm.DB, selectCounter, year, month string) (MostPerDay, error) {
	var mostPerDay MostPerDay
	tx := db.
		Model(models.Checkin{}).
		Select(selectCounter +
			"strftime('%Y-%m-%d', checkins.checkin_at) as date," +
			"strftime('%Y', checkins.checkin_at) as year," +
			"strftime('%m', checkins.checkin_at) as month").
		Group("date").
		Order("count DESC").
		Order("date ASC")

	if len(month) > 0 {
		tx.Where("year = ? AND month = ?", year, month)
	} else {
		tx.Where("year = ?", year)
	}

	res := tx.First(&mostPerDay)

	return mostPerDay, res.Error
}
