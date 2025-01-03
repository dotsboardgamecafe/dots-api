package model

import (
	"context"
	"dots-api/lib/utils"
	"dots-api/services/api/request"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

type HallOfFameEnt struct {
	UserCode            string `db:"user_code"`
	UserFullName        string `db:"user_fullname"`
	UserName            string `db:"user_name"`
	UserImgUrl          string `db:"user_img_url"`
	TournamentBannerUrl string `db:"tournament_banner_url"`
	TournamentName      string `db:"tournament_name"`
	CafeName            string `db:"cafe_name"`
	CafeAddress         string `db:"cafe_address"`
}

type MonthlyTopAchieverEnt struct {
	Ranking         int    `db:"rank"`
	UserCode        string `db:"user_code"`
	UserFullName    string `db:"user_fullname"`
	UserName        string `db:"user_name"`
	UserImgUrl      string `db:"user_img_url"`
	CafeCode        string `db:"cafe_code"`
	Location        string `db:"location"`
	TotalPoint      int    `db:"total_point"`
	TotalGamePlayed int    `db:"total_game_played"`
}

func (c *Contract) GetHallOfFameList(db *pgxpool.Pool, ctx context.Context, param request.HallOfFameParam) ([]HallOfFameEnt, request.HallOfFameParam, error) {
	var (
		err        error
		list       []HallOfFameEnt
		paramQuery []interface{}

		query = `SELECT 
			u.user_code,
			COALESCE(u.username, '') AS user_name,
			u.fullname AS user_fullname,
			u.image_url AS user_img_url,
			g.image_url AS tournament_banner_url,
			t.name AS tournament_name,
			c.address AS cafe_address,
			c.name AS cafe_name
		FROM users u
			JOIN tournament_participants tp ON u.id = tp.user_id
			JOIN tournaments t ON t.id = tp.tournament_id
			JOIN games g ON t.game_id = g.id
			JOIN cafes c ON c.id = g.cafe_id
		WHERE tp.position = 1
		`
	)

	// Populate Search
	paramQuery, query = generateFilterQueryHallOfFame(param, query)

	// Limit
	paramQuery = append(paramQuery, param.Limit)
	query += fmt.Sprintf(" LIMIT $%d ", len(paramQuery))

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, param, c.errHandler("model.GetHallOfFameList", err, utils.ErrGetHallOfFame)
	}

	defer rows.Close()
	for rows.Next() {
		var data HallOfFameEnt
		err = rows.Scan(
			&data.UserCode,
			&data.UserName,
			&data.UserFullName,
			&data.UserImgUrl,
			&data.TournamentBannerUrl,
			&data.TournamentName,
			&data.CafeAddress,
			&data.CafeName,
		)

		if err != nil {
			return list, param, c.errHandler("model.GetHallOfFameList", err, utils.ErrScanHallOfFame)
		}

		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) GetUniqueGame(db *pgxpool.Pool, ctx context.Context, param request.MonthlyTopAchieverParam) ([]MonthlyTopAchieverEnt, request.MonthlyTopAchieverParam, error) {
	var (
		err        error
		list       []MonthlyTopAchieverEnt
		paramQuery []interface{}
		where      []string
		query      = `
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
			ROW_NUMBER() OVER(ORDER BY COUNT(1) DESC) AS rank,
			u.user_code,
			COALESCE(u.username, '') AS user_name,
			u.fullname AS user_fullname,
			u.image_url AS user_img_url,
			0 AS total_point,
			COUNT(DISTINCT up.game_id) AS total_played_game
		FROM
			unique_participants up
		LEFT JOIN
			users u ON u.id = up.user_id
		LEFT JOIN
			games g ON g.id = up.game_id 
		LEFT JOIN 
			cafes c ON c.id = g.cafe_id
		`
	)

	// MONTH-YEAR
	if param.Month > 0 && param.Year > 0 {
		paramQuery = append(paramQuery, fmt.Sprintf("%d-%02d", param.Year, param.Month))
		where = append(where, fmt.Sprintf("TO_CHAR(up.created_date, 'YYYY-MM') <= $%d", len(paramQuery)))
	}

	// CAFE CITY
	if len(param.CafeCity) > 0 {
		paramQuery = append(paramQuery, strings.ToLower(param.CafeCity))
		where = append(where, fmt.Sprintf("LOWER(c.city) = $%d", len(paramQuery)))
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	// Limit
	paramQuery = append(paramQuery, param.Limit)
	query += fmt.Sprintf(" GROUP BY u.user_code, u.username, u.fullname, u.image_url LIMIT $%d ", len(paramQuery))

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, param, c.errHandler("model.GetMonthlyTopAchiever", err, utils.ErrGetMonthlyTopAchiever)
	}

	defer rows.Close()
	for rows.Next() {
		var data MonthlyTopAchieverEnt
		err = rows.Scan(
			&data.Ranking,
			&data.UserCode,
			&data.UserName,
			&data.UserFullName,
			&data.UserImgUrl,
			&data.TotalPoint,
			&data.TotalGamePlayed,
		)

		if err != nil {
			return list, param, c.errHandler("model.GetMonthlyTopAchiever", err, utils.ErrScanGetMonthlyTopAchiever)
		}

		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) GetMostVP(db *pgxpool.Pool, ctx context.Context, param request.MonthlyTopAchieverParam) ([]MonthlyTopAchieverEnt, request.MonthlyTopAchieverParam, error) {
	var (
		err                    error
		list                   []MonthlyTopAchieverEnt
		paramQuery             []interface{}
		queryGetRoomTotalPoint = `
		SELECT up.user_id, SUM(up.point) AS total_point
		FROM users_points up 
			JOIN rooms_participants rp ON up.user_id = rp.user_id AND rp.status = 'active'
			JOIN rooms r ON r.id = rp.room_id
		WHERE r.room_code = up.source_code
		`
		queryGetTournamentTotalPoint = `
		SELECT up.user_id, SUM(up.point) AS total_point
		FROM users_points up 
			JOIN tournament_participants rp ON up.user_id = rp.user_id AND rp.status = 'active'
			JOIN tournaments r ON r.id = rp.tournament_id
		WHERE r.tournament_code = up.source_code
		`
		queryGetNonGameTotalPoint = `
		SELECT up.user_id, SUM(up.point) AS total_point
		FROM users_points up 
		WHERE up.data_source NOT IN ('room', 'tournament')
		`
		query = `
		SELECT
			ROW_NUMBER() OVER(ORDER BY SUM(tp.total_point) DESC) AS rank,
			u.user_code,
			COALESCE(u.username, '') AS user_name,
			u.fullname AS user_fullname,
			u.image_url AS user_img_url,
			SUM(tp.total_point) AS total_point,
			0 AS total_played_game
		FROM users u 
		`
	)

	// MONTH-YEAR
	if param.Month > 0 && param.Year > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Month)

		orWhere = append(orWhere, fmt.Sprintf("EXTRACT('month' FROM up.created_date) = $%d", len(paramQuery)))

		paramQuery = append(paramQuery, param.Year)
		orWhere = append(orWhere, fmt.Sprintf("EXTRACT('year' FROM up.created_date) = $%d", len(paramQuery)))

		queryGetRoomTotalPoint += " AND " + strings.Join(orWhere, " AND ")
		queryGetTournamentTotalPoint += " AND " + strings.Join(orWhere, " AND ")
		queryGetNonGameTotalPoint += " AND " + strings.Join(orWhere, " AND ")
	}

	// CAFE CITY
	if len(param.CafeCity) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, strings.ToLower(param.CafeCity))

		orWhere = append(orWhere, fmt.Sprintf("LOWER(r.location_city) = $%d", len(paramQuery)))

		queryGetRoomTotalPoint += " AND " + strings.Join(orWhere, " AND ")
		queryGetTournamentTotalPoint += " AND " + strings.Join(orWhere, " AND ")
	}

	query += ` JOIN (` + queryGetRoomTotalPoint + " GROUP BY up.user_id UNION ALL " + queryGetTournamentTotalPoint + " GROUP BY up.user_id UNION ALL " + queryGetNonGameTotalPoint + " GROUP BY up.user_id " +
		` ) AS tp ON u.id = tp.user_id GROUP BY u.id`

	// Limit
	paramQuery = append(paramQuery, param.Limit)
	query += fmt.Sprintf(" LIMIT $%d ", len(paramQuery))

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, param, c.errHandler("model.GetMonthlyTopAchiever", err, utils.ErrGetMonthlyTopAchiever)
	}

	defer rows.Close()
	for rows.Next() {
		var data MonthlyTopAchieverEnt
		err = rows.Scan(
			&data.Ranking,
			&data.UserCode,
			&data.UserName,
			&data.UserFullName,
			&data.UserImgUrl,
			&data.TotalPoint,
			&data.TotalGamePlayed,
		)

		if err != nil {
			return list, param, c.errHandler("model.GetMonthlyTopAchiever", err, utils.ErrScanGetMonthlyTopAchiever)
		}

		list = append(list, data)
	}

	return list, param, nil
}

// Private Function
func generateFilterQueryHallOfFame(param request.HallOfFameParam, query string) ([]interface{}, string) {
	var (
		where      []string
		paramQuery []interface{}
	)

	// YEAR
	if param.Year > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Year)
		orWhere = append(orWhere, fmt.Sprintf("EXTRACT('year' FROM t.created_date) = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// Append All Where Conditions
	if len(where) > 0 {
		query += " AND " + strings.Join(where, " AND ")
	}

	query += ` ORDER BY t.created_date DESC`

	return paramQuery, query
}
