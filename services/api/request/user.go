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
	UpdateProfileUserReq struct {
		FullName    string `json:"full_name"`
		ImageUrl    string `json:"image_url"`
		PhoneNumber string `json:"phone_number"`
	}

	UpdatePasswordReq struct {
		OldPassword     string `json:"old_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"required"`
	}

	UpdateStatusUserReq struct {
		Status string `json:"status" validate:"required"`
	}

	UpdateUserReq struct {
		FullName    string `json:"full_name"`
		ImageUrl    string `json:"image_url"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
		UserName    string `json:"username"`
		Status      string `json:"status"`
	}

	UserParam struct {
		Page       int      `json:"page"`
		Limit      int      `json:"limit"`
		Offset     int      `json:"offset"`
		Count      int      `json:"count"`
		Sort       string   `json:"sort"`
		Order      string   `json:"order"`
		Keyword    string   `json:"keyword"`
		Status     string   `json:"status"`
		LatestTier []string `json:"latest_tier"`
	}
)

func (param *UserParam) ParseUser(values url.Values) error {
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
		if order[0] == "name" {
			order[0] = "username"
		}

		arrStr := new(array.ArrStr)
		if exist, _ := arrStr.InArray(order[0], []string{"id", "email", "username", "fullname", "status", "created_date"}); exist {
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

	if latestTier, ok := values["latest_tier"]; ok && len(latestTier) > 0 {
		param.LatestTier = strings.Split(strings.ToLower(latestTier[0]), ",")
	}

	if limit, ok := values["limit"]; ok && len(limit) > 0 {
		if l, err := strconv.Atoi(limit[0]); err == nil {
			param.Limit = l
		}
	}

	param.Offset = (param.Page - 1) * param.Limit

	return nil
}
