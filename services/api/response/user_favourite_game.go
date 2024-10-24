package response

type UserFavouriteGameRes struct {
	UserId                  int64  `json:"user_id"`
	UserCode                string `json:"user_code"`
	GameCategoryName        string `json:"game_category_name"`
	GameCategoryDescription string `json:"game_category_description"`
	GameCategoryImageUrl    string `json:"game_category_image_url"`
	TotalPlay               int64  `json:"total_play"`
}
