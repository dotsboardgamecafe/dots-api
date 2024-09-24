package response

type GameTypeRes struct {
	GameTypeCode string `json:"code"`
	Name         string `json:"name"`
	CreatedDate  string `json:"created_date,omitempty"`
}
