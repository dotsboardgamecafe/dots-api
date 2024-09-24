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

// GetPermissionListAct ...
func (h *Contract) GetPermissionListAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		param = request.PermissionParam{}
	)

	// Define urlQuery and Parse
	err = param.ParsePermission(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetPermissionList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, ToResponsePermissionList(data), param)
}

// AddPermissionAct ...
func (h *Contract) AddPermissionAct(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req = request.PermissionReq{}
		ctx = context.TODO()
		m   = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusPermission, req.Status) {
		h.SendBadRequest(w, "wrong status value for admin(active|inactive)")
		return
	}

	if !utils.Contains(utils.HTTPMethodList, req.RouteMethod) {
		h.SendBadRequest(w, "wrong route method value for permission(GET|POST|PUT|DELETE|OPTIONS)")
		return
	}

	// Generate Random Code
	code := utils.GeneratePrefixCode(utils.PermissionPrefix)
	err = m.AddPermission(h.DB, ctx, code, req.Name, req.RoutePattern, req.RouteMethod, req.Description, req.Status)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

// GetPermissionDetailAct ...
func (h *Contract) GetPermissionDetailAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	data, err := m.GetPermissionByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, response.PermissionRes{
		PermissionCode: data.PermissionCode,
		Name:           data.Name,
		RoutePattern:   data.RoutePattern,
		RouteMethod:    data.RouteMethod,
		Description:    data.Description,
		Status:         data.Status,
	}, nil)
}

// UpdatePermissionAct ...
func (h *Contract) UpdatePermissionAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		req  = request.PermissionUpdateReq{}
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusPermission, req.Status) {
		h.SendBadRequest(w, "wrong status value for permission(active|inactive)")
		return
	}

	if !utils.Contains(utils.HTTPMethodList, req.RouteMethod) {
		h.SendBadRequest(w, "wrong route method value for permission(GET|POST|PUT|DELETE|OPTIONS)")
		return
	}

	//validate permission code
	_, err = m.GetPermissionByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.UpdatePermissionByCode(h.DB, ctx, code, req.Name, req.RoutePattern, req.RouteMethod, req.Description, req.Status)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

// DeletePermissionAct ...
func (h *Contract) DeletePermissionAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	err = m.DeletePermissionByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	h.SendSuccess(w, nil, nil)

}

func ToResponsePermissionList(data []model.PermissionEnt) []response.PermissionRes {
	var res = make([]response.PermissionRes, 0)

	// Populate response
	for _, v := range data {
		res = append(res, response.PermissionRes{
			PermissionId:   v.ID,
			PermissionCode: v.PermissionCode,
			Name:           v.Name,
			RoutePattern:   v.RoutePattern,
			RouteMethod:    v.RouteMethod,
			Description:    v.Description,
			Status:         v.Status,
		})
	}

	return res
}
