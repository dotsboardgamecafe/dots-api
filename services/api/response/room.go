package response

type RoomRes struct {
	GameMasterCode     string               `json:"game_master_code"`
	GameMasterName     string               `json:"game_master_name"`
	GameMasterImgUrl   string               `json:"game_master_img_url"`
	GameCode           string               `json:"game_code"`
	GameName           string               `json:"game_name"`
	GameImgUrl         string               `json:"game_img_url"`
	CafeCode           string               `json:"cafe_code"`
	CafeName           string               `json:"cafe_name"`
	CafeAddress        string               `json:"cafe_address"`
	RoomCode           string               `json:"room_code"`
	RoomType           string               `json:"room_type"`
	RoomBannerUrl      string               `json:"room_banner_url"`
	Name               string               `json:"name"`
	Description        string               `json:"description"`
	Difficulty         string               `json:"difficulty"`
	StartDate          string               `json:"start_date"`
	EndDate            string               `json:"end_date"`
	StartTime          string               `json:"start_time"`
	EndTime            string               `json:"end_time"`
	MaximumParticipant int                  `json:"maximum_participant"`
	DayPastEndDate     float64              `json:"day_past_end_date"`
	BookingPrice       float64              `json:"booking_price"`
	RewardPoint        int                  `json:"reward_point"`
	InstagramLink      string               `json:"instagram_link"`
	Status             string               `json:"status"`
	CurrentUsedSlot    int                  `json:"current_used_slot"`
	RoomParticipant    []RoomParticipantRes `json:"room_participants"`
	HaveJoined         bool                 `json:"have_joined"`
}

type RoomListRes struct {
	CafeId             int64   `json:"cafe_id"`
	CafeCode           string  `json:"cafe_code"`
	CafeName           string  `json:"cafe_name"`
	CafeAddress        string  `json:"cafe_address"`
	RoomCode           string  `json:"room_code"`
	RoomType           string  `json:"room_type"`
	RoomImgUrl         string  `json:"room_img_url"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	Instruction        string  `json:"instruction"`
	Difficulty         string  `json:"difficulty"`
	StartDate          string  `json:"start_date"`
	EndDate            string  `json:"end_date"`
	StartTime          string  `json:"start_time"`
	EndTime            string  `json:"end_time"`
	MaximumParticipant int     `json:"maximum_participant"`
	CurrentUsedSlot    int     `json:"current_used_slot"`
	InstagramLink      string  `json:"instagram_link"`
	DayPastEndDate     float64 `json:"day_past_end_date"`
	Status             string  `json:"status"`
	BookingPrice       float64 `json:"booking_price"`
	GameMasterName     string  `json:"game_master_name"`
	GameMasterImageUrl string  `json:"game_master_image_url"`
	GameCode           string  `json:"game_code"`
	GameName           string  `json:"game_name"`
	GameImgUrl         string  `json:"game_img_url"`
}

type BookingRes struct {
	InvoiceUrl string `json:"invoice_url"`
	ExpiredAt  string `json:"expired_at"`
}
