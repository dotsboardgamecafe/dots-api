package request

import (
	"dots-api/lib/array"
	"net/url"
	"strconv"
	"strings"
)

type (
	UserGameCollectionParam struct {
		Page    int    `json:"page"`
		Limit   int    `json:"limit"`
		Offset  int    `json:"offset"`
		Count   int    `json:"count"`
		Sort    string `json:"sort"`
		Order   string `json:"order"`
		MaxPage int    `json:"max_page"`
	}
)

func (param *UserGameCollectionParam) ParseUserGameCollection(values url.Values) error {

	param.Page = 1
	param.Limit = 10
	param.Sort = "desc"
	param.Order = "game_id"
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
		if exist, _ := arrStr.InArray(order[0], []string{"game_id"}); exist {
			param.Order = order[0]
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
