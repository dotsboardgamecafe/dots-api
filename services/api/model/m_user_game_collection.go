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
		// where      []string

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
	query += " ORDER BY " + param.Order + " " + param.Sort + " "

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

func (c *Contract) AddUserGameCollections(db *pgxpool.Pool, ctx context.Context, payload request.UserGameCollectionAddPayload) error {
	var (
		err error
	)

	res, err := db.Exec(ctx, `
		INSERT INTO users_game_collections(user_id, game_id, created_date)
			SELECT
				(SELECT id FROM users WHERE user_code = $1)::int as user_id,
				(SELECT id FROM games WHERE game_code = $2)::int as game_id,
				CURRENT_TIMESTAMP as created_date
	`, payload.UserCode, payload.GameCode)

	if err != nil {
		return c.errHandler("model.AddUserGameCollection", err, utils.ErrAddingUserGameCollection)
	}

	if res.RowsAffected() == 0 {
		return c.errHandler("model.AddUserGameCollection", errors.New(utils.ErrAddingUserGameCollection), utils.ErrAddingUserGameCollection)
	}

	return nil
}

func (c *Contract) CheckUserGameCollectionExists(db *pgxpool.Pool, ctx context.Context, payload request.UserGameCollectionAddPayload) error {
	var (
		err        error
		res        bool
		query      string
		paramQuery []interface{}
	)

	paramQuery = append(paramQuery, payload.UserCode)
	paramQuery = append(paramQuery, payload.GameCode)

	query = `
		SELECT EXISTS(
			SELECT
				u.id,
				u.user_code,
				g.id,
				g.game_code AS game_code,
				g."name" AS game_name,
				g.image_url AS game_image_url
			FROM games g
			LEFT JOIN users_game_collections ugc ON ugc.game_id = g.id
			LEFT JOIN users u ON ugc.user_id = u.id
			WHERE u.user_code = $1 AND g.game_code = $2
			LIMIT 1
		) as exists
	`

	err = db.QueryRow(ctx, query, paramQuery...).Scan(&res)
	if err != nil {
		return c.errHandler("model.CheckUserGameCollection", err, utils.ErrUserGameCollectionExists)
	}

	if res {
		return c.errHandler("model.CheckUserGameCollection", errors.New(utils.ErrUserGameCollectionExists), utils.ErrUserGameCollectionExists)
	}

	return nil
}
