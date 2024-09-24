package response

type UserGameCollectionRes struct {
	UserId       int64  `json:"user_id"`
	UserCode     string `json:"user_code"`
	GameCode     string `json:"game_code"`
	GameId       int64  `json:"game_id"`
	GameName     string `json:"game_name"`
	GameImageUrl string `json:"game_image_url"`
}
