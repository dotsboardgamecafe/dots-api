package response

type TournamentParticipantRes struct {
	UserCode       string `json:"user_code"`
	UserName       string `json:"user_name"`
	UserImgUrl     string `json:"user_image_url"`
	StatusWinner   bool   `json:"status_winner"`
	Status         string `json:"status"`
	AdditionalInfo string `json:"additional_info"`
	Position       int    `json:"position"`
	RewardPoint    int    `json:"reward_point"`
}
