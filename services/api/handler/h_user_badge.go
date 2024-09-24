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

// GetUserBadgeListAct ...
func (h *Contract) GetUserBadgeListAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.UserBadgeRes, 0)
		code  = chi.URLParam(r, "code")
		param = request.UserBadgeParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseUserBadge(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetUserBadgeList(h.DB, ctx, code, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.UserBadgeRes{
			BadgeId:       v.BadgeId,
			UserId:        v.UserId.Int64,
			BadgeName:     v.BadgeName,
			BadgeImageURL: v.BadgeImageURL,
			BadgeCode:     v.BadgeCode,
			BadgeCategory: v.BadgeCategory,
			Description:   v.Description.String,
			VPPoint:       v.VPPoint,
			IsClaim:       v.IsClaim.Bool,
			CreatedDate:   v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
			IsBadgeOwned:  v.IsBadgeOwned.Bool,
			NeedToClaim:   v.NeedToClaim.Bool,
		})
	}

	h.SendSuccess(w, res, param)
}

// GetUserBadgeByBadgeCodeAct ...
func (h *Contract) GetUserBadgeByBadgeCodeAct(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		ctx       = context.TODO()
		m         = model.Contract{App: h.App}
		res       = response.UserBadgeRes{}
		code      = chi.URLParam(r, "code")
		badgeCode = chi.URLParam(r, "badge-code")
	)

	v, err := m.GetUserBadgeByBadgeCode(h.DB, ctx, code, badgeCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	res = response.UserBadgeRes{
		BadgeId:       v.BadgeId,
		UserId:        v.UserId.Int64,
		BadgeName:     v.BadgeName,
		BadgeImageURL: v.BadgeImageURL,
		BadgeCode:     v.BadgeCode,
		BadgeCategory: v.BadgeCategory,
		Description:   v.Description.String,
		VPPoint:       v.VPPoint,
		IsClaim:       v.IsClaim.Bool,
		CreatedDate:   v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		IsBadgeOwned:  v.IsBadgeOwned.Bool,
		NeedToClaim:   v.NeedToClaim.Bool,
	}

	h.SendSuccess(w, res, nil)
}

// UpdateUserBadgeByBadgeCodeAct ...
func (h *Contract) UpdateUserBadgeByBadgeCodeAct(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		ctx       = context.TODO()
		m         = model.Contract{App: h.App}
		req       request.IsClaimBadgeReq
		code      = chi.URLParam(r, "code")
		badgeCode = chi.URLParam(r, "badge-code")
	)

	// Db tx start
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	userId, err := m.GetUserIdByUserCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Get Badge ID by Code
	badges, err := m.GetBadgeDetailByCode(h.DB, ctx, badgeCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.UpdateUserBadge(tx, ctx, int64(userId), badges.Id, req.IsClaim)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Add user point for all the winners
	err = m.AddUserPoint(tx, ctx, userId, utils.Badge, badgeCode, int(badges.VPPoint))
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Db tx commit
	err = tx.Commit(ctx)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		tx.Rollback(ctx)
		return
	}

	h.SendSuccess(w, nil, nil)
}
