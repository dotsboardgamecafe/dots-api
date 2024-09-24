package request

import (
	"dots-api/lib/utils"
	"fmt"
	"net/url"
	"strconv"
)

type (
	UserTransactionParam struct {
		Page    int    `json:"page"`
		MaxPage int    `json:"max_page"`
		Limit   int    `json:"limit"`
		Offset  int    `json:"offset"`
		Count   int    `json:"count"`
		Status  string `json:"status"`
	}
)

func (param *UserTransactionParam) ParseUserTransactionParam(values url.Values) error {
	param.Page = 1
	param.Limit = 10
	param.Status = ""
	param.Offset = 0

	if page, ok := values["page"]; ok && len(page) > 0 {
		if p, err := strconv.Atoi(page[0]); err == nil && p > 1 {
			param.Page = p
		}
	}

	if status, ok := values["status"]; ok && len(status) > 0 {
		if !utils.Contains(utils.XenditTransactionStatus, status[0]) {
			return fmt.Errorf("%s", "wrong status value for transaction(PENDING|PAID|SETTLED|EXPIRED")
		}
		param.Status = status[0]
	}

	if limit, ok := values["limit"]; ok && len(limit) > 0 {
		if l, err := strconv.Atoi(limit[0]); err == nil {
			param.Limit = l
		}
	}

	param.Offset = (param.Page - 1) * param.Limit

	return nil
}
