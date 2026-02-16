package request

import (
	"dots-api/lib/array"
	"net/url"
	"strconv"
	"strings"
)

type UserPointHistoryParam struct {
	Page       int    `json:"page"`
	MaxPage    int    `json:"max_page"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	Count      int    `json:"count"`
	Sort       string `json:"sort"`
	Order      string `json:"order"`
	SourceType string `json:"source_type"`
}

func (param *UserPointHistoryParam) Parse(values url.Values) error {
	param.Page = 1
	param.Limit = 10
	param.Sort = "asc"
	param.Order = "created_date"
	param.Offset = 0

	if page, ok := values["page"]; ok && len(page) > 0 {
		if p, err := strconv.Atoi(page[0]); err == nil && p > 1 {
			param.Page = p
		}
	}

	if sort, ok := values["sort"]; ok && len(sort) > 0 && strings.ToLower(sort[0]) == "asc" {
		param.Sort = "asc"
	}

	if order, ok := values["order"]; ok && len(order) > 0 {
		arrStr := new(array.ArrStr)
		orderList := []string{"created_date", "point"}

		if exist, _ := arrStr.InArray(order[0], orderList); exist {
			param.Order = order[0]
		} else {
			param.Order = "created_date"
		}
	}

	if sourceType, ok := values["source_type"]; ok && len(sourceType) > 0 {
		arrStr := new(array.ArrStr)
		sourceTypeList := []string{
			"room", "room_play", "room_paid", "tournament", "tournament_play",
			"tournament_paid", "redeem", "badge", "game", "profile", "point", "tier",
		}
		if exist, _ := arrStr.InArray(sourceType[0], sourceTypeList); exist {
			param.SourceType = sourceType[0]
		}
	}

	if limit, ok := values["limit"]; ok && len(limit) > 0 {
		if l, err := strconv.Atoi(limit[0]); err == nil {
			param.Limit = l
		}
	}

	param.Offset = (param.Page - 1) * param.Limit

	return nil
}

type UserPointAdjustmentRequest struct {
	AdjustmentType string  `json:"adjustment_type" validate:"required,oneof=add subtract"`
	Point          int     `json:"point" validate:"required,min=1,max=10000"`
	Description    *string `json:"description,omitempty"`
}
