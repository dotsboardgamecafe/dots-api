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

type UserGameCollectionResp struct {
	UserId       sql.NullInt64  `db:"user_id"`
	UserCode     sql.NullString `db:"user_code"`
	GameId       sql.NullInt64  `db:"game_id"`
	GameCode     sql.NullString `db:"game_code"`
	GameName     sql.NullString `db:"game_name"`
	GameImageUrl sql.NullString `db:"game_image_url"`
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
			g.image_url AS game_image_url
		FROM games g
		LEFT JOIN tournaments t ON g.id = t.game_id
		LEFT JOIN tournament_participants tp ON tp.tournament_id = t.id AND tp.status = 'active'
		left JOIN rooms r ON g.id = r.game_id
		left JOIN rooms_participants rp ON rp.room_id = r.id AND rp.status = 'active'
		LEFT JOIN users u ON rp.user_id = u.id OR tp.user_id = u.id
		WHERE u.user_code = $1
		GROUP BY
			u.id,
			u.user_code,
			g.id,
			g.game_code,
			g."name",
			g.image_url 
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
		err = rows.Scan(&data.UserId, &data.UserCode, &data.GameId, &data.GameCode, &data.GameName, &data.GameImageUrl)
		if err != nil {
			return list, param, c.errHandler("model.GetUserGameCollections", err, utils.ErrScanningListUserGameCollection)
		}
		list = append(list, data)
	}

	return list, param, nil
}
