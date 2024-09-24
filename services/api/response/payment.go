package response

type MessageResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}
