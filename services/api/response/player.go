package response

// Hall of Fame
type HallOfFameRes struct {
	UserCode            string `json:"user_code"`
	UserName            string `json:"user_name"`
	UserFullName        string `json:"user_fullname"`
	UserImgUrl          string `json:"user_img_url"`
	TournamentBannerUrl string `json:"tournament_banner_url"`
	TournamentName      string `json:"tournament_name"`
	CafeName            string `json:"cafe_name"`
	CafeAddress         string `json:"cafe_address"`
}

// Monthly Top Achiever
type MonthlyTopAchiever struct {
	Ranking         int    `json:"rank"`
	UserCode        string `json:"user_code"`
	UserFullName    string `json:"user_fullname"`
	UserName        string `json:"user_name"`
	UserImgUrl      string `json:"user_img_url"`
	Location        string `json:"location"`
	TotalPoint      int    `json:"total_point,omitempty"`
	TotalGamePlayed int    `json:"total_game_played,omitempty"`
}
