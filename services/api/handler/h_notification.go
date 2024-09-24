package handler

import (
	"dots-api/bootstrap"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// GetNotificationList handles HTTP request to get list of notifications
func (h *Contract) GetNotificationList(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = r.Context()
		m              = model.Contract{App: h.App}
		userIdentifier = bootstrap.GetIdentifierCodeFromToken(ctx, r)
		res            = make([]response.NotificationResp, 0)
		param          = request.NotificationParam{}
	)

	// Define URL query and Parse
	err = param.ParseNotification(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetListNotifications(h.DB, ctx, param, userIdentifier)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.NotificationResp{
			ReceiverSource:   v.ReceiverSource,
			ReceiverCode:     v.ReceiverCode,
			NotificationCode: v.NotificationCode,
			TransactionCode:  v.TransactionCode,
			Type:             v.Type,
			Title:            v.Title.String,
			Description:      v.Description.String,
			StatusRead:       v.StatusRead,
			ImageUrl:         v.ImageUrl.String,
			CreatedDate:      v.CreatedDate.Format(time.RFC3339),
			UpdatedDate:      v.UpdatedDate.Time.Format(time.RFC3339),
		})
	}

	h.SendSuccess(w, res, param)
}

// UpdateNotificationIsSeen handles HTTP request to update the is_seen status of a notification
func (h *Contract) UpdateNotificationIsSeenAct(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		req       = request.NotificationIsSeenReq{}
		notifCode = chi.URLParam(r, "code")
		ctx       = r.Context()
		m         = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	err = m.UpdateNotificationIsSeenByNotificationCode(h.DB, ctx, notifCode, req.IsSeen)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}
