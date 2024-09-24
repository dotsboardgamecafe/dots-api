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
	AdminReq struct {
		Email       string `json:"email" validate:"required,max=100"`
		Name        string `json:"name" validate:"required,max=100"`
		UserName    string `json:"username" validate:"max=15"`
		Password    string `json:"password" validate:"required,max=100"`
		Status      string `json:"status" validate:"required,max=50"`
		PhoneNumber string `json:"phone_number" validate:"max=15"`
		ImageUrl    string `json:"image_url"`
	}

	AdminUpdateReq struct {
		Email       string `json:"email" validate:"max=100"`
		Name        string `json:"name" validate:"max=100"`
		UserName    string `json:"username" validate:"max=15"`
		Password    string `json:"password" validate:"max=100"`
		Status      string `json:"status" validate:"max=50"`
		PhoneNumber string `json:"phone_number" validate:"max=15"`
		ImageUrl    string `json:"image_url"`
	}

	UpdateStatusAdminReq struct {
		Status string `json:"status" validate:"required"`
	}

	AdminParam struct {
		Page    int    `json:"page"`
		Limit   int    `json:"limit"`
		Offset  int    `json:"offset"`
		Count   int    `json:"count"`
		Sort    string `json:"sort"`
		Order   string `json:"order"`
		Keyword string `json:"keyword"`
		Status  string `json:"status"`
	}
)

func (param *AdminParam) ParseAdmin(values url.Values) error {
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
		if exist, _ := arrStr.InArray(order[0], []string{"id", "email", "name", "status", "created_date"}); exist {
			param.Order = order[0]
		}
	}

	if status, ok := values["status"]; ok && len(status) > 0 {
		if !utils.Contains(utils.StatusAdmin, status[0]) {
			return fmt.Errorf("%s", "wrong status value for admin(active|inactive")
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
