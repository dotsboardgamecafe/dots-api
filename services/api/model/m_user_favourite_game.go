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

		query = `
			WITH 
			unique_participants AS (
					SELECT DISTINCT r.game_id, rp.user_id, rp.created_date
					FROM rooms r 
					INNER JOIN rooms_participants rp ON r.id = rp.room_id AND rp.status = 'active'
					UNION
					SELECT DISTINCT t.game_id, tp.user_id, tp.created_date
					FROM tournaments t
					INNER JOIN tournament_participants tp ON t.id = tp.tournament_id AND tp.status = 'active'
					UNION
					SELECT DISTINCT game_id, user_id, created_date FROM users_game_collections
			)
			SELECT
					u.id AS user_id,
					u.user_code,
					gc.category_name AS game_category_name,
					gc.category_description,
					gc.category_image_url,
					COUNT(up.game_id) AS total_play
				FROM unique_participants up
				LEFT JOIN games g ON g.id = up.game_id
				LEFT JOIN games_categories gc ON gc.game_id = g.id
				LEFT JOIN users u ON u.id = up.user_id
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
