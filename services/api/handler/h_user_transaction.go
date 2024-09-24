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

func (h *Contract) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		m              = model.Contract{App: h.App}
		userIdentifier = chi.URLParam(r, "code")
		res            = make([]response.UserBookingSummaryRes, 0)
		param          = request.UserTransactionParam{}
	)

	// Define urlQuery and Parse
	err = param.ParseUserTransactionParam(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetUserTransactionList(h.DB, ctx, userIdentifier, param)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.UserBookingSummaryRes{
			DataSource:       v.DataSource,
			TransactionCode:  v.TransactionCode,
			GameName:         v.GameName,
			GameImgUrl:       v.GameImgUrl,
			Price:            v.Price,
			AwardedUserPoint: v.AwardedUserPoint,
			PaymentMethod:    v.PaymentMethod,
			Status:           v.Status,
			CreatedDate:      v.CreatedDate.Format(time.RFC3339),
		})
	}

	h.SendSuccess(w, res, param)
}

func (h *Contract) GetUserTransactionDetail(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		userIdentifier = chi.URLParam(r, "code")
		trxCode        = chi.URLParam(r, "trx_code")
		m              = model.Contract{App: h.App}
	)

	// Get Transaction Detail
	data, err := m.GetTransactionByCode(h.DB, ctx, userIdentifier, trxCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, response.UserBookingSummaryRes{
		DataSource:       data.DataSource,
		TransactionCode:  data.TransactionCode,
		GameName:         data.GameName,
		Price:            data.Price,
		AwardedUserPoint: data.AwardedUserPoint,
		PaymentMethod:    data.PaymentMethod,
		Status:           data.Status,
		CreatedDate:      data.CreatedDate.Format(utils.DATE_TIME_FORMAT),
	}, nil)
}
