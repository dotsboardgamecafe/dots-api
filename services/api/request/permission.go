package request

import (
	"dots-api/lib/array"
	"dots-api/lib/utils"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type (
	PermissionReq struct {
		Name         string `json:"name" validate:"required"`
		RoutePattern string `json:"route_pattern" validate:"required"`
		RouteMethod  string `json:"route_method" validate:"required"`
		Description  string `json:"description" validate:"required"`
		Status       string `json:"status" validate:"required"`
	}

	PermissionUpdateReq struct {
		Name         string `json:"name" validate:"required"`
		RoutePattern string `json:"route_pattern" validate:"required"`
		RouteMethod  string `json:"route_method" validate:"required"`
		Description  string `json:"description" validate:"required"`
		Status       string `json:"status" validate:"required"`
	}

	PermissionParam struct {
		Page    int    `json:"page"`
		Limit   int    `json:"limit"`
		Offset  int    `json:"offset"`
		Count   int    `json:"count"`
		MaxPage int    `json:"max_page"`
		Sort    string `json:"sort"`
		Order   string `json:"order"`
		Keyword string `json:"keyword"`
		Status  string `json:"status"`
	}
)

func (param *PermissionParam) ParsePermission(values url.Values) error {
	param.Keyword = ""
	param.Page = 1
	param.Limit = 10
	param.Sort = "desc"
	param.Order = "id"
	param.Status = ""
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
		if exist, _ := arrStr.InArray(order[0], []string{"id", "name", "status", "created_date"}); exist {
			param.Order = order[0]
		}
	}

	if status, ok := values["status"]; ok && len(status) > 0 {
		if !utils.Contains(utils.StatusPermission, status[0]) {
			return fmt.Errorf("%s", "wrong status value for permission(active|inactive")
		}
		param.Status = status[0]
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
