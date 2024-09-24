package response

type UserBookingSummaryRes struct {
	Id               int64   `json:"id,omitempty"`
	UserId           int64   `json:"user_id,omitempty"`
	DataSource       string  `json:"data_source"`
	TransactionCode  string  `json:"transaction_code"`
	GameName         string  `json:"game_name"`
	GameImgUrl       string  `json:"game_img_url,omitempty"`
	Price            float64 `json:"final_price_amount"`
	AwardedUserPoint int     `json:"awarded_user_point"`
	PaymentMethod    string  `json:"payment_method"`
	Status           string  `json:"status"`
	CreatedDate      string  `json:"created_date"`
	UpdatedDate      string  `json:"updated_date"`
}
