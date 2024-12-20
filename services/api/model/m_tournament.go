package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"dots-api/services/api/request"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	TournamentsEnt struct {
		TournamentId    int64           `db:"tournament_id"`
		GameId          int64           `db:"game_id"`
		GameCode        string          `db:"game_code"`
		GameName        string          `db:"game_name"`
		GameType        string          `db:"game_type"`
		GameImgUrl      string          `db:"game_img_url"`
		CafeCode        string          `db:"cafe_code"`
		CafeName        string          `db:"cafe_name"`
		CafeAddress     string          `db:"cafe_address"`
		TournamentCode  string          `db:"tournament_code"`
		ImageUrl        sql.NullString  `db:"image_url"`
		PrizesImgUrl    sql.NullString  `db:"prizes_img_url"`
		Name            sql.NullString  `db:"name"`
		TournamentRules string          `db:"tournament_rules"`
		Level           string          `db:"level"`
		StartDate       sql.NullTime    `db:"start_date"`
		EndDate         sql.NullTime    `db:"end_date"`
		StartTime       time.Time       `db:"start_time"`
		EndTime         time.Time       `db:"end_time"`
		ParticipantVP   int64           `db:"participant_vp"`
		BookingPrice    float64         `db:"booking_price"`
		PlayerSlot      int64           `db:"player_slot"`
		CurrentUsedSlot int64           `db:"current_used_slot"`
		Status          string          `db:"status"`
		DayPastEndDate  sql.NullFloat64 `db:"days_past_end_date"`
		CreatedDate     time.Time       `db:"created_date"`
		UpdatedDate     sql.NullTime    `db:"updated_date"`
		DeletedDate     sql.NullTime    `db:"deleted_date"`
		LocationCity    string          `db:"location_city"`
	}

	NonWinnerEntity struct {
		UserId           int64  `db:"user_id"`
		TournamentCode   string `db:"tournament_code"`
		ParticipantPoint int    `db:"participant_vp"`
	}
)

func (c *Contract) GetTournamentList(db *pgxpool.Pool, ctx context.Context, param request.TournamentParam) ([]TournamentsEnt, request.TournamentParam, error) {
	var (
		err        error
		list       []TournamentsEnt
		paramQuery []interface{}
		totalData  int
		where      []string

		query = `SELECT 
		games.game_code, games.game_type, games.name, games.image_url, 
		cafes.cafe_code, cafes.name as cafe_name, cafes.address as cafe_address, 
		tournaments.id as tournament_id, tournaments.tournament_code, tournaments.image_url, tournaments.name, tournaments.prizes_img_url, 
		tournaments.tournament_rules, tournaments.level, tournaments.start_date, tournaments.end_date, tournaments.booking_price, 
		tournaments.start_time, tournaments.end_time, tournaments.booking_price, 
		tournaments.player_slot, tournaments.participant_vp, 
		COALESCE(tp.count_participants, 0) AS current_used_slot,
		tournaments.status, DATE_PART('day', NOW() - tournaments.end_date) AS days_past_end_date, 
		tournaments.created_date, tournaments.updated_date, tournaments.deleted_date
		FROM tournaments
		LEFT JOIN games ON games.id = tournaments.game_id  
		LEFT JOIN cafes ON cafes.id = games.cafe_id  
		LEFT JOIN (
			SELECT 
				COUNT(tournament_participants.user_id) AS count_participants, 
				tournament_participants.tournament_id
				FROM tournament_participants
				WHERE tournament_participants.status != 'cancel'
				GROUP BY tournament_participants.tournament_id
		) AS tp 
		ON tp.tournament_id = tournaments.id
	`
	)

	// TOURNAMENT NAME (KEYWORD)
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		orWhere = append(orWhere, fmt.Sprintf("tournaments.name iLIKE $%d", len(paramQuery)))
		where = append(where, "("+strings.Join(orWhere, " AND ")+")")
	}

	// CAFE CITY
	if len(param.CafeCity) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, strings.ToLower(param.CafeCity))
		orWhere = append(orWhere, fmt.Sprintf("LOWER(tournaments.location_city) = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// GAME CODE
	if len(param.GameCode) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.GameCode)
		orWhere = append(orWhere, fmt.Sprintf("games.game_code = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// TOURNAMENT STATUS
	if len(param.Status) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Status)
		orWhere = append(orWhere, fmt.Sprintf("tournaments.status = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// TOURNAMENT DATE
	if len(param.TournamentDate) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.TournamentDate)
		orWhere = append(orWhere, fmt.Sprintf("$%d BETWEEN tournaments.start_date AND tournaments.end_date", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// Append All Where Conditions
	if len(where) > 0 {
		query += " WHERE ( tournaments.deleted_date IS NULL AND games.status = 'active' AND games.deleted_date IS NULL ) and ( cafes.status = 'active' AND cafes.deleted_date IS NULL ) AND (" + strings.Join(where, " AND ") + ")"
	} else {
		query += " WHERE ( tournaments.deleted_date IS NULL AND games.status = 'active' AND games.deleted_date IS NULL ) and ( cafes.status = 'active' AND cafes.deleted_date IS NULL ) "
	}

	{
		newQcount := `SELECT COUNT(1) FROM ( ` + query + ` ) AS total`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetTournamentList", err, utils.ErrCountingListTournament)
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

	// Build ORDER BY clause
	orderByClauses := []string{}
	for i, order := range param.Order {
		sortOrder := "asc"
		if i < len(param.Sort) {
			sortOrder = param.Sort[i]
		}
		orderByClauses = append(orderByClauses, fmt.Sprintf("%s %s", order, sortOrder))
	}
	query += " ORDER BY " + strings.Join(orderByClauses, ", ")

	paramQuery = append(paramQuery, param.Offset)
	query += fmt.Sprintf(" OFFSET $%d ", len(paramQuery))

	paramQuery = append(paramQuery, param.Limit)
	query += fmt.Sprintf(" LIMIT $%d ", len(paramQuery))

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, param, c.errHandler("model.GetTournamentList", err, utils.ErrGettingListTournament)
	}

	defer rows.Close()
	for rows.Next() {
		var data TournamentsEnt
		err = rows.Scan(
			&data.GameCode, &data.GameType, &data.GameName, &data.ImageUrl,
			&data.CafeCode, &data.CafeName, &data.CafeAddress,
			&data.TournamentId, &data.TournamentCode, &data.ImageUrl, &data.Name, &data.PrizesImgUrl,
			&data.TournamentRules, &data.Level, &data.StartDate, &data.EndDate, &data.BookingPrice,
			&data.StartTime, &data.EndTime, &data.BookingPrice,
			&data.PlayerSlot, &data.ParticipantVP, &data.CurrentUsedSlot,
			&data.Status, &data.DayPastEndDate,
			&data.CreatedDate, &data.UpdatedDate, &data.DeletedDate,
		)
		if err != nil {

			return list, param, c.errHandler("model.GetTournamentList", err, utils.ErrScanningListTournament)
		}
		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) GetTournamentByCode(db *pgxpool.Pool, ctx context.Context, code string) (TournamentsEnt, error) {
	var (
		err   error
		data  TournamentsEnt
		query = `
		SELECT 
		tournaments.game_id, games.game_code, games.game_type, games.name, games.image_url, 
		cafes.cafe_code, cafes.name as cafe_name, cafes.address as cafe_address, 
		tournaments.id as tournament_id,tournaments.tournament_code, tournaments.image_url, tournaments.name, tournaments.prizes_img_url, 
		tournaments.tournament_rules, tournaments.level, tournaments.start_date, tournaments.end_date, tournaments.booking_price, 
		tournaments.start_time, tournaments.end_time, tournaments.booking_price, 
		tournaments.player_slot, tournaments.participant_vp, tournaments.status,  DATE_PART('day', NOW() - tournaments.end_date) AS days_past_end_date, 
		tournaments.created_date, tournaments.updated_date, tournaments.deleted_date, 
		COALESCE(tp.count_participants, 0) AS current_used_slot
		FROM tournaments
		LEFT JOIN games ON games.id = tournaments.game_id AND games.status = 'active' AND games.deleted_date IS NULL 
		LEFT JOIN cafes ON cafes.id = games.cafe_id AND cafes.status = 'active' AND cafes.deleted_date IS NULL 
		LEFT JOIN (
			SELECT 
				COUNT(tournament_participants.user_id) AS count_participants, 
				tournament_participants.tournament_id
				FROM tournament_participants
				WHERE tournament_participants.status != 'cancel'
				GROUP BY tournament_participants.tournament_id
		) AS tp 
			ON tp.tournament_id = tournaments.id
    	WHERE tournaments.tournament_code = $1 AND tournaments.deleted_date IS NULL`
	)

	err = db.QueryRow(ctx, query, code).Scan(
		&data.GameId, &data.GameCode, &data.GameType, &data.GameName, &data.ImageUrl,
		&data.CafeCode, &data.CafeName, &data.CafeAddress,
		&data.TournamentId, &data.TournamentCode, &data.ImageUrl, &data.Name, &data.PrizesImgUrl,
		&data.TournamentRules, &data.Level, &data.StartDate, &data.EndDate, &data.BookingPrice,
		&data.StartTime, &data.EndTime, &data.BookingPrice,
		&data.PlayerSlot, &data.ParticipantVP, &data.Status, &data.DayPastEndDate,
		&data.CreatedDate, &data.UpdatedDate, &data.DeletedDate, &data.CurrentUsedSlot,
	)

	if err != nil {
		return data, c.errHandler("model.GetTournamentByCode", err, utils.ErrGettingTournamentByCode)
	}

	return data, nil
}

func (c *Contract) AddTournament(tx pgx.Tx, ctx context.Context, gameId int64, tournamentCode, imageUrl, name, tournamentRules, level, prizeImageUrl string, bookingPrice float64, startDate, endDate, startTime, endTime interface{}, playerSlot, participantVP int64, status string, locationCity string) (int64, error) {
	sql := `INSERT INTO tournaments(
		game_id, tournament_code, image_url, prizes_img_url, name, tournament_rules, level, start_date, end_date, start_time, end_time, player_slot, booking_price, participant_vp, status, created_date, location_city
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	RETURNING id`

	var id int64
	err := tx.QueryRow(ctx, sql, gameId, tournamentCode, imageUrl, prizeImageUrl, name, tournamentRules, level, startDate, endDate, startTime, endTime, playerSlot, bookingPrice, participantVP, status, time.Now().In(time.UTC), locationCity).Scan(&id)
	if err != nil {
		return 0, c.errHandler("model.AddTournament", err, utils.ErrAddingTournament)
	}

	return id, nil
}

func (c *Contract) UpdateTournamentByCode(tx pgx.Tx, ctx context.Context, gameId int64, code, imageUrl, name, tournament_rules, level, status, prizeImageUrl string, bookingPrice float64, startDate, endDate, startTime, endTime interface{}, playerSlot, participantVP int64, locationCity string) error {
	var (
		err error
		sql = `
		UPDATE tournaments 
		SET game_id = $1, name = $2, image_url = $3, tournament_rules = $4, level = $5, status = $6, prizes_img_url = $7, start_date = $8, end_date = $9, start_time = $10, end_time = $11, player_slot = $12, booking_price = $13, participant_vp = $14, updated_date = $15, location_city = $16
		WHERE tournament_code = $17`
	)
	_, err = tx.Exec(ctx, sql, gameId, name, imageUrl, tournament_rules, level, status, prizeImageUrl, startDate, endDate, startTime, endTime, playerSlot, bookingPrice, participantVP, time.Now().UTC(), locationCity, code)
	if err != nil {
		return c.errHandler("model.UpdateTournament", err, utils.ErrUpdatingTournament)
	}

	return nil
}

func (c *Contract) DeleteTournamentByCode(db *pgxpool.Pool, ctx context.Context, code string) error {
	var (
		err error
		sql = `
		UPDATE tournaments 
		SET updated_date = NOW(), deleted_date=$1 
		WHERE tournament_code=$2`
	)
	_, err = db.Exec(ctx, sql, time.Now().In(time.UTC), code)
	if err != nil {
		return c.errHandler("model.DeleteTournamentByCode", err, utils.ErrDeletingTournament)
	}

	return nil
}

func (c *Contract) CountParticipantTournamentByTournamentId(db *pgxpool.Pool, ctx context.Context, id int64) (int, error) {
	var (
		err              error
		totalParticipant int

		queryGetTotalParticipant = `select count(id)  from tournament_participants tp 
		where tp.status!= 'cancel' and tournament_id  = $1
		group by tournament_id `
	)

	err = db.QueryRow(ctx, queryGetTotalParticipant, id).Scan(
		&totalParticipant,
	)
	if err != nil && err != pgx.ErrNoRows {
		return totalParticipant, c.errHandler("model.CountParticipantTournamentByTournamentId", err, utils.ErrCountParticipantTournamentByTournamentId)
	}

	return totalParticipant, nil
}

func (c *Contract) GetTournamentIdByCode(db *pgxpool.Pool, ctx context.Context, code string) (int64, error) {
	var (
		err   error
		id    int64
		query = `SELECT id FROM tournaments WHERE tournament_code=$1`
	)
	err = db.QueryRow(ctx, query, code).Scan(&id)
	if err != nil {
		return id, c.errHandler("model.GetTournamentIdByCode", err, utils.ErrGettingTournamentByCode)
	}
	return id, nil
}

func (c *Contract) GetRemainingOfNonWinnerPlayers(db *pgxpool.Pool, ctx context.Context, code string) ([]NonWinnerEntity, error) {
	var (
		err   error
		list  []NonWinnerEntity
		query = `SELECT 
			tp.user_id AS user_id, 
			t.tournament_code AS tournament_code, 
			t.participant_vp AS participant_vp
		FROM tournaments t JOIN tournament_participants tp ON t.id = tp.tournament_id
		WHERE tp.position = 0 AND t.tournament_code = $1;`
	)

	rows, err := db.Query(ctx, query, code)
	if err != nil {
		return list, c.errHandler("model.GetRemainingOfNonWinnerPlayers", err, utils.ErrGetRemainingOfNonWinnerPlayers)
	}

	defer rows.Close()
	for rows.Next() {
		var data NonWinnerEntity
		err = rows.Scan(
			&data.UserId, &data.TournamentCode, &data.ParticipantPoint,
		)
		if err != nil {
			return list, c.errHandler("model.GetRemainingOfNonWinnerPlayers", err, utils.ErrGetRemainingOfNonWinnerPlayers)
		}

		list = append(list, data)
	}

	return list, nil
}

func (c *Contract) UpdateTournamentStatus(db *pgxpool.Pool, ctx context.Context, tournamentCode, status string) error {
	sql := `
		UPDATE tournaments
		SET updated_date = NOW(), status = $1
		WHERE tournament_code = $2 AND deleted_date IS NULL
	`

	_, err := db.Exec(ctx, sql, status, tournamentCode)
	if err != nil {
		return c.errHandler("model.UpdateTournamentStatus", err, utils.ErrUpdatingTournamentStatus)
	}

	return nil
}

func (c *Contract) UpdateTournamentStatusTrx(tx pgx.Tx, ctx context.Context, roomCode, status string) error {
	sql := `
		UPDATE tournaments
		SET status = $1, updated_date = NOW()
		WHERE tournament_code = $2 AND deleted_date IS NULL
	`

	_, err := tx.Exec(ctx, sql, status, roomCode)
	if err != nil {
		return c.errHandler("model.UpdateTournamentStatusTrx", err, utils.ErrUpdatingTournamentStatus)
	}

	return nil
}
