package request

import (
	"dots-api/lib/array"
	"net/url"
	"strconv"
	"strings"
)

type (
	SettingReq struct {
		SetGroup     string `json:"set_group" validate:"required,max=50"`
		SetLabel     string `json:"set_label" validate:"required,max=100"`
		SetOrder     int    `json:"set_order" validate:"required,max=100"`
		ContentType  string `json:"content_type" validate:"required,max=10"`
		ContentValue string `json:"content_value" validate:"required"`
		IsActive     bool   `json:"is_active"`
	}

	SettingParam struct {
		Page     int    `json:"page"`
		Limit    int    `json:"limit"`
		Offset   int    `json:"offset"`
		Count    int    `json:"count"`
		Sort     string `json:"sort"`
		Order    string `json:"order"`
		Keyword  string `json:"keyword"`
		IsActive string `json:"is_active"`
		SetGroup string `json:"set_group"`
	}
)

func (param *SettingParam) ParseSetting(values url.Values) error {
	param.Keyword = ""
	param.Page = 1
	param.Limit = 10
	param.Sort = "desc"
	param.Order = "id"
	param.IsActive = ""
	param.SetGroup = ""
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
		if exist, _ := arrStr.InArray(order[0], []string{"id", "set_label", "set_order"}); exist {
			param.Order = order[0]
		}
	}

	if isActive, ok := values["is_active"]; ok && len(isActive) > 0 {
		param.IsActive = isActive[0]
	}

	if setGroup, ok := values["set_group"]; ok && len(setGroup) > 0 {
		param.SetGroup = setGroup[0]
	}

	if keyword, ok := values["keyword"]; ok && len(keyword) > 0 {
		param.Keyword = keyword[0]
	}

	if limit, ok := values["limit"]; ok && len(limit) > 0 {
		if l, err := strconv.Atoi(limit[0]); err == nil {
			param.Limit = l
		}
	}

	param.Offset = (param.Page - 1) * param.Limit

	return nil
}
