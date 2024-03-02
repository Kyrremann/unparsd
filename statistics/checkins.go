package statistics

import (
	"gorm.io/gorm"
)

type DailyCheckins struct {
	Year           int `json:"year"`
	Month          int `json:"month"`
	Day            int `json:"day"`
	Checkins       int `json:"checkins"`
	UniqueCheckins int `json:"unique_checkins"`
}

func CheckinsByDay(db *gorm.DB) ([]DailyCheckins, error) {
	var checkinsByDay []DailyCheckins
	tx := db.
		Table("checkins").
		Select("strftime('%Y-%m-%d', checkins.checkin_at) as date," +
			"strftime('%d', checkins.checkin_at) as day," +
			"strftime('%m', checkins.checkin_at) as month," +
			"strftime('%Y', checkins.checkin_at) as year," +
			"count(checkins.id) as checkins," +
			"count(DISTINCT(beers.id)) as unique_checkins").
		Joins("INNER JOIN beers ON checkins.beer_id = beers.id").
		Group("date")

	res := tx.Find(&checkinsByDay)

	return checkinsByDay, res.Error
}
