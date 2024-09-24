package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	payment "dots-api/lib/xendit"
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

type OriginUserTransactionEnt struct {
	Id              int64        `db:"id"`
	UserId          int64        `db:"user_id"`
	UserCode        string       `db:"user_code"`
	UserFullname    string       `db:"user_fullname"`
	UserEmail       string       `db:"user_email"`
	DataSource      string       `db:"data_source"`
	SourceCode      string       `db:"source_code"`
	TransactionCode string       `db:"transaction_code"`
	AggregatorCode  string       `db:"aggregator_code"`
	Price           float64      `db:"price"`
	PaymentMethod   string       `db:"payment_method"`
	PaymentLink     string       `db:"payment_link"`
	RespPayload     string       `db:"resp_payload"`
	Status          string       `db:"status"`
	CreatedDate     time.Time    `db:"created_date"`
	UpdatedDate     sql.NullTime `db:"updated_date"`
	ExpiredDate     sql.NullTime `db:"expired_date"`
}

type UserTransactionEnt struct {
	Id               int64        `db:"id"`
	UserId           int64        `db:"user_id"`
	DataSource       string       `db:"data_source"`
	TransactionCode  string       `db:"transaction_code"`
	AggregatorCode   string       `db:"aggregator_code"`
	GameCode         string       `db:"game_code"`
	GameName         string       `db:"game_name"`
	GameImgUrl       string       `db:"game_img_url"`
	Price            float64      `db:"final_price_amount"`
	AwardedUserPoint int          `db:"awarded_user_point"`
	PaymentMethod    string       `db:"payment_method"`
	PaymentLink      string       `db:"payment_link"`
	RespPayload      string       `db:"resp_payload"`
	Status           string       `db:"status"`
	CreatedDate      time.Time    `db:"created_date"`
	UpdatedDate      sql.NullTime `db:"updated_date"`
	ExpiredDate      sql.NullTime `db:"expired_date"`
}

func (c *Contract) GetUserTransactionList(db *pgxpool.Pool, ctx context.Context, userCode string, param request.UserTransactionParam) ([]UserTransactionEnt, request.UserTransactionParam, error) {
	var (
		err        error
		list       []UserTransactionEnt
		paramQuery []interface{}
		totalData  int

		query = `SELECT users_transactions.id AS id,
			transaction_code,
			user_id,
			users_transactions.data_source,
			games.game_code AS game_code,
			games.name AS game_name,
			games.image_url AS game_img_url,
			users_transactions.price  AS final_price_amount,
			payment_method,
			users_transactions.status,
			users_transactions.created_date,
			users_transactions.updated_date,
			users_transactions.expired_date
    	FROM users_transactions
			JOIN users ON users.id = users_transactions.user_id
			LEFT JOIN rooms ON rooms.room_code  = users_transactions.source_code and users_transactions.data_source = 'room'
			LEFT JOIN tournaments ON tournaments.tournament_code  = users_transactions.source_code and users_transactions.data_source = 'tournament'
			JOIN games ON games.id = rooms.game_id OR games.id = tournaments.game_id
		WHERE users.user_code = $1`
	)

	// Populate Search
	paramQuery, query = generateTransactionFilterByQuery(param, query, userCode)

	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS total`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetUserTransactionList", err, utils.ErrCountingTotalUserTransaction)
		}
		param.Count = totalData
	}

	// Select Max Page
	if param.Count > param.Limit && param.Page > int(param.Count/param.Limit) {
		param.MaxPage = int(math.Ceil(float64(param.Count) / float64(param.Limit)))
	} else {
		param.MaxPage = int(param.Count / param.Limit)
	}

	// Limit and Offset
	param.Offset = (param.Page - 1) * param.Limit
	query += " ORDER BY users_transactions.id DESC "

	paramQuery = append(paramQuery, param.Offset)
	query += fmt.Sprintf(" OFFSET $%d ", len(paramQuery))

	paramQuery = append(paramQuery, param.Limit)
	query += fmt.Sprintf(" LIMIT $%d ", len(paramQuery))

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, param, c.errHandler("model.GetUserTransactionList", err, utils.ErrGettingListUserTransaction)
	}

	defer rows.Close()
	for rows.Next() {
		var data UserTransactionEnt
		err = rows.Scan(
			&data.Id,
			&data.TransactionCode,
			&data.UserId,
			&data.DataSource,
			&data.GameCode,
			&data.GameName,
			&data.GameImgUrl,
			&data.Price,
			&data.PaymentMethod,
			&data.Status,
			&data.CreatedDate,
			&data.UpdatedDate,
			&data.ExpiredDate,
		)

		if err != nil {
			return list, param, c.errHandler("model.GetUserTransactionList", err, utils.ErrScanningListUserTransaction)
		}

		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) GetTransactionByCode(db *pgxpool.Pool, ctx context.Context, userCode string, trxCode string) (UserTransactionEnt, error) {
	var (
		err   error
		data  UserTransactionEnt
		query = `SELECT users_transactions.id AS id,
			transaction_code,
			users_transactions.data_source,
			games.game_code AS game_code,
			games.name AS game_name,
			users_transactions.price  AS final_price_amount,
			COALESCE(users_points.point, 0) AS awarded_user_point,
			payment_method,
			payment_link,
			users_transactions.status,
			users_transactions.created_date,
			users_transactions.updated_date,
			users_transactions.expired_date
		FROM users_transactions
			JOIN users ON users.id = users_transactions.user_id
			LEFT JOIN rooms ON rooms.room_code  = users_transactions.source_code and users_transactions.data_source = 'room'
			LEFT JOIN tournaments ON tournaments.tournament_code  = users_transactions.source_code and users_transactions.data_source = 'tournament'
			JOIN games ON games.id = rooms.game_id OR games.id = tournaments.game_id
			LEFT JOIN users_points ON users_transactions.source_code = users_points.source_code AND users_points.user_id = users.id
    	WHERE users.user_code = $1 AND transaction_code = $2`
	)

	err = db.QueryRow(ctx, query, userCode, trxCode).Scan(
		&data.Id,
		&data.TransactionCode,
		&data.DataSource,
		&data.GameCode,
		&data.GameName,
		&data.Price,
		&data.AwardedUserPoint,
		&data.PaymentMethod,
		&data.PaymentLink,
		&data.Status,
		&data.CreatedDate,
		&data.UpdatedDate,
		&data.ExpiredDate,
	)

	if err != nil {
		return data, c.errHandler("model.GetTransactionByCode", err, utils.ErrGettingtUserTransactionDetailCode)
	}

	return data, nil
}

func (c *Contract) CreateOneTimeInvoice(tx pgx.Tx, ctx context.Context, userId int64, dataSource string, sourceCode string, price float64, titleSubs string, email string) (int64, string, string, time.Time, error) {
	var (
		err             error
		paymentLink     string
		orderId         int64
		lastInsertId    int64
		userFullname    string
		userPhoneNumber string
		xenditKey       = c.Config.GetString("xendit.api_key")
		currentTime     = time.Now().In(time.UTC)

		sqlInsert = `INSERT INTO users_transactions (
			user_id, data_source, source_code, transaction_code, aggregator_code,
			price, payment_method, payment_link, status, resp_payload, created_date, updated_date, expired_date
	   ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8 , $9, $10, $11, $12, $13) RETURNING id`

		sqlGetUserData = `SELECT fullname, phone_number FROM users WHERE id = $1;`
	)

	// Generate OrderCode
	orderCode := utils.GeneratePrefixCode(utils.TransactionPrefix)

	_ = tx.QueryRow(ctx, sqlGetUserData, userId).Scan(&userFullname, &userPhoneNumber)

	// Forward Xendit Payment
	resp, errX := payment.XenditClient{Key: xenditKey}.CreateInvoice(int64(price), orderCode, email, titleSubs, userFullname, userPhoneNumber)
	if errX != nil {
		fmt.Println("CreateOneTimeInvoice [Create]", errX)

		return orderId, orderCode, paymentLink, currentTime, c.errHandler("model.CreateOrder", errX, "Error create one time purchase Xendit")
	}

	marResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("CreateOneTimeInvoice [Marshall]", err)

		return orderId, orderCode, paymentLink, currentTime, err
	}

	// Remap the xendit response to order entity
	invoiceUrl := resp.GetInvoiceUrl()

	err = tx.QueryRow(ctx, sqlInsert,
		userId,
		dataSource,
		sourceCode,
		orderCode,
		*resp.Id,
		price,
		"",
		invoiceUrl,
		resp.Status.String(),
		string(marResp),
		currentTime,
		currentTime,
		resp.ExpiryDate,
	).Scan(&lastInsertId)

	if err != nil {
		return lastInsertId, orderCode, invoiceUrl, resp.ExpiryDate, c.errHandler("model.CreateOneTimeInvoice", err, utils.ErrCreatingOneInvoice)
	}

	return lastInsertId, orderCode, invoiceUrl, resp.ExpiryDate, nil
}

func (c *Contract) GetInvoiceTrxByCode(db *pgxpool.Pool, ctx context.Context, aggregatorCode string) (OriginUserTransactionEnt, error) {
	var (
		err  error
		data OriginUserTransactionEnt
		sql  = `
		SELECT 
			ut.id, ut.user_id, ut.data_source, ut.source_code, ut.transaction_code, ut.aggregator_code,
			ut.price, ut.payment_method, ut.payment_link, ut.status,
			ut.created_date, ut.updated_date,
			u.user_code, u.fullname AS user_fullname, u.email AS user_email
		FROM users_transactions ut
		JOIN users u ON ut.user_id = u.id
		WHERE ut.aggregator_code = $1
	`
	)
	err = db.QueryRow(ctx, sql, aggregatorCode).Scan(
		&data.Id, &data.UserId, &data.DataSource, &data.SourceCode,
		&data.TransactionCode, &data.AggregatorCode, &data.Price,
		&data.PaymentMethod, &data.PaymentLink, &data.Status,
		&data.CreatedDate, &data.UpdatedDate,
		&data.UserCode, &data.UserFullname, &data.UserEmail)

	if err != nil {
		return data, c.errHandler("model.GetInvoiceTrxByCode", err, utils.ErrGetInvoiceByAggregatorCode)
	}

	return data, nil
}

func (c *Contract) UpdateInvoiceTrx(tx pgx.Tx, ctx context.Context, aggregatorCode string, paymentMethod string, status string, response string) error {
	var (
		err error
	)

	query := `UPDATE users_transactions 
	SET payment_method = $1, status = $2, resp_payload = $3, updated_date = $4
	WHERE aggregator_code = $5`
	_, err = tx.Exec(ctx, query, paymentMethod, status, response, time.Now().In(time.UTC), aggregatorCode)
	if err != nil {
		return c.errHandler("model.UpdateInvoiceTrx", err, utils.ErrUpdateInvoiceByAggregatorCode)
	}

	return nil
}

func (c *Contract) GetTotalBookingAmountByUserID(db *pgxpool.Pool, ctx context.Context, userID int64) (int, error) {
	var (
		totalAmount    pgtype.Numeric
		totalAmountInt int
	)

	query := `SELECT COALESCE(SUM(price), 0) AS total_booking 
		FROM users_transactions 
		WHERE user_id = $1 AND status = 'PAID';`

	err := db.QueryRow(ctx, query, userID).Scan(&totalAmount)
	if err != nil {
		return 0, fmt.Errorf("error executing query: %w", err)
	}

	// Convert pgtype.Numeric to int
	err = totalAmount.AssignTo(&totalAmountInt)
	if err != nil {
		return 0, fmt.Errorf("error converting totalAmount to int64: %w", err)
	}

	return totalAmountInt, nil
}

// Private function
func generateTransactionFilterByQuery(param request.UserTransactionParam, query string, userCode string) ([]interface{}, string) {
	var (
		where      []string
		paramQuery []interface{}
	)

	paramQuery = append(paramQuery, userCode)

	// TRX STATUS
	if len(param.Status) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Status)
		orWhere = append(orWhere, fmt.Sprintf("users_transactions.status = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// Append All Where Conditions
	if len(where) > 0 {
		query += " AND " + strings.Join(where, " AND ")
	}

	return paramQuery, query
}
