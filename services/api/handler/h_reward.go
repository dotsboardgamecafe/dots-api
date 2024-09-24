package handler

import (
	"context"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// GetRewardsListAct
func (h *Contract) GetRewardsListAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.RewardRes, 0)
		param = request.RewardParam{}
	)

	// Parse URL query parameters
	err = param.ParseReward(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetRewardList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		tierRes := response.TierRes{
			TierCode:    v.Tier.TierCode,
			Name:        v.Tier.Name,
			Description: v.Tier.Description.String,
			MinPoint:    v.Tier.MinPoint,
			MaxPoint:    v.Tier.MaxPoint,
			Status:      v.Tier.Status.String,
			CreatedDate: v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
			UpdatedDate: v.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
		}
		res = append(res, response.RewardRes{
			Tier:         tierRes,
			Name:         v.Name,
			ImageUrl:     v.ImageUrl,
			CategoryType: v.CategoryType,
			RewardCode:   v.RewardCode,
			Description:  v.Description.String,
			Status:       v.Status,
			ExpiredDate:  v.ExpiredDate.Time.Format(utils.DATE_TIME_FORMAT),
			CreatedDate:  v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
			UpdatedDate:  v.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
		})
	}

	h.SendSuccess(w, res, param)
}

// GetRewardDetailAct
func (h *Contract) GetRewardDetailAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	data, err := m.GetRewardByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	tierRes := response.TierRes{
		TierCode:    data.Tier.TierCode,
		Name:        data.Tier.Name,
		Description: data.Tier.Description.String,
		MinPoint:    data.Tier.MinPoint,
		MaxPoint:    data.Tier.MaxPoint,
		Status:      data.Tier.Status.String,
		CreatedDate: data.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		UpdatedDate: data.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
	}
	res := response.RewardRes{
		RewardCode:   data.RewardCode,
		Tier:         tierRes,
		Name:         data.Name,
		ImageUrl:     data.ImageUrl,
		CategoryType: data.CategoryType,
		Description:  data.Description.String,
		Status:       data.Status,
		ExpiredDate:  data.ExpiredDate.Time.Format(utils.DATE_TIME_FORMAT),
		CreatedDate:  data.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		UpdatedDate:  data.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
	}

	h.SendSuccess(w, res, nil)
}

// AddRewardAct
func (h *Contract) AddRewardAct(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		req         request.RewardReq
		ctx         = context.TODO()
		m           = model.Contract{App: h.App}
		expiredDate interface{}
		rewardCode  = utils.GeneratePrefixCode(utils.RewardPrefix)
	)

	// Binding and validation
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusReward, req.Status) {
		h.SendBadRequest(w, "wrong status value for reward(active|inactive")
		return
	}

	tierId, err := m.GetTierIdByCode(h.DB, ctx, req.TierCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if req.ExpiredDate != "" {
		// Convert time strings to time.Time format
		expiredDate, err = time.Parse(time.DateOnly, req.ExpiredDate)
		if err != nil {
			h.SendBadRequest(w, utils.ErrInvalidExpiredDate)
			return
		}
	}

	// Add new reward
	err = m.AddReward(h.DB, ctx, tierId, req.Name, req.ImageUrl, req.CategoryType, rewardCode, req.Description, req.Status, expiredDate)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

// UpdateRewardAct
func (h *Contract) UpdateRewardAct(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		req         request.RewardReq
		code        = chi.URLParam(r, "code") // Get reward code from URL
		ctx         = context.TODO()
		m           = model.Contract{App: h.App}
		expiredDate interface{}
	)

	// Binding and validation
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	tierId, err := m.GetTierIdByCode(h.DB, ctx, req.TierCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if req.ExpiredDate != "" {
		// Convert time strings to time.Time format
		expiredDate, err = time.Parse(time.DateOnly, req.ExpiredDate)
		if err != nil {
			h.SendBadRequest(w, utils.ErrInvalidExpiredDate)
			return
		}
	}

	// Update existing reward
	err = m.UpdateReward(h.DB, ctx, code, tierId, req.Name, req.ImageUrl, req.CategoryType, req.Description, req.Status, expiredDate)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

// DeleteRewardAct
func (h *Contract) DeleteRewardAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code") // Get reward code from URL
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	// Delete existing reward
	err = m.DeleteReward(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}
