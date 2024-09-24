package request

import (
	"net/url"
	"strconv"
	"strings"
)

type RewardReq struct {
	TierCode     string `json:"tier_code"`
	Name         string `json:"name"`
	ImageUrl     string `json:"image_url"`
	Description  string `json:"description"`
	CategoryType string `json:"category_type"`
	Status       string `json:"status"`
	ExpiredDate  string `json:"expired_date"`
}

type RewardParam struct {
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
	Offset       int    `json:"offset"`
	Count        int    `json:"count"`
	Sort         string `json:"sort"`
	Order        string `json:"order"`
	MaxPage      int    `json:"max_page"`
	Keyword      string `json:"keyword"`
	Status       string `json:"status"`
	Name         string `json:"name"`
	CategoryType string `json:"category_type"`
	ExpiredDate  string `json:"expired_date"`
	TierName     string `json:"tier_name"`
}

func (param *RewardParam) ParseReward(values url.Values) error {
	param.Keyword = ""
	param.Page = 1
	param.Limit = 10
	param.Sort = "desc"
	param.Order = "created_date"
	param.Status = ""
	param.Offset = 0
	param.CategoryType = ""
	param.ExpiredDate = ""
	param.TierName = ""

	if page, ok := values["page"]; ok && len(page) > 0 {
		if p, err := strconv.Atoi(page[0]); err == nil && p > 0 {
			param.Page = p
		}
	}

	if sort, ok := values["sort"]; ok && len(sort) > 0 && strings.ToLower(sort[0]) == "asc" {
		param.Sort = "asc"
	}

	if order, ok := values["order"]; ok && len(order) > 0 {
		if order[0] == "created_date" || order[0] == "reward_name" || order[0] == "status" {
			param.Order = order[0]
		}
	}

	if status, ok := values["status"]; ok && len(status) > 0 {
		param.Status = status[0]
	}

	if categoryType, ok := values["category_type"]; ok && len(categoryType) > 0 {
		param.CategoryType = categoryType[0]
	}

	if expiredDate, ok := values["expired_date"]; ok && len(expiredDate) > 0 {
		param.ExpiredDate = expiredDate[0]
	}

	if tierName, ok := values["tier_name"]; ok && len(tierName) > 0 {
		param.TierName = tierName[0]
	}

	if keyword, ok := values["keyword"]; ok && len(keyword) > 0 {
		param.Keyword = keyword[0]
	}

	if limit, ok := values["limit"]; ok && len(limit) > 0 {
		if l, err := strconv.Atoi(limit[0]); err == nil && l > 0 {
			param.Limit = l
		}
	}

	param.Offset = (param.Page - 1) * param.Limit

	return nil
}
