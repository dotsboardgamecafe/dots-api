package request

import (
	"net/url"
	"strconv"
)

type (
	UserRedeemRequest struct {
		InvoiceCode string `json:"invoice_code" validate:"required,max=30"`
	}

	UserRedeemHistoryParam struct {
		Page    int `json:"page"`
		MaxPage int `json:"max_page"`
		Limit   int `json:"limit"`
		Offset  int `json:"offset"`
		Count   int `json:"count"`
	}

	UserClaimedHistoryParam struct {
		Page    int    `json:"page"`
		Keyword string `json:"keyword"`
		MaxPage int    `json:"max_page"`
		Limit   int    `json:"limit"`
		Offset  int    `json:"offset"`
		Count   int    `json:"count"`
	}
)

func (param *UserRedeemHistoryParam) ParseUserRedeemHistory(values url.Values) error {
	param.Page = 1
	param.Limit = 10
	param.Offset = 0

	if page, ok := values["page"]; ok && len(page) > 0 {
		if p, err := strconv.Atoi(page[0]); err == nil && p > 1 {
			param.Page = p
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

func (param *UserClaimedHistoryParam) ParseUserRedeemHistory(values url.Values) error {
	param.Page = 1
	param.Limit = 10
	param.Offset = 0

	if page, ok := values["page"]; ok && len(page) > 0 {
		if p, err := strconv.Atoi(page[0]); err == nil && p > 1 {
			param.Page = p
		}
	}

	if limit, ok := values["limit"]; ok && len(limit) > 0 {
		if l, err := strconv.Atoi(limit[0]); err == nil {
			param.Limit = l
		}
	}

	if keyword, ok := values["keyword"]; ok && len(keyword) > 0 {
		param.Keyword = keyword[0]
	}

	param.Offset = (param.Page - 1) * param.Limit

	return nil
}
