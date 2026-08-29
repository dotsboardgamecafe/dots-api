package handler

import (
	"context"
	"dots-api/bootstrap"
	"dots-api/lib/onesignal"
	"dots-api/lib/rabbit"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// GetTournamentList
func (h *Contract) GetTournamentList(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.TournamentRes, 0)
		param = request.TournamentParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseTournament(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetTournamentList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		dataTournamentBadge, err := m.GetTournamentBadgeByTournamentCode(h.DB, ctx, v.TournamentCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		// Populate resTournamentBadge
		resTournamentBadge := make([]response.BadgeRes, 0)
		for _, badge := range dataTournamentBadge {
			var badgeRuleResSlice []response.BadgeRuleRes
			badgeRules, err := m.GetBadgeRuleList(h.DB, ctx, badge.Id)
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

			resBadge := response.BadgeRes{
				BadgeCode:     badge.BadgeCode,
				BadgeCategory: badge.BadgeCategory,
				Name:          badge.Name,
				ImageURL:      badge.ImageURL,
				BadgeRules:    badgeRuleResSlice,
				VPPoint:       badge.VPPoint,
				Status:        badge.Status,
				CreatedDate:   badge.CreatedDate.Format(utils.DATE_TIME_FORMAT),
				UpdatedDate:   badge.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
				DeletedDate:   badge.DeletedDate.Time.Format(utils.DATE_TIME_FORMAT),
			}
			resTournamentBadge = append(resTournamentBadge, resBadge)
		}

		res = append(res, response.TournamentRes{
			GameCode:         v.GameCode,
			GameName:         v.GameName,
			GameImgUrl:       v.GameImgUrl,
			CafeCode:         v.CafeCode,
			CafeName:         v.CafeName,
			CafeAddress:      v.CafeAddress,
			TournamentCode:   v.TournamentCode,
			ImageUrl:         v.ImageUrl.String,
			PrizesImgUrl:     v.PrizesImgUrl.String,
			Name:             v.Name.String,
			TournamentRules:  v.TournamentRules,
			BookingPrice:     v.BookingPrice,
			Difficulty:       v.Level,
			StartDate:        v.StartDate.Time.Format(utils.DATE_FORMAT),
			EndDate:          v.EndDate.Time.Format(utils.DATE_FORMAT),
			StartTime:        v.StartTime.Format(utils.TIME_FORMAT),
			EndTime:          v.EndTime.Format(utils.TIME_FORMAT),
			PlayerSlot:       v.PlayerSlot,
			ParticipantVP:    v.ParticipantVP,
			Status:           v.Status,
			DayPastEndDate:   v.DayPastEndDate.Float64,
			CurrentUsedSlot:  v.CurrentUsedSlot,
			CreatedDate:      v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
			UpdatedDate:      v.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
			TournamentBadges: resTournamentBadge,
		})
	}

	h.SendSuccess(w, res, param)
}

// GetTournamentDetailAct
func (h *Contract) GetTournamentDetailAct(w http.ResponseWriter, r *http.Request) {
	var (
		err                      error
		code                     = chi.URLParam(r, "code")
		ctx                      = context.TODO()
		m                        = model.Contract{App: h.App}
		resTournamentParticipant = make([]response.TournamentParticipantRes, 0)
		resTournamentBadge       = make([]response.BadgeRes, 0)
		userCode                 = bootstrap.GetIdentifierCodeFromToken(ctx, r)
		haveJoinedTournament     = false
	)

	user, err := m.GetUserByUserCode(h.DB, ctx, userCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	dataTournament, err := m.GetTournamentByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	dataTournamentParticipant, err := m.GetAllParticipantByTournamentCode(h.DB, ctx, dataTournament.TournamentCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	dataTournamentBadge, err := m.GetTournamentBadgeByTournamentCode(h.DB, ctx, dataTournament.TournamentCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Check if already booked
	participant, err := m.GetOneTournamentParticipant(h.DB, ctx, dataTournament.TournamentId, int64(user.ID))
	if err != nil && err.Error() != utils.EmptyData {
		h.SendBadRequest(w, err.Error())
		return
	}

	if participant.Id > 0 && participant.Status == "active" {
		haveJoinedTournament = true
	}

	// Populate resTournamentParticipant
	for _, participant := range dataTournamentParticipant {
		resParticipant := response.TournamentParticipantRes{
			UserCode:       participant.UserCode,
			UserName:       participant.UserName,
			UserImgUrl:     participant.UserImgUrl,
			UserStyle:      response.UserStyleRes{Color: participant.UserStyle.Color},
			StatusWinner:   participant.StatusWinner,
			Status:         participant.Status,
			AdditionalInfo: participant.AdditionalInfo.String,
			Position:       participant.Position,
			RewardPoint:    int(participant.RewardPoint.Int64),
			LatestTier:     participant.LatestTier.String,
		}
		resTournamentParticipant = append(resTournamentParticipant, resParticipant)
	}

	// Populate resTournamentBadge
	for _, badge := range dataTournamentBadge {
		var badgeRuleResSlice []response.BadgeRuleRes
		badgeRules, err := m.GetBadgeRuleList(h.DB, ctx, badge.Id)
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

		resBadge := response.BadgeRes{
			BadgeCode:     badge.BadgeCode,
			BadgeCategory: badge.BadgeCategory,
			Name:          badge.Name,
			ImageURL:      badge.ImageURL,
			BadgeRules:    badgeRuleResSlice,
			VPPoint:       badge.VPPoint,
			Status:        badge.Status,
			CreatedDate:   badge.CreatedDate.Format(utils.DATE_TIME_FORMAT),
			UpdatedDate:   badge.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
			DeletedDate:   badge.DeletedDate.Time.Format(utils.DATE_TIME_FORMAT),
		}
		resTournamentBadge = append(resTournamentBadge, resBadge)
	}

	// Populate response
	tournamentRes := response.TournamentRes{
		GameCode:               dataTournament.GameCode,
		GameName:               dataTournament.GameName,
		GameImgUrl:             dataTournament.GameImgUrl,
		CafeCode:               dataTournament.CafeCode,
		CafeName:               dataTournament.CafeName,
		CafeAddress:            dataTournament.CafeAddress,
		TournamentCode:         dataTournament.TournamentCode,
		PrizesImgUrl:           dataTournament.PrizesImgUrl.String,
		ImageUrl:               dataTournament.ImageUrl.String,
		BookingPrice:           dataTournament.BookingPrice,
		Name:                   dataTournament.Name.String,
		TournamentRules:        dataTournament.TournamentRules,
		Difficulty:             dataTournament.Level,
		StartDate:              dataTournament.StartDate.Time.Format(utils.DATE_FORMAT),
		EndDate:                dataTournament.EndDate.Time.Format(utils.DATE_FORMAT),
		StartTime:              dataTournament.StartTime.Format(utils.TIME_FORMAT),
		EndTime:                dataTournament.EndTime.Format(utils.TIME_FORMAT),
		PlayerSlot:             dataTournament.PlayerSlot,
		ParticipantVP:          dataTournament.ParticipantVP,
		Status:                 dataTournament.Status,
		CurrentUsedSlot:        dataTournament.CurrentUsedSlot,
		DayPastEndDate:         dataTournament.DayPastEndDate.Float64,
		CreatedDate:            dataTournament.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		UpdatedDate:            dataTournament.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
		DeletedDate:            dataTournament.DeletedDate.Time.Format(utils.DATE_TIME_FORMAT),
		TournamentParticipants: resTournamentParticipant,
		TournamentBadges:       resTournamentBadge,
		HaveJoined:             haveJoinedTournament,
	}

	h.SendSuccess(w, tournamentRes, nil)
}

// AddTournament
func (h *Contract) AddTournamentAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		req            = request.TournamentReq{}
		ctx            = context.TODO()
		tournamentCode = utils.GeneratePrefixCode(utils.TournamentPrefix)
		m              = model.Contract{App: h.App}
		tournamentId   int64
		startDate      interface{}
		endDate        interface{}
		startTime      interface{}
		endTime        interface{}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusTournament, req.Status) {
		h.SendBadRequest(w, "wrong status value for tournament(active|inactive")
		return
	}

	// Start a transaction
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// For transaction
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		tx.Commit(ctx)
	}()

	gameId, err := m.GetGameIdByCode(h.DB, ctx, req.GameCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	gameEnt, err := m.GetGameByCode(h.DB, ctx, req.GameCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Validate LocationCity
	locationCity, err := m.GetCafeLocationCityByCode(h.DB, ctx, req.LocationCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if req.StartDate != "" {
		// Convert start date string to time.Time
		startDate, err = time.Parse(time.DateOnly, req.StartDate)
		if err != nil {
			h.SendBadRequest(w, utils.ErrInvalidStartDateFormat)
			return
		}
	}

	if req.EndDate != "" {
		// Convert end date string to time.Time
		endDate, err = time.Parse(time.DateOnly, req.EndDate)
		if err != nil {
			h.SendBadRequest(w, utils.ErrInvalidEndDateFormat)
			return
		}
	}

	if req.StartTime != "" {
		// Convert start time string to time.Time
		startTime, err = time.Parse(time.TimeOnly, req.StartTime)
		if err != nil {
			h.SendBadRequest(w, utils.ErrInvalidStartTimeFormat)
			return
		}
	}

	if req.EndTime != "" {
		// Convert end time string to time.Time
		endTime, err = time.Parse(time.TimeOnly, req.EndTime)
		if err != nil {
			h.SendBadRequest(w, utils.ErrInvalidEndTimeFormat)
			return
		}
	}

	// Add tournament to the database
	tournamentId, err = m.AddTournament(
		tx,
		ctx,
		gameId,
		tournamentCode,
		req.ImageUrl,
		req.Name,
		req.TournamentRules,
		req.Level,
		req.PrizesImgUrl,
		req.BookingPrice,
		startDate,
		endDate,
		startTime,
		endTime,
		req.PlayerSlot,
		req.ParticipantVP,
		req.Status,
		locationCity,
	)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Iterate through req.BadgeCodes and insert tournament badges
	for _, badgeCode := range req.BadgeCodes {
		badgeId, err := m.GetBadgeIdByCode(h.DB, ctx, badgeCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
		_, err = m.InsertTournamentBadge(tx, ctx, tournamentId, badgeId)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
	}

	go func() {
		dataListUser, err := m.GetAllUsers(m.DB, ctx)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		onesignal := onesignal.New(m.App)
		for _, user := range dataListUser {

			// Generate Notification code
			notifCode := utils.GeneratePrefixCode(utils.NotifPrefix)

			description := response.NotificationTournamentResp{
				StartDate:   startDate.(time.Time).Format("2006-01-02"),
				StartTime:   startTime.(time.Time).Format("15:04:05"),
				EndTime:     endTime.(time.Time).Format("15:04:05"),
				CafeName:    gameEnt.CafeName,
				GameName:    gameEnt.Name,
				CafeAddress: gameEnt.CafeAddress,
				Level:       req.Level,
			}

			descriptionJSON, err := json.Marshal(description)
			if err != nil {
				h.WriteLog("h.AddTournament: ", err)
				continue
			}

			OSDescription := utils.TournamentReminderPushNotificationDescription + "\n\n" + "Tournament name: " + req.Name + "\n" + "Date: " + startDate.(time.Time).Format("02 January 2006") + "\n" + "Location: " + description.CafeName
			_, err = onesignal.CreateOSNotifications(user.XPlayer, req.Name, OSDescription, utils.Tournament)
			if err != nil {
				h.WriteLog("h.AddTournament: ", err)
				continue
			}

			err = m.AddNotification(m.DB, ctx, notifCode, "user", user.UserCode, tournamentCode, utils.TournamentReminder, req.Name, descriptionJSON, req.ImageUrl)
			if err != nil {
				h.WriteLog("h.AddTournament: ", err)
			}
		}
	}()

	h.SendSuccess(w, nil, nil)
}

// UpdateTournament
func (h *Contract) UpdateTournamentAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		req            = request.TournamentReq{}
		tournamentCode = chi.URLParam(r, "code")
		ctx            = context.TODO()
		startDate      interface{}
		endDate        interface{}
		startTime      interface{}
		endTime        interface{}
		m              = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if isTournamentStatusAttrInvalid(w, h, req.Status) {
		return
	}

	// Start a transaction
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// For transaction
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		tx.Commit(ctx)
	}()

	if req.StartDate != "" {
		// Convert start date string to time.Time
		startDate, err = time.Parse(time.DateOnly, req.StartDate)
		if err != nil {
			h.SendBadRequest(w, utils.ErrInvalidStartDateFormat)
			return
		}
	}

	if req.EndDate != "" {
		// Convert end date string to time.Time
		endDate, err = time.Parse(time.DateOnly, req.EndDate)
		if err != nil {
			h.SendBadRequest(w, utils.ErrInvalidEndDateFormat)
			return
		}
	}

	if req.StartTime != "" {
		// Convert start time string to time.Time
		startTime, err = time.Parse(time.TimeOnly, req.StartTime)
		if err != nil {
			h.SendBadRequest(w, utils.ErrInvalidStartTimeFormat)
			return
		}
	}

	if req.EndTime != "" {
		// Convert end time string to time.Time
		endTime, err = time.Parse(time.TimeOnly, req.EndTime)
		if err != nil {
			h.SendBadRequest(w, utils.ErrInvalidEndTimeFormat)
			return
		}
	}

	tournament, err := m.GetTournamentByCode(h.DB, ctx, tournamentCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Validate LocationCity
	locationCity, err := m.GetCafeLocationCityByCode(h.DB, ctx, req.LocationCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	gameId, err := m.GetGameIdByCode(h.DB, ctx, req.GameCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if req.Status == "inactive" {
		// Get current total participant
		totalParticipant, err := m.CountParticipantTournamentByTournamentId(h.DB, ctx, tournament.TournamentId)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		if totalParticipant > 0 {
			h.SendBadRequest(w, "Unable to set inactive status because of existing participants")
			return
		}
	}

	// Update tournament
	err = m.UpdateTournamentByCode(
		tx, ctx, gameId, tournamentCode, req.ImageUrl,
		req.Name, req.TournamentRules, req.Level, req.Status,
		req.PrizesImgUrl, req.BookingPrice, startDate, endDate,
		startTime, endTime, req.PlayerSlot,
		req.ParticipantVP, locationCity,
	)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Delete existing badges
	err = m.DeleteTournamentBadge(tx, ctx, tournament.TournamentId)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Insert new badges
	for _, badgeCode := range req.BadgeCodes {
		badgeID, err := m.GetBadgeIdByCode(h.DB, ctx, badgeCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		_, err = m.InsertTournamentBadge(tx, ctx, tournament.TournamentId, badgeID)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
	}

	go removeUpcomingRoomNotification(&m, ctx, tournamentCode, utils.TournamentReminder)
	go func() {
		dataListUser, err := m.GetAllUsers(m.DB, ctx)
		if err != nil {
			h.WriteLog("h.UpdateTournamentAct: ", err)
			return
		}

		onesignal := onesignal.New(h.App)
		for _, user := range dataListUser {

			// Generate Notification code
			notifCode := utils.GeneratePrefixCode(utils.NotifPrefix)

			description := response.NotificationTournamentResp{
				StartDate:   startDate.(time.Time).Format("2006-01-02"),
				StartTime:   startTime.(time.Time).Format("15:04:05"),
				EndTime:     endTime.(time.Time).Format("15:04:05"),
				CafeName:    tournament.CafeName,
				GameName:    tournament.GameName,
				CafeAddress: tournament.CafeAddress,
				Level:       req.Level,
			}

			descriptionJSON, err := json.Marshal(description)
			if err != nil {
				h.WriteLog("h.UpdateTournamentAct: ", err)
				continue
			}

			OSDescription := utils.TournamentReminderPushNotificationDescription + "\n\n" + "Tournament name: " + req.Name + "\n" + "Date: " + startDate.(time.Time).Format("02 January 2006") + "\n" + "Location: " + description.CafeName
			_, err = onesignal.CreateOSNotifications(user.XPlayer, req.Name, OSDescription, utils.Tournament)
			if err != nil {
				h.WriteLog("h.UpdateTournamentAct: ", err)
				continue
			}

			err = m.AddNotification(m.DB, ctx, notifCode, "user", user.UserCode, tournamentCode, utils.TournamentReminder, req.Name, descriptionJSON, req.ImageUrl)
			if err != nil {
				h.WriteLog("h.UpdateTournamentAct: ", err)
			}
		}
	}()

	h.SendSuccess(w, nil, nil)
}

// SetWinnerAct
func (h *Contract) SetWinnerTournamentAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		reqs           = request.SetWinnerTournamentReq{}
		tournamentCode = chi.URLParam(r, "code")
		ctx            = context.TODO()
		m              = model.Contract{App: h.App}
		queueHost      = m.Config.GetString("queue.rabbitmq.host")
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &reqs); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	// Start a transaction
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// For transaction
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		tx.Commit(ctx)
	}()

	tournamentData, err := m.GetTournamentByCode(h.DB, ctx, tournamentCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	tournamentId := tournamentData.TournamentId
	tournamentParticipantPoint := int(tournamentData.ParticipantVP)

	params := tournamentParams{TournamentStatus: tournamentData.Status}
	if isTournamentClosed(w, h, params) {
		return
	}

	// Set tournament status to Closed
	err = m.UpdateTournamentStatusTrx(tx, ctx, tournamentCode, utils.TournamentStatus["CLOSED"])
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	var (
		participantIds []int64
	)
	for _, req := range reqs.TournamentParticipant {
		var statusWinner bool

		userId, err := m.GetUserIdByUserCode(h.DB, ctx, req.UserCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		tournamentEnt, err := m.GetOneTournamentParticipant(h.DB, ctx, tournamentId, int64(userId))
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		if tournamentEnt.Status != "active" {
			h.SendBadRequest(w, fmt.Sprintf("Tournament participant %s is inactive", req.UserCode))
			return
		}

		// Processing all users who have badge and winner vp points
		// Else, for remaining players whom didn't win the game got participant point
		if req.BadgeCode != "" {
			badgeData, err := m.GetBadgeDetailByCode(h.DB, ctx, req.BadgeCode)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}

			// delete if exist and replace with newer
			err = m.DeleteUserBadge(tx, ctx, int64(userId), badgeData.Id)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}

			err = m.AddUserBadgeTx(tx, ctx, int64(userId), badgeData.Id)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}

			if req.Position > 0 {
				statusWinner = true
			}

			// Update participant tournament info
			err = m.UpdateTournamentParticipant(
				tx, ctx, tournamentId, int64(userId), statusWinner, req.Position, tournamentEnt.Status, tournamentEnt.AdditionalInfo.String, int64(tournamentParticipantPoint), tournamentEnt.TransactionCode.String)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}

			// Add user point for all the winners
			err = m.AddUserPoint(tx, ctx, userId, utils.UserPointType["TOURNAMENT_PLAY"], tournamentCode, tournamentParticipantPoint)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}

			go func(userID int64, position int) { // Publisher badge
				queueData := rabbit.QueueDataPayload(
					rabbit.QueueUserBadge,
					rabbit.QueueUserBadgeReq(
						utils.SpesificBoardGameCategory,
						userID,
					),
				)
				err = rabbit.PublishQueue(ctx, queueHost, queueData)
				if err != nil {
					log.Printf("Error : %s", err)
				}

				if position == 1 {
					queueData := rabbit.QueueDataPayload(
						rabbit.QueueUserBadge,
						rabbit.QueueUserBadgeReq(
							utils.TournamentWon,
							userID,
						),
					)
					err = rabbit.PublishQueue(ctx, queueHost, queueData)
					if err != nil {
						log.Printf("Error : %s", err)
					}
				}
			}(userId, req.Position)
		} else {
			userId, err := m.GetUserIdByUserCode(h.DB, ctx, req.UserCode)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}

			tournamentEnt, err := m.GetOneTournamentParticipant(h.DB, ctx, tournamentId, int64(userId))
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}

			if tournamentEnt.Status != "active" {
				h.SendBadRequest(w, fmt.Sprintf("Tournament participant %s is inactive", req.UserCode))
				return
			}

			// Update participant tournament info
			err = m.UpdateTournamentParticipant(
				tx, ctx, tournamentId, int64(userId), statusWinner, req.Position, tournamentEnt.Status, tournamentEnt.AdditionalInfo.String, int64(tournamentParticipantPoint), tournamentEnt.TransactionCode.String)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}

			err = m.AddUserPoint(tx, ctx, userId, utils.UserPointType["TOURNAMENT_PLAY"], tournamentCode, tournamentParticipantPoint)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}
		}

		participantIds = append(participantIds, int64(userId))
	}

	for _, id := range participantIds {
		if err := m.AddUserGameCollections(h.DB, ctx, id, tournamentData.GameId); err == nil {
			m.AddUserPoint(tx, ctx, id, utils.UserPointType["GAME_COLLECTION"], tournamentData.GameCode, utils.FirstTimePlayed.Int())

			queueHost := m.Config.GetString("queue.rabbitmq.host")
			queueData := rabbit.QueueDataPayload(
				rabbit.QueueUserBadge,
				rabbit.QueueUserBadgeReq(
					utils.PlayingGames,
					id,
				),
			)
			err = rabbit.PublishQueue(ctx, queueHost, queueData)
			if err != nil {
				log.Printf("Error : %s", err)
			}

			// Publisher badge
			queueData = rabbit.QueueDataPayload(
				rabbit.QueueUserBadge,
				rabbit.QueueUserBadgeReq(
					utils.TimeLimit,
					id,
				),
			)
			err = rabbit.PublishQueue(ctx, queueHost, queueData)
			if err != nil {
				log.Printf("Error : %s", err)
			}

			// Publisher badge
			queueData = rabbit.QueueDataPayload(
				rabbit.QueueUserBadge,
				rabbit.QueueUserBadgeReq(
					utils.PlayingGames,
					id,
				),
			)
			err = rabbit.PublishQueue(ctx, queueHost, queueData)
			if err != nil {
				log.Printf("Error : %s", err)
			}
		}
	}

	go removeUpcomingRoomNotification(&m, ctx, tournamentCode, utils.TournamentReminder)

	h.SendSuccess(w, nil, nil)
}

func (h *Contract) DeleteTournamentAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		tournamentCode = chi.URLParam(r, "code")
		ctx            = context.TODO()
		m              = model.Contract{App: h.App}
	)

	tournament, err := m.GetTournamentByCode(h.DB, ctx, tournamentCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	params := tournamentParams{TournamentStatus: tournament.Status}
	if isTournamentClosed(w, h, params) {
		return
	}

	if isThereAnyTournamentParticipants(w, h, int(tournament.CurrentUsedSlot)) {
		return
	}

	err = m.DeleteTournamentByCode(h.DB, ctx, tournamentCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	go removeUpcomingRoomNotification(&m, ctx, tournamentCode, utils.TournamentReminder)

	h.SendSuccess(w, nil, nil)
}

// BookingTournament
func (h *Contract) BookingTournamentAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		m              = model.Contract{App: h.App}
		tournamentCode = chi.URLParam(r, "code")
		userCode       = bootstrap.GetIdentifierCodeFromToken(ctx, r)
	)

	user, err := m.GetUserByUserCode(h.DB, ctx, userCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Get room by code
	trnm, err := m.GetTournamentByCode(h.DB, ctx, tournamentCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	params := tournamentParams{TournamentStatus: trnm.Status, IsBooking: true}
	if isTournamentClosed(w, h, params) {
		return
	}

	if isTournamentInactive(w, h, params) {
		return
	}

	// Get current total participant
	totalParticipant, err := m.CountParticipantTournamentByTournamentId(h.DB, ctx, trnm.TournamentId)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if totalParticipant >= int(trnm.PlayerSlot) {
		h.SendBadRequest(w, "Sorry, this tournament is fully booked")
		return
	}

	// Check if already booked
	participant, err := m.GetOneTournamentParticipant(h.DB, ctx, trnm.TournamentId, int64(user.ID))
	if err != nil && err.Error() != utils.EmptyData {
		h.SendBadRequest(w, err.Error())
		return
	}

	if participant.Id > 0 && participant.Status != "cancel" {
		//check if status pending and already exist transaction
		if participant.Status == "pending" {
			order, err := m.GetTransactionByCode(h.DB, ctx, userCode, participant.TransactionCode.String)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}

			h.SendSuccess(w, response.BookingRes{
				InvoiceUrl: order.PaymentLink,
				ExpiredAt:  order.ExpiredDate.Time.Local().Format(utils.DATE_TIME_FORMAT),
			}, nil)
			return
		}

		h.SendBadRequest(w, "You have already booked this tournament")
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

	var (
		statusParticipant, transactionCode, invoiceUrl string
		expiredAt                                      time.Time
		earnedPoint                                    int64
	)

	statusParticipant = "pending"
	if trnm.BookingPrice > 0 {
		earnedPoint = int64(utils.CalculateUserRedeemPoint(trnm.BookingPrice))
		//call xendit
		_, transactionCode, invoiceUrl, expiredAt, err = m.CreateOneTimeInvoice(tx, ctx, int64(user.ID), utils.UserPointType["TOURNAMENT_TYPE"], tournamentCode, trnm.BookingPrice, fmt.Sprintf("INVOICE-%s-%s", userCode, tournamentCode), user.Email.String)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
	} else {
		earnedPoint = int64(utils.CalculateUserRedeemPoint(trnm.BookingPrice))
		//call xendit
		_, transactionCode, invoiceUrl, expiredAt, err = m.CreateFreeInvoice(tx, ctx, int64(user.ID), utils.UserPointType["TOURNAMENT_TYPE"], tournamentCode, trnm.BookingPrice, fmt.Sprintf("INVOICE-%s-%s", userCode, tournamentCode), user.Email.String)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		statusParticipant = "active"

		// Add user point from price (to log transaction event though it was free)
		err = m.AddUserPoint(tx, ctx, int64(user.ID), utils.UserPointType["TOURNAMENT_PAID"], transactionCode, 0)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		// Add user point from vp point participation (used for activity)
		err = m.AddUserPoint(tx, ctx, int64(user.ID), utils.UserPointType["TOURNAMENT_TYPE"], trnm.TournamentCode, 0)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		go func(ctx context.Context, userID int64) {
			// Publisher badge
			queueData := rabbit.QueueDataPayload(
				rabbit.QueueUserBadge,
				rabbit.QueueUserBadgeReq(
					utils.TimeLimit,
					int64(user.ID),
				),
			)
			queueHost := m.Config.GetString("queue.rabbitmq.host")
			err = rabbit.PublishQueue(ctx, queueHost, queueData)
			if err != nil {
				log.Printf("Error : %s", err)
			}
		}(ctx, int64(user.ID))
	}

	// check if exist
	if participant.Id > 0 && participant.Status != "active" {
		//update status
		err = m.UpdateTournamentParticipant(tx, ctx, trnm.TournamentId, int64(user.ID), false, participant.Position, statusParticipant, participant.AdditionalInfo.String, earnedPoint, transactionCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
	} else {
		//add participant
		err = m.InsertOneTournamentParticipant(tx, ctx, trnm.TournamentId, int64(user.ID), false, participant.Position, statusParticipant, participant.AdditionalInfo.String, earnedPoint, transactionCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
	}

	h.SendSuccess(w, response.BookingRes{
		InvoiceUrl: invoiceUrl,
		ExpiredAt:  expiredAt.Local().Format(utils.DATE_TIME_FORMAT),
	}, nil)
}

func (h *Contract) UpdateTournamentStatus(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		req            = request.UpdateStatusTournamentReq{}
		m              = model.Contract{App: h.App}
		tournamentCode = chi.URLParam(r, "code")
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err.Error())
		return
	}

	tournament, err := m.GetTournamentByCode(h.DB, ctx, tournamentCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	params := tournamentParams{TournamentStatus: tournament.Status}
	if isTournamentClosed(w, h, params) {
		return
	}

	if isTournamentStatusAttrInvalid(w, h, req.Status) {
		return
	}

	err = m.UpdateTournamentStatus(h.DB, ctx, tournamentCode, req.Status)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if req.Status == utils.StatusTournament[0] {
		go func() {
			dataListUser, err := m.GetAllUsers(m.DB, ctx)
			if err != nil {
				h.WriteLog("h.UpdateTournamentStatus: ", err)
				return
			}

			onesignal := onesignal.New(h.App)
			for _, user := range dataListUser {

				// Generate Notification code
				notifCode := utils.GeneratePrefixCode(utils.NotifPrefix)

				description := response.NotificationTournamentResp{
					StartDate:   tournament.StartDate.Time.Format("2006-01-02"),
					StartTime:   tournament.StartTime.Format("15:04:05"),
					EndTime:     tournament.EndTime.Format("15:04:05"),
					CafeName:    tournament.CafeName,
					GameName:    tournament.GameName,
					CafeAddress: tournament.CafeAddress,
					Level:       tournament.Level,
				}

				descriptionJSON, err := json.Marshal(description)
				if err != nil {
					h.WriteLog("h.UpdateTournamentStatus: ", err)
					continue
				}

				OSDescription := utils.TournamentReminderPushNotificationDescription + "\n\n" + "Tournament name: " + tournament.Name.String + "\n" + "Date: " + tournament.StartDate.Time.Format("02 January 2006") + "\n" + "Location: " + description.CafeName
				_, err = onesignal.CreateOSNotifications(user.XPlayer, tournament.Name.String, OSDescription, utils.Tournament)
				if err != nil {
					h.WriteLog("h.UpdateTournamentStatus: ", err)
					continue
				}

				err = m.AddNotification(m.DB, ctx, notifCode, "user", user.UserCode, tournamentCode, utils.TournamentReminder, tournament.Name.String, descriptionJSON, tournament.ImageUrl.String)
				if err != nil {
					h.WriteLog("h.UpdateTournamentStatus: ", err)
				}
			}
		}()
	} else {
		go removeUpcomingRoomNotification(&m, ctx, tournamentCode, utils.TournamentReminder)
	}

	h.SendSuccess(w, nil, nil)
}

// Set Participant status as "Inactive" to prevent them to gain
// points and badges from the room in case they are not present
// when the sessions end or asks for refund.
func (h *Contract) RemoveTournamentParticipantAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		m              = model.Contract{App: h.App}
		tournamentCode = chi.URLParam(r, "code")
		req            = request.TournamentParticipantRemoveReq{}
	)

	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err.Error())
		return
	}

	tournament, err := m.GetTournamentByCode(h.DB, ctx, tournamentCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	params := tournamentParams{TournamentStatus: tournament.Status}
	if isTournamentClosed(w, h, params) {
		return
	}

	user, err := m.GetUserByUserCode(h.DB, ctx, req.UserCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	participant, err := m.GetOneTournamentParticipant(h.DB, ctx, tournament.TournamentId, int64(user.ID))
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if participant.Status == "inactive" {
		h.SendBadRequest(w, utils.ErrDeactiveParticipantIsInactive)
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

	trx, err := m.GetTransactionBySourceCode(h.DB, ctx, user.UserCode, utils.UserPointType["TOURNAMENT_TYPE"], tournamentCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.DeleteTournamentParticipant(tx, ctx, tournament.TournamentId, int64(user.ID))
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.RemoveUserPoint(tx, ctx, int64(user.ID), utils.UserPointType["TOURNAMENT_TYPE"], tournamentCode, 0)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if participant.RewardPoint.Int64 > 0 {
		err = m.RemoveUserPoint(tx, ctx, int64(user.ID), utils.UserPointType["TOURNAMENT_PLAY"], tournamentCode, int(participant.RewardPoint.Int64))
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
	}

	if tournament.BookingPrice <= 0 {
		err = m.RemoveUserPoint(tx, ctx, int64(user.ID), utils.UserPointType["TOURNAMENT_PAID"], trx.TransactionCode, 0)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		err = m.UpdateInvoiceTrx(tx, ctx, trx.AggregatorCode, "", utils.PaymentStatus["EXPIRED"], "")
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
	}

	h.SendSuccess(w, nil, nil)
}

type tournamentParams struct {
	TournamentStatus string
	IsBooking        bool
}

func isTournamentClosed(w http.ResponseWriter, h *Contract, params tournamentParams) bool {
	if params.TournamentStatus != "closed" {
		return false
	}

	message := "Modifications are not allowed on closed tournaments (deletion or update)"
	if params.IsBooking && params.TournamentStatus == "closed" {
		message = "Booking is not allowed for closed tournaments"
	}

	h.SendBadRequest(w, message)
	return true
}

func isTournamentInactive(w http.ResponseWriter, h *Contract, params tournamentParams) bool {
	if params.TournamentStatus != "inactive" {
		return false
	}

	h.SendBadRequest(w, "Booking is not allowed for inactive tournaments")
	return true
}

func isTournamentStatusAttrInvalid(w http.ResponseWriter, h *Contract, tournamentStatus string) bool {
	if utils.Contains(utils.StatusTournament, tournamentStatus) {
		return false
	}

	h.SendBadRequest(w, "Invalid tournament status. Please use 'active' or 'inactive'")
	return true
}

func isThereAnyTournamentParticipants(w http.ResponseWriter, h *Contract, currentParticipant int) bool {
	if currentParticipant == 0 {
		return false
	}

	message := fmt.Sprintf("Cannot delete tournament because there are %d participants in this tournament", currentParticipant)
	h.SendBadRequest(w, message)
	return true
}
