package request

import (
	"dots-api/lib/utils"

	"fmt"
	"net/url"
	"strconv"
)

type (
	HallOfFameParam struct {
		Limit int `json:"limit"`
		Count int `json:"count"`
		Year  int `json:"year"`
	}

	MonthlyTopAchieverParam struct {
		Limit    int    `json:"limit"`
		Count    int    `json:"count"`
		Category string `json:"category"`
		Month    int    `json:"month" validate:"min=1,max=12"`
		Year     int    `json:"year"`
		CafeCity string `json:"cafe_city"`
	}
)

func (param *HallOfFameParam) ParseHallOfFame(values url.Values) error {
	param.Limit = 10

	if year, ok := values["year"]; ok && len(year) > 0 {
		if y, err := strconv.Atoi(year[0]); err == nil && y > 1 {
			param.Year = y
		}
	}

	if limit, ok := values["limit"]; ok && len(limit) > 0 {
		if l, err := strconv.Atoi(limit[0]); err == nil {
			param.Limit = l
		}
	}

	return nil
}

func (param *MonthlyTopAchieverParam) ParseMonthlyTopAchiever(values url.Values) error {
	// Default limit 10; category -> VP
	param.Category = "vp"
	param.Limit = 10

	// Get period of month and year
	month, monthParamExist := values["month"]
	year, yearParamExist := values["year"]
	isPeriodParamExist := (monthParamExist && len(month) > 0) && (yearParamExist && len(year) > 0)

	if isPeriodParamExist {
		if m, err := strconv.Atoi(month[0]); err == nil && m > 0 && m < 13 {
			param.Month = m
		}

		if y, err := strconv.Atoi(year[0]); err == nil {
			param.Year = y
		}
	}

	if categoryType, ok := values["category"]; ok && len(categoryType) > 0 {
		if !utils.Contains(utils.MonthlyTopAchieverCategory, categoryType[0]) {
			return fmt.Errorf("%s", "wrong type value for category(vp|unique_game)")
		}
		param.Category = categoryType[0]
	}

	if cafeCity, ok := values["cafe_city"]; ok && len(cafeCity) > 0 {
		param.CafeCity = cafeCity[0]
	}

	if limit, ok := values["limit"]; ok && len(limit) > 0 {
		if l, err := strconv.Atoi(limit[0]); err == nil {
			param.Limit = l
		}
	}

	return nil
}
