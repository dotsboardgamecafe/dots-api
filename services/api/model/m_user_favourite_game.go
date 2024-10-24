package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"dots-api/services/api/request"
	"fmt"
	"math"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserFavouriteGameResp struct {
	UserId                  sql.NullInt64  `db:"user_id"`
	UserCode                sql.NullString `db:"user_code"`
	GameCategoryName        sql.NullString `db:"game_category_name"`
	GameCategoryDescription sql.NullString `db:"game_category_description"`
	GameCategoryImageUrl    sql.NullString `db:"game_category_image_url"`
	TotalPlay               sql.NullInt64  `db:"total_play"`
}

func (c *Contract) GetUserFavouriteGames(db *pgxpool.Pool, ctx context.Context, userCode string, param request.UserFavouriteGameParam) ([]UserFavouriteGameResp, request.UserFavouriteGameParam, error) {
	var (
		err        error
		list       []UserFavouriteGameResp
		paramQuery []interface{}
		totalData  int
		// where      []string

		query = `SELECT
					u.id AS user_id,
					u.user_code,
					gc.category_name AS game_category_name,
					gc.category_description,
					gc.category_image_url,
					COUNT(gc.id) AS total_play
				FROM games g
				JOIN games_categories gc ON gc.game_id = g.id
				LEFT JOIN rooms r ON r.game_id = g.id
				LEFT JOIN rooms_participants rp ON rp.room_id = r.id AND rp.status = 'active'
				LEFT JOIN tournaments t ON t.game_id = g.id
				LEFT JOIN tournament_participants tp ON tp.tournament_id = t.id AND tp.status = 'active'
				LEFT JOIN users u ON rp.user_id = u.id OR tp.user_id = u.id
				WHERE u.user_code = $1
				GROUP BY
					u.id,
					u.user_code,
					gc.category_name,
					gc.category_description,
					gc.category_image_url
				ORDER BY total_play DESC `
	)

	paramQuery = append(paramQuery, userCode)

	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS data`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetUserFavouriteGames", err, utils.ErrCountingListUserFavouriteGame)
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

	paramQuery = append(paramQuery, param.Offset)
	query += fmt.Sprintf("offset $%d ", len(paramQuery))

	paramQuery = append(paramQuery, param.Limit)
	query += fmt.Sprintf("limit $%d ", len(paramQuery))

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, param, c.errHandler("model.GetUserFavouriteGames", err, utils.ErrGettingListUserFavouriteGame)
	}

	defer rows.Close()
	for rows.Next() {
		var data UserFavouriteGameResp
		err = rows.Scan(&data.UserId, &data.UserCode, &data.GameCategoryName, &data.GameCategoryDescription, &data.GameCategoryImageUrl, &data.TotalPlay)
		if err != nil {
			return list, param, c.errHandler("model.GetUserFavouriteGames", err, utils.ErrScanningListUserFavouriteGame)
		}
		list = append(list, data)
	}

	return list, param, nil
}
