package handler

import (
	"context"
	"dots-api/bootstrap"
	"dots-api/lib/rabbit"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// GetRoomList
func (h *Contract) GetRoomList(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.RoomListRes, 0)
		param = request.RoomParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseRoom(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetRoomList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.RoomListRes{
			CafeId:             v.CafeId,
			CafeCode:           v.CafeCode,
			CafeName:           v.CafeName,
			CafeAddress:        v.CafeAddress,
			RoomCode:           v.RoomCode,
			RoomType:           v.RoomType,
			RoomImgUrl:         v.RoomImgUrl,
			Name:               v.Name,
			Description:        v.Description,
			Instruction:        v.Instruction,
			Difficulty:         v.Difficulty,
			StartDate:          v.StartDate.Time.Format(utils.DATE_FORMAT),
			EndDate:            v.EndDate.Time.Format(utils.DATE_FORMAT),
			StartTime:          v.StartTime.Format(utils.TIME_FORMAT),
			EndTime:            v.EndTime.Format(utils.TIME_FORMAT),
			MaximumParticipant: v.MaximumParticipant,
			CurrentUsedSlot:    v.CurrentUsedSlot,
			InstagramLink:      v.InstagramLink,
			Status:             v.Status,
			DayPastEndDate:     v.DayPastEndDate.Float64,
			BookingPrice:       v.BookingPrice,
			GameMasterName:     v.GameMasterName.String,
			GameMasterImageUrl: v.GameMasterImageUrl.String,
			GameCode:           v.GameCode,
			GameName:           v.GameName,
			GameImgUrl:         v.GameImgUrl,
		})
	}

	h.SendSuccess(w, res, param)
}

// GetRoomByCode
func (h *Contract) GetRoomByCode(w http.ResponseWriter, r *http.Request) {
	var (
		err                error
		code               = chi.URLParam(r, "code")
		ctx                = context.TODO()
		m                  = model.Contract{App: h.App}
		resRoomParticipant = make([]response.RoomParticipantRes, 0)
		userCode           = bootstrap.GetIdentifierCodeFromToken(ctx, r)
		haveJoinedRoom     = false
	)

	// Get Room Info Detail
	roomInfo, err := m.GetRoomByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Check if already booked & status active
	participant, err := m.GetParticipantByRoomCodeAndUserCode(h.DB, ctx, code, userCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if len(participant.UserCode) > 0 && participant.Status == "active" {
		haveJoinedRoom = true
	}

	// After that, get the list of all participants on selected room code
	participantInfo, err := m.GetAllParticipantByRoomCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate participant
	for _, participant := range participantInfo {
		resParticipant := response.RoomParticipantRes{
			UserCode:       participant.UserCode,
			UserName:       participant.UserName,
			UserImgUrl:     participant.UserImgUrl,
			StatusWinner:   participant.StatusWinner,
			Status:         participant.Status,
			AdditionalInfo: participant.AdditionalInfo.String,
			Position:       participant.Position,
			RewardPoint:    int(participant.RewardPoint.Int64),
			LatestTier:     participant.LatestTier.String,
		}
		resRoomParticipant = append(resRoomParticipant, resParticipant)
	}

	// Populate response
	h.SendSuccess(w, response.RoomRes{
		GameMasterCode:     roomInfo.GameMasterCode.String,
		GameMasterName:     roomInfo.GameMasterName.String,
		GameMasterImgUrl:   roomInfo.GameMasterImgUrl.String,
		GameCode:           roomInfo.GameCode,
		GameName:           roomInfo.GameName,
		GameImgUrl:         roomInfo.GameImgUrl,
		CafeCode:           roomInfo.CafeCode,
		CafeName:           roomInfo.CafeName,
		CafeAddress:        roomInfo.CafeAddress,
		RoomCode:           roomInfo.RoomCode,
		RoomType:           roomInfo.RoomType,
		Name:               roomInfo.Name,
		Description:        roomInfo.Description,
		Difficulty:         roomInfo.Difficulty,
		StartDate:          roomInfo.StartDate.Time.Format(utils.DATE_FORMAT),
		EndDate:            roomInfo.EndDate.Time.Format(utils.DATE_FORMAT),
		StartTime:          roomInfo.StartTime.Format(utils.TIME_FORMAT),
		EndTime:            roomInfo.EndTime.Format(utils.TIME_FORMAT),
		MaximumParticipant: roomInfo.MaximumParticipant,
		BookingPrice:       roomInfo.BookingPrice,
		RewardPoint:        roomInfo.RewardPoint,
		InstagramLink:      roomInfo.InstagramLink,
		Status:             roomInfo.Status,
		DayPastEndDate:     roomInfo.DayPastEndDate.Float64,
		RoomBannerUrl:      roomInfo.BannerRoomUrl,
		CurrentUsedSlot:    roomInfo.CurrentUsedSlot,
		RoomParticipant:    resRoomParticipant,
		HaveJoined:         haveJoinedRoom,
	}, nil)
}

// AddRoom
func (h *Contract) AddRoom(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		req       = request.RoomReq{}
		ctx       = context.TODO()
		m         = model.Contract{App: h.App}
		startDate interface{}
		endDate   interface{}
		startTime interface{}
		endTime   interface{}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusRoom, req.Status) {
		h.SendBadRequest(w, "wrong status value for rooms(active|inactive)")
		return
	}

	if !utils.Contains(utils.RoomType, req.RoomType) {
		h.SendBadRequest(w, "wrong type value for rooms(normal|special_event)")
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

	// Validate LocationCity
	locationCity, err := m.GetCafeLocationCityByCode(h.DB, ctx, req.LocationCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Get Admin Id
	gameMasterId, err := m.GetAdminIdByCode(h.DB, ctx, req.GameMasterCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Get Game Id
	gameId, err := m.GetGameIdByCode(h.DB, ctx, req.GameCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	// Generate Random Code
	code := utils.GeneratePrefixCode(utils.RoomPrefix)
	err = m.AddRoom(
		h.DB,
		ctx,
		gameMasterId,
		gameId,
		code,
		req.RoomType,
		req.Name,
		req.Description,
		startDate,
		endDate,
		startTime,
		endTime,
		float64(req.BookingPrice),
		req.RewardPoint,
		req.InstagramLink,
		req.Status,
		req.Difficulty,
		req.Instruction,
		req.MaximumParticipant,
		req.ImageURL,
		locationCity)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

// UpdateRoom
func (h *Contract) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		req       = request.RoomReq{}
		code      = chi.URLParam(r, "code")
		ctx       = context.TODO()
		m         = model.Contract{App: h.App}
		startDate interface{}
		endDate   interface{}
		startTime interface{}
		endTime   interface{}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	// Validate LocationCity
	locationCity, err := m.GetCafeLocationCityByCode(h.DB, ctx, req.LocationCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Get Admin Id
	gameMasterId, err := m.GetAdminIdByCode(h.DB, ctx, req.GameMasterCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Get Game Id
	gameId, err := m.GetGameIdByCode(h.DB, ctx, req.GameCode)
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

	err = m.UpdateRoom(h.DB, ctx, code, gameMasterId, gameId, code, req.RoomType, req.Name, req.Description, startDate, endDate,
		startTime, endTime, float64(req.BookingPrice), req.RewardPoint, req.InstagramLink, req.Status, req.Difficulty, req.Instruction, req.MaximumParticipant, req.ImageURL, locationCity)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

func (h *Contract) BookingRoom(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		ctx      = context.TODO()
		m        = model.Contract{App: h.App}
		roomCode = chi.URLParam(r, "code")
		userCode = bootstrap.GetIdentifierCodeFromToken(ctx, r)
	)

	room, err := m.GetRoomByCode(h.DB, ctx, roomCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	params := roomParams{RoomStatus: room.Status, IsBooking: true}
	if isRoomClosed(w, h, params) {
		return
	}

	if isRoomInactive(w, h, params) {
		return
	}

	user, err := m.GetUserByUserCode(h.DB, ctx, userCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Get current total participant
	totalParticipant, err := m.CountParticipantRoomByRoomId(h.DB, ctx, room.RoomId)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if totalParticipant >= room.MaximumParticipant {
		h.SendBadRequest(w, "Sorry, this room is fully booked")
		return
	}

	// Check if already booked
	participant, err := m.GetParticipantByRoomCodeAndUserCode(h.DB, ctx, roomCode, userCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if len(participant.UserCode) > 0 && participant.Status != "cancel" {
		// Check if status pending and already exist transaction
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

		h.SendBadRequest(w, "You have already booked this room")
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

	//call xendit
	_, transactionCode, invoiceUrl, expiredAt, err := m.CreateOneTimeInvoice(tx, ctx, int64(user.ID), utils.UserPointType["ROOM_TYPE"], roomCode, room.BookingPrice, fmt.Sprintf("INVOICE-%s-%s", userCode, roomCode), user.Email.String)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	statusParticipant := "pending"
	earnedPoint := int64(utils.CalculateUserRedeemPoint(room.BookingPrice))
	// check if exist
	if len(participant.UserCode) > 0 && participant.Status != "active" {
		//update status
		err = m.UpdateRoomParticipant(tx, ctx, room.RoomId, int64(user.ID), false, participant.Position, statusParticipant, participant.AdditionalInfo.String, earnedPoint, transactionCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
	} else {
		//add participant
		err = m.InsertOneRoomParticipant(tx, ctx, room.RoomId, int64(user.ID), statusParticipant, earnedPoint, transactionCode)
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

// Set Winner
func (h *Contract) SetWinnerRoomAct(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		reqs     = request.SetWinnerRoomReq{}
		roomCode = chi.URLParam(r, "code")
		ctx      = context.TODO()
		m        = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &reqs); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

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

	room, err := m.GetRoomByCode(h.DB, ctx, roomCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	participants, err := m.GetAllParticipantByRoomCode(h.DB, ctx, roomCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	params := roomParams{RoomStatus: room.Status}
	if isRoomClosed(w, h, params) {
		return
	}

	// Set room status to Closed
	err = m.UpdateRoomStatusTrx(tx, ctx, roomCode, utils.RoomStatus["CLOSED"])
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	for _, req := range reqs.RoomParticipant {
		var statusWinner bool
		roomId, err := m.GetRoomIdByCode(h.DB, ctx, roomCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		userId, err := m.GetUserIdByUserCode(h.DB, ctx, req.UserCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		roomParticipantEnt, err := m.GetOneRoomParticipant(h.DB, ctx, roomId, int64(userId))
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		if req.Position > 0 {
			statusWinner = true
		}

		err = m.UpdateRoomParticipant(tx, ctx, roomId, int64(userId), statusWinner, roomParticipantEnt.Position, roomParticipantEnt.Status, "member", int64(roomParticipantEnt.RewardPoint.Int64), roomParticipantEnt.TransactionCode.String)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		// Publisher badge
		queueData := rabbit.QueueDataPayload(
			rabbit.QueueUserBadge,
			rabbit.QueueUserBadgeReq(
				utils.SpesificBoardGameCategory,
				int64(userId),
			),
		)
		queueHost := m.Config.GetString("queue.rabbitmq.host")
		err = rabbit.PublishQueue(ctx, queueHost, queueData)
		if err != nil {
			log.Printf("Error : %s", err)
		}
	}

	for _, participant := range participants {
		_ = m.AddUserGameCollections(h.DB, ctx, participant.UserId, room.GameId)
	}

	h.SendSuccess(w, nil, nil)
}

func (h *Contract) UpdateRoomStatus(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		ctx      = context.TODO()
		req      = request.UpdateStatusRoomReq{}
		m        = model.Contract{App: h.App}
		roomCode = chi.URLParam(r, "code")
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err.Error())
		return
	}

	room, err := m.GetRoomByCode(h.DB, ctx, roomCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	params := roomParams{RoomStatus: room.Status}
	if isRoomClosed(w, h, params) {
		return
	}

	if isRoomStatusAttrInvalid(w, h, req.Status) {
		return
	}

	err = m.UpdateRoomStatus(h.DB, ctx, roomCode, req.Status)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

func (h *Contract) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		ctx      = context.TODO()
		m        = model.Contract{App: h.App}
		roomCode = chi.URLParam(r, "code")
	)

	room, err := m.GetRoomByCode(h.DB, ctx, roomCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	params := roomParams{RoomStatus: room.Status}
	if isRoomClosed(w, h, params) {
		return
	}

	if isThereAnyRoomParticipants(w, h, room.CurrentUsedSlot) {
		return
	}

	err = m.DeleteRoomByCode(h.DB, ctx, roomCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

type roomParams struct {
	RoomStatus string
	IsBooking  bool
}

func isRoomClosed(w http.ResponseWriter, h *Contract, params roomParams) bool {
	if params.RoomStatus != "closed" {
		return false
	}

	message := "Modifications are not allowed on closed rooms (deletion or update)"
	if params.IsBooking && params.RoomStatus == "closed" {
		message = "Booking is not allowed for closed rooms"
	}

	h.SendBadRequest(w, message)
	return true
}

func isRoomInactive(w http.ResponseWriter, h *Contract, params roomParams) bool {
	if params.RoomStatus != "inactive" {
		return false
	}

	h.SendBadRequest(w, "Booking is not allowed for inactive rooms")
	return true
}

func isRoomStatusAttrInvalid(w http.ResponseWriter, h *Contract, roomStatus string) bool {
	if utils.Contains(utils.StatusRoom, roomStatus) {
		return false
	}

	h.SendBadRequest(w, "Invalid room status. Please use 'active' or 'inactive'")
	return true
}

func isThereAnyRoomParticipants(w http.ResponseWriter, h *Contract, currentParticipant int) bool {
	if currentParticipant == 0 {
		return false
	}

	message := fmt.Sprintf("Cannot delete room because there are %d participants in this room", currentParticipant)
	h.SendBadRequest(w, message)
	return true
}
