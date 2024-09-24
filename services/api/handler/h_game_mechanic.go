package handler

import (
	"context"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Contract) GetGameMechanicList(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.GameMechanicRes, 0)
		param = request.GameMechanicParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseGameMechanic(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, err := m.ListOfGameMechanics(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.GameMechanicRes{
			GameMechanicCode: v.SettingCode,
			Name:             v.ContentValue,
			CreatedDate:      v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		})
	}

	h.SendSuccess(w, res, param)
}

func (h *Contract) GetDetailGameMechanic(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
		code = chi.URLParam(r, "code")
	)

	data, err := m.GetSettingByCode(h.DB, ctx, code, "game_mechanic")
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, response.GameMechanicRes{
		GameMechanicCode: data.SettingCode,
		Name:             data.ContentValue,
	}, nil)
}

func (h *Contract) AddGameMechanic(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req = request.GameMechanicReq{}
		ctx = context.TODO()
		m   = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	// Generate Random Code
	code := utils.GeneratePrefixCode(utils.SettingPrefix)
	setLabel := utils.ConvertToSnakeCase(req.Name)
	err = m.AddSetting(h.DB, ctx, code, "game_mechanic", setLabel, 1, "string", req.Name, true)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, "Game Mechanic berhasil ditambah", nil)
}

func (h *Contract) UpdateGameMechanic(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		req  = request.GameMechanicReq{}
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	data, err := m.GetSettingByCode(h.DB, ctx, code, "game_mechanic")
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	previousName := data.ContentValue
	key := utils.ConvertToSnakeCase(req.Name)
	err = m.UpdateGameMechanic(h.DB, ctx, code, key, req.Name, previousName)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, "Game Mechanic berhasil di-update", nil)
}

func (h *Contract) DeleteGameMechanic(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	data, err := m.GetSettingByCode(h.DB, ctx, code, "game_mechanic")
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.DeleteGameMechanic(h.DB, ctx, code, data.ContentValue)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, "Game Mechanic berhasil dihapus", nil)
}
