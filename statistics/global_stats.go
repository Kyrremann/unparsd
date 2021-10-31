package statistics

import (
	"strconv"

	"gorm.io/gorm"
)

// TODO: Days drinking
type GlobalStats struct {
	Checkins     int            `json:"checkins"`
	UniqueBeers  int            `json:"unique_beers"`
	StartDate    string         `json:"start_date"`
	DaysDrinking int            `json:"days_drinking"`
	Periodes     []PeriodeStats `gorm:"-" json:"years"`
}

// TODO: most_per_day
type PeriodeStats struct {
	Checkins         int            `json:"checkins"`
	Breweries        int            `json:"breweries"`
	BreweryCountries int            `json:"brewery_countries"`
	Venues           int            `json:"venues"`
	VenueCountries   int            `json:"venue_countries"`
	Beers            int            `json:"beers"`
	MaxAbv           float64        `json:"max_abv"`
	AvgAbv           float64        `json:"avg_abv"`
	Styles           int            `json:"styles"`
	StartDate        string         `json:"start_date"`
	Month            *int           `json:"month,omitempty"`
	Year             int            `json:"year"`
	Months           []PeriodeStats `gorm:"-" json:"months,omitempty"`
}

func periodeStats(db *gorm.DB, groupBy, year string) ([]PeriodeStats, error) {
	var periodeStats []PeriodeStats
	tx := db.
		Table("checkins").
		Select("count(checkins.id) as checkins," +
			"count(DISTINCT(breweries.country)) as brewery_countries," +
			"count(DISTINCT(beers.brewery_id)) as breweries," +
			"count(DISTINCT(beers.id)) as beers," +
			"count(DISTINCT(beers.type)) as styles," +
			"max(beers.abv) as max_abv," +
			"ROUND(avg(beers.abv), 2) as avg_abv," +
			"count(DISTINCT(checkins.venue_id)) as venues," +
			"count(DISTINCT(venues.country)) as venue_countries," +
			"date(min(checkins.checkin_at)) as start_date," +
			"strftime('%m', checkins.checkin_at) as month," +
			"strftime('%Y', checkins.checkin_at) as year").
		Joins("INNER JOIN breweries ON beers.brewery_id == breweries.id").
		Joins("INNER JOIN beers ON checkins.beer_id == beers.id").
		Joins("INNER JOIN venues ON checkins.venue_id == venues.name").
		Group(groupBy)

	if len(year) > 0 {
		tx.Where("year", year)
	}

	res := tx.Find(&periodeStats)

	return periodeStats, res.Error
}

func monthlyStats(db *gorm.DB, year int) ([]PeriodeStats, error) {
	return periodeStats(db, "month", strconv.Itoa(year))
}

func yearlyStats(db *gorm.DB) ([]PeriodeStats, error) {
	yearly, err := periodeStats(db, "year", "")
	if err != nil {
		return nil, err
	}

	for i, year := range yearly {
		monthly, err := monthlyStats(db, year.Year)
		if err != nil {
			return nil, err
		}
		yearly[i].Months = monthly
		yearly[i].Month = nil
	}

	return yearly, nil
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

	return globalStat, nil
}
