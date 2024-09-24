package response

type GameMechanicRes struct {
	GameMechanicCode string `json:"code"`
	Name             string `json:"name"`
	CreatedDate      string `json:"created_date,omitempty"`
}
