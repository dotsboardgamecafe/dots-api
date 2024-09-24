package handler

import (
	"context"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetCafeListAct ...
func (h *Contract) GetCafeListAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.CafeRes, 0)
		param = request.CafeParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseCafe(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetCafeList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.CafeRes{
			CafeCode:    v.CafeCode,
			Name:        v.Name,
			Address:     v.Address,
			Description: v.Description,
			Status:      v.Status,
			Province:    v.Province,
			City:        v.City,
		})
	}

	h.SendSuccess(w, res, param)
}

// AddCafeAct ...
func (h *Contract) AddCafeAct(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req = request.CafeReq{}
		ctx = context.TODO()
		m   = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusCafe, req.Status) {
		h.SendBadRequest(w, "wrong status value for cafe(active|inactive")
		return
	}

	//validate province
	provinces, err := m.GetSettingList(h.DB, ctx, request.SettingParam{
		SetGroup: "province",
	})
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	provinceNames := []string{}
	for _, v := range provinces {
		provinceNames = append(provinceNames, v.ContentValue)
	}

	if !utils.Contains(provinceNames, req.Province) {
		h.SendBadRequest(w, fmt.Sprintf("wrong province value for cafe(%v)", provinceNames))
		return
	}

	//validate city
	cities, err := m.GetSettingList(h.DB, ctx, request.SettingParam{
		SetGroup: "city",
	})
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	cityNames := []string{}
	for _, v := range cities {
		cityNames = append(cityNames, v.ContentValue)
	}

	if !utils.Contains(cityNames, req.City) {
		h.SendBadRequest(w, fmt.Sprintf("wrong province value for cafe(%v)", cityNames))
		return
	}

	// Generate Random Code
	code := utils.GeneratePrefixCode(utils.CafePrefix)
	err = m.AddCafe(h.DB, ctx, code, req.Name, req.Address, req.Description, req.Status, req.Province, req.City)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, response.CafeRes{
		CafeCode:    code,
		Name:        req.Name,
		Address:     req.Address,
		Description: req.Description,
		Status:      req.Status,
		Province:    req.Province,
		City:        req.City,
	}, nil)
}

// GetCafeDetailAct ...
func (h *Contract) GetCafeDetailAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	data, err := m.GetCafeByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, response.CafeRes{
		CafeCode:    data.CafeCode,
		Name:        data.Name,
		Address:     data.Address,
		Description: data.Description,
		Status:      data.Status,
		Province:    data.Province,
		City:        data.City,
	}, nil)
}

// UpdateCafeAct ...
func (h *Contract) UpdateCafeAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		req  = request.CafeReq{}
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusCafe, req.Status) {
		h.SendBadRequest(w, "wrong status value for cafe(active|inactive")
		return
	}

	err = m.UpdateCafeByCode(h.DB, ctx, code, req.Name, req.Address, req.Description, req.Status, req.Province, req.City)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, response.CafeRes{
		CafeCode:    code,
		Name:        req.Name,
		Address:     req.Address,
		Description: req.Description,
		Status:      req.Status,
	}, nil)
}

// DeleteCafeAct ...
func (h *Contract) DeleteCafeAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	err = m.DeleteCafeByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	h.SendSuccess(w, nil, nil)

}
