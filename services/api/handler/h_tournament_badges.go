package handler

import (
	"context"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetBadgeDetailAct ...
func (h *Contract) GetTournamentBadgeDetailAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
		res  = make([]response.BadgeRes, 0)
	)

	badges, err := m.GetBadgeDetailByParentCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response with badge details and rules
	for _, v := range badges {
		badgeRules, err := m.GetBadgeRuleList(h.DB, ctx, v.Id)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		// Create a slice to store badge rules
		var badgeRuleResSlice []response.BadgeRuleRes

		// Populate the badge rules slice
		for _, rule := range badgeRules {
			badgeRuleRes := response.BadgeRuleRes{
				BadgeRuleCode: rule.BadgeRuleCode,
				BadgeId:       rule.BadgeId,
				KeyCondition:  rule.KeyCondition,
				ValueType:     rule.ValueType,
				Value:         rule.Value,
			}
			badgeRuleResSlice = append(badgeRuleResSlice, badgeRuleRes)
		}

		res = append(res, response.BadgeRes{
			BadgeCode:     v.BadgeCode,
			BadgeCategory: v.BadgeCategory,
			Name:          v.Name,
			ImageURL:      v.ImageURL,
			VPPoint:       v.VPPoint,
			Description:   v.Description.String,
			BadgeRules:    badgeRuleResSlice,
			Status:        v.Status,
			CreatedDate:   v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
			UpdatedDate:   v.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
			DeletedDate:   v.DeletedDate.Time.Format(utils.DATE_TIME_FORMAT),
		})
	}

	h.SendSuccess(w, res, nil)
}

// AddBadgeAct ...
func (h *Contract) AddTournamentBadgeAct(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		req        = request.TournamentBadgeListReq{}
		ctx        = context.TODO()
		m          = model.Contract{App: h.App}
		parentCode = utils.GeneratePrefixCode(utils.ParentBadgePrefix)
		badgeID    int64
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
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

	for _, v := range req.TournamentBadges {
		// Generate Badge Code
		badgeCode := utils.GeneratePrefixCode(utils.BadgePrefix)

		// Add Badge
		badgeID, err = m.AddBadge(tx, ctx, badgeCode, v.BadgeCategory, v.Name, v.ImageURL, v.VPPoint, v.Status, v.Description, parentCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		// Generate Badge Rule Code
		badgeRuleCode := utils.GeneratePrefixCode(utils.BadgeRulePrefix)

		var tournamentCategory request.TournamentCategory
		valueJSON, err := json.Marshal(v.BadgeRule.Value)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
		err = json.Unmarshal(valueJSON, &tournamentCategory)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
		err = m.AddBadgeRule(tx, ctx, badgeRuleCode, badgeID, v.BadgeRule.KeyCondition, v.BadgeRule.ValueType, tournamentCategory)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

// UpdateTournamentBadgeAct updates an existing badge in the database.
func (h *Contract) UpdateTournamentBadgeAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		req  request.TournamentBadgeListReq
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	// Bind and validate request
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	_, err = m.GetBadgeDetailByParentCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
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
		} else {
			tx.Commit(ctx)
		}
	}()

	for _, v := range req.TournamentBadges {
		if !utils.Contains(utils.StatusBadges, v.Status) {
			h.SendBadRequest(w, "wrong status value for badge (active|inactive)")
			return
		}

		// Update Badge
		err = m.UpdateBadgeByCode(tx, ctx, v.BadgeCategory, v.Name, v.Description, v.ImageURL, v.Status, v.VPPoint, v.BadgeCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		badgeID, err := m.GetBadgeIdByCode(h.DB, ctx, v.BadgeCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		// Delete Badge Rule
		err = m.DeleteBadgeRule(tx, ctx, badgeID)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		// Generate Badge Rule Code
		badgeRuleCode := utils.GeneratePrefixCode(utils.BadgeRulePrefix)

		var tournamentCategory request.TournamentCategory
		valueJSON, err := json.Marshal(v.BadgeRule.Value)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
		err = json.Unmarshal(valueJSON, &tournamentCategory)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		// Add Badge Rule
		err = m.AddBadgeRule(tx, ctx, badgeRuleCode, badgeID, v.BadgeRule.KeyCondition, v.BadgeRule.ValueType, tournamentCategory)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}
