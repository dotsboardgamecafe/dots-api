package handler

import (
	"context"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetGameListAct ...
func (h *Contract) GetGameListAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.GameRes, 0)
		param = request.GameParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseGame(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetGameList(h.DB, ctx, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		// For mapping data game master
		var gameMasterRes []response.AdminRes

		gameMasters, err := m.GetGameAdmins(h.DB, ctx, v.Id)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		for _, v := range gameMasters {
			gameMasterRes = append(gameMasterRes, response.AdminRes{
				AdminCode:   v.AdminCode,
				Email:       v.Email,
				Name:        v.Name,
				UserName:    v.UserName,
				Status:      v.Status,
				ImageURL:    v.ImageURL,
				PhoneNumber: v.PhoneNumber,
			})
		}

		res = append(res, response.GameRes{
			CafeCode:           v.CafeCode,
			CafeName:           v.CafeName,
			GameCode:           v.GameCode,
			GameType:           v.GameType,
			Location:           v.Location,
			Name:               v.Name,
			ImageUrl:           v.ImageUrl,
			CollectionUrl:      response.BuildCollectionURLResp(v.CollectionUrl),
			Description:        v.Description,
			Status:             v.Status,
			Duration:           v.Duration,
			Difficulty:         v.Difficulty.String,
			Level:              v.Level,
			MinimalParticipant: v.MinimalParticipant.Int64,
			MaximumParticipant: v.MaximumParticipant,
			GameCategories:     response.BuildGameCategoryResp(v.GameCategories.String),
			GameMasters:        gameMasterRes,
			NumberOfPopularity: v.NumberOfPopularity,
			// GameCharacteristic: response.BuildGameCharacteristicResp(v.GameCharacteristic.String),
		})
	}

	h.SendSuccess(w, res, param)
}

// AddGameAct ...
func (h *Contract) AddGameAct(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req = request.GameReq{}
		ctx = context.TODO()
		m   = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusGame, req.Status) {
		h.SendBadRequest(w, "wrong status value for Game(active|inactive")
		return
	}

	// Get Cafe Id
	cafeId, err := m.GetCafeIdByCode(h.DB, ctx, req.CafeCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Generate Random Code
	code := utils.GeneratePrefixCode(utils.GamePrefix)
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Convert the array to JSON
	collectionUrl, err := json.Marshal(req.CollectionUrl)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	// adminId, err := m.GetAdminIdByCode(h.DB, ctx, req.AdminCode)
	// if err != nil {
	// 	h.SendBadRequest(w, err.Error())
	// 	tx.Rollback(ctx)
	// 	return
	// }

	gameId, err := m.AddGame(tx, ctx, cafeId, code, req.GameType, req.Name, req.ImageUrl, string(collectionUrl), req.Description, req.Difficulty, req.Status, req.Level, req.MinimalParticipant, req.MaximumParticipant, req.Duration)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		tx.Rollback(ctx)
		return
	}

	admins, err := m.GetAdminIdsByCode(h.DB, ctx, req.AdminCodes)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		tx.Rollback(ctx)
		return
	}

	err = m.SyncGameAdmins(tx, ctx, gameId, admins)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		tx.Rollback(ctx)
		return
	}

	// Insert game category
	for _, v := range req.GameCategories {
		err = m.InsertOneGameCategory(tx, ctx, gameId, v.CategoryName)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			tx.Rollback(ctx)
			return
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		tx.Rollback(ctx)
		return
	}

	data, err := m.GetGameByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, response.GameRes{
		CafeCode:           data.CafeCode,
		CafeName:           data.CafeName,
		GameCode:           data.GameCode,
		GameType:           data.GameType,
		Name:               data.Name,
		ImageUrl:           data.ImageUrl,
		CollectionUrl:      response.BuildCollectionURLResp(data.CollectionUrl),
		Description:        data.Description,
		Difficulty:         data.Difficulty.String,
		Level:              data.Level,
		Duration:           data.Duration,
		AdminCode:          data.AdminCode.String,
		MinimalParticipant: data.MinimalParticipant.Int64,
		MaximumParticipant: data.MaximumParticipant,
		Status:             data.Status,
		GameCategories:     response.BuildGameCategoryResp(data.GameCategories.String),
	}, nil)
}

// GetGameDetailAct ...
func (h *Contract) GetGameDetailAct(w http.ResponseWriter, r *http.Request) {
	var (
		err           error
		code          = chi.URLParam(r, "code")
		ctx           = context.TODO()
		m             = model.Contract{App: h.App}
		gameMasterRes []response.AdminRes
		dataPlayerRes []response.UsersHavePlayedGameHistoryRes
	)

	data, err := m.GetGameByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	gameMasters, err := m.GetGameAdmins(h.DB, ctx, data.Id)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	for _, v := range gameMasters {
		gameMasterRes = append(gameMasterRes, response.AdminRes{
			AdminCode:   v.AdminCode,
			Email:       v.Email,
			Name:        v.Name,
			UserName:    v.UserName,
			Status:      v.Status,
			ImageURL:    v.ImageURL,
			PhoneNumber: v.PhoneNumber,
		})
	}

	// Retrieve the list of players who have played the game and the total count
	playersData, totalPlayers, err := m.GetUsersHavePlayedGameHistory(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Build the response with player data
	for _, playerData := range playersData {
		playerResponse := response.UsersHavePlayedGameHistoryRes{
			GameId:    playerData.GameId.Int64,
			GameName:  playerData.GameName.String,
			UserCode:  playerData.UserCode.String,
			UserName:  playerData.UserName.String,
			UserImage: playerData.UserImage.String,
		}

		dataPlayerRes = append(dataPlayerRes, playerResponse)
	}

	h.SendSuccess(w, response.GameDetailRes{
		CafeCode:                  data.CafeCode,
		CafeName:                  data.CafeName,
		GameCode:                  data.GameCode,
		GameType:                  data.GameType,
		Location:                  data.Location,
		Name:                      data.Name,
		ImageUrl:                  data.ImageUrl,
		CollectionUrl:             response.BuildCollectionURLResp(data.CollectionUrl),
		Description:               data.Description,
		Status:                    data.Status,
		NumberOfPopularity:        data.NumberOfPopularity,
		Duration:                  data.Duration,
		Difficulty:                data.Difficulty.String,
		Level:                     data.Level,
		MinimalParticipant:        data.MinimalParticipant.Int64,
		MaximumParticipant:        data.MaximumParticipant,
		GameCategories:            response.BuildGameCategoryResp(data.GameCategories.String),
		GameRelated:               response.BuildGameRelatedResp(data.GameRelated.String),
		GameRooms:                 response.BuildGameAvailableRoomResResp(data.GameRoomAvailables.String),
		GameMasters:               gameMasterRes,
		UserHavePlayedGameHistory: dataPlayerRes,
		TotalPlayer:               totalPlayers,
	}, nil)
}

// UpdateGameAct ...
func (h *Contract) UpdateGameAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		req  = request.GameReq{}
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if !utils.Contains(utils.StatusGame, req.Status) {
		h.SendBadRequest(w, "wrong status value for Game(active|inactive")
		return
	}

	// Get Cafe Id
	cafeId, err := m.GetCafeIdByCode(h.DB, ctx, req.CafeCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Get Game Id
	gameId, err := m.GetGameIdByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Db tx start
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Convert the array to JSON
	collectionUrl, err := json.Marshal(req.CollectionUrl)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	err = m.UpdateGameByCode(tx, ctx, cafeId, code, req.GameType, req.Name, req.ImageUrl, string(collectionUrl), req.Description, req.Difficulty, req.Status, req.Level, req.MinimalParticipant, req.MaximumParticipant, req.Duration)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		tx.Rollback(ctx)
		return
	}

	admins, err := m.GetAdminIdsByCode(h.DB, ctx, req.AdminCodes)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		tx.Rollback(ctx)
		return
	}

	err = m.SyncGameAdmins(tx, ctx, gameId, admins)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		tx.Rollback(ctx)
		return
	}

	// Delete and re-insert game category
	err = m.DeleteGameCategory(tx, ctx, gameId)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		tx.Rollback(ctx)
		return
	}
	for _, v := range req.GameCategories {
		err = m.InsertOneGameCategory(tx, ctx, gameId, v.CategoryName)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			tx.Rollback(ctx)
			return
		}
	}

	// Delete and re-insert game characteristic
	err = m.DeleteGameCharacteristic(tx, ctx, gameId)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		tx.Rollback(ctx)
		return
	}

	// Db tx commit
	err = tx.Commit(ctx)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		tx.Rollback(ctx)
		return
	}

	data, err := m.GetGameByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, response.GameRes{
		CafeCode:       data.CafeCode,
		CafeName:       data.CafeName,
		GameCode:       data.GameCode,
		GameType:       data.GameType,
		Name:           data.Name,
		ImageUrl:       data.ImageUrl,
		CollectionUrl:  response.BuildCollectionURLResp(data.CollectionUrl),
		Description:    data.Description,
		Status:         data.Status,
		GameCategories: response.BuildGameCategoryResp(data.GameCategories.String),
		GameRelated:    response.BuildGameRelatedResp(data.GameRelated.String),
	}, nil)
}

func (h *Contract) DeleteGameAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	// Get the game ID by its code
	id, err := m.GetGameIdByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Check if the game is being used in any rooms or tournaments
	isExist, err := m.CheckExistGameUsed(h.DB, ctx, id)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}
	if isExist {
		h.SendBadRequest(w, utils.ErrForbiddenDeleteGame)
		return
	}

	// Delete the game by its ID
	err = m.DeleteGameById(h.DB, ctx, id)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Send success response
	h.SendSuccess(w, nil, nil)
}

// GetGameQRCodeAct ...
func (h *Contract) GetGameQRCodeAct(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code = chi.URLParam(r, "code")
		ctx  = context.TODO()
		m    = model.Contract{App: h.App}
	)

	data, err := m.GetGameByCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	qr, err := m.GetGameQRCodeByCode(ctx, data.GameCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, qr, nil)
}
