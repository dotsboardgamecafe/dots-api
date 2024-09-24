package handler

import (
	"dots-api/services/api/response"
	"net/http"
)

func (h *Contract) SuccessCallback(w http.ResponseWriter, r *http.Request) {
	h.SendSuccess(w, response.MessageResponse{
		Message: "Success",
		Status:  true,
	}, nil)
}

func (h *Contract) FailureCallback(w http.ResponseWriter, r *http.Request) {
	h.SendSuccess(w, response.MessageResponse{
		Message: "Fail",
		Status:  false,
	}, nil)
}
