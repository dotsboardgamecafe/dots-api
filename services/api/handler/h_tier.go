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

// GetTiersListAct
func (h *Contract) GetTiersListAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.TierRes, 0)
		param = request.TierParam{}
	)

	// Parse URL query parameters
	err = param.ParseTier(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetTiersList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.TierRes{
			TierCode:    v.TierCode,
			Name:        v.Name,
			Description: v.Description.String,
			MinPoint:    v.MinPoint,
			MaxPoint:    v.MaxPoint,
			Status:      v.Status.String,
			CreatedDate: v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
			UpdatedDate: v.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
		})
	}

	h.SendSuccess(w, res, param)
}

// GetTiersDetailAct
func (h *Contract) GetTiersDetailAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	data, err := m.GetTierByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	res := response.TierRes{
		TierCode:    data.TierCode,
		Name:        data.Name,
		Description: data.Description.String,
		MinPoint:    data.MinPoint,
		MaxPoint:    data.MaxPoint,
		Status:      data.Status.String,
		CreatedDate: data.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		UpdatedDate: data.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
	}

	// Populate response
	h.SendSuccess(w, res, nil)
}

// AddTierAct
func (h *Contract) AddTierAct(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		req      request.TierReq
		ctx      = context.TODO()
		m        = model.Contract{App: h.App}
		tierCode = utils.GeneratePrefixCode(utils.TierPrefix)
	)

	// Binding and validation
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}
	if !utils.Contains(utils.StatusTier, req.Status) {
		h.SendBadRequest(w, "wrong status value for tier(active|inactive")
		return
	}
	err = m.InsertTier(h.DB, ctx, tierCode, req.TierName, req.Description, req.MinPoint, req.MaxPoint, req.Status)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

// UpdateTierAct
func (h *Contract) UpdateTierAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		req  request.TierReq
		code = chi.URLParam(r, "code") // Get tier code from URL
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	// Binding and validation
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}
	if !utils.Contains(utils.StatusTier, req.Status) {
		h.SendBadRequest(w, "wrong status value for tier(active|inactive")
		return
	}
	err = m.UpdateTier(h.DB, ctx, code, req.TierName, req.Description, req.MinPoint, req.MaxPoint, req.Status)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

// DeleteTierAct
func (h *Contract) DeleteTierAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code") // Get tier code from URL
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	err = m.DeleteTier(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}
