package model

import (
	"context"
	"dots-api/lib/utils"
	"dots-api/services/api/request"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func (c *Contract) ListOfGameTypes(db *pgxpool.Pool, ctx context.Context, param request.GameMechanicParam) ([]SettingEnt, error) {
	var (
		err        error
		list       []SettingEnt
		where      []string
		paramQuery []interface{}
		totalData  int

		query = `SELECT id, setting_code, content_value, created_date FROM settings WHERE set_group = 'game_type'`
	)

	// Populate Search
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		orWhere = append(orWhere, fmt.Sprintf("set_label iLIKE $%d", len(paramQuery)))
		orWhere = append(orWhere, fmt.Sprintf("content_value iLIKE $%d", len(paramQuery)))
		where = append(where, "("+strings.Join(orWhere, " OR ")+")")
	}

	// Append All Where Conditions
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	{
		newQcount := `SELECT COUNT(1) FROM ( ` + query + ` ) AS data`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, c.errHandler("model.ListOfGameTypes", err, utils.ErrCountingListSetting)
		}
		param.Count = totalData
	}

	// Limit and Offset
	if param.Page > 0 && param.Limit > 0 {
		// Select Max Page
		if param.Count > param.Limit && param.Page > int(param.Count/param.Limit) {
			param.Page = int(math.Ceil(float64(param.Count) / float64(param.Limit)))
		}

		param.Offset = (param.Page - 1) * param.Limit
		query += " ORDER BY id " + param.Sort + " "

		paramQuery = append(paramQuery, param.Offset)
		query += fmt.Sprintf("offset $%d ", len(paramQuery))

		paramQuery = append(paramQuery, param.Limit)
		query += fmt.Sprintf("limit $%d ", len(paramQuery))
	}

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, c.errHandler("model.ListOfGameTypes", err, utils.ErrGettingListSetting)
	}

	defer rows.Close()
	for rows.Next() {
		var data SettingEnt
		err = rows.Scan(&data.Id, &data.SettingCode, &data.ContentValue, &data.CreatedDate)
		if err != nil {
			return list, c.errHandler("model.ListOfGameTypes", err, utils.ErrScanningListSetting)
		}
		list = append(list, data)
	}

	return list, nil
}

func (c *Contract) UpdateGameType(db *pgxpool.Pool, ctx context.Context, settingCode string, typeKey string, typeName string, prevTypeName string) error {
	var (
		err               error
		sqlUpdateGameType = `UPDATE settings 
			SET set_key = $1, content_value = $2, updated_date = $3
		WHERE setting_code = $4`
		sqlUpdateRelatedGameCategories = `UPDATE games 
			SET game_type = $1
		WHERE game_type = $2`
	)

	fail := func(err error) error {
		return fmt.Errorf("UpdateGameType: %v", err)
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

	_, err = tx.Exec(ctx, sqlUpdateGameType,
		typeKey, typeName, time.Now().In(time.UTC), settingCode)
	if err != nil {
		return fail(err)
	}

	_, err = tx.Exec(ctx, sqlUpdateRelatedGameCategories, typeName, prevTypeName)
	if err != nil {
		return fail(err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fail(err)
	}

	return nil
}

func (c *Contract) DeleteGameType(db *pgxpool.Pool, ctx context.Context, settingCode string, mechanicName string) error {
	var (
		err   error
		query = "DELETE FROM settings WHERE setting_code = $1"
	)

	_, err = db.Exec(ctx, query, settingCode)
	if err != nil {
		return c.errHandler("model.DeleteGameType", err, utils.ErrGettingListSetting)
	}

	return nil
}

func (c *Contract) IsGameTypeExistsOnGames(db *pgxpool.Pool, ctx context.Context, gameType string) (int, error) {
	var (
		err       error
		totalData int
		query     = "SELECT COUNT(1) AS total FROM games WHERE game_type = $1"
	)

	err = db.QueryRow(ctx, query, gameType).Scan(&totalData)
	if err != nil {
		return 0, c.errHandler("model.IsGameTypeExistsOnGames", err, utils.ErrGettingListSetting)
	}

	return totalData, nil
}
