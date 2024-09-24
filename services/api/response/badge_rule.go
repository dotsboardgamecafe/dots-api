package response

type BadgeRuleRes struct {
	BadgeRuleCode     string      `json:"badge_rule_code"`
	BadgeId           int64       `json:"badge_id"`
	CategoryBadgeRule string      `json:"category_badge_rule"`
	KeyCondition      string      `json:"name"`
	ValueType         string      `json:"image_url"`
	Value             interface{} `json:"value"`
}
