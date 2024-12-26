package response

import "encoding/json"

type GameRes struct {
	CafeCode           string                 `json:"cafe_code"`
	CafeName           string                 `json:"cafe_name"`
	GameCode           string                 `json:"game_code"`
	GameType           string                 `json:"game_type"`
	Name               string                 `json:"name"`
	Location           string                 `json:"location"`
	ImageUrl           string                 `json:"image_url"`
	CollectionUrl      []string               `json:"collection_url"`
	Description        string                 `json:"description"`
	Status             string                 `json:"status"`
	Difficulty         string                 `json:"difficulty"`
	Level              float64                `json:"level"`
	Duration           int64                  `json:"duration"`
	AdminCode          string                 `json:"admin_code"`
	MinimalParticipant int64                  `json:"minimal_participant"`
	MaximumParticipant int64                  `json:"maximum_participant"`
	GameCategories     []GameCategoryRes      `json:"game_categories"`
	GameRelated        []GameRelatedRes       `json:"game_related"`
	GameRooms          []GameAvailableRoomRes `json:"game_rooms"`
	GameMasters        []AdminRes             `json:"game_masters"`
	NumberOfPopularity int64                  `json:"number_of_popularity"`
}

type GameDetailRes struct {
	CafeCode                  string                          `json:"cafe_code"`
	CafeName                  string                          `json:"cafe_name"`
	GameCode                  string                          `json:"game_code"`
	GameType                  string                          `json:"game_type"`
	Name                      string                          `json:"name"`
	Location                  string                          `json:"location"`
	ImageUrl                  string                          `json:"image_url"`
	CollectionUrl             []string                        `json:"collection_url"`
	Description               string                          `json:"description"`
	Status                    string                          `json:"status"`
	NumberOfPopularity        int64                           `json:"number_of_popularity"`
	Difficulty                string                          `json:"difficulty"`
	Level                     float64                         `json:"level"`
	Duration                  int64                           `json:"duration"`
	AdminCode                 string                          `json:"admin_code"`
	MinimalParticipant        int64                           `json:"minimal_participant"`
	MaximumParticipant        int64                           `json:"maximum_participant"`
	GameCategories            []GameCategoryRes               `json:"game_categories"`
	GameRelated               []GameRelatedRes                `json:"game_related"`
	GameRooms                 []GameAvailableRoomRes          `json:"game_rooms"`
	GameMasters               []AdminRes                      `json:"game_masters"`
	UserHavePlayedGameHistory []UsersHavePlayedGameHistoryRes `json:"user_have_played_game_history"`
	TotalPlayer               int64                           `json:"total_player"`
}

type UsersHavePlayedGameHistoryRes struct {
	GameId      int64  `json:"game_id"`
	GameName    string `json:"game_name"`
	UserCode    string `json:"user_code"`
	UserName    string `json:"username"`
	UserImage   string `json:"user_image"`
	CreatedDate string `json:"created_date,omitempty"`
}

type GameQRCodeRes struct {
	Source string `json:"qrcode"`
}

func BuildCollectionURLResp(data string) []string {
	var resp []string
	_ = json.Unmarshal([]byte(data), &resp)
	return resp
}
