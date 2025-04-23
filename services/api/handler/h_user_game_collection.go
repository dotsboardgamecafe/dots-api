package handler

import (
	"context"
	"dots-api/lib/rabbit"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetUserGameCollectionAct ...
func (h *Contract) GetUserGameCollectionAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		res   = make([]response.UserGameCollectionRes, 0)
		code  = chi.URLParam(r, "code")
		param = request.UserGameCollectionParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseUserGameCollection(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetUserGameCollections(h.DB, ctx, code, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.UserGameCollectionRes{
			UserId:       v.UserId.Int64,
			UserCode:     v.UserCode.String,
			GameCode:     v.GameCode.String,
			GameId:       v.GameId.Int64,
			GameName:     v.GameName.String,
			GameImageUrl: v.GameImageUrl.String,
			CreatedDate:  v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		})
	}

	h.SendSuccess(w, res, param)
}

// AddUserGameCollectionAct ...
func (h *Contract) AddUserGameCollectionAct(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		ctx     = context.TODO()
		m       = model.Contract{App: h.App}
		code    = chi.URLParam(r, "code")
		payload = request.UserGameCollectionAddPayload{}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &payload); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	game, err := m.GetGameByCode(h.DB, ctx, payload.GameCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	userID, err := m.GetUserIdByUserCode(h.DB, ctx, code)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.AddUserGameCollections(h.DB, ctx, userID, game.Id)
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
		}
		tx.Commit(ctx)
	}()

	// Add user point that count as activity, smh
	err = m.AddUserPoint(tx, ctx, userID, utils.UserPointType["GAME_COLLECTION"], game.GameCode, utils.FirstTimePlayed.Int())
	if err != nil {
		h.Log.File().Error("h.AddUserGameCollectionAct: "+err.Error(), code, game.GameCode)
	}

	go func(ctx context.Context, userID int64) {
		// Publisher badge
		queueData := rabbit.QueueDataPayload(
			rabbit.QueueUserBadge,
			rabbit.QueueUserBadgeReq(
				utils.PlayingGames,
				userID,
			),
		)
		queueHost := m.Config.GetString("queue.rabbitmq.host")
		err = rabbit.PublishQueue(ctx, queueHost, queueData)
		if err != nil {
			log.Printf("Error : %s", err)
		}
	}(ctx, userID)

	h.SendSuccess(
		w,
		response.UserGameCollectionRes{
			UserId:       userID,
			UserCode:     code,
			GameCode:     payload.GameCode,
			GameId:       game.Id,
			GameName:     game.Name,
			GameImageUrl: game.ImageUrl,
		},
		nil,
	)
}
