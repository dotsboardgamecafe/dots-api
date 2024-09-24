package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"time"

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
		INSERT INTO notifications(receiver_source, receiver_code, transaction_code, notification_code, "type", "title", "description", status_read, created_date)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := db.Exec(ctx, sql,
		receiverSource, receiverCode, transactionCode, notifCode, nType, title, description, false, time.Now().In(time.UTC),
	)
	if err != nil {
		return c.errHandler("model.AddNotification", err, utils.ErrAddingNotification)
	}
	return nil
}
