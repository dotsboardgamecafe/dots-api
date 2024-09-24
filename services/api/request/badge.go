package request

import (
	"net/url"
	"strconv"
	"strings"
)

type BadgeReq struct {
	BadgeCategory string         `json:"badge_category" validate:"required,max=100"`
	Name          string         `json:"name"`
	ImageURL      string         `json:"image_url"`
	Status        string         `json:"status"`
	VPPoint       int64          `json:"vp_point"`
	Description   string         `json:"description"`
	BadgeRule     []BadgeRuleReq `json:"badge_rule"`
}

type TournamentBadgeListReq struct {
	TournamentBadges []TournamentBadgeReq `json:"tournament_badges"`
}

type TournamentBadgeReq struct {
	BadgeCode     string       `json:"badge_code"`
	BadgeCategory string       `json:"badge_category" validate:"required,max=100"`
	Name          string       `json:"name"`
	ImageURL      string       `json:"image_url"`
	Status        string       `json:"status"`
	Description   string       `json:"description"`
	VPPoint       int64        `json:"vp_point"`
	BadgeRule     BadgeRuleReq `json:"badge_rule"`
}

type UpdateBadgeReq struct {
	BadgeCategory string               `json:"badge_category" validate:"required,max=100"`
	Name          string               `json:"name"`
	ImageURL      string               `json:"image_url"`
	Status        string               `json:"status"`
	VPPoint       int64                `json:"vp_point"`
	Description   string               `json:"description"`
	BadgeRule     []UpdateBadgeRuleReq `json:"badge_rule"`
}

type IsClaimBadgeReq struct {
	IsClaim bool `json:"is_claim"`
}

type BadgeParam struct {
	Page          int    `json:"page"`
	Limit         int    `json:"limit"`
	Offset        int    `json:"offset"`
	Count         int    `json:"count"`
	MaxPage       int    `json:"max_page"`
	Sort          string `json:"sort"`
	Order         string `json:"order"`
	Keyword       string `json:"keyword"`
	Status        string `json:"status"`
	BadgeCategory string `json:"badge_category"`
}

func (param *BadgeParam) ParseBadge(values url.Values) error {
	param.Keyword = ""
	param.Page = 1
	param.Limit = 10
	param.Sort = "desc"
	param.Order = "created_date"
	param.Status = ""
	param.Offset = 0
	param.BadgeCategory = ""

	if page, ok := values["page"]; ok && len(page) > 0 {
		if p, err := strconv.Atoi(page[0]); err == nil && p > 1 {
			param.Page = p
		}
	}

	if sort, ok := values["sort"]; ok && len(sort) > 0 && strings.ToLower(sort[0]) == "asc" {
		param.Sort = "asc"
	}

	if order, ok := values["order"]; ok && len(order) > 0 {
		if order[0] == "created_date" || order[0] == "name" || order[0] == "status" {
			param.Order = order[0]
		}
	}

	if status, ok := values["status"]; ok && len(status) > 0 {
		param.Status = status[0]
	}

	if keyword, ok := values["keyword"]; ok && len(keyword) > 0 {
		param.Keyword = keyword[0]
	}

	if badgeCategory, ok := values["badge_category"]; ok && len(badgeCategory) > 0 {
		param.BadgeCategory = badgeCategory[0]
	}

	if limit, ok := values["limit"]; ok && len(limit) > 0 {
		if l, err := strconv.Atoi(limit[0]); err == nil {
			param.Limit = l
		}
	}

	param.Offset = (param.Page - 1) * param.Limit

	return nil
}
