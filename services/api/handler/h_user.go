package handler

import (
	"context"
	"dots-api/bootstrap"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// GetUserListAct ...
func (h *Contract) GetUserListAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.UserRes, 0)
		param = request.UserParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseUser(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetUserList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		stats := response.UserStatsRes{
			VP:              v.LatestPoint,
			RoomNormalCount: int(v.RoomNormalCount.Int64),
			RoomEventCount:  int(v.RoomEventCount.Int64),
			TournamentCount: int(v.TournamentCount.Int64),
			BadgeCount:      int(v.BadgeCount.Int64),
			GameCount:       int(v.GameCount.Int64),
		}
		res = append(res, response.UserRes{
			UserCode:           v.UserCode,
			Email:              v.Email.String,
			UserName:           v.UserName.String,
			PhoneNumber:        v.PhoneNumber,
			FullName:           v.FullName,
			DateOfBirth:        v.DateOfBirth.String,
			Gender:             v.Gender.String,
			ImageURL:           v.ImageURL.String,
			UserStats:          stats,
			LatestPoint:        v.LatestPoint,
			LatestTier:         v.LatestTierName,
			Password:           v.Password,
			XPlayer:            v.XPlayer,
			StatusVerification: v.StatusVerification,
			Status:             v.Status,
			UserStyle:          response.UserStyleRes{Color: v.UserStyle.Color},
			TotalSpent:         v.TotalSpent,
			CreatedDate:        v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
			UpdatedDate:        v.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
		})
	}

	h.SendSuccess(w, res, param)
}

// GetUserDetailAct ...
func (h *Contract) GetUserDetailAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	data, err := m.GetUserByUserCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, response.UserRes{
		UserCode:           data.UserCode,
		Email:              data.Email.String,
		UserName:           data.UserName.String,
		PhoneNumber:        data.PhoneNumber,
		FullName:           data.FullName,
		DateOfBirth:        data.DateOfBirth.String,
		Gender:             data.Gender.String,
		ImageURL:           data.ImageURL.String,
		LatestPoint:        data.LatestPoint,
		LatestTier:         data.LatestTierName,
		Password:           data.Password,
		XPlayer:            data.XPlayer,
		StatusVerification: data.StatusVerification,
		Status:             data.Status,
		UserStyle:          response.UserStyleRes{Color: data.UserStyle.Color},
		TotalSpent:         data.TotalSpent,
		CreatedDate:        data.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		UpdatedDate:        data.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
	}, nil)
}

func (h *Contract) GetUserProfileAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		userIdentifier = bootstrap.GetIdentifierCodeFromToken(ctx, r)
		m              = model.Contract{App: h.App}
		dataUser       model.UserEnt
		res            = response.UserProfileRes{}
		TierBenefits   = make([]response.TierWithBenefitRes, 0)
	)

	dataUser, err = m.GetUserByUserCode(h.DB, ctx, userIdentifier)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// After that, get the list of all rewards based on user latest tier
	rewardList, err := m.GetTierWithReward(h.DB, ctx, dataUser.LatestTierId)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate TierBenefits
	for _, v := range rewardList {
		TierBenefits = append(TierBenefits, response.TierWithBenefitRes{
			RewardCode:        v.RewardCode,
			RewardName:        v.RewardName.String,
			RewardImageUrl:    v.RewardImageUrl.String,
			RewardDescription: v.RewardDescription.String,
		})
	}

	// Get tier range min-max point
	TierMinRangePoint := response.TierRangePointRes{
		MinPoint: dataUser.TierMinRangePoint,
		MaxPoint: dataUser.TierMaxRangePoint,
	}

	// Populate response
	res = response.UserProfileRes{
		UserCode:           dataUser.UserCode,
		Email:              dataUser.Email.String,
		UserName:           dataUser.UserName.String,
		PhoneNumber:        dataUser.PhoneNumber,
		FullName:           dataUser.FullName,
		DateOfBirth:        dataUser.DateOfBirth.String,
		Gender:             dataUser.Gender.String,
		ImageURL:           dataUser.ImageURL.String,
		LatestPoint:        dataUser.LatestPoint,
		LatestTier:         dataUser.LatestTierName,
		Password:           dataUser.Password,
		XPlayer:            dataUser.XPlayer,
		TierRangePoint:     &TierMinRangePoint,
		TierBenefits:       TierBenefits,
		StatusVerification: dataUser.StatusVerification,
		Status:             dataUser.Status,
		UserStyle:          response.UserStyleRes{Color: dataUser.UserStyle.Color},
		MemberSince:        dataUser.CreatedDate.Format(utils.YEAR_FORMAT),
		CreatedDate:        dataUser.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		UpdatedDate:        dataUser.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
	}

	h.SendSuccess(w, res, nil)
}

// GetUserProfileByCodeAct
func (h *Contract) GetUserProfileByCodeAct(w http.ResponseWriter, r *http.Request) {
	var (
		err          error
		ctx          = context.TODO()
		code         = chi.URLParam(r, "code")
		m            = model.Contract{App: h.App}
		dataUser     model.UserEnt
		res          = response.UserProfileRes{}
		TierBenefits = make([]response.TierWithBenefitRes, 0)
	)

	dataUser, err = m.GetUserByUserCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// After that, get the list of all rewards based on user latest tier
	rewardList, err := m.GetTierWithReward(h.DB, ctx, dataUser.LatestTierId)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate TierBenefits
	for _, v := range rewardList {
		TierBenefits = append(TierBenefits, response.TierWithBenefitRes{
			RewardCode:        v.RewardCode,
			RewardName:        v.RewardName.String,
			RewardImageUrl:    v.RewardImageUrl.String,
			RewardDescription: v.RewardDescription.String,
		})
	}

	// Get tier range min-max point
	TierMinRangePoint := response.TierRangePointRes{
		MinPoint: dataUser.TierMinRangePoint,
		MaxPoint: dataUser.TierMaxRangePoint,
	}

	// Populate response
	res = response.UserProfileRes{
		UserCode:           dataUser.UserCode,
		Email:              dataUser.Email.String,
		UserName:           dataUser.UserName.String,
		PhoneNumber:        dataUser.PhoneNumber,
		FullName:           dataUser.FullName,
		DateOfBirth:        dataUser.DateOfBirth.String,
		Gender:             dataUser.Gender.String,
		ImageURL:           dataUser.ImageURL.String,
		LatestPoint:        dataUser.LatestPoint,
		LatestTier:         dataUser.LatestTierName,
		XPlayer:            dataUser.XPlayer,
		TierRangePoint:     &TierMinRangePoint,
		TierBenefits:       TierBenefits,
		StatusVerification: dataUser.StatusVerification,
		Status:             dataUser.Status,
		UserStyle:          response.UserStyleRes{Color: dataUser.UserStyle.Color},
		MemberSince:        dataUser.CreatedDate.Format(utils.YEAR_FORMAT),
		CreatedDate:        dataUser.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		UpdatedDate:        dataUser.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
	}

	h.SendSuccess(w, res, nil)
}

func (h *Contract) UpdateUserProfileAct(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		ctx         = context.TODO()
		req         = request.UpdateProfileUserReq{}
		m           = model.Contract{App: h.App}
		userCode    = bootstrap.GetIdentifierCodeFromToken(ctx, r)
		fullName    string
		userName    string
		dateOfBirth string
		gender      string
		imageUri    string
		phoneNUmber string
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	// Check if username exceeds max character limit (12)
	if len(req.UserName) > 12 {
		h.SendBadRequest(w, utils.ErrUsernameExceedsMaxCharacterLimit)
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

	currentData, err := m.GetUserByUserCode(h.DB, ctx, userCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Re-assign profile data manually
	fullName = req.FullName
	if reflect.TypeOf(fullName) == nil || fullName == "" {
		fullName = currentData.FullName
	}

	userName = req.UserName
	if reflect.TypeOf(userName) == nil || userName == "" {
		userName = currentData.UserName.String
	}

	if userName != currentData.UserName.String {
		err = m.CheckIfUsernameExists(h.DB, ctx, req.UserName)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
	}

	imageUri = req.ImageUrl
	if reflect.TypeOf(imageUri) == nil || imageUri == "" {
		imageUri = currentData.ImageURL.String
	}

	phoneNUmber = req.PhoneNumber
	if reflect.TypeOf(phoneNUmber) == nil || phoneNUmber == "" {
		phoneNUmber = currentData.PhoneNumber
	}

	if req.DateOfBirth == "" && !currentData.DateOfBirth.Valid {
		h.SendBadRequest(w, "you have not set your date of birth")
		return
	}

	if req.DateOfBirth == "" {
		dateOfBirth = currentData.DateOfBirth.String
	} else {
		dateOfBirth = req.DateOfBirth
	}

	if req.Gender == "" && !currentData.Gender.Valid {
		h.SendBadRequest(w, "you have not set your gender")
		return
	}

	if req.Gender == "male" || req.Gender == "female" {
		gender = req.Gender
	} else {
		gender = currentData.Gender.String
	}

	err = m.UpdateUserProfile(h.DB, ctx, userCode, fullName, userName, gender, dateOfBirth, imageUri, phoneNUmber)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if currentData.ImageURL.String == "" && imageUri != "" {
		m.AddUserPoint(tx, ctx, int64(currentData.ID), utils.UserPointType["PROFILE"], userCode, utils.UpdateProfile.Int())
	}

	h.SendSuccess(w, nil, nil)
}

func (h *Contract) UpdateUserStatusAct(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		ctx      = context.TODO()
		req      = request.UpdateStatusUserReq{}
		m        = model.Contract{App: h.App}
		userCode = chi.URLParam(r, "code")
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusUser, req.Status) {
		h.SendBadRequest(w, "wrong status value for user(active|inactive)")
		return
	}

	err = m.UpdateUserStatus(h.DB, ctx, userCode, req.Status)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

func (h *Contract) UpdateUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		ctx      = context.TODO()
		req      = request.UpdateUserReq{}
		m        = model.Contract{App: h.App}
		userCode = chi.URLParam(r, "code")
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	// Check if username exceeds max character limit (12)
	if len(req.UserName) > 12 {
		h.SendBadRequest(w, utils.ErrUsernameExceedsMaxCharacterLimit)
		return
	}

	currentData, err := m.GetUserByUserCode(h.DB, ctx, userCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if req.UserName != currentData.UserName.String {
		err = m.CheckIfUsernameExists(h.DB, ctx, req.UserName)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}
	}

	if !utils.Contains(utils.StatusUser, req.Status) {
		h.SendBadRequest(w, "wrong status value for user(active|inactive)")
		return
	}

	err = m.UpdateUser(h.DB, ctx, userCode, req.FullName, req.DateOfBirth.Format(utils.DATE_FORMAT), req.Gender, req.ImageUrl, req.PhoneNumber, req.Email, req.UserName, req.Status)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

func (h *Contract) UpdatePasswordUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		userIdentifier = bootstrap.GetIdentifierCodeFromToken(ctx, r)
		m              = model.Contract{App: h.App}
		req            = request.UpdatePasswordReq{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	err = m.UpdatePasswordUser(h.DB, ctx, userIdentifier, req.OldPassword, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

func (h *Contract) GetAllPlayerActivities(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
		code = chi.URLParam(r, "code")
		res  = make([]response.PlayerActivitiesRes, 0)
	)

	data, err := m.GetPlayerAndOtherActivities(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		if userPointType, ok := utils.GetUserPointTypeByValue(v.DataSource); ok {
			v.TitleDescription = strings.Replace(utils.UserPointTypeText[userPointType], "[username]", v.UserName, 1)
			v.TitleDescription = strings.Replace(v.TitleDescription, "[name]", v.ItemName, 1)
			v.TitleDescription = strings.Replace(v.TitleDescription, "[vp]", fmt.Sprint(v.Point), 1)
		}

		res = append(res, response.PlayerActivitiesRes{
			DataSource:       v.DataSource,
			TitleDescription: v.TitleDescription,
			SourceCode:       v.SourceCode,
			UserName:         v.UserName,
			UserCode:         v.UserCode,
			UserImageUrl:     v.UserImageUrl,
			UserStyle:        response.UserStyleRes{Color: v.UserStyle.Color},
			ItemName:         v.ItemName,
			ItemCode:         v.ItemCode,
			ItemImgUrl:       v.ItemImgUrl,
			Point:            v.Point,
			CreatedDate:      v.CreatedDate.Format(time.RFC3339),
		})
	}

	// Populate response
	h.SendSuccess(w, res, nil)
}

func (h *Contract) GetUserPointActivities(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
		code = chi.URLParam(r, "code")
		res  = make([]response.PlayerActivitiesRes, 0)
	)

	data, err := m.GetUserPointActivities(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.PlayerActivitiesRes{
			TitleDescription: v.TitleDescription,
			DataSource:       v.DataSource,
			SourceCode:       v.SourceCode,
			UserName:         v.UserName,
			UserStyle:        response.UserStyleRes{Color: v.UserStyle.Color},
			Point:            v.Point,
			CreatedDate:      v.CreatedDate.Format(time.RFC3339),
		})
	}

	// Populate response
	h.SendSuccess(w, res, nil)
}

func (h *Contract) GetUserPointHistories(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		code  = chi.URLParam(r, "code")
		res   = make([]response.UserPointHistoryResponse, 0)
		param = request.UserPointHistoryParam{}
	)

	// Define urlQuery and Parse
	_ = param.Parse(r.URL.Query())

	data, param, err := m.GetUserPointHistories(h.DB, ctx, code, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	tz, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		h.Log.File().Error("handler.GetUserPointHistories: failed to load timezone", err)
		tz = time.UTC
	}

	for _, v := range data {
		res = append(res, response.UserPointHistoryResponse{
			SourceId:       v.SourceId,
			SourceUserCode: v.SourceUserCode,
			SourceType:     v.SourceType,
			SourceCode:     v.SourceCode,
			SourceName:     v.SourceName,
			Point:          v.Point,
			CreatedDate:    v.CreatedDate.In(tz).Format(utils.DATE_TIME_FORMAT),
		})
	}

	// Populate response
	h.SendSuccess(w, res, param)
}

func (h *Contract) AdjustUserPointAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
		code = chi.URLParam(r, "code")
		req  = request.UserPointAdjustmentRequest{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	user, err := m.GetUserByUserCode(h.DB, ctx, code)
	if err != nil {
		h.SendInternalServerErr(w, err.Error())
		return
	}

	tx, err := h.DB.Begin(ctx)
	if err != nil {
		h.SendInternalServerErr(w, err.Error())
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
		tx.Commit(ctx)
	}()

	dataSource := utils.UserPointType["POINT_ADD"]
	sourceCode := "P-ADD"
	point := req.Point

	if req.AdjustmentType == "subtract" {
		dataSource = utils.UserPointType["POINT_SUB"]
		sourceCode = "P-SUB"
		point = -req.Point
	}

	m.AddUserPoint(tx, ctx, int64(user.ID), dataSource, sourceCode, point)

	h.SendSuccess(w, nil, nil)
}

// DeleteUserAct ...
func (h *Contract) DeleteUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		code     = chi.URLParam(r, "code")
		ctx      = context.TODO()
		m        = model.Contract{App: h.App}
		userCode = bootstrap.GetIdentifierCodeFromToken(ctx, r)
		role     = bootstrap.GetIdentifierRoleFromToken(ctx, r)
	)

	if role == utils.User && userCode != code {
		h.SendForbidden(w, utils.ErrForbiddenAuth)
	}

	err = m.DeleteUserByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	h.SendSuccess(w, nil, nil)

}

func (h *Contract) StoreUserStyleAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
		req  = request.StoreUserStyleReq{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	user, err := m.GetUserByUserCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Store user style
	err = m.StoreUserStyle(h.DB, ctx, int64(user.ID), model.UserStyle{Color: req.Color})
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}
