package handler

import (
	"context"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetUserGameHistoryAct ...
func (h *Contract) GetUserGameHistoryAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.UserGameHistoryRes, 0)
		code  = chi.URLParam(r, "code")
		param = request.UserGameHistoryParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseUserGameHistory(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetUserGameHistories(h.DB, ctx, code, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.UserGameHistoryRes{
			UserId:       v.UserId.Int64,
			UserCode:     v.UserCode.String,
			GameId:       v.GameId.Int64,
			GameName:     v.GameName.String,
			GameImageUrl: v.GameImageUrl.String,
		})
	}

	h.SendSuccess(w, res, param)
}
