package handler

import (
	"context"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetUserFavouriteGameAct ...
func (h *Contract) GetUserFavouriteGameAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.UserFavouriteGameRes, 0)
		code  = chi.URLParam(r, "code")
		param = request.UserFavouriteGameParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseUserFavouriteGame(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetUserFavouriteGames(h.DB, ctx, code, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.UserFavouriteGameRes{
			UserId:                  v.UserId.Int64,
			UserCode:                v.UserCode.String,
			GameCategoryId:          v.GameCategoryId.Int64,
			GameCategoryName:        v.GameCategoryName.String,
			GameCategoryDescription: v.GameCategoryDescription.String,
			GameCategoryImageUrl:    v.GameCategoryImageUrl.String,
			TotalPlay:               v.TotalPlay.Int64,
		})
	}

	h.SendSuccess(w, res, param)
}
