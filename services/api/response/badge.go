package response

type BadgeRes struct {
	BadgeCode     string         `json:"badge_code"`
	BadgeCategory string         `json:"badge_category"`
	Name          string         `json:"name"`
	ImageURL      string         `json:"image_url"`
	BadgeRules    []BadgeRuleRes `json:"badge_rules"`
	VPPoint       int64          `json:"vp_point"`
	Description   string         `json:"description"`
	Status        string         `json:"status"`
	ParentCode    string         `json:"parent_code"`
	CreatedDate   string         `json:"created_date"`
	UpdatedDate   string         `json:"updated_date"`
	DeletedDate   string         `json:"deleted_date"`
}
