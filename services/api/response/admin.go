package response

type AdminRes struct {
	AdminCode   string `json:"admin_code"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	UserName    string `json:"user_name"`
	Status      string `json:"status"`
	ImageURL    string `json:"image_url"`
	PhoneNumber string `json:"phone_number"`
}
