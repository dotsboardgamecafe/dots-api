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

// GetSettingListAct ...
func (h *Contract) GetSettingListAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.SettingRes, 0)
		param = request.SettingParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseSetting(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, err := m.GetSettingList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.SettingRes{
			SettingCode:  v.SettingCode,
			SetGroup:     v.SetGroup,
			SetKey:       v.SetKey,
			SetLabel:     v.SetLabel,
			SetOrder:     v.SetOrder,
			ContentType:  v.ContentType,
			ContentValue: v.ContentValue,
			IsActive:     v.IsActive,
		})
	}

	h.SendSuccess(w, res, param)
}

// AddSettingAct ...
func (h *Contract) AddSettingAct(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req = request.SettingReq{}
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
	err = m.AddSetting(h.DB, ctx, code, req.SetGroup, req.SetLabel, req.SetOrder, req.ContentType, req.ContentValue, req.IsActive)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

// GetSettingDetailAct ...
func (h *Contract) GetSettingDetailAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	data, err := m.GetSettingByCode(h.DB, ctx, code, "")
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, response.SettingRes{
		SettingCode:  data.SettingCode,
		SetGroup:     data.SetGroup,
		SetKey:       data.SetKey,
		SetLabel:     data.SetLabel,
		SetOrder:     data.SetOrder,
		ContentType:  data.ContentType,
		ContentValue: data.ContentValue,
		IsActive:     data.IsActive,
	}, nil)

}

// UpdateSettingAct ...
func (h *Contract) UpdateSettingAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		req  = request.SettingReq{}
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	err = m.UpdateSetting(h.DB, ctx, code, req.SetLabel, req.SetOrder, req.ContentValue, req.IsActive)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}
