package handler

import (
	"context"
	"dots-api/bootstrap"
	"dots-api/lib/rabbit"
	"dots-api/lib/upload"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetBadgeListAct ...
func (h *Contract) GetBadgeListAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.BadgeRes, 0)
		param = request.BadgeParam{}
	)

	// Parse URL query parameters
	err = param.ParseBadge(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Retrieve badge list from the badgesbase
	badges, param, err := m.GetBadgeList(h.DB, ctx, param)
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
			BadgeRules:    badgeRuleResSlice,
			Description:   v.Description.String,
			VPPoint:       v.VPPoint,
			Status:        v.Status,
			ParentCode:    v.ParentCode.String,
			CreatedDate:   v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
			UpdatedDate:   v.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
			DeletedDate:   v.DeletedDate.Time.Format(utils.DATE_TIME_FORMAT),
		})
	}

	h.SendSuccess(w, res, param)
}

// GetUnownedBadgeUserListAct ...
func (h *Contract) GetUnownedBadgeUserListAct(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		ctx      = context.TODO()
		m        = model.Contract{App: h.App}
		res      = make([]response.BadgeRes, 0)
		param    = request.UnownedBadgeParam{}
		userCode = chi.URLParam(r, "user_code")
	)

	// Parse URL query parameters
	err = param.ParseBadge(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Retrieve badge list from the badgesbase
	badges, param, err := m.GetUnownedBadgeListByUserCode(h.DB, ctx, userCode, param)
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
			BadgeRules:    badgeRuleResSlice,
			Description:   v.Description.String,
			VPPoint:       v.VPPoint,
			Status:        v.Status,
			ParentCode:    v.ParentCode.String,
			CreatedDate:   v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
			UpdatedDate:   v.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
			DeletedDate:   v.DeletedDate.Time.Format(utils.DATE_TIME_FORMAT),
		})
	}

	h.SendSuccess(w, res, param)
}

// GetBadgeDetailAct ...
func (h *Contract) GetBadgeDetailAct(w http.ResponseWriter, r *http.Request) {
	var (
		err               error
		code              = chi.URLParam(r, "code")
		ctx               = context.TODO()
		m                 = model.Contract{App: h.App}
		badgeRuleResSlice []response.BadgeRuleRes
	)

	badges, err := m.GetBadgeDetailByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	badgeRules, err := m.GetBadgeRuleList(h.DB, ctx, badges.Id)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

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

	h.SendSuccess(w, response.BadgeRes{
		BadgeCode:     badges.BadgeCode,
		BadgeCategory: badges.BadgeCategory,
		Name:          badges.Name,
		ImageURL:      badges.ImageURL,
		BadgeRules:    badgeRuleResSlice,
		Description:   badges.Description.String,
		VPPoint:       badges.VPPoint,
		Status:        badges.Status,
		CreatedDate:   badges.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		UpdatedDate:   badges.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
		DeletedDate:   badges.DeletedDate.Time.Format(utils.DATE_TIME_FORMAT),
	}, nil)
}

// AddBadgeAct ...
func (h *Contract) AddBadgeAct(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		req       = request.BadgeReq{}
		ctx       = context.TODO()
		m         = model.Contract{App: h.App}
		badgeCode = utils.GeneratePrefixCode(utils.BadgePrefix)
		badgeID   int64
		isGift    bool
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusBadges, req.Status) {
		h.SendBadRequest(w, "wrong status value for badge(active|inactive")
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

	// Add Badge
	badgeID, err = m.AddBadge(tx, ctx, badgeCode, req.BadgeCategory, req.Name, req.ImageURL, req.VPPoint, req.Status, req.Description, "")
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	isGift = req.BadgeCategory == utils.BadgeCategoryGift.String()
	// Add Badge Rules
	if !isGift {
		for _, rule := range req.BadgeRule {
			badgeRuleCode := utils.GeneratePrefixCode(utils.BadgeRulePrefix)
			if rule.KeyCondition == utils.SpesificBoardGameCategory {
				var spesificBoardGameCategory request.SpesificBoardGameCategory
				valueJSON, err := json.Marshal(rule.Value)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = json.Unmarshal(valueJSON, &spesificBoardGameCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = m.AddBadgeRule(tx, ctx, badgeRuleCode, badgeID, rule.KeyCondition, rule.ValueType, spesificBoardGameCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
			} else if rule.KeyCondition == utils.TimeLimit {
				var timeLimitCategory request.TimeLimitCategory
				valueJSON, err := json.Marshal(rule.Value)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = json.Unmarshal(valueJSON, &timeLimitCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = m.AddBadgeRule(tx, ctx, badgeRuleCode, badgeID, rule.KeyCondition, rule.ValueType, timeLimitCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
			} else if rule.KeyCondition == utils.TotalSpend {
				var totalSpendCategory int
				valueJSON, err := json.Marshal(rule.Value)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = json.Unmarshal(valueJSON, &totalSpendCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = m.AddBadgeRule(tx, ctx, badgeRuleCode, badgeID, rule.KeyCondition, rule.ValueType, totalSpendCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
			} else if rule.KeyCondition == utils.Tournament {
				var tournamentCategory request.TournamentCategory
				valueJSON, err := json.Marshal(rule.Value)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = json.Unmarshal(valueJSON, &tournamentCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = m.AddBadgeRule(tx, ctx, badgeRuleCode, badgeID, rule.KeyCondition, rule.ValueType, tournamentCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
			} else if rule.KeyCondition == utils.TournamentWon {
				var totalTournamentWon int
				valueJSON, err := json.Marshal(rule.Value)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = json.Unmarshal(valueJSON, &totalTournamentWon)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = m.AddBadgeRule(tx, ctx, badgeRuleCode, badgeID, rule.KeyCondition, rule.ValueType, totalTournamentWon)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
			} else if rule.KeyCondition == utils.PlayingGames {
				var totalPlayingDifferentGames int
				valueJSON, err := json.Marshal(rule.Value)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = json.Unmarshal(valueJSON, &totalPlayingDifferentGames)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = m.AddBadgeRule(tx, ctx, badgeRuleCode, badgeID, rule.KeyCondition, rule.ValueType, totalPlayingDifferentGames)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
			}
		}

		// check if status active will send publish check badge available
		if req.Status == "active" && !isGift {
			// Publisher badge
			queueData := rabbit.QueueDataPayload(
				rabbit.QueueBadges,
				rabbit.QueueBadgeReq(
					badgeCode,
				),
			)
			queueHost := m.Config.GetString("queue.rabbitmq.host")
			err = rabbit.PublishQueue(ctx, queueHost, queueData)
			if err != nil {
				log.Printf("Error : %s", err)
			}
		}
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

// UpdateBadgeAct updates an existing badge in the database.
func (h *Contract) UpdateBadgeAct(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		req            = request.UpdateBadgeReq{}
		badgeCode      = chi.URLParam(r, "code")
		ctx            = context.TODO()
		m              = model.Contract{App: h.App}
		isGift    bool = false
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusBadges, req.Status) {
		h.SendBadRequest(w, "wrong status value for badge(active|inactive")
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

	// Get Badge ID by Code
	badgeID, err := m.GetBadgeIdByCode(h.DB, ctx, badgeCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Update Badge
	err = m.UpdateBadgeByCode(tx, ctx, req.BadgeCategory, req.Name, req.Description, req.ImageURL, req.Status, req.VPPoint, badgeCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if len(req.BadgeRule) != 0 {
		// Delete Badge Rule
		err = m.DeleteBadgeRule(tx, ctx, badgeID)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
	}

	isGift = req.BadgeCategory == utils.BadgeCategoryGift.String()
	if !isGift {
		for _, v := range req.BadgeRule {
			// Add Badge Rule
			badgeRuleCode := utils.GeneratePrefixCode(utils.BadgeRulePrefix)
			switch v.KeyCondition {
			case utils.SpesificBoardGameCategory:
				var spesificBoardGameCategory request.SpesificBoardGameCategory
				valueJSON, err := json.Marshal(v.Value)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = json.Unmarshal(valueJSON, &spesificBoardGameCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = m.AddBadgeRule(tx, ctx, badgeRuleCode, badgeID, v.KeyCondition, v.ValueType, spesificBoardGameCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}

			case utils.TimeLimit:
				var timeLimitCategory request.TimeLimitCategory
				valueJSON, err := json.Marshal(v.Value)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = json.Unmarshal(valueJSON, &timeLimitCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = m.AddBadgeRule(tx, ctx, badgeRuleCode, badgeID, v.KeyCondition, v.ValueType, timeLimitCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}

			case utils.TotalSpend:
				var totalSpendCategory int
				valueJSON, err := json.Marshal(v.Value)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = json.Unmarshal(valueJSON, &totalSpendCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = m.AddBadgeRule(tx, ctx, badgeRuleCode, badgeID, v.KeyCondition, v.ValueType, totalSpendCategory)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}

			case utils.TournamentWon:
				var totalTournamentWon int
				valueJSON, err := json.Marshal(v.Value)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = json.Unmarshal(valueJSON, &totalTournamentWon)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = m.AddBadgeRule(tx, ctx, badgeRuleCode, badgeID, v.KeyCondition, v.ValueType, totalTournamentWon)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}

			case utils.PlayingGames:
				var totalPlayingDifferentGames int
				valueJSON, err := json.Marshal(v.Value)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = json.Unmarshal(valueJSON, &totalPlayingDifferentGames)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}
				err = m.AddBadgeRule(tx, ctx, badgeRuleCode, badgeID, v.KeyCondition, v.ValueType, totalPlayingDifferentGames)
				if err != nil {
					h.SendBadRequest(w, err.Error())
					return
				}

			default:
				h.SendBadRequest(w, "Invalid KeyCondition")
				return
			}
		}

		// check if status active will send publish check badge available
		if req.Status == "active" {
			// Publisher badge
			queueData := rabbit.QueueDataPayload(
				rabbit.QueueBadges,
				rabbit.QueueBadgeReq(
					badgeCode,
				),
			)
			queueHost := m.Config.GetString("queue.rabbitmq.host")
			err = rabbit.PublishQueue(ctx, queueHost, queueData)
			if err != nil {
				log.Printf("Error : %s", err)
			}
		}
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

// DeleteBadgeAct ...
func (h *Contract) DeleteBadgeAct(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		badgeCode = chi.URLParam(r, "code")
		ctx       = context.TODO()
		m         = model.Contract{App: h.App}
	)

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

	// Get Badge ID by Code
	badgeID, err := m.GetBadgeIdByCode(h.DB, ctx, badgeCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.DeleteBadge(tx, ctx, badgeID)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.DeleteBadgeRule(tx, ctx, badgeID)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

// GiftBadgeToUserAct...
func (h *Contract) GiftBadgeToUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		badgeCode = chi.URLParam(r, "code")
		req       request.AddBadgeToUserReq
		ctx       = context.TODO()
		m         = model.Contract{App: h.App}
	)

	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	badge, err := m.GetBadgeDetailByCode(h.DB, ctx, badgeCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	userId, err := m.GetUserIdByUserCode(h.DB, ctx, req.UserCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	adminId, err := m.GetAdminIdByCode(h.DB, ctx, bootstrap.GetIdentifierCodeFromToken(ctx, r))
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.AddBadgeToUser(h.DB, ctx, badge.Id, userId, adminId)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

func (h *Contract) ImportCsvToUpdateBadgesAct(w http.ResponseWriter, r *http.Request) {
	var (
		ctx        = context.TODO()
		info       = new(upload.Info)
		name       = "file"
		allowedExt = []string{"csv"}
		m          = model.Contract{App: h.App}
	)

	info.MaxSize = 10
	file, _, err := info.MultipartHandler(w, r, name, allowedExt)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	defer file.Close()

	// Reset the file pointer to the beginning after the MultipartHandler has read it
	if _, err := file.Seek(0, 0); err != nil {
		h.SendBadRequest(w, utils.ErrToReadCSVFile)
		return
	}

	var records []struct {
		BadgeCode string
		VPPoint   string
	}

	fields := map[string]string{
		"badge_code": "BadgeCode",
		"vp_point":   "VPPoint",
	}

	err = utils.ParseCSVToStruct(file, &records, fields)
	if err != nil {
		h.SendBadRequest(w, utils.ErrFailedParseCSVDataToStruct)
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

	for _, v := range records {
		// Get Badge ID by Code
		badge, err := m.GetBadgeDetailByCode(h.DB, ctx, v.BadgeCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		vp, err := strconv.ParseInt(v.VPPoint, 10, 64)
		if err != nil {
			h.SendBadRequest(w, utils.ErrImportingBadgeData)
			return
		}

		// Update Badge
		err = m.UpdateBadgeByCode(tx, ctx, badge.BadgeCategory, badge.Name, badge.Description.String, badge.ImageURL, badge.Status, vp, badge.BadgeCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
	}

	for _, v := range records {
		badge, err := m.GetBadgeDetailByCode(h.DB, ctx, v.BadgeCode)
		if err != nil {
			continue
		}

		// check if status active will send publish check badge available
		if badge.Status == "active" {
			// Publisher badge
			queueData := rabbit.QueueDataPayload(
				rabbit.QueueBadges,
				rabbit.QueueBadgeReq(
					badge.BadgeCode,
				),
			)
			queueHost := m.Config.GetString("queue.rabbitmq.host")
			err = rabbit.PublishQueue(ctx, queueHost, queueData)
			if err != nil {
				log.Printf("Error : %s", err)
			}
		}
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

func (h *Contract) DistributeBadgeToUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		param = request.BadgeParam{}
	)

	// Parse URL query parameters
	err = param.ParseBadge(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Retrieve badge list from the badgesbase
	badges, param, err := m.GetBadgeList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	queueHost := m.Config.GetString("queue.rabbitmq.host")
	// Populate response with badge details and rules
	states := []map[string]interface{}{}
	for _, v := range badges {
		if v.Status != "active" {
			continue
		}

		rules, err := m.GetBadgeRuleList(h.DB, ctx, v.Id)
		if err != nil {
			continue
		}

		isTournamentWinner := slices.IndexFunc(rules, func(rule model.BadgeRuleEnt) bool {
			return rule.KeyCondition == utils.TournamentWinner
		})

		if isTournamentWinner != -1 {
			continue
		}

		// Publisher badge
		queueData := rabbit.QueueDataPayload(
			rabbit.QueueBadges,
			rabbit.QueueBadgeReq(
				v.BadgeCode,
			),
		)
		err = rabbit.PublishQueue(ctx, queueHost, queueData)
		if err != nil {
			log.Printf("Error : %s", err)
			continue
		}
		states = append(states, map[string]interface{}{"badge_code": v.BadgeCode, "status": "sent"})
	}

	h.SendSuccess(w, map[string]interface{}{"message": "success", "data": states}, param)
}
