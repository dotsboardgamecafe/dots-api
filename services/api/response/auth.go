package response

type LoginUserRes struct {
	Token       string          `json:"token"`
	UserCode    string          `json:"user_code"`
	ImageURL    string          `json:"image_url"`
	FullName    string          `json:"fullname"`
	DateOfBirth string          `json:"date_of_birth"`
	Gender      string          `json:"gender"`
	PhoneNumber string          `json:"phone_number"`
	Email       string          `json:"email"`
	ExpiredAt   string          `json:"expired_at"`
	ActorType   string          `json:"actor_type"`
	CreatedDate string          `json:"created_date"`
	RoleId      int             `json:"role_id"`
	RoleCode    string          `json:"role_code"`
	Permissions []PermissionRes `json:"permissions"`
}

type RegisterUserRes struct {
	UserCode    string `json:"user_code"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email"`
	CreatedDate string `json:"created_date"`
}

type VerifyPasswordRes struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}
