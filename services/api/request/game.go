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
	GameReq struct {
		CafeCode           string            `json:"cafe_code" validate:"required,max=50"`
		GameType           string            `json:"game_type" validate:"required,max=100"`
		Name               string            `json:"name" validate:"required,max=100"`
		ImageUrl           string            `json:"image_url" validate:"max=500"`
		CollectionUrl      []string          `json:"collection_url" validate:"max=500"`
		Description        string            `json:"description"`
		Status             string            `json:"status" validate:"required,max=50"`
		Difficulty         string            `json:"difficulty"`
		Level              float64           `json:"level"`
		Duration           int64             `json:"duration"`
		MinimalParticipant int64             `json:"minimal_participant"`
		MaximumParticipant int64             `json:"maximum_participant" validate:"required"`
		AdminCode          string            `json:"admin_code"`
		GameCategories     []GameCategoryReq `json:"game_categories"`
	}

	GameParam struct {
		Page               int      `json:"page"`
		Limit              int      `json:"limit"`
		Offset             int      `json:"offset"`
		Count              int      `json:"count"`
		Sort               string   `json:"sort"`
		SortKey            string   `json:"sort_key"`
		Keyword            string   `json:"keyword"`
		Status             string   `json:"status"`
		Level              float64  `json:"level"`
		NumberOfPlayers    int      `json:"number_of_players"`
		MinimalParticipant int      `json:"minimal_participant"`
		MaximumParticipant int      `json:"maximum_participant"`
		MinDuration        int      `json:"min_duration"`
		MaxDuration        int      `json:"max_duration"`
		MaxPage            int      `json:"max_page"`
		CafeCode           string   `json:"cafe_code"`
		GameType           []string `json:"game_type"`
		GameCategoryName   []string `json:"game_category_name"`
		Difficulty         string   `json:"difficulty"`
		Location           []string `json:"location"`
	}
)

func (param *GameParam) ParseGame(values url.Values) error {
	param.Keyword = ""
	param.Page = 1
	param.Limit = 10
	param.Level = 0
	param.MinDuration = 0
	param.MaxDuration = 0
	param.MinimalParticipant = 0
	param.MaximumParticipant = 0
	param.Sort = "desc"
	param.SortKey = "g.created_date"
	param.Status = ""
	param.Difficulty = ""
	param.Offset = 0

	if page, ok := values["page"]; ok && len(page) > 0 {
		if p, err := strconv.Atoi(page[0]); err == nil && p > 1 {
			param.Page = p
		}
	}

	if level, ok := values["level"]; ok && len(level) > 0 {
		if l, err := strconv.ParseFloat(level[0], 64); err == nil && l > 1 {
			param.Level = l
		}
	}

	if sort, ok := values["sort"]; ok && len(sort) > 0 && strings.ToLower(sort[0]) == "asc" {
		param.Sort = "asc"
	}

	if sortKey, ok := values["sort_key"]; ok && len(sortKey) > 0 {
		arrStr := new(array.ArrStr)
		if exist, _ := arrStr.InArray(sortKey[0], []string{"name", "level", "created_date", "number_of_popularity"}); exist {
			if sortKey[0] != "number_of_popularity" {
				sortKey[0] = "g." + sortKey[0]
			}
			param.SortKey = sortKey[0]
		}
	}

	if status, ok := values["status"]; ok && len(status) > 0 {
		if !utils.Contains(utils.StatusGame, status[0]) {
			return fmt.Errorf("%s", "wrong status value for game(active|inactive)")
		}
		param.Status = status[0]
	}

	if cafeCode, ok := values["cafe_code"]; ok && len(cafeCode) > 0 {
		param.CafeCode = cafeCode[0]
	}

	if difficulty, ok := values["difficulty"]; ok && len(difficulty) > 0 {
		param.Difficulty = difficulty[0]
	}

	if gameType, ok := values["game_type"]; ok && len(gameType) > 0 {
		param.GameType = strings.Split(gameType[0], ",")
	}

	if gameCategoryName, ok := values["game_category_name"]; ok && len(gameCategoryName) > 0 {
		param.GameCategoryName = strings.Split(strings.ToLower(gameCategoryName[0]), ",")
	}

	if location, ok := values["location"]; ok && len(location) > 0 {
		param.Location = strings.Split(strings.ToLower(location[0]), ",")
	}

	if keyword, ok := values["keyword"]; ok && len(keyword) > 0 {
		param.Keyword = keyword[0]
	}

	if limit, ok := values["limit"]; ok && len(limit) > 0 {
		if l, err := strconv.Atoi(limit[0]); err == nil {
			param.Limit = l
		}
	}

	if minDuration, ok := values["min_duration"]; ok && len(minDuration) > 0 {
		if l, err := strconv.Atoi(minDuration[0]); err == nil {
			param.MinDuration = l
		}
	}

	if maxDuration, ok := values["max_duration"]; ok && len(maxDuration) > 0 {
		if l, err := strconv.Atoi(maxDuration[0]); err == nil {
			param.MaxDuration = l
		}
	}

	if NumberOfPlayers, ok := values["number_of_players"]; ok && len(NumberOfPlayers) > 0 {
		if l, err := strconv.Atoi(NumberOfPlayers[0]); err == nil {
			param.NumberOfPlayers = l
		}
	}

	if minParticipant, ok := values["minimal_participant"]; ok && len(minParticipant) > 0 {
		if l, err := strconv.Atoi(minParticipant[0]); err == nil {
			param.MinimalParticipant = l
		}
	}

	if maxParticipant, ok := values["maximum_participant"]; ok && len(maxParticipant) > 0 {
		if l, err := strconv.Atoi(maxParticipant[0]); err == nil {
			param.MaximumParticipant = l
		}
	}

	param.Offset = (param.Page - 1) * param.Limit

	return nil
}
