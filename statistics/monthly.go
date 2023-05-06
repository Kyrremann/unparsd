package statistics

import (
	"math"
	"time"

	"github.com/kyrremann/unparsd/models"
	"gorm.io/gorm"
)

type Monthly struct {
	Checkins    int
	StartDay    int
	StartMonth  string
	Year        int
	BeersPerDay float64
}

func GetMonthlyBannerSummary(db *gorm.DB) ([]Monthly, error) {
	var monthly []Monthly
	res := db.
		Model(models.Checkin{}).
		Select("count(checkins.id) as checkins," +
			"strftime('%d', date(min(checkins.checkin_at))) as start_day," +
			"strftime('%m', date(min(checkins.checkin_at))) as start_month," +
			"strftime('%Y', checkins.checkin_at) as year").
		Group("year").
		Find(&monthly)

	if res.Error != nil {
		return nil, res.Error
	}

	for i, m := range monthly {
		daysInYear := time.Date(m.Year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
		monthly[i].BeersPerDay = math.Round((float64(m.Checkins)/float64(daysInYear))*100.00) / 100.00
		month, err := getMonthAsString(m.StartMonth)
		if err != nil {
			return nil, err
		}
		monthly[i].StartMonth = month
	}

	return monthly, res.Error
}
