package response

type UserRes struct {
	UserCode           string `json:"user_code"`
	Email              string `json:"email"`
	UserName           string `json:"username"`
	PhoneNumber        string `json:"phone_number"`
	FullName           string `json:"fullname"`
	DateOfBirth        string `json:"date_of_birth"`
	Gender             string `json:"gender"`
	ImageURL           string `json:"image_url"`
	LatestPoint        int    `json:"latest_point"`
	LatestTier         string `json:"latest_tier"`
	Password           string `json:"password"`
	XPlayer            string `json:"x_player"`
	StatusVerification bool   `json:"status_verification"`
	Status             string `json:"status"`
	TotalSpent         int    `json:"total_spent"`
	CreatedDate        string `json:"created_date"`
	UpdatedDate        string `json:"updated_date"`
	DeletedDate        string `json:"deleted_date,omitempty"`
}

type UserProfileRes struct {
	UserCode           string               `json:"user_code"`
	Email              string               `json:"email"`
	UserName           string               `json:"username"`
	PhoneNumber        string               `json:"phone_number"`
	FullName           string               `json:"fullname"`
	DateOfBirth        string               `json:"date_of_birth"`
	Gender             string               `json:"gender"`
	ImageURL           string               `json:"image_url"`
	LatestPoint        int                  `json:"latest_point"`
	LatestTier         string               `json:"latest_tier"`
	TierRangePoint     *TierRangePointRes   `json:"tier_range_point,omitempty"`
	TierBenefits       []TierWithBenefitRes `json:"tier_benefits"`
	MemberSince        string               `json:"member_since"`
	Password           string               `json:"password"`
	XPlayer            string               `json:"x_player"`
	StatusVerification bool                 `json:"status_verification"`
	Status             string               `json:"status"`
	CreatedDate        string               `json:"created_date"`
	UpdatedDate        string               `json:"updated_date"`
	DeletedDate        string               `json:"deleted_date"`
}

type PlayerActivitiesRes struct {
	UserName         string `json:"username"`
	TitleDescription string `json:"title_description"`
	GameImgUrl       string `json:"game_image_url,omitempty"`
	GameName         string `json:"game_name,omitempty"`
	GameCode         string `json:"game_code,omitempty"`
	DataSource       string `json:"data_source"`
	SourceCode       string `json:"source_code"`
	Point            int    `json:"point"`
	CreatedDate      string `json:"created_date"`
}
