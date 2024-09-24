package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TournamentsEnt struct {
	GameCode        string         `db:"game_code"`
	GameName        string         `db:"game_name"`
	GameType        string         `db:"game_type"`
	GameImgUrl      string         `db:"game_img_url"`
	CafeCode        string         `db:"cafe_code"`
	CafeName        string         `db:"cafe_name"`
	CafeAddress     string         `db:"cafe_address"`
	TournamentCode  string         `db:"tournament_code"`
	ImageUrl        sql.NullString `db:"image_url"`
	PrizesImgUrl    sql.NullString `db:"prizes_img_url"`
	Name            sql.NullString `db:"name"`
	TournamentRules string         `db:"tournament_rules"`
	Level           string         `db:"level"`
	StartDate       sql.NullTime   `db:"start_date"`
	EndDate         sql.NullTime   `db:"end_date"`
	StartTime       time.Time      `db:"start_time"`
	EndTime         time.Time      `db:"end_time"`
	ParticipantVP   int64          `db:"participant_vp"`
	BookingPrice    int64          `db:"booking_price"`
	PlayerSlot      int64          `db:"player_slot"`
	CurrentUsedSlot int64          `db:"current_used_slot"`
	Status          string         `db:"status"`
	CreatedDate     time.Time      `db:"created_date"`
	UpdatedDate     sql.NullTime   `db:"updated_date"`
	DeletedDate     sql.NullTime   `db:"deleted_date"`
}

func (c *Contract) GetTournamentList(db *pgxpool.Pool, ctx context.Context, startDate string) ([]TournamentsEnt, error) {
	var (
		tournaments []TournamentsEnt
		query       = `
		SELECT 
		games.game_code, games.game_type, games.name, games.image_url, 
		cafes.cafe_code, cafes.name as cafe_name, cafes.address as cafe_address, 
		tournaments.tournament_code, tournaments.image_url, tournaments.name, tournaments.prizes_img_url, 
		tournaments.tournament_rules, tournaments.level, tournaments.start_date, tournaments.end_date, 
		tournaments.start_time, tournaments.end_time, tournaments.booking_price, 
		tournaments.player_slot, tournaments.participant_vp, 
		(
			SELECT COUNT(1) AS total_current_user 
			FROM tournaments JOIN tournament_participants ON tournament_participants.tournament_id = tournaments.id
			WHERE tournament_participants.status != 'cancel'
		) AS current_used_slot,
		tournaments.status,  
		tournaments.created_date, tournaments.updated_date, tournaments.deleted_date
		FROM tournaments
		LEFT JOIN games ON games.id = tournaments.game_id AND games.status = 'active' AND games.deleted_date IS NULL 
		LEFT JOIN cafes ON cafes.id = games.cafe_id AND cafes.status = 'active' AND cafes.deleted_date IS NULL 
    WHERE tournaments.start_date = $1 AND tournaments.deleted_date IS NULL
    `
	)

	rows, err := db.Query(ctx, query, startDate)
	if err != nil {
		return nil, c.errHandler("model.GetListTournament", err, utils.ErrFetchingTournamentByStartDate)
	}
	defer rows.Close()

	for rows.Next() {
		var data TournamentsEnt
		err = rows.Scan(
			&data.GameCode, &data.GameType, &data.GameName, &data.ImageUrl,
			&data.CafeCode, &data.CafeName, &data.CafeAddress,
			&data.TournamentCode, &data.ImageUrl, &data.Name, &data.PrizesImgUrl,
			&data.TournamentRules, &data.Level, &data.StartDate, &data.EndDate,
			&data.StartTime, &data.EndTime, &data.BookingPrice,
			&data.PlayerSlot, &data.ParticipantVP, &data.CurrentUsedSlot,
			&data.Status,
			&data.CreatedDate, &data.UpdatedDate, &data.DeletedDate,
		)
		if err != nil {
			return nil, c.errHandler("model.GetListTournament", err, utils.ErrScanningListTournament)
		}
		tournaments = append(tournaments, data)
	}

	if err := rows.Err(); err != nil {
		return nil, c.errHandler("model.GetListTournament", err, utils.ErrGettingListTournament)
	}

	return tournaments, nil
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

func (c *Contract) GetListTournamentCodes(db *pgxpool.Pool, ctx context.Context) ([]string, error) {
	var (
		err             error
		tournamentCodes []string
		query           = `SELECT tournament_code FROM tournaments WHERE end_date < NOW() AND deleted_date IS NULL`
	)

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, c.errHandler("model.GetListTournamentCodes", err, utils.ErrGettingListTournamentCode)
	}
	defer rows.Close()

	for rows.Next() {
		var tournamentCode string
		if err = rows.Scan(&tournamentCode); err != nil {
			return nil, c.errHandler("model.GetListTournamentCodes", err, utils.ErrGettingListTournamentCode)
		}
		tournamentCodes = append(tournamentCodes, tournamentCode)
	}

	if rows.Err() != nil {
		return nil, c.errHandler("model.GetListTournamentCodes", rows.Err(), utils.ErrGettingListTournamentCode)
	}

	return tournamentCodes, nil
}
