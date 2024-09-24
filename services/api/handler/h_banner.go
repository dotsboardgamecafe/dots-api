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

// GetBannerListAct ...
func (h *Contract) GetBannerListAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.BannerRes, 0)
		param = request.BannerParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseBanner(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetBannerList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.BannerRes{
			BannerCode:  v.BannerCode,
			BannerType:  v.BannerType,
			Title:       v.Title,
			Description: v.Description,
			ImageURL:    v.ImageURL,
			Status:      v.Status,
			CreatedDate: v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
			UpdatedDate: v.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
		})
	}

	h.SendSuccess(w, res, param)
}

// GetBannerDetailAct ...
func (h *Contract) GetBannerDetailAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	data, err := m.GetBannerByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, response.BannerRes{
		BannerCode:  data.BannerCode,
		BannerType:  data.BannerType,
		Title:       data.Title,
		Description: data.Description,
		ImageURL:    data.ImageURL,
		Status:      data.Status,
		CreatedDate: data.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		UpdatedDate: data.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
	}, nil)
}

// AddBannerAct ...
func (h *Contract) AddBannerAct(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req = request.BannerReq{}
		ctx = context.TODO()
		m   = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusBanner, req.Status) {
		h.SendBadRequest(w, "wrong status value for banner (publish|unpublish)")
		return
	}

	// Generate Random Code
	code := utils.GeneratePrefixCode(utils.BannerPrefix)
	err = m.AddBanner(h.DB, ctx, code, req.Title, req.Description, req.BannerType, req.ImageURL, req.Status)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

// UpdateBannerAct ...
func (h *Contract) UpdateBannerAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		req  = request.BannerReq{}
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusBanner, req.Status) {
		h.SendBadRequest(w, "wwrong status value for banner (publish|unpublish)")
		return
	}

	err = m.UpdateBannerByCode(h.DB, ctx, code, req.BannerType, req.Title, req.Description, req.ImageURL, req.Status)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

// DeleteBannerAct ...
func (h *Contract) DeleteBannerAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	err = m.DeleteBannerByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	h.SendSuccess(w, nil, nil)

}
