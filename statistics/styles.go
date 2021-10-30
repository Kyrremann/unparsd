package statistics

import (
	"github.com/kyrremann/unparsd/models"
	"github.com/kyrremann/unparsd/parsing"
	"gorm.io/gorm"
)

type DistinctStyle struct {
	Type     string
	Distinct int
	Total    int
}

type style struct {
	Type  string
	Total int
}

func contains(list []string, el string) bool {
	for _, v := range list {
		if v == el {
			return true
		}
	}

	return false
}

func intersection(a, b []string) []string {
	var c []string

	for _, el := range a {
		if !contains(b, el) {
			c = append(c, el)
		}
	}

	return c
}

func MissingStyles(db *gorm.DB) ([]string, error) {
	var allStyles []string
	err := parsing.ParseJSON("../fixture/all_styles.json", &allStyles)
	if err != nil {
		return nil, err
	}

	var styles []string
	res := db.Model(&models.Beer{}).Distinct("Type").Pluck("Type", &styles)
	if res.Error != nil {
		return nil, res.Error
	}

	return intersection(allStyles, styles), nil
}

func DistinctStyles(db *gorm.DB) ([]DistinctStyle, error) {
	var styles []DistinctStyle
	var distinctive []style
	res := db.Model(&models.Beer{}).Select("Type, count(Type) as total").Group("Type").Find(&distinctive)
	if res.Error != nil {
		return nil, res.Error
	}

	for _, d := range distinctive {
		styles = append(styles, DistinctStyle{Type: d.Type, Distinct: d.Total})
	}

	var checkins []style
	res = db.Model(&models.Checkin{}).Select("checkins.beer").Joins("Beer").Select("Beer.Type, count(Beer.Type) as total").Group("Beer.Type").Find(&checkins)
	if res.Error != nil {
		return nil, res.Error
	}

	for _, c := range checkins {
		for i, style := range styles {
			if style.Type == c.Type {
				style.Total = c.Total
				styles[i] = style
				break
			}
		}
	}

	return styles, nil
}
