package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"dots-api/services/api/request"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserGameHistoryResp struct {
	UserId             sql.NullInt64  `db:"user_id"`
	UserCode           sql.NullString `db:"user_code"`
	GameId             sql.NullInt64  `db:"game_id"`
	GameName           sql.NullString `db:"game_name"`
	GameImageUrl       sql.NullString `db:"game_image_url"`
	GameDuration       sql.NullInt64  `db:"game_duration"`
	GameDifficulty     sql.NullString `db:"game_difficulty"`
	GameType           sql.NullString `db:"game_type"`
	GamePlayerSlot     sql.NullInt64  `db:"player_slot"`
	GameMasterId       sql.NullInt64  `db:"game_master_id"`
	GameMasterCode     sql.NullString `db:"game_master_code"`
	GameMasterName     sql.NullString `db:"game_master_name"`
	GameMasterImageUrl sql.NullString `db:"game_master_image_url"`
	GamePlayType       sql.NullString `db:"game_play_type"`
}

type UsersHavePlayedGameHistoryEnt struct {
	GameId      sql.NullInt64  `db:"game_id"`
	GameName    sql.NullString `db:"game_name"`
	UserCode    sql.NullString `db:"user_code"`
	UserName    sql.NullString `db:"username"`
	UserImage   sql.NullString `db:"user_image"`
	CreatedDate time.Time      `db:"created_date"`
}

func (c *Contract) GetUserGameHistories(db *pgxpool.Pool, ctx context.Context, userCode string, param request.UserGameHistoryParam) ([]UserGameHistoryResp, request.UserGameHistoryParam, error) {
	var (
		err        error
		list       []UserGameHistoryResp
		paramQuery []interface{}
		totalData  int
		// where      []string
		query = `select
					gdata.user_id,
					u.user_code,
					gdata.game_id,
					g."name" as game_name,
					g.image_url as game_image_url,
					g.duration as game_duration,
					g.difficulty as game_difficulty,
					g.game_type,
					gdata.player_slot,
					gdata.game_master_id,
					a.admin_code as game_master_code,
					a."name" as game_master_name,
					a.image_url as game_master_image_url,
					gdata.game_play_type
				from (
					select rp.user_id, r.game_id, r.game_master_id, r.maximum_participant as player_slot,  'non-tournament' as game_play_type
					from rooms  r
					left join rooms_participants rp on rp.room_id = r.id  
					union all
					select tp.user_id, t.game_id, null as game_master_id, t.player_slot, 'tournament' as game_play_type
					from tournaments  t 
					left join tournament_participants tp on tp.tournament_id = t.id 
				) as gdata
				left join users u on gdata.user_id = u.id
				left join games g on gdata.game_id = g.id
				left join admins a on gdata.game_master_id=a.id
				where u.user_code=$1 `
	)

	paramQuery = append(paramQuery, userCode)
	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS data`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetUserGameHistorys", err, utils.ErrCountingListUserGameHistory)
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
		return list, param, c.errHandler("model.GetUserGameHistories", err, utils.ErrGettingListUserGameHistory)
	}

	defer rows.Close()
	for rows.Next() {
		var data UserGameHistoryResp
		err = rows.Scan(&data.UserId, &data.UserCode, &data.GameId, &data.GameName, &data.GameImageUrl, &data.GameDuration, &data.GameDifficulty, &data.GameType, &data.GamePlayerSlot, &data.GameMasterId, &data.GameMasterCode, &data.GameMasterName, &data.GameMasterImageUrl, &data.GamePlayType)
		if err != nil {
			return list, param, c.errHandler("model.GetUserGameHistories", err, utils.ErrScanningListUserGameHistory)
		}
		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) CountDifferentGamesByUserID(db *pgxpool.Pool, ctx context.Context, userId int64) (int, error) {
	var (
		err       error
		gameCount int
	)

	query := `
		SELECT COUNT(DISTINCT gdata.game_id) as game_count
		FROM (
			SELECT rp.user_id, r.game_id
			FROM rooms r
			LEFT JOIN rooms_participants rp ON rp.room_id = r.id  WHERE rp.status = 'active'
			UNION ALL
			SELECT tp.user_id, t.game_id
			FROM tournaments t 
			LEFT JOIN tournament_participants tp ON tp.tournament_id = t.id where tp.status = 'active'  
		) AS gdata
		LEFT JOIN users u ON gdata.user_id = u.id
		WHERE u.id = $1 
	`

	err = db.QueryRow(ctx, query, userId).Scan(&gameCount)
	if err != nil {
		return 0, c.errHandler("model.CountDifferentGamesByUserID", err, utils.ErrCountingDifferentGames)
	}

	return gameCount, nil
}

func (c *Contract) GetUsersHavePlayedGameHistory(db *pgxpool.Pool, ctx context.Context, gameCode string) ([]UsersHavePlayedGameHistoryEnt, int64, error) {
	var (
		err       error
		list      []UsersHavePlayedGameHistoryEnt
		totalData int64
		query     = `WITH 
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
		SELECT DISTINCT
			g.id AS game_id, 
			g."name" AS game_name, 
			u.user_code, 
			u.username, 
			u.image_url AS user_image,
			up.	created_date
		FROM unique_participants up

		LEFT JOIN users as u ON up.user_id = u.id

		LEFT JOIN games as g ON up.game_id = g.id 
		WHERE g.game_code = $1
		`
	)

	// Count total records
	countQuery := `SELECT COUNT(*) FROM (` + query + `) AS data`
	err = db.QueryRow(ctx, countQuery, gameCode).Scan(&totalData)
	if err != nil {
		return list, totalData, c.errHandler("model.GetUserGameHistorys", err, utils.ErrCountingListUserGameHistory)
	}

	query += `ORDER BY up.created_date desc LIMIT 3`
	// Execute query
	rows, err := db.Query(ctx, query, gameCode)
	if err != nil {
		return list, totalData, c.errHandler("model.GetUserGameHistories", err, utils.ErrGettingListUserGameHistory)
	}
	defer rows.Close()
	// Process results
	for rows.Next() {
		var data UsersHavePlayedGameHistoryEnt
		err = rows.Scan(&data.GameId, &data.GameName, &data.UserCode, &data.UserName, &data.UserImage, &data.CreatedDate)
		if err != nil {
			return list, totalData, c.errHandler("model.GetUserGameHistories", err, utils.ErrScanningListUserGameHistory)
		}
		list = append(list, data)
	}

	return list, totalData, nil
}
