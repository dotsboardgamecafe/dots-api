package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"dots-api/services/api/request"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserGameCollectionResp struct {
	UserId       sql.NullInt64  `db:"user_id"`
	UserCode     sql.NullString `db:"user_code"`
	GameId       sql.NullInt64  `db:"game_id"`
	GameCode     sql.NullString `db:"game_code"`
	GameName     sql.NullString `db:"game_name"`
	GameImageUrl sql.NullString `db:"game_image_url"`
	CreatedDate  time.Time      `db:"created_date"`
}

func (c *Contract) GetUserGameCollections(db *pgxpool.Pool, ctx context.Context, userCode string, param request.UserGameCollectionParam) ([]UserGameCollectionResp, request.UserGameCollectionParam, error) {
	var (
		err        error
		list       []UserGameCollectionResp
		paramQuery []interface{}
		totalData  int

		query = `
			SELECT
				u.id,
				u.user_code,
				g.id,
				g.game_code AS game_code,
				g."name" AS game_name,
				g.image_url AS game_image_url,
				ugc.created_date
			FROM games g
			LEFT JOIN users_game_collections ugc ON ugc.game_id = g.id
			LEFT JOIN users u ON ugc.user_id = u.id
			WHERE u.user_code = $1
		`
	)

	paramQuery = append(paramQuery, userCode)

	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS data`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetUserGameCollections", err, utils.ErrCountingListUserGameCollection)
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
	query += " ORDER BY " + param.SortKey + " " + param.Sort + " "

	paramQuery = append(paramQuery, param.Offset)
	query += fmt.Sprintf("offset $%d ", len(paramQuery))

	paramQuery = append(paramQuery, param.Limit)
	query += fmt.Sprintf("limit $%d ", len(paramQuery))

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, param, c.errHandler("model.GetUserGameCollections", err, utils.ErrGettingListUserGameCollection)
	}

	defer rows.Close()
	for rows.Next() {
		var data UserGameCollectionResp
		err = rows.Scan(&data.UserId, &data.UserCode, &data.GameId, &data.GameCode, &data.GameName, &data.GameImageUrl, &data.CreatedDate)
		if err != nil {
			return list, param, c.errHandler("model.GetUserGameCollections", err, utils.ErrScanningListUserGameCollection)
		}
		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) CountUserGameCollectionsByUserID(db *pgxpool.Pool, ctx context.Context, userId int64) (int, error) {
	var (
		err   error
		count int
	)

	err = db.QueryRow(ctx, `SELECT COUNT(*) FROM users_game_collections WHERE user_id = $1`, userId).Scan(&count)
	if err != nil {
		return 0, c.errHandler("model.CountUserGameCollectionsByUserID", err, utils.ErrCountingUserGameCollection)
	}

	return count, nil
}

func (c *Contract) AddUserGameCollections(db *pgxpool.Pool, ctx context.Context, userId, gameId int64) error {
	var (
		err   error
		query = `
			INSERT INTO 
				users_game_collections(user_id, game_id, created_date)
				VALUES ($1, $2, CURRENT_TIMESTAMP)
		`
	)

	_, err = db.Exec(ctx, query, userId, gameId)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return c.errHandler("model.AddUserGameCollection", errors.New(utils.ErrUserGameCollectionExists), utils.ErrUserGameCollectionExists)
			}
		}

		return c.errHandler("model.AddUserGameCollection", err, utils.ErrAddingUserGameCollection)
	}

	return nil
}
