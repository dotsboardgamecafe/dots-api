package response

type UserGameHistoryRes struct {
	UserId             int64   `json:"user_id"`
	UserCode           string  `json:"user_code"`
	GameId             int64   `json:"game_id"`
	GameName           string  `json:"game_name"`
	GameImageUrl       string  `json:"game_image_url"`
	GameDuration       int64   `json:"game_duration"`
	GameDifficulty     float64 `json:"game_difficulty"`
	GameType           string  `json:"game_type"`
	GamePlayerSlot     int64   `json:"player_slot"`
	GameMasterId       int64   `json:"game_master_id"`
	GameMasterCode     string  `json:"game_master_code"`
	GameMasterName     string  `json:"game_master_name"`
	GameMasterImageUrl string  `json:"game_master_image_url"`
	GamePlayType       string  `json:"game_play_type"`
}
