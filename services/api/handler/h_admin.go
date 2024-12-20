package handler

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetAdminListAct ...
func (h *Contract) GetAdminListAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		param = request.AdminParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseAdmin(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetAdminList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, ToResponseAdminList(data), param)
}

// AddAdminAct ...
func (h *Contract) AddAdminAct(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req = request.AdminReq{}
		ctx = context.TODO()
		m   = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusAdmin, req.Status) {
		h.SendBadRequest(w, "wrong status value for admin(active|inactive")
		return
	}

	// Check if username exist (only if username set)
	if req.UserName != "" {
		isExist, _ := m.IsAdminUsernameExist(h.DB, ctx, req.UserName)

		if isExist {
			h.SendBadRequest(w, utils.ErrUsernameAlreadyRegistered)
			return
		}
	}

	//validate email
	dataAdmin, err := m.GetAdminByEmail(h.DB, ctx, req.Email)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if len(dataAdmin.AdminCode) > 0 {
		h.SendBadRequest(w, "email already used")
		return
	}

	//validate phone
	dataAdmin, err = m.GetAdminByPhoneNumber(h.DB, ctx, req.PhoneNumber)
	if err != nil {
		if err != sql.ErrNoRows {
			h.SendBadRequest(w, err.Error())
			return
		}
	}

	if len(dataAdmin.AdminCode) > 0 {
		h.SendBadRequest(w, "phone_number already used")
		return
	}

	// Generate Random Code
	code := utils.GeneratePrefixCode(utils.AdminPrefix)
	err = m.AddAdmin(h.DB, ctx, code, req.Email, req.Name, req.UserName, req.Password, req.Status, req.Role, req.PhoneNumber, req.ImageUrl)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

// GetAdminDetailAct ...
func (h *Contract) GetAdminDetailAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	data, err := m.GetAdminByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, response.AdminRes{
		AdminCode:   data.AdminCode,
		Email:       data.Email,
		Name:        data.Name,
		UserName:    data.UserName,
		Status:      data.Status,
		ImageURL:    data.ImageURL,
		PhoneNumber: data.PhoneNumber,
		Role:        data.Role,
	}, nil)
}

// UpdateAdminAct ...
func (h *Contract) UpdateAdminAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		req  = request.AdminUpdateReq{}
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusAdmin, req.Status) {
		h.SendBadRequest(w, "wrong status value for admin(active|inactive")
		return
	}

	// Check if username exist (only if username set)
	if req.UserName != "" {
		isExist, _ := m.IsAdminUsernameExist(h.DB, ctx, req.UserName)

		if isExist {
			h.SendBadRequest(w, utils.ErrUsernameAlreadyRegistered)
			return
		}
	}

	err = m.UpdateAdminByCode(h.DB, ctx, code, req.Email, req.Name, req.UserName, req.Password, req.Status, req.Role, req.PhoneNumber, req.ImageUrl)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

// DeleteAdminAct ...
func (h *Contract) DeleteAdminAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	err = m.DeleteAdminByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	h.SendSuccess(w, nil, nil)

}

func (h *Contract) UpdateAdminStatus(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		ctx       = context.TODO()
		req       = request.UpdateStatusAdminReq{}
		m         = model.Contract{App: h.App}
		adminCode = chi.URLParam(r, "code")
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err.Error())
		return
	}

	if !utils.Contains(utils.StatusAdmin, req.Status) {
		h.SendBadRequest(w, "wrong status value for admin(active|inactive)")
		return
	}

	err = m.UpdateAdminStatus(h.DB, ctx, adminCode, req.Status)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

func ToResponseAdminList(data []model.AdminEnt) []response.AdminRes {
	var res = make([]response.AdminRes, 0)

	// Populate response
	for _, v := range data {
		res = append(res, response.AdminRes{
			AdminCode:   v.AdminCode,
			Email:       v.Email,
			UserName:    v.UserName,
			Name:        v.Name,
			Status:      v.Status,
			ImageURL:    v.ImageURL,
			PhoneNumber: v.PhoneNumber,
			Role:        v.Role,
		})
	}

	return res
}
