package response

type CafeRes struct {
	CafeCode    string `json:"cafe_code"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Province    string `json:"province"`
	City        string `json:"city"`
}
