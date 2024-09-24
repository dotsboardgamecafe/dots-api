package request

import (
	"dots-api/lib/array"
	"net/url"
	"strconv"
	"strings"
)

type TournamentReq struct {
	GameCode        string   `json:"game_code"`
	ImageUrl        string   `json:"image_url"`
	Name            string   `json:"name" validate:"required,max=100"`
	TournamentRules string   `json:"tournament_rules"`
	Level           string   `json:"level"`
	StartDate       string   `json:"start_date"`
	EndDate         string   `json:"end_date"`
	StartTime       string   `json:"start_time"`
	EndTime         string   `json:"end_time"`
	PlayerSlot      int64    `json:"player_slot"`
	BookingPrice    float64  `json:"booking_price"`
	BadgeCodes      []string `json:"badge_codes"`
	PrizesImgUrl    string   `json:"prizes_img_url"`
	Status          string   `json:"status" validate:"max=10"`
	ParticipantVP   int64    `json:"participant_vp"`
	LocationCode    string   `json:"location_code" validate:"max=50"`
}

type SetWinnerTournamentReq struct {
	TournamentParticipant []TournamentParticipant `json:"tournament_participant"`
}

type UpdateStatusTournamentReq struct {
	Status string `json:"status" validate:"required"`
}

type TournamentParticipant struct {
	Position  int    `json:"position"`
	BadgeCode string `json:"badge_code"`
	UserCode  string `json:"user_code"`
}

type TournamentParam struct {
	Page           int      `json:"page"`
	Limit          int      `json:"limit"`
	Offset         int      `json:"offset"`
	MaxPage        int      `json:"max_page"`
	Count          int      `json:"count"`
	Sort           []string `json:"sort"`
	Order          []string `json:"order"`
	Keyword        string   `json:"keyword"`
	Status         string   `json:"status"`
	Name           string   `json:"name"`
	CafeCity       string   `json:"location"`
	GameCode       string   `json:"game_code"`
	TournamentDate string   `json:"tournament_date"`
}

func (param *TournamentParam) ParseTournament(values url.Values) error {
	param.Keyword = ""
	param.Page = 1
	param.Limit = 10
	param.Sort = []string{"desc"}
	param.Order = []string{"tournaments.created_date"}
	param.Status = ""
	param.Offset = 0
	param.CafeCity = ""
	param.GameCode = ""
	param.TournamentDate = ""

	if page, ok := values["page"]; ok && len(page) > 0 {
		if p, err := strconv.Atoi(page[0]); err == nil && p > 1 {
			param.Page = p
		}
	}

	if sort, ok := values["sort"]; ok && len(sort) > 0 {
		param.Sort = strings.Split(sort[0], ",")
	}

	if order, ok := values["order"]; ok && len(order) > 0 {
		arrStr := new(array.ArrStr)
		orders := strings.Split(order[0], ",")
		validOrders := []string{}
		for _, o := range orders {
			if exist, _ := arrStr.InArray(o, []string{"tournaments.created_date", "tournaments.name", "tournaments.status", "days_past_end_date"}); exist {
				validOrders = append(validOrders, o)
			}
		}
		if len(validOrders) > 0 {
			param.Order = validOrders
		}
	}

	if status, ok := values["status"]; ok && len(status) > 0 {
		param.Status = status[0]
	}

	if cafeCity, ok := values["location"]; ok && len(cafeCity) > 0 {
		param.CafeCity = cafeCity[0]
	}

	if gameCode, ok := values["game_code"]; ok && len(gameCode) > 0 {
		param.GameCode = gameCode[0]
	}

	if tournamentDate, ok := values["tournament_date"]; ok && len(tournamentDate) > 0 {
		param.TournamentDate = tournamentDate[0]
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
