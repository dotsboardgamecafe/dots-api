package response

type (
	SuccessRedirectRes struct {
		Message string `json:"message"`
	}

	FailedRedirectRes struct {
		Message string `json:"message"`
	}
)
