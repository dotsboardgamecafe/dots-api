package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"dots-api/services/api/request"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	UserRedeemHistoryEnt struct {
		Id                 int64        `db:"id"`
		UserId             int64        `db:"user_id"`
		CustomId           string       `db:"custom_id"`
		InvoiceCode        string       `db:"invoice_code"`
		InvoiceAmount      int          `db:"invoice_amount"`
		InvoiceDescription string       `db:"invoice_description"`
		CreatedDate        time.Time    `db:"created_date"`
		UpdatedDate        sql.NullTime `db:"updated_date"`
	}

	UserRedeemHistoryDetailRes struct {
		Id                 int64        `db:"id"`
		UserId             int64        `db:"user_id"`
		UserCode           string       `db:"user_code"`
		CustomId           string       `db:"custom_id"`
		PointEarned        int          `db:"point_earned"`
		InvoiceCode        string       `db:"invoice_code"`
		InvoiceAmount      int          `db:"invoice_amount"`
		InvoiceDescription string       `db:"invoice_description"`
		InvoiceInfo        string       `db:"redeem_information"`
		RequestedPlatform  string       `db:"requested_platform"`
		CreatedDate        time.Time    `db:"created_date"`
		UpdatedDate        sql.NullTime `db:"updated_date"`
	}

	UserRedeemPayload struct {
		UserCode           string
		PointEarned        int
		CustomId           string
		InvoiceCode        string
		InvoiceAmount      float64
		InvoiceDescription string
		Information        []byte
		RequestedPlatform  string
	}

	UserClaimedInvoice struct {
		OrderId          int64            `json:"order_id"`
		OrderNo          string           `json:"order_no"`
		OrderStatus      string           `json:"order_status"`
		OrderCreatedTime string           `json:"created_time"`
		OrderTotalAmount string           `json:"total_amount"`
		OrderTotalQty    int              `json:"total_item_qty"`
		OrderItems       []UserOrderItems `json:"order_items"`
	}

	UserOrderItems struct {
		Qty                int     `json:"qty"`
		ProductPrice       float64 `json:"price"`
		ProductId          int64   `json:"product_id"`
		ProductSKU         string  `json:"sku"`
		ProductName        string  `json:"name"`
		CategoryName       string  `json:"category_name"`
		ClassificationName string  `json:"klasifikasi"`
	}
)

func (c *Contract) GetUserRedeemHistories(db *pgxpool.Pool, ctx context.Context, param request.UserRedeemHistoryParam, userCode string) ([]UserRedeemHistoryDetailRes, request.UserRedeemHistoryParam, error) {
	var (
		err        error
		list       []UserRedeemHistoryDetailRes
		paramQuery []interface{}
		totalData  int

		query = `SELECT
			u.id AS user_id,
			user_code,
			custom_id,
			users_points.point AS point_earned,
			invoice_code,
			invoice_amount,
			description AS invoice_description,
			user_redeem_histories.created_date AS created_date,
			user_redeem_histories.updated_date AS updated_date
		FROM user_redeem_histories
			JOIN users u ON u.id = user_redeem_histories.user_id
			JOIN users_points ON user_redeem_histories.custom_id = users_points.source_code
		WHERE u.user_code = $1
		ORDER BY user_redeem_histories.id DESC LIMIT 5`
	)

	paramQuery = append(paramQuery, userCode)

	// Populate Search
	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS total`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetUserRedeemHistories", err, utils.ErrCountingListUserRedeemHistory)
		}
		param.Count = totalData
	}

	// Select Max Page
	// if param.Count > param.Limit && param.Page > int(param.Count/param.Limit) {
	// 	param.MaxPage = int(math.Ceil(float64(param.Count) / float64(param.Limit)))
	// }

	// Limit and Offset
	// param.Offset = (param.Page - 1) * param.Limit
	// paramQuery = append(paramQuery, param.Offset)
	// query += fmt.Sprintf(" OFFSET $%d ", len(paramQuery))

	// paramQuery = append(paramQuery, param.Limit)
	// query += fmt.Sprintf(" LIMIT $%d ", len(paramQuery))

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, param, c.errHandler("model.GetUserRedeemHistories", err, utils.ErrGettingListUserRedeemHistory)
	}

	defer rows.Close()
	for rows.Next() {
		var data UserRedeemHistoryDetailRes
		err = rows.Scan(
			&data.UserId, &data.UserCode,
			&data.CustomId, &data.PointEarned,
			&data.InvoiceCode, &data.InvoiceAmount, &data.InvoiceDescription,
			&data.CreatedDate, &data.UpdatedDate,
		)
		if err != nil {
			return list, param, c.errHandler("model.GetUserRedeemHistories", err, utils.ErrScanningListUserRedeemHistory)
		}
		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) GetClaimedInvoiceHistories(db *pgxpool.Pool, ctx context.Context, param request.UserClaimedHistoryParam, userCode string) ([]UserRedeemHistoryDetailRes, request.UserClaimedHistoryParam, error) {
	var (
		err        error
		list       []UserRedeemHistoryDetailRes
		where      []string
		paramQuery []interface{}
		totalData  int

		query = `
		SELECT
			u.id AS user_id,
			user_code,
			custom_id,
			invoice_code,
			invoice_amount,
			description AS invoice_description,
			redeem_information AS invoice_information,
			user_redeem_histories.created_date AS created_date,
			user_redeem_histories.updated_date AS updated_date
		FROM user_redeem_histories
			JOIN users u ON u.id = user_redeem_histories.user_id
		WHERE u.user_code = $1
	`
	)

	paramQuery = append(paramQuery, userCode)

	// Populate Search
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		orWhere = append(orWhere, fmt.Sprintf("invoice_code iLIKE $%d", len(paramQuery)))
		where = append(where, " AND ("+strings.Join(orWhere, " OR ")+")")
	}

	// Append All Where Conditions
	if len(where) > 0 {
		query += " " + strings.Join(where, " AND ")
	}

	// Populate Search
	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS total`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetClaimedInvoiceHistories", err, utils.ErrCountingListUserRedeemHistory)
		}
		param.Count = totalData
	}

	// Select Max Page
	if param.Count > param.Limit && param.Page > int(param.Count/param.Limit) {
		param.MaxPage = int(math.Ceil(float64(param.Count) / float64(param.Limit)))
	}

	// Limit and Offset
	param.Offset = (param.Page - 1) * param.Limit
	paramQuery = append(paramQuery, param.Offset)
	query += fmt.Sprintf(" OFFSET $%d ", len(paramQuery))

	paramQuery = append(paramQuery, param.Limit)
	query += fmt.Sprintf(" LIMIT $%d ", len(paramQuery))

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, param, c.errHandler("model.GetClaimedInvoiceHistories", err, utils.ErrGettingListUserRedeemHistory)
	}

	defer rows.Close()
	for rows.Next() {
		var data UserRedeemHistoryDetailRes
		err = rows.Scan(
			&data.UserId, &data.UserCode,
			&data.CustomId, &data.InvoiceCode, &data.InvoiceAmount,
			&data.InvoiceDescription, &data.InvoiceInfo,
			&data.CreatedDate, &data.UpdatedDate,
		)
		if err != nil {
			return list, param, c.errHandler("model.GetClaimedInvoiceHistories", err, utils.ErrScanningListUserRedeemHistory)
		}
		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) GetUserRedeemDetail(db *pgxpool.Pool, ctx context.Context, userCode string, invoiceCode string) (UserRedeemHistoryDetailRes, error) {
	var (
		err  error
		data UserRedeemHistoryDetailRes

		query = `SELECT
			u.id AS user_id,
			user_code,
			custom_id,
			users_points.point AS point_earned,
			invoice_code,
			invoice_amount,
			description AS invoice_description,
			user_redeem_histories.created_date AS created_date,
			user_redeem_histories.updated_date AS updated_date
		FROM user_redeem_histories
			JOIN users u ON u.id = user_redeem_histories.user_id
			JOIN users_points ON user_redeem_histories.custom_id = users_points.source_code
		WHERE u.user_code = $1 AND user_redeem_histories.invoice_code = $2`
	)

	err = db.QueryRow(ctx, query, userCode, invoiceCode).Scan(
		&data.UserId, &data.UserCode,
		&data.CustomId, &data.PointEarned,
		&data.InvoiceCode, &data.InvoiceAmount, &data.InvoiceDescription,
		&data.CreatedDate, &data.UpdatedDate,
	)

	if err != nil {
		return data, c.errHandler("model.GetUserRedeemDetail", err, utils.ErrGettingtUserRedeemHistoryDetailCode)
	}

	return data, nil
}

func (c *Contract) RedeemInvoice(db *pgxpool.Pool, ctx context.Context, userId int64, userRedeemData UserRedeemPayload) (earnedPoint int, err error) {
	currentDateTime := time.Now().In(time.UTC)

	// Create a helper function for preparing failure results.
	fail := func(err error) (int, error) {
		return 0, fmt.Errorf("RedeemInvoice: %v", err)
	}

	// Begin transaction (tx)
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fail(err)
	}

	// Error & rollback handling
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_ = tx.Rollback(closeCtx)
	}()

	// Insert new row for user_redeem_histories
	_, err = tx.Exec(ctx, `INSERT INTO user_redeem_histories(
		user_id,
		custom_id,
		invoice_code,
		invoice_amount,
		description,
		redeem_information,
		requested_platform,
		created_date,
		updated_date
	) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		userId,
		userRedeemData.CustomId,
		userRedeemData.InvoiceCode,
		userRedeemData.InvoiceAmount,
		userRedeemData.InvoiceDescription,
		userRedeemData.Information,
		userRedeemData.RequestedPlatform,
		currentDateTime,
		currentDateTime,
	)
	if err != nil {
		return fail(err)
	}

	// Insert new row for users_points
	err = c.AddUserPoint(tx, ctx,
		userId,
		utils.UserPointType["REDEEM_TYPE"],
		userRedeemData.CustomId,
		userRedeemData.PointEarned)
	if err != nil {
		return fail(err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fail(err)
	}

	return userRedeemData.PointEarned, nil
}

func (c *Contract) SyncRedeemInformation(db *pgxpool.Pool, ctx context.Context, userId int64, invoiceCode string, newestInvoiceInfo []byte) (err error) {
	sql := `UPDATE user_redeem_histories
		SET redeem_information = $1
		WHERE user_id = $2 AND invoice_code = $3;`

	_, err = db.Exec(ctx, sql, newestInvoiceInfo, userId, invoiceCode)
	if err != nil {
		return c.errHandler("model.SyncRedeemInformation", err, utils.ErrFailedUpdateRedeemInfo)
	}

	return nil
}

func (c *Contract) IsInvoiceCodeExist(db *pgxpool.Pool, ctx context.Context, InvoiceCode string) (bool, error) {
	var (
		isExist bool
		query   = `SELECT EXISTS(
			SELECT 1 FROM user_redeem_histories WHERE invoice_code = $1
		)`
	)

	err := db.QueryRow(ctx, query, InvoiceCode).Scan(&isExist)
	if err != nil {
		return false, c.errHandler("model.IsInvoiceCodeExist", err, "Error checking the existence of user invoice")
	}

	return isExist, nil
}

func (c *Contract) GetTotalInvoiceAmountByUserID(db *pgxpool.Pool, ctx context.Context, userID int64) (int, error) {
	var (
		totalAmount    pgtype.Numeric
		totalAmountInt int
	)

	query := `SELECT SUM(invoice_amount) FROM user_redeem_histories WHERE user_id = $1`

	err := db.QueryRow(ctx, query, userID).Scan(&totalAmount)
	if err != nil {
		return 0, fmt.Errorf("error executing query: %w", err)
	}

	// Check if the result is null
	if totalAmount.Int == nil {
		return 0, nil
	}

	// Convert pgtype.Numeric to int
	err = totalAmount.AssignTo(&totalAmountInt)
	if err != nil {
		return 0, fmt.Errorf("error converting totalAmount to int64: %w", err)
	}

	return totalAmountInt, nil
}

func (redeemHistory *UserRedeemHistoryDetailRes) ParseInvoiceInformation() UserClaimedInvoice {
	var invoiceDetail UserClaimedInvoice
	json.Unmarshal([]byte(redeemHistory.InvoiceInfo), &invoiceDetail)

	return invoiceDetail
}
