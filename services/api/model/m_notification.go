package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"dots-api/services/api/request"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type NotificationEnt struct {
	Id               int64          `db:"id"`
	ReceiverSource   string         `db:"receiver_source"`
	ReceiverCode     string         `db:"receiver_code"`
	NotificationCode string         `db:"notification_code"`
	TransactionCode  string         `db:"transaction_code"`
	Type             string         `db:"type"`
	Title            sql.NullString `db:"title"`
	Description      sql.NullString `db:"description"`
	StatusRead       bool           `db:"status_read"`
	ImageUrl         sql.NullString `db:"image_url"`
	CreatedDate      time.Time      `db:"created_date"`
	UpdatedDate      sql.NullTime   `db:"updated_date"`
}

// AddNotification adds a new notification to the database
func (c *Contract) AddNotification(db *pgxpool.Pool, ctx context.Context, notifCode, receiverSource, receiverCode, transactionCode, nType, title string, description interface{}, imageUrl string) error {
	sql := `
		INSERT INTO notifications(receiver_source, receiver_code, transaction_code, notification_code, "type", "title", "description", status_read, image_url, created_date)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := db.Exec(ctx, sql,
		receiverSource, receiverCode, transactionCode, notifCode, nType, title, description, false, imageUrl, time.Now().In(time.UTC),
	)
	if err != nil {
		return c.errHandler("model.AddNotification", err, utils.ErrAddingNotification)
	}
	return nil
}

// AddNotification With Transactional
func (c *Contract) AddNotificationWithTx(tx pgx.Tx, ctx context.Context, notifCode, receiverSource, receiverCode, transactionCode, nType, title string, description interface{}, imageUrl string) error {
	sql := `
		INSERT INTO notifications(receiver_source, receiver_code, transaction_code, notification_code, "type", "title", "description", status_read, image_url, created_date)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := tx.Exec(ctx, sql,
		receiverSource, receiverCode, transactionCode, notifCode, nType, title, description, false, imageUrl, time.Now().In(time.UTC),
	)
	if err != nil {
		return c.errHandler("model.AddNotificationWithTx", err, utils.ErrAddingNotification)
	}
	return nil
}

// UpdateNotificationIsSeenByNotificationCode
func (c *Contract) UpdateNotificationIsSeenByNotificationCode(db *pgxpool.Pool, ctx context.Context, code string, isSeen bool) error {
	sql := `
		UPDATE notifications
		SET status_read = $1
		WHERE notification_code = $2
	`
	_, err := db.Exec(ctx, sql, isSeen, code)
	if err != nil {
		return c.errHandler("model.UpdateNotificationIsSeenByNotificationCode", err, utils.ErrUpdatingNotification)
	}
	return nil
}

// GetListNotifications retrieves a list of notifications based on given parameters with pagination
func (c *Contract) GetListNotifications(db *pgxpool.Pool, ctx context.Context, param request.NotificationParam, code string) ([]NotificationEnt, request.NotificationParam, error) {
	var (
		err                        error
		list                       []NotificationEnt
		where                      []string
		paramArgs                  []interface{}
		totalData, totalUnreadData int
	)

	query := `SELECT 
		id, receiver_source, receiver_code, notification_code, transaction_code, "type", "title", "description", status_read, image_url, created_date, updated_date
	FROM notifications WHERE receiver_code = '` + code + "' AND ( DATE(created_date) BETWEEN CURRENT_DATE - INTERVAL '1 month' AND CURRENT_DATE ) "

	// Populate Search
	if len(param.ReceiverSource) > 0 {
		paramArgs = append(paramArgs, param.ReceiverSource)
		where = append(where, "receiver_source = $"+fmt.Sprint(len(paramArgs)))
	}
	if len(param.ReceiverCode) > 0 {
		paramArgs = append(paramArgs, param.ReceiverCode)
		where = append(where, "receiver_code = $"+fmt.Sprint(len(paramArgs)))
	}

	// Append All Where Conditions
	if len(where) > 0 {
		query += " AND" + strings.Join(where, " AND ")
	}
	// Count total data notification unread
	countUnreadQuery := `SELECT COUNT(*) FROM ( ` + query + `AND status_read = false ) AS data`
	err = db.QueryRow(ctx, countUnreadQuery, paramArgs...).Scan(&totalUnreadData)
	if err != nil {
		return list, param, c.errHandler("model.GetListNotifications", err, utils.ErrCountingListNotification)
	}

	// Count total data
	countQuery := `SELECT COUNT(*) FROM ( ` + query + ` ) AS data`
	err = db.QueryRow(ctx, countQuery, paramArgs...).Scan(&totalData)
	if err != nil {
		return list, param, c.errHandler("model.GetListNotifications", err, utils.ErrCountingListNotification)
	}

	// Update Count Param
	param.Count = totalData
	param.CountUnread = totalUnreadData

	// Select Max Page
	if param.Count > param.Limit && param.Page > int(param.Count/param.Limit) {
		param.Page = int(math.Ceil(float64(param.Count) / float64(param.Limit)))
	}

	// Limit and Offset
	param.Offset = (param.Page - 1) * param.Limit
	query += " ORDER BY " + param.Order + " " + param.Sort + " "
	query += fmt.Sprintf("LIMIT $%d OFFSET $%d ", len(paramArgs)+1, len(paramArgs)+2)
	paramArgs = append(paramArgs, param.Limit, param.Offset)

	rows, err := db.Query(ctx, query, paramArgs...)
	if err != nil {
		return list, param, c.errHandler("model.GetListNotifications", err, utils.ErrGettingListNotification)
	}
	defer rows.Close()

	for rows.Next() {
		var data NotificationEnt
		err = rows.Scan(&data.Id, &data.ReceiverSource, &data.ReceiverCode, &data.NotificationCode, &data.TransactionCode, &data.Type, &data.Title, &data.Description, &data.StatusRead, &data.ImageUrl, &data.CreatedDate, &data.UpdatedDate)
		if err != nil {
			return list, param, c.errHandler("model.GetListNotifications", err, utils.ErrScanningListNotification)
		}
		list = append(list, data)
	}

	return list, param, nil
}
