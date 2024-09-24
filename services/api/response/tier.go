package response

type TierRes struct {
	TierCode    string `json:"tier_code"`
	Name        string `json:"name"`
	MinPoint    int64  `json:"min_point"`
	MaxPoint    int64  `json:"max_point"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedDate string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
	DeletedDate string `json:"deleted_date"`
}

type TierRangePointRes struct {
	MinPoint int `json:"min_point"`
	MaxPoint int `json:"max_point"`
}

type TierWithBenefitRes struct {
	RewardCode        string `json:"reward_code"`
	RewardName        string `json:"reward_name"`
	RewardImageUrl    string `json:"reward_img_url"`
	RewardDescription string `json:"reward_description"`
}
