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

func (c *Contract) ListOfGameMechanics(db *pgxpool.Pool, ctx context.Context, param request.GameMechanicParam) ([]SettingEnt, error) {
	var (
		err        error
		list       []SettingEnt
		where      []string
		paramQuery []interface{}
		totalData  int

		query = `SELECT id, setting_code, content_value, created_date FROM settings WHERE set_group = 'game_mechanic'`
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
			return list, c.errHandler("model.ListOfGameMechanics", err, utils.ErrCountingListSetting)
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
		return list, c.errHandler("model.ListOfGameMechanics", err, utils.ErrGettingListSetting)
	}

	defer rows.Close()
	for rows.Next() {
		var data SettingEnt
		err = rows.Scan(&data.Id, &data.SettingCode, &data.ContentValue, &data.CreatedDate)
		if err != nil {
			return list, c.errHandler("model.ListOfGameMechanics", err, utils.ErrScanningListSetting)
		}
		list = append(list, data)
	}

	return list, nil
}

func (c *Contract) UpdateGameMechanic(db *pgxpool.Pool, ctx context.Context, settingCode string, mechanicKey string, mechanicName string, prevMechanicName string) error {
	var (
		err                   error
		sqlUpdateGameMechanic = `UPDATE settings 
			SET set_key = $1, content_value = $2, updated_date = $3
		WHERE setting_code = $4`
		sqlUpdateRelatedGameCategories = `UPDATE games_categories 
			SET category_name = $1
		WHERE category_name = $2`
	)

	fail := func(err error) error {
		return fmt.Errorf("UpdateGameMechanic: %v", err)
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

	_, err = tx.Exec(ctx, sqlUpdateGameMechanic,
		mechanicKey, mechanicName, time.Now().In(time.UTC), settingCode)
	if err != nil {
		return fail(err)
	}

	_, err = tx.Exec(ctx, sqlUpdateRelatedGameCategories, mechanicName, prevMechanicName)
	if err != nil {
		return fail(err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fail(err)
	}

	return nil
}

func (c *Contract) DeleteGameMechanic(db *pgxpool.Pool, ctx context.Context, settingCode string, mechanicName string) error {
	var (
		err                          error
		sqlDeleteGameMechanic        = "DELETE FROM settings WHERE setting_code = $1"
		sqlDeleteRelatedGameMechanic = "DELETE FROM games_categories WHERE category_name = $1"
	)

	fail := func(err error) error {
		return fmt.Errorf("DeleteGameMechanic: %v", err)
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

	_, err = tx.Exec(ctx, sqlDeleteGameMechanic, settingCode)
	if err != nil {
		return fail(err)
	}

	_, err = tx.Exec(ctx, sqlDeleteRelatedGameMechanic, mechanicName)
	if err != nil {
		return fail(err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fail(err)
	}

	return nil
}
