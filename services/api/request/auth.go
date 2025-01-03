package request

type LoginUserReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterUserReq struct {
	Fullname        string `json:"fullname" validate:"required"`
	DateOfBirth     string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
	Gender          string `json:"gender" validate:"required,oneof=male female"`
	Email           string `json:"email" validate:"required"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	Username        string `json:"username" validate:"required,max=10"`
}

type RequestVerifyEmailReq struct {
	Email string `json:"email" validate:"required"`
	Type  string `json:"type" validate:"required"`
}

type VerifyEmailReq struct {
	Email string `json:"email" validate:"required"`
	Token string `json:"token" validate:"required"`
}

type ResetPasswordReq struct {
	NewPassword     string `json:"new_password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type VerifyPasswordReq struct {
	Password string `json:"password" validate:"required"`
}
