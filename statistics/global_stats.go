package statistics

import (
	"math"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type GlobalStats struct {
	Checkins      int                    `json:"checkins"`
	UniqueBeers   int                    `json:"unique_beers"`
	StartDate     string                 `json:"start_date"`
	DaysDrinking  int                    `gorm:"-" json:"days_drinking"`
	BeersPerDay   float64                `gorm:"-" json:"beers_per_day"`
	Periods       map[string]PeriodStats `gorm:"-" json:"years"`
	GeneratedDate time.Time
}

type MostPerDay struct {
	Count int    `json:"count"`
	Date  string `json:"date"`
}

type PeriodStats struct {
	Checkins              int           `json:"checkins"`
	UniqueBreweries       int           `json:"unique_breweries"`
	BreweryCountries      int           `json:"brewery_countries"`
	UniqueVenues          int           `json:"unique_venues"`
	VenueCountries        int           `json:"venue_countries"`
	UniqueBeers           int           `json:"unique_beers"`
	MaxAbv                float64       `json:"max_abv"`
	AvgAbv                float64       `json:"avg_abv"`
	Styles                int           `json:"styles"`
	StartDate             string        `json:"start_date"`
	Month                 string        `json:"month,omitempty"`
	Year                  string        `json:"year"`
	Months                []PeriodStats `gorm:"-" json:"months,omitempty"`
	MostCheckinsPerDay    MostPerDay    `gorm:"-" json:"most_checkins_per_day"`
	MostUniqueBeersPerDay MostPerDay    `gorm:"-" json:"most_unique_beers_per_day"`
	BeersPerDay           float64       `gorm:"-" json:"beers_per_day"`
}

func getMonthAsString(month string) (string, error) {
	intMonth, err := strconv.Atoi(month)
	if err != nil {
		return "", err
	}
	return time.Month(intMonth).String(), nil
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

func daysTillNowInYear() int {
	return time.Now().YearDay()
}

func daysInYear(year int) int {
	return time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
}

func periodStats(db *gorm.DB, groupBy, year string) ([]PeriodStats, error) {
	var periodStats []PeriodStats
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

	res := tx.Find(&periodStats)

	return periodStats, res.Error
}

func monthlyStats(db *gorm.DB, year string) ([]PeriodStats, error) {
	monthly, err := periodStats(db, "month", year)
	if err != nil {
		return nil, err
	}

	for i, ps := range monthly {
		daysInMonth, err := daysInMonth(year, ps.Month)
		if err != nil {
			return nil, err
		}
		monthly[i].BeersPerDay = math.Round((float64(ps.Checkins)/float64(daysInMonth))*100.00) / 100.00

		stringMonth, err := getMonthAsString(ps.Month)
		if err != nil {
			return nil, err
		}
		monthly[i].Month = stringMonth
	}

	return monthly, nil
}

func yearlyStats(db *gorm.DB) ([]PeriodStats, error) {
	yearly, err := periodStats(db, "year", "")
	if err != nil {
		return nil, err
	}

	for i, ps := range yearly {
		intYear, err := strconv.Atoi(ps.Year)
		if err != nil {
			return nil, err
		}

		monthly, err := monthlyStats(db, ps.Year)
		if err != nil {
			return nil, err
		}
		yearly[i].Months = monthly
		yearly[i].Month = ""

		var days int
		if time.Now().Year() == intYear {
			days = daysTillNowInYear()
		} else {
			days = daysInYear(intYear)
		}
		yearly[i].BeersPerDay = math.Round((float64(ps.Checkins)/float64(days))*100.00) / 100.00
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

	periods, err := yearlyStats(db)
	if err != nil {
		return GlobalStats{}, err
	}

	globalStat.Periods = make(map[string]PeriodStats)
	for _, ps := range periods {
		globalStat.Periods[ps.Year] = ps
	}

	daysDrinking, err := daysSince(globalStat.StartDate)
	if err != nil {
		return GlobalStats{}, err
	}
	globalStat.DaysDrinking = daysDrinking
	globalStat.BeersPerDay = math.Round((float64(globalStat.Checkins)/float64(daysDrinking))*100.00) / 100.00
	globalStat.GeneratedDate = time.Now()

	return globalStat, nil
}
