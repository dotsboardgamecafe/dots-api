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
	RoomReq struct {
		GameMasterCode     string  `json:"game_master_code" validate:"required,max=50"`
		GameCode           string  `json:"game_code" validate:"required,max=50"`
		LocationCode       string  `json:"location_code" validate:"required,max=50"`
		RoomType           string  `json:"room_type" validate:"required,max=50"`
		Name               string  `json:"name" validate:"required,max=100"`
		Description        string  `json:"description"`
		Instruction        string  `json:"instruction" validate:"max=500"`
		Difficulty         string  `json:"difficulty" validate:"max=50"`
		StartDate          string  `json:"start_date" validate:"required"`
		EndDate            string  `json:"end_date" validate:"required"`
		StartTime          string  `json:"start_time" validate:"required"`
		EndTime            string  `json:"end_time" validate:"required"`
		MaximumParticipant int     `json:"maximum_participant" validate:"required"`
		BookingPrice       float64 `json:"booking_price" validate:"required,min=10000"`
		RewardPoint        int     `json:"reward_point"`
		InstagramLink      string  `json:"instagram_link" validate:"max=500"`
		ImageURL           string  `json:"image_url" validate:"max=500"`
		Status             string  `json:"status" validate:"required,max=50"`
	}

	BookingRoomReq struct {
		Price      float64 `json:"price" validate:"required,numeric"`
		DataSource string  `json:"data_source" validate:"required"`
		SourceCode string  `json:"source_code" validate:"required"`
	}

	SetWinnerRoomReq struct {
		RoomParticipant []RoomParticipant `json:"room_participant"`
	}

	RoomParticipant struct {
		Position int    `json:"position"`
		UserCode string `json:"user_code"`
	}

	UpdateStatusRoomReq struct {
		Status string `json:"status" validate:"required"`
	}

	RoomParam struct {
		Page     int      `json:"page"`
		MaxPage  int      `json:"max_page"`
		Limit    int      `json:"limit"`
		Offset   int      `json:"offset"`
		Count    int      `json:"count"`
		Sort     []string `json:"sort"`
		Order    []string `json:"order"`
		Keyword  string   `json:"keyword"`
		Status   string   `json:"status"`
		CafeCity string   `json:"location"`
		RoomType string   `json:"room_type"`
	}
)

func (param *RoomParam) ParseRoom(values url.Values) error {
	param.Keyword = ""
	param.Page = 1
	param.Limit = 10
	param.Sort = []string{"desc"}
	param.Order = []string{"rooms.created_date"}
	param.Status = ""
	param.Offset = 0

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
			if exist, _ := arrStr.InArray(o, []string{"rooms.id", "rooms.name", "rooms.description", "rooms.status", "rooms.created_date", "days_past_end_date"}); exist {
				validOrders = append(validOrders, o)
			}
		}
		if len(validOrders) > 0 {
			param.Order = validOrders
		}
	}

	if status, ok := values["status"]; ok && len(status) > 0 {
		if !utils.Contains(utils.StatusRoom, status[0]) {
			return fmt.Errorf("%s", "wrong status value for room(publish|unpublish)")
		}
		param.Status = status[0]
	}

	if roomType, ok := values["room_type"]; ok && len(roomType) > 0 {
		if !utils.Contains(utils.RoomType, roomType[0]) {
			return fmt.Errorf("%s", "wrong type value for room(normal|special_event)")
		}
		param.RoomType = roomType[0]
	}

	if cafeCity, ok := values["location"]; ok && len(cafeCity) > 0 {
		param.CafeCity = cafeCity[0]
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
