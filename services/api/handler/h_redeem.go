package handler

import (
	"context"
	"dots-api/bootstrap"
	POS "dots-api/lib/point_of_sale"
	"dots-api/lib/point_of_sale/sub_modules"
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

/*
APP SIDE
- Listing all user redeems history
- Get redeem detail
- Redeem an Invoice
*/
func (h *Contract) GetRedeemHistory(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		m              = model.Contract{App: h.App}
		res            = make([]response.UserReedemHistoryRes, 0)
		param          = request.UserRedeemHistoryParam{}
		userIdentifier = bootstrap.GetIdentifierCodeFromToken(ctx, r)
	)

	// Define urlQuery and Parse
	err = param.ParseUserRedeemHistory(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetUserRedeemHistories(h.DB, ctx, param, userIdentifier)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		createdDatePlusXMinutes := utils.AddTime(v.CreatedDate, (10 * time.Minute))
		isNew := createdDatePlusXMinutes.After(time.Now())

		res = append(res, response.UserReedemHistoryRes{
			UserCode:           v.UserCode,
			CustomId:           v.CustomId,
			PointEarned:        v.PointEarned,
			InvoiceCode:        fmt.Sprintf("#%s", v.InvoiceCode),
			InvoiceAmount:      v.InvoiceAmount,
			InvoiceDescription: v.InvoiceDescription,
			IsInvoiceNew:       &isNew,
			ClaimedDate:        v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		})
	}

	h.SendSuccess(w, res, nil)
}

func (h *Contract) GetRedeemDetail(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		userIdentifier = bootstrap.GetIdentifierCodeFromToken(ctx, r)
		invoiceCode    = chi.URLParam(r, "invoice_code")
		m              = model.Contract{App: h.App}
	)

	data, err := m.GetUserRedeemDetail(h.DB, ctx, userIdentifier, invoiceCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, response.UserReedemHistoryRes{
		UserCode:           data.UserCode,
		CustomId:           data.CustomId,
		PointEarned:        data.PointEarned,
		InvoiceCode:        fmt.Sprintf("#%s", data.InvoiceCode),
		InvoiceAmount:      data.InvoiceAmount,
		InvoiceDescription: data.InvoiceDescription,
		ClaimedDate:        data.CreatedDate.Format(utils.DATE_TIME_FORMAT),
	}, nil)
}

func (h *Contract) Redeem(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		req            request.UserRedeemRequest
		ctx            = context.TODO()
		userIdentifier = bootstrap.GetIdentifierCodeFromToken(ctx, r)
		m              = model.Contract{App: h.App}
		redeemCode     = utils.GeneratePrefixCode(utils.RedeemPrefix)
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	userId, err := m.GetUserIdByUserCode(h.DB, ctx, userIdentifier)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Hit API OLSERA POS (to fetch invoice amount and invoice products)
	PointOfSaleSystem, err := POS.GetPointOfSale("Olsera")
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// [Staging Only] Add conditional checker to allow redeem the same invoice
	if !PointOfSaleSystem.AllowToRedeemTheSameInvoice() {
		isExist, err := m.IsInvoiceCodeExist(h.DB, ctx, req.InvoiceCode)

		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		// Throw error if invoice code exist under this user
		if isExist {
			h.SendUnprocessableEntity(w, fmt.Sprintf("Invoice #%s sudah digunakan", req.InvoiceCode))
			return
		}
	}

	// Fetching Invoice
	_, invoiceDetail, err := PointOfSaleSystem.GetInvoiceCodeFromList(req.InvoiceCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Assign Payload
	redeemedPayload := generatePayload(redeemCode, utils.RedeemPlatform["APP"], *invoiceDetail)

	// Save & Calculate Earned Point
	earnedPoint, err := m.RedeemInvoice(h.DB, ctx, userId, *redeemedPayload)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Publisher badge
	queueData := rabbit.QueueDataPayload(
		rabbit.QueueUserBadge,
		rabbit.QueueUserBadgeReq(
			utils.TotalSpend,
			userId,
		),
	)
	queueHost := m.Config.GetString("queue.rabbitmq.host")
	err = rabbit.PublishQueue(ctx, queueHost, queueData)
	if err != nil {
		log.Printf("Error : %s", err)
	}

	// Populate response
	h.SendSuccess(w, response.UserReedemHistoryRes{
		UserCode:    userIdentifier,
		PointEarned: earnedPoint,
		InvoiceCode: fmt.Sprintf("#%s", redeemedPayload.InvoiceCode),
	}, nil)
}

/*
CMS SIDE
- Listing all claimed invoices
- Claim an Invoice
*/
func (h *Contract) GetInvoicesClaimedHistory(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		ctx      = context.TODO()
		m        = model.Contract{App: h.App}
		res      = make([]response.UserClaimHistoryRes, 0)
		param    = request.UserClaimedHistoryParam{}
		userCode = chi.URLParam(r, "user_code")
	)

	// Define urlQuery and Parse
	err = param.ParseUserRedeemHistory(r.URL.Query())
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	data, param, err := m.GetClaimedInvoiceHistories(h.DB, ctx, param, userCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, v := range data {
		res = append(res, response.UserClaimHistoryRes{
			InvoiceCode:   fmt.Sprintf("#%s", v.InvoiceCode),
			InvoiceAmount: v.InvoiceAmount,
			InvoiceItems:  v.ParseInvoiceInformation().OrderItems,
			ClaimedTime:   v.CreatedDate.Format(utils.DATE_TIME_FORMAT),
			ClaimedDate:   v.CreatedDate.Format(utils.DATE_DAY_FORMAT),
		})
	}

	h.SendSuccess(w, res, param)
}

func (h *Contract) Claim(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		req        request.UserRedeemRequest
		ctx        = context.TODO()
		userCode   = chi.URLParam(r, "user_code")
		m          = model.Contract{App: h.App}
		redeemCode = utils.GeneratePrefixCode(utils.RedeemPrefix)
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	userId, err := m.GetUserIdByUserCode(h.DB, ctx, userCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Hit API OLSERA POS (to fetch invoice amount and invoice products)
	PointOfSaleSystem, err := POS.GetPointOfSale("Olsera")
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// [Staging Only] Add conditional checker to allow redeem the same invoice
	if !PointOfSaleSystem.AllowToRedeemTheSameInvoice() {
		isExist, err := m.IsInvoiceCodeExist(h.DB, ctx, req.InvoiceCode)

		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		// Throw error if invoice code exist under this user
		if isExist {
			h.SendUnprocessableEntity(w, fmt.Sprintf("Invoice #%s sudah digunakan", req.InvoiceCode))
			return
		}
	}

	// Fetching Invoice
	_, invoiceDetail, err := PointOfSaleSystem.GetInvoiceCodeFromList(req.InvoiceCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Assign Payload
	redeemedPayload := generatePayload(redeemCode, utils.RedeemPlatform["CMS"], *invoiceDetail)

	// Save & Calculate Earned Point
	earnedPoint, err := m.RedeemInvoice(h.DB, ctx, userId, *redeemedPayload)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Publisher badge
	queueData := rabbit.QueueDataPayload(
		rabbit.QueueUserBadge,
		rabbit.QueueUserBadgeReq(
			utils.TotalSpend,
			userId,
		),
	)
	queueHost := m.Config.GetString("queue.rabbitmq.host")
	err = rabbit.PublishQueue(ctx, queueHost, queueData)
	if err != nil {
		log.Printf("Error : %s", err)
	}

	// Populate response
	h.SendSuccess(w, response.UserReedemHistoryRes{
		UserCode:    userCode,
		PointEarned: earnedPoint,
		InvoiceCode: fmt.Sprintf("#%s", redeemedPayload.InvoiceCode),
	}, nil)
}

func (h *Contract) SyncInvoice(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		req      request.UserRedeemRequest
		ctx      = context.TODO()
		userCode = chi.URLParam(r, "user_code")
		m        = model.Contract{App: h.App}
	)

	// Binding and Validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	userId, err := m.GetUserIdByUserCode(h.DB, ctx, userCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Hit API OLSERA POS (to fetch invoice amount and invoice products)
	PointOfSaleSystem, err := POS.GetPointOfSale("Olsera")
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Fetching Invoice
	_, invoiceDetail, err := PointOfSaleSystem.GetInvoiceCodeFromList(req.InvoiceCode)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	newestInvoiceInfo, _ := invoiceDetail.SavedInformationToJSONString()

	err = m.SyncRedeemInformation(h.DB, ctx, userId, req.InvoiceCode, newestInvoiceInfo)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, fmt.Sprintf("Invoice #%s berhasil di-sync", req.InvoiceCode), nil)
}

/*
PRIVATE FUNCTION
*/
func generatePayload(redeemCode string, requestedPlatform string, invoiceDetail sub_modules.CloseOrderDetail) *model.UserRedeemPayload {
	parsedTotalAmount, _ := invoiceDetail.GetTotalAmount()
	listOfProducts := invoiceDetail.GetLineOfProducts()
	invoiceInfo, _ := invoiceDetail.SavedInformationToJSONString()

	redeemPayload := new(model.UserRedeemPayload)
	redeemPayload.CustomId = redeemCode
	redeemPayload.InvoiceCode = invoiceDetail.Data.OrderNo
	redeemPayload.InvoiceAmount = parsedTotalAmount
	redeemPayload.InvoiceDescription = listOfProducts
	redeemPayload.PointEarned = utils.CalculateUserRedeemPoint(parsedTotalAmount)
	redeemPayload.Information = invoiceInfo
	redeemPayload.RequestedPlatform = requestedPlatform

	return redeemPayload
}
