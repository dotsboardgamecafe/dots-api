package model

import (
	"context"
	"dots-api/lib/utils"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	GameAdminEnt struct {
		GameId  int64 `json:"game_id"`
		AdminId int64 `json:"admin_id"`
	}
)

func (c *Contract) GetGameAdmins(db *pgxpool.Pool, ctx context.Context, gameId int64) ([]AdminEnt, error) {
	var (
		err  error
		data []AdminEnt
		sql  = `SELECT 
		a.admin_code, a.email, a.name, a.username, a.status, a.image_url, a.phone_number, r.name as role
		FROM admins a 
		left join games_admins ga on ga.admin_id = a.id
		left join roles r on r.id = a.role_id
		WHERE ga.game_id = $1`
	)

	rows, err := db.Query(ctx, sql, gameId)
	if err != nil {
		return data, c.errHandler("model.GetGameAdmins", err, utils.ErrGettingGameAdmins)
	}

	for rows.Next() {
		var admin AdminEnt
		if err := rows.Scan(&admin.AdminCode, &admin.Email, &admin.Name, &admin.UserName, &admin.Status, &admin.ImageURL, &admin.PhoneNumber, &admin.Role); err != nil {
			return data, c.errHandler("model.GetGameAdmins", err, utils.ErrScanningGameAdmins)
		}
		data = append(data, admin)
	}

	return data, nil
}

func (c *Contract) SyncGameAdmins(tx pgx.Tx, ctx context.Context, gameId int64, adminIds []int64) error {
	var (
		err    error
		query  = `DELETE FROM games_admins WHERE game_id = $1`
		query1 = `INSERT INTO games_admins(game_id, admin_id) SELECT $1, id FROM admins WHERE id = ANY($2)`
	)

	_, err = tx.Exec(ctx, query, gameId)
	if err != nil {
		return c.errHandler("model.SyncGameAdmins", err, utils.ErrSyncingGameAdmins)
	}

	if len(adminIds) == 0 {
		return nil
	}

	_, err = tx.Exec(ctx, query1, gameId, adminIds)
	if err != nil {
		return c.errHandler("model.SyncGameAdmins", err, utils.ErrSyncingGameAdmins)
	}

	return nil
}
