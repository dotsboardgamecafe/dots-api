package handler

import (
	"context"
	"dots-api/lib/onesignal"
	"dots-api/lib/rabbit"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
)

func (h *Contract) TransactionCallback(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		m              = model.Contract{App: h.App}
		req            = request.InvoiceCallbackRequest{}
		xPlayer        string
		bannerImageUri string
	)

	if err = h.Bind(r, &req); err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	if err = h.Validator.Driver.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	// Marshal request payload
	response, err := json.Marshal(req)
	if err != nil {
		return
	}

	trx, err := m.GetInvoiceTrxByCode(h.DB, ctx, req.ID)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Check if the payment price same with order price
	if int64(trx.Price) != int64(req.Amount) {
		h.SendBadRequest(w, "Paid amount not match with order amount")
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

	err = m.UpdateInvoiceTrx(tx, ctx, req.ID, req.PaymentMethod, req.Status, string(response))
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	statusParticipant := "cancel"
	// Push Notification SUCCESS PAYMENT
	if req.Status == "PAID" {
		statusParticipant = "active"
	}

	// Logic handle status trnsaction
	switch trx.DataSource {
	case utils.UserPointType["ROOM_TYPE"]:
		// execute update status participant room
		participant, err := m.GetParticipantByRoomCodeAndUserCode(h.DB, ctx, trx.SourceCode, trx.UserCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		//update status participant
		err = m.UpdateRoomParticipant(tx, ctx, participant.RoomId, participant.UserId, participant.StatusWinner, participant.Position, statusParticipant, participant.AdditionalInfo.String, participant.RewardPoint.Int64, participant.TransactionCode.String)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		if req.Status == "PAID" {
			// Add user point from price
			earnedPoint := utils.CalculateUserRedeemPoint(trx.Price)
			err = m.AddUserPoint(tx, ctx, participant.UserId, trx.DataSource, trx.SourceCode, earnedPoint)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}

			// Add user point from vp point participation
			err = m.AddUserPoint(tx, ctx, participant.UserId, trx.DataSource, trx.SourceCode, participant.ParticipationPoint)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}
		}

		// Publisher badge
		queueData := rabbit.QueueDataPayload(
			rabbit.QueueUserBadge,
			rabbit.QueueUserBadgeReq(
				utils.TimeLimit,
				trx.UserId,
			),
		)

		queueHost := m.Config.GetString("queue.rabbitmq.host")
		err = rabbit.PublishQueue(ctx, queueHost, queueData)
		if err != nil {
			log.Printf("Error : %s", err)
		}

		xPlayer = participant.UserXPlayer
		bannerImageUri = participant.RoomBannerUri

	case utils.UserPointType["TOURNAMENT_TYPE"]:
		// execute update status participant tournament

		// Get tournament by code
		trnm, err := m.GetTournamentByCode(h.DB, ctx, trx.SourceCode)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		// Check if already booked
		participant, err := m.GetParticipantByTournamentCodeAndUserCode(h.DB, ctx, trnm.TournamentCode, trx.UserCode)
		if err != nil && err.Error() != utils.EmptyData {
			h.SendBadRequest(w, err.Error())
			return
		}

		//update status participant
		err = m.UpdateTournamentParticipant(tx, ctx, participant.TournamentId, participant.UserId, participant.StatusWinner, participant.Position, statusParticipant, participant.AdditionalInfo.String, participant.RewardPoint.Int64, participant.TransactionCode.String)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		if req.Status == "PAID" {
			// Add user point
			earnedPoint := utils.CalculateUserRedeemPoint(trx.Price)
			err = m.AddUserPoint(tx, ctx, participant.UserId, trx.DataSource, trx.SourceCode, earnedPoint)
			if err != nil {
				h.SendBadRequest(w, err.Error())
				return
			}
		}

		xPlayer = participant.UserXPlayer
		bannerImageUri = participant.TournamentBannerUri
	default:
		h.SendBadRequest(w, "undefined type transaction")
		return
	}

	// Publisher badge total spent on status PAID
	if req.Status == "PAID" {
		queueData := rabbit.QueueDataPayload(
			rabbit.QueueUserBadge,
			rabbit.QueueUserBadgeReq(
				utils.TotalSpend,
				trx.UserId,
			),
		)

		queueHost := m.Config.GetString("queue.rabbitmq.host")
		err = rabbit.PublishQueue(ctx, queueHost, queueData)
		if err != nil {
			log.Printf("Error : %s", err)
		}
	}

	// Publisher badge
	queueData := rabbit.QueueDataPayload(
		rabbit.QueueUserBadge,
		rabbit.QueueUserBadgeReq(
			utils.SpesificBoardGameCategory,
			trx.UserId,
		),
	)

	queueHost := m.Config.GetString("queue.rabbitmq.host")
	err = rabbit.PublishQueue(ctx, queueHost, queueData)
	if err != nil {
		log.Printf("Error : %s", err)
	}

	// Notification Handler
	sendNotification(h.DB, ctx, m, req.Status, trx, xPlayer, bannerImageUri)

	h.SendSuccess(w, nil, nil)
}

// Send Notification via email & PN - private function
func sendNotification(db *pgxpool.Pool, ctx context.Context, m model.Contract, paymentStatus string, data model.OriginUserTransactionEnt, xPlayer string, bannerImageUri string) {
	// email := mail.New(m.App)
	onesignal := onesignal.New(m.App)

	if paymentStatus == "PAID" {
		notifCode := utils.GeneratePrefixCode(utils.NotifPrefix)

		description := "Pembayaran Anda telah berhasil diproses. Anda telah berhasil masuk room / tournament!."

		descriptionJSON, err := json.Marshal(description)
		if err != nil {
			log.Printf("Error : %s", err)
		}

		// Insert data into db
		err = m.AddNotification(db, ctx, notifCode, "user", data.UserCode, data.TransactionCode, utils.SuccessPaymentType, utils.SuccessPaymentTitle, descriptionJSON, bannerImageUri)
		if err != nil {
			log.Printf("Error : %s", err)
		}

		// COMMENT OUT FOR NOW - Use Xendit Payment Email Service
		// ======================================================
		// EmailData := mail.SuccessPaymentData{
		// 	Name:          data.UserFullname,
		// 	TransactionId: data.AggregatorCode,
		// 	PaymentDate:   data.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		// 	Amount:        int64(data.Price),
		// }

		// err = email.SendMail(mail.SuccessfulPayment, mail.MailSubj[mail.SuccessfulPayment], data.UserEmail, EmailData)
		// if err != nil {
		// 	log.Printf("Error : %s", err)
		// }

		_, err = onesignal.CreateOSNotifications(xPlayer, utils.SuccessPaymentTitle, description, utils.Transaction)
		if err != nil {
			log.Printf("Error : %s", err)
		}

		return
	}

	// Generate Notification code
	notifCode := utils.GeneratePrefixCode(utils.NotifPrefix)

	var (
		notifTitle       = utils.ExpiredPaymentTitle
		notifType        = utils.ExpiredPaymentType
		notifDescription = utils.ExpiredPaymentDescription
	)

	if paymentStatus != utils.PaymentStatus["EXPIRED"] {
		notifTitle = utils.FailPaymentTitle
		notifType = utils.FailPaymentType
		notifDescription = utils.FailPaymentDescription
	}

	descriptionJSON, err := json.Marshal(notifDescription)
	if err != nil {
		log.Printf("Error : %s", err)
	}

	// Insert data into db
	err = m.AddNotification(m.DB, ctx, notifCode, "user", data.UserCode, data.TransactionCode, notifType, notifTitle, descriptionJSON, bannerImageUri)
	if err != nil {
		log.Printf("Error : %s", err)
	}

	// COMMENT OUT FOR NOW - Use Xendit Payment Email Service
	// ======================================================
	// EmailData := mail.FailedPaymentData{
	// 	Name: data.UserFullname,
	// }

	// err = email.SendMail(mail.SuccessfulPayment, mail.MailSubj[mail.SuccessfulPayment], data.UserEmail, EmailData)
	// if err != nil {
	// 	log.Printf("Error : %s", err)
	// }

	return
}
