package response

type UserBadgeRes struct {
	UserId        int64  `json:"user_id"`
	BadgeId       int64  `json:"badge_id"`
	BadgeName     string `json:"badge_name"`
	BadgeImageURL string `json:"badge_image_url"`
	BadgeCode     string `json:"badge_code"`
	VPPoint       int    `json:"vp_point"`
	Description   string `json:"description"`
	BadgeCategory string `json:"badge_category"`
	IsClaim       bool   `json:"is_claim"`
	CreatedDate   string `json:"created_date"`
	IsBadgeOwned  bool   `json:"is_badge_owned"`
	NeedToClaim   bool   `json:"need_to_claim"`
}
