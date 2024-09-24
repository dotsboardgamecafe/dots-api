package response

type NotificationResp struct {
	ReceiverSource   string `json:"receiver_source"`
	ReceiverCode     string `json:"receiver_code"`
	NotificationCode string `json:"notification_code"`
	TransactionCode  string `json:"transaction_code"`
	Type             string `json:"type"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	ImageUrl         string `json:"image_url"`
	StatusRead       bool   `json:"status_read"`
	CreatedDate      string `json:"created_date"`
	UpdatedDate      string `json:"updated_date,omitempty"`
}

type NotificationTournamentResp struct {
	StartDate   string `json:"start_date"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	CafeName    string `json:"cafe_name"`
	CafeAddress string `json:"cafe_address"`
	GameName    string `json:"game_name"`
	Level       string `json:"level"`
}

type NotificationRoomResp struct {
	StartDate   string `json:"start_date"`
	CafeName    string `json:"cafe_name"`
	CafeAddress string `json:"cafe_address"`
	GameName    string `json:"game_name"`
	Level       string `json:"level"`
}
