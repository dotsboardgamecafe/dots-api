package response

type TournamentRes struct {
	GameCode               string                     `json:"game_code"`
	GameName               string                     `json:"game_name"`
	GameImgUrl             string                     `json:"game_img_url"`
	CafeCode               string                     `json:"cafe_code"`
	CafeName               string                     `json:"cafe_name"`
	CafeAddress            string                     `json:"cafe_address"`
	TournamentCode         string                     `json:"tournament_code"`
	PrizesImgUrl           string                     `json:"prizes_img_url"`
	ImageUrl               string                     `json:"image_url"`
	Name                   string                     `json:"name"`
	TournamentRules        string                     `json:"tournament_rules"`
	Difficulty             string                     `json:"difficulty"`
	StartDate              string                     `json:"start_date"`
	EndDate                string                     `json:"end_date"`
	StartTime              string                     `json:"start_time"`
	EndTime                string                     `json:"end_time"`
	PlayerSlot             int64                      `json:"player_slot"`
	BookingPrice           float64                    `json:"booking_price"`
	ParticipantVP          int64                      `json:"participant_vp"`
	Status                 string                     `json:"status"`
	DayPastEndDate         float64                    `json:"day_past_end_date"`
	CurrentUsedSlot        int64                      `json:"current_used_slot"`
	TournamentParticipants []TournamentParticipantRes `json:"tournament_participants"`
	TournamentBadges       []BadgeRes                 `json:"tournament_badges"`
	CreatedDate            string                     `json:"created_date"`
	UpdatedDate            string                     `json:"updated_date"`
	DeletedDate            string                     `json:"deleted_date"`
	HaveJoined             bool                       `json:"have_joined"`
}
