package handler

import (
	"context"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"net/http"
)

// GetHallOfFame
func (h *Contract) GetHallOfFame(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.HallOfFameRes, 0)
		param = request.HallOfFameParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseHallOfFame(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetHallOfFameList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.HallOfFameRes{
			UserName:            v.UserName,
			UserFullName:        v.UserFullName,
			UserImgUrl:          v.UserImgUrl,
			TournamentBannerUrl: v.TournamentBannerUrl,
			TournamentName:      v.TournamentName,
			CafeName:            v.CafeName,
			CafeAddress:         v.CafeAddress,
		})
	}

	h.SendSuccess(w, res, param)
}

// GetMonthlyTopAchiever
func (h *Contract) GetMonthlyTopAchiever(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.MonthlyTopAchiever, 0)
		param = request.MonthlyTopAchieverParam{}
		data  = make([]model.MonthlyTopAchieverEnt, 0)
	)

	// Define urlQuery and Parse
	err = param.ParseMonthlyTopAchiever(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if param.Category == "vp" {
		data, param, err = m.GetMostVP(h.DB, ctx, param)
	} else if param.Category == "unique_game" {
		data, param, err = m.GetUniqueGame(h.DB, ctx, param)
	}

	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		locationCity := v.Location
		if locationCity == "" {
			locationCity = param.CafeCity
		}

		res = append(res, response.MonthlyTopAchiever{
			Ranking:         v.Ranking,
			UserFullName:    v.UserFullName,
			UserName:        v.UserName,
			UserImgUrl:      v.UserImgUrl,
			Location:        locationCity,
			TotalPoint:      v.TotalPoint,
			TotalGamePlayed: v.TotalGamePlayed,
		})
	}

	h.SendSuccess(w, res, param)
}
