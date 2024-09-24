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

// GetRoleListAct ...
func (h *Contract) GetRoleListAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		param = request.RoleParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseRole(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetRoleList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, ToResponseRoleList(data), param)
}

// AddRoleAct ...
func (h *Contract) AddRoleAct(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req = request.RoleReq{}
		ctx = context.TODO()
		m   = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusRole, req.Status) {
		h.SendBadRequest(w, "wrong status value for admin(active|inactive")
		return
	}

	// Start a transaction
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
		tx.Commit(ctx)
	}()

	// Generate Random Code
	code := utils.GeneratePrefixCode(utils.RolePrefix)
	roleId, err := m.AddRole(tx, ctx, code, req.Name, req.Description, req.Status)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	//insert data role permission
	if len(req.PermissionIds) > 0 {
		//delete existing data role permission if any
		err = m.DeleteRolePermissionByRoleId(tx, ctx, roleId)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		//insert data role permission if any
		for _, v := range req.PermissionIds {
			err = m.AddRolePermission(tx, ctx, roleId, v)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}
		}
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

// GetRoleDetailAct ...
func (h *Contract) GetRoleDetailAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	role, err := m.GetRoleByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	//get date permission
	permissions, err := m.GetRolePermissionByRoleId(h.DB, ctx, role.ID)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, response.RoleRes{
		RoleCode:    role.RoleCode,
		Name:        role.Name,
		Description: role.Description,
		Status:      role.Status,
		Permissions: ToResponsePermissionList(permissions),
	}, nil)
}

// UpdateRoleAct ...
func (h *Contract) UpdateRoleAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		req  = request.RoleUpdateReq{}
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusRole, req.Status) {
		h.SendBadRequest(w, "wrong status value for permission(active|inactive)")
		return
	}

	// Start a transaction
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
		tx.Commit(ctx)
	}()

	//validate role code
	role, err := m.GetRoleByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.UpdateRoleByCode(tx, ctx, code, req.Name, req.Description, req.Status)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	//insert data role permission
	if len(req.PermissionIds) > 0 {
		//delete existing data role permission if any
		err = m.DeleteRolePermissionByRoleId(tx, ctx, role.ID)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		//insert data role permission if any
		for _, v := range req.PermissionIds {
			err = m.AddRolePermission(tx, ctx, role.ID, v)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}
		}
	}

	h.SendSuccess(w, nil, nil)
}

// DeleteRoleAct ...
func (h *Contract) DeleteRoleAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	err = m.DeleteRoleByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	h.SendSuccess(w, nil, nil)

}

func ToResponseRoleList(data []model.RoleEnt) []response.RoleRes {
	var res = make([]response.RoleRes, 0)

	// Populate response
	for _, v := range data {
		res = append(res, response.RoleRes{
			RoleCode:    v.RoleCode,
			Name:        v.Name,
			Description: v.Description,
			Status:      v.Status,
		})
	}

	return res
}
