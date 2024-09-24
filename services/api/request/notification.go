package request

import (
	"net/url"
	"strconv"
	"strings"
)

type NotificationIsSeenReq struct {
	IsSeen bool `json:"is_seen"`
}

type NotificationParam struct {
	Page           int    `json:"page"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	Count          int    `json:"count"`
	CountUnread    int    `json:"count_unread"`
	Sort           string `json:"sort"`
	Order          string `json:"order"`
	Keyword        string `json:"keyword"`
	Status         string `json:"status"`
	ReceiverSource string `json:"receiver_source"`
	ReceiverCode   string `json:"receiver_code"`
}

func (param *NotificationParam) ParseNotification(values url.Values) error {
	param.Keyword = ""
	param.Page = 1
	param.Limit = 10
	param.Sort = "desc"
	param.Order = "notifications.id"
	param.Status = ""
	param.Offset = 0
	param.ReceiverSource = ""
	param.ReceiverCode = ""

	if page, ok := values["page"]; ok && len(page) > 0 {
		if p, err := strconv.Atoi(page[0]); err == nil && p > 1 {
			param.Page = p
		}
	}

	if sort, ok := values["sort"]; ok && len(sort) > 0 && strings.ToLower(sort[0]) == "asc" {
		param.Sort = "asc"
	}

	if order, ok := values["order"]; ok && len(order) > 0 {
		if order[0] == "created_date" || order[0] == "name" {
			param.Order = order[0]
		}
	}

	if status, ok := values["status"]; ok && len(status) > 0 {
		param.Status = status[0]
	}

	if receiverSource, ok := values["receiver_source"]; ok && len(receiverSource) > 0 {
		param.ReceiverSource = receiverSource[0]
	}

	if receiverCode, ok := values["receiver_code"]; ok && len(receiverCode) > 0 {
		param.ReceiverCode = receiverCode[0]
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
