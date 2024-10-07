package request

type BadgeRuleReq struct {
	BadgeRuleCategory string      `json:"badge_rule_category"`
	KeyCondition      string      `json:"key_condition"`
	ValueType         string      `json:"value_type"`
	Value             interface{} `json:"value"`
}

type UpdateBadgeRuleReq struct {
	BadgeRuleCode     string      `json:"badge_rule_code"`
	BadgeRuleCategory string      `json:"badge_rule_category"`
	KeyCondition      string      `json:"key_condition"`
	ValueType         string      `json:"value_type"`
	Value             interface{} `json:"value"`
}

type TimeLimitCategory struct {
	Category  string `json:"category"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type SpesificBoardGameCategory struct {
	GameCode     []string `json:"game_code"`
	NeedGM       bool     `json:"need_gm"`
	TotalPlayed  int64    `json:"total_played"`
	BookingPrice int64    `json:"booking_price"`
}

type TournamentCategory struct {
	Position int `json:"position"`
}
