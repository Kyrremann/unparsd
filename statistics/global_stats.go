package statistics

import (
	"math"
	"strconv"
	"time"

	"github.com/kyrremann/unparsd/models"
	"gorm.io/gorm"
)

type GlobalStats struct {
	Checkins     int            `json:"checkins"`
	UniqueBeers  int            `json:"unique_beers"`
	StartDate    string         `json:"start_date"`
	DaysDrinking int            `gorm:"-" json:"days_drinking"`
	BeersPerDay  float64        `gorm:"-" json:"beers_per_day"`
	Periodes     []PeriodeStats `gorm:"-" json:"years"`
}

type MostPerDay struct {
	Count int    `json:"count"`
	Date  string `json:"date"`
}

type PeriodeStats struct {
	Checkins              int            `json:"checkins"`
	UniqueBreweries       int            `json:"unique_breweries"`
	BreweryCountries      int            `json:"brewery_countries"`
	UniqueVenues          int            `json:"unique_venues"`
	VenueCountries        int            `json:"venue_countries"`
	UniqueBeers           int            `json:"unique_beers"`
	MaxAbv                float64        `json:"max_abv"`
	AvgAbv                float64        `json:"avg_abv"`
	Styles                int            `json:"styles"`
	StartDate             string         `json:"start_date"`
	Month                 *string        `json:"month,omitempty"`
	Year                  string         `json:"year"`
	Months                []PeriodeStats `gorm:"-" json:"months,omitempty"`
	MostCheckinsPerDay    MostPerDay     `gorm:"-" json:"most_checkins_per_day"`
	MostUniqueBeersPerDay MostPerDay     `gorm:"-" json:"most_unique_beers_per_day"`
	BeersPerDay           float64        `gorm:"-" json:"beers_per_day"`
}

func daysInMonth(year, month string) (int, error) {
	intYear, err := strconv.Atoi(year)
	if err != nil {
		return 0, err
	}

	intMonth, err := strconv.Atoi(month)
	if err != nil {
		return 0, err
	}

	return time.Date(intYear, time.Month(intMonth)+1, 0, 0, 0, 0, 0, time.UTC).Day(), nil
}

func daysInYear(year string) (int, error) {
	intYear, err := strconv.Atoi(year)
	if err != nil {
		return 0, err
	}

	return time.Date(intYear, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay(), nil
}

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

func periodeStats(db *gorm.DB, groupBy, year string) ([]PeriodeStats, error) {
	var periodeStats []PeriodeStats
	tx := db.
		Table("checkins").
		Select("count(checkins.id) as checkins," +
			"count(DISTINCT(breweries.country)) as brewery_countries," +
			"count(DISTINCT(beers.brewery_id)) as unique_breweries," +
			"count(DISTINCT(beers.id)) as unique_beers," +
			"count(DISTINCT(beers.type)) as styles," +
			"ROUND(max(beers.abv), 2) as max_abv," +
			"ROUND(avg(beers.abv), 2) as avg_abv," +
			"count(DISTINCT(checkins.venue_name)) as unique_venues," +
			"count(DISTINCT(venues.country)) as venue_countries," +
			"date(min(checkins.checkin_at)) as start_date," +
			"strftime('%m', checkins.checkin_at) as month," +
			"strftime('%Y', checkins.checkin_at) as year").
		Joins("LEFT JOIN beers ON checkins.beer_id == beers.id").
		Joins("LEFT JOIN breweries ON beers.brewery_id == breweries.id").
		Joins("LEFT JOIN venues ON checkins.venue_name == venues.name").
		Group(groupBy)

	if len(year) > 0 {
		tx.Where("year", year)
	}

	res := tx.Find(&periodeStats)

	return periodeStats, res.Error
}

func monthlyStats(db *gorm.DB, year string) ([]PeriodeStats, error) {
	monthly, err := periodeStats(db, "month", year)
	if err != nil {
		return nil, err
	}

	for i, ps := range monthly {
		mostCheckinsPerDay, err := MostCheckinsPerDay(db, year, *ps.Month)
		if err != nil {
			return nil, err
		}

		monthly[i].MostCheckinsPerDay = mostCheckinsPerDay
		mostUniqueBeersPerDay, err := MostUniqueBeersPerDay(db, year, *ps.Month)
		if err != nil {
			return nil, err
		}
		monthly[i].MostUniqueBeersPerDay = mostUniqueBeersPerDay

		daysInMonth, err := daysInMonth(year, *ps.Month)
		if err != nil {
			return nil, err
		}
		monthly[i].BeersPerDay = math.Round((float64(ps.Checkins)/float64(daysInMonth))*100.00) / 100.00
	}

	return monthly, nil
}

func yearlyStats(db *gorm.DB) ([]PeriodeStats, error) {
	yearly, err := periodeStats(db, "year", "")
	if err != nil {
		return nil, err
	}

	for i, ps := range yearly {
		monthly, err := monthlyStats(db, ps.Year)
		if err != nil {
			return nil, err
		}
		yearly[i].Months = monthly
		yearly[i].Month = nil

		mostCheckinsPerDay, err := MostCheckinsPerDay(db, ps.Year, "")
		if err != nil {
			return nil, err
		}

		yearly[i].MostCheckinsPerDay = mostCheckinsPerDay
		mostUniqueBeersPerDay, err := MostUniqueBeersPerDay(db, ps.Year, "")
		if err != nil {
			return nil, err
		}
		yearly[i].MostUniqueBeersPerDay = mostUniqueBeersPerDay

		daysInYear, err := daysInYear(ps.Year)
		if err != nil {
			return nil, err
		}
		yearly[i].BeersPerDay = math.Round((float64(ps.Checkins)/float64(daysInYear))*100.00) / 100.00
	}

	return yearly, nil
}

func daysSince(dateString string) (int, error) {
	date, err := time.Parse("2006-02-03", dateString)
	if err != nil {
		return 0, err
	}

	return int(time.Since(date).Hours() / 24), nil
}

func AllMyStats(db *gorm.DB) (GlobalStats, error) {
	var globalStat GlobalStats
	res := db.
		Table("checkins").
		Select("count(checkins.id) as checkins," +
			"count(DISTINCT(beers.id)) as unique_beers," +
			"date(min(checkins.checkin_at)) as start_date").
		Joins("INNER JOIN beers ON checkins.beer_id == beers.id").
		Find(&globalStat)

	if res.Error != nil {
		return GlobalStats{}, res.Error
	}

	if res.Error != nil {
		return GlobalStats{}, res.Error
	}

	periodes, err := yearlyStats(db)
	if err != nil {
		return GlobalStats{}, err
	}
	globalStat.Periodes = periodes

	daysDrinking, err := daysSince(globalStat.StartDate)
	if err != nil {
		return GlobalStats{}, err
	}
	globalStat.DaysDrinking = daysDrinking
	globalStat.BeersPerDay = math.Round((float64(globalStat.Checkins)/float64(daysDrinking))*100.00) / 100.00

	return globalStat, nil
}
