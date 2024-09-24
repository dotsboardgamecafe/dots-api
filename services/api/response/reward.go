package response

type RewardRes struct {
	Tier         TierRes
	Name         string `json:"name"`
	ImageUrl     string `json:"image_url"`
	CategoryType string `json:"category_type"`
	RewardCode   string `json:"reward_code"`
	VoucherCode  string `json:"voucher_code"`
	Status       string `json:"status"`
	Description  string `json:"description"`
	ExpiredDate  string `json:"expired_date"`
	CreatedDate  string `json:"created_date"`
	UpdatedDate  string `json:"updated_date"`
}
