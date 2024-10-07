package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	TournamentParticipantEnt struct {
		Id              int64          `db:"id"`
		TournamentId    int64          `db:"tournament_id"`
		UserId          int64          `db:"user_id"`
		StatusWinner    bool           `db:"status_winner"`
		Status          string         `db:"status"`
		Position        int            `db:"position"`
		AdditionalInfo  sql.NullString `db:"additional_info"`
		RewardPoint     sql.NullInt64  `db:"reward_point"`
		TransactionCode sql.NullString `db:"transaction_code"`
	}

	TournamentParticipantRespEnt struct {
		Id                  int64          `db:"id"`
		TournamentId        int64          `db:"tournament_id"`
		UserCode            string         `db:"user_code"`
		UserName            string         `db:"user_name"`
		UserImgUrl          string         `db:"user_image_url"`
		UserXPlayer         string         `db:"user_x_player"`
		StatusWinner        bool           `db:"status_winner"`
		Status              string         `db:"status"`
		AdditionalInfo      sql.NullString `db:"additional_info"`
		LatestTier          sql.NullString `db:"latest_tier"`
		Position            int            `db:"position"`
		RewardPoint         sql.NullInt64  `db:"reward_point"`
		TransactionCode     sql.NullString `db:"transaction_code"`
		UserId              int64          `db:"user_id"`
		TournamentBannerUri string         `db:"tournament_banner_uri"`
	}
)

func (c *Contract) GetAllParticipantByTournamentCode(db *pgxpool.Pool, ctx context.Context, code string) ([]TournamentParticipantRespEnt, error) {
	var (
		err   error
		list  []TournamentParticipantRespEnt
		query = `SELECT 
			tp.id,
			tp.tournament_id,
			u.user_code,
			COALESCE(u.username, '') AS user_name,
			u.image_url AS user_image_url,
			u.x_player AS user_x_player, 
			tp.status_winner,
			tp.status,
			tp.position,
			tp.additional_info,
			tp.reward_point, 
			tr.name as latest_tier_name 
			FROM tournament_participants tp
				LEFT JOIN tournaments t ON tp.tournament_id = t.id
				LEFT JOIN users u ON tp.user_id = u.id
				LEFT JOIN tiers tr ON tr.id = u.latest_tier_id 
			WHERE tournament_code = $1 AND tp.status = 'active' `
	)

	rows, err := db.Query(ctx, query, code)
	if err != nil {
		return list, c.errHandler("model.GetAllParticipantByTournamentCode", err, utils.ErrGettingAllParticipantByTournamentCode)
	}

	defer rows.Close()
	for rows.Next() {
		var data TournamentParticipantRespEnt
		err = rows.Scan(
			&data.Id, &data.TournamentId, &data.UserCode,
			&data.UserName, &data.UserImgUrl, &data.UserXPlayer,
			&data.StatusWinner, &data.Status, &data.Position,
			&data.AdditionalInfo, &data.RewardPoint, &data.LatestTier,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				return list, nil
			}
			return list, c.errHandler("model.GetAllParticipantByTournamentCode", err, utils.ErrScanningListTournament)
		}
		list = append(list, data)
	}

	return list, nil
}

func (c *Contract) GetOneTournamentParticipant(db *pgxpool.Pool, ctx context.Context, participantID, userID int64) (TournamentParticipantEnt, error) {
	var (
		participant TournamentParticipantEnt
		query       = `
            SELECT id, tournament_id, user_id, status_winner, status, position, 
                   additional_info, reward_point, transaction_code
            FROM tournament_participants
            WHERE tournament_id = $1 AND user_id = $2
        `
	)

	err := db.QueryRow(ctx, query, participantID, userID).Scan(
		&participant.Id,
		&participant.TournamentId,
		&participant.UserId,
		&participant.StatusWinner,
		&participant.Status,
		&participant.Position,
		&participant.AdditionalInfo,
		&participant.RewardPoint,
		&participant.TransactionCode,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return participant, c.errHandler("model.GetOneTournamentParticipant", err, utils.EmptyData)
		}
		return participant, c.errHandler("model.GetOneTournamentParticipant", err, utils.ErrFetchingTournamentParticipant)
	}

	return participant, nil
}

func (c *Contract) CountTournamentParticipantByUserId(db *pgxpool.Pool, ctx context.Context, tournamentID, userID int64) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM tournament_participants
		WHERE tournament_id = $1 AND user_id = $2 AND status = 'active'
	`

	err := db.QueryRow(ctx, query, tournamentID, userID).Scan(&count)
	if err != nil {
		return 0, c.errHandler("model.CountTournamentParticipant", err, utils.ErrFetchingTournamentParticipant)
	}

	return count, nil
}

func (c *Contract) CountTournamentParticipantByUserIdAndGameIdAndIsGameMasterAndBookingPrice(db *pgxpool.Pool, ctx context.Context, userId, gameId int64, bookingPrice float64) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
			FROM tournament_participants tp
		LEFT JOIN tournaments t ON tp.tournament_id = t.id
		LEFT JOIN games g ON g.id = t.game_id
		LEFT JOIN users u ON tp.user_id = u.id
		WHERE u.id = $1 AND t.game_id = $2 AND tp.status = 'active' AND t.booking_price >= $3
	`

	err := db.QueryRow(ctx, query, userId, gameId, bookingPrice).Scan(&count)
	if err != nil {
		return 0, c.errHandler("model.CountTournamentParticipantByUserIdAndGameIdAndIsGameMaster", err, utils.ErrCountParticipantRoomByUserIdAndGameIdAndIsGameMaster)
	}

	return count, nil
}

func (c *Contract) CountTournamentParticipantByUserIdAndStartDateAndEndDate(db *pgxpool.Pool, ctx context.Context, userId int64, startDate, endDate string) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM tournament_participants tp
		LEFT JOIN tournaments t ON tp.tournament_id = t.id
		LEFT JOIN games g ON g.id = t.game_id
		LEFT JOIN users u ON tp.user_id = u.id
		WHERE u.id = $1 AND tp.status = 'active' AND (t.start_date BETWEEN $2 AND $3)
	`

	err := db.QueryRow(ctx, query, userId, startDate, endDate).Scan(&count)
	if err != nil {
		return 0, c.errHandler("model.CountTournamentParticipantByUserIdAndStartDateAndEndDate", err, utils.ErrCountParticipantTournamentByStartDateAndEndDate)
	}

	return count, nil
}

func (c *Contract) CountTournamentParticipantByUserIdAndStartDateAndLifeTime(db *pgxpool.Pool, ctx context.Context, userId int64, startDate string, endDate string) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM tournament_participants tp
		LEFT JOIN tournaments t ON tp.tournament_id = t.id
		LEFT JOIN games g ON g.id = t.game_id
		LEFT JOIN users u ON tp.user_id = u.id
		WHERE u.id = $1 AND tp.status = 'active' AND (t.start_date BETWEEN $2 AND $3)  AND t.end_date = $4
	`

	err := db.QueryRow(ctx, query, userId, startDate, startDate, nil).Scan(&count)
	if err != nil {
		return 0, c.errHandler("model.CountRoomParticipantByUserIdAndStartDateAndEndDate", err, utils.ErrCountParticipantRoomByStartDateAndEndDate)
	}

	return count, nil
}

func (c *Contract) GetParticipantByTournamentCodeAndUserCode(db *pgxpool.Pool, ctx context.Context, tournamentCode, userCode string) (TournamentParticipantRespEnt, error) {
	var (
		err  error
		data TournamentParticipantRespEnt

		queryGetRoomGameDetail = `SELECT 
		u.user_code AS user_code,
		u.username AS user_name,
		u.image_url AS user_image_url,
		u.x_player AS user_x_player, 
		tp.status_winner,
		tp.status,
		tp.transaction_code,
		tp.tournament_id,
		tp.user_id,
		t.image_url AS tournament_banner_uri
		FROM tournament_participants tp
			JOIN tournaments t ON tp.tournament_id = t.id
			JOIN users u ON tp.user_id = u.id
		WHERE t.tournament_code = $1 and u.user_code =$2 `
	)

	err = db.QueryRow(ctx, queryGetRoomGameDetail, tournamentCode, userCode).Scan(
		&data.UserCode, &data.UserName, &data.UserImgUrl, &data.UserXPlayer,
		&data.StatusWinner, &data.Status, &data.TransactionCode,
		&data.TournamentId, &data.UserId, &data.TournamentBannerUri,
	)

	if err != nil && err != pgx.ErrNoRows {
		return data, c.errHandler("model.GetParticipantByTournamentCodeAndUserCode", err, utils.ErrGettingTournamentByCodeAndUserCode)
	}

	return data, nil
}

func (c *Contract) CountTournamentWinnerByUserId(db *pgxpool.Pool, ctx context.Context, userId int64) (int, error) {
	var (
		err   error
		total int

		query = `SELECT COUNT(*)
		FROM tournament_participants tp
			JOIN tournaments t ON tp.tournament_id = t.id
			JOIN users u ON tp.user_id = u.id
		WHERE u.id = $1 AND tp.status_winner = $2`
	)

	err = db.QueryRow(ctx, query, userId, true).Scan(&total)
	if err != nil && err != pgx.ErrNoRows {
		return 0, c.errHandler("model.CountTournamentWinnerByUserId", err, utils.ErrCountParticipantWonTournament)
	}

	return total, nil
}

func (c *Contract) InsertOneTournamentParticipant(tx pgx.Tx, ctx context.Context, tournamentId, userId int64, statusWinner bool, position int, status string, additionalInfo string, rewardPoint int64, transactionCode string) error {
	var (
		err   error
		query = `INSERT INTO tournament_participants(tournament_id, user_id, status_winner, position, status, additional_info, reward_point, transaction_code) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`
	)
	_, err = tx.Exec(ctx, query, tournamentId, userId, statusWinner, position, status, additionalInfo, rewardPoint, transactionCode)
	if err != nil {
		return c.errHandler("model.InsertOneTournamentParticipant", err, utils.ErrAddingTournamentParticipant)
	}
	return nil
}

func (c *Contract) UpdateTournamentParticipant(tx pgx.Tx, ctx context.Context, tournamentId, userId int64, statusWinner bool, position int, status string, additionalInfo string, rewardPoint int64, transactionCode string) error {
	var (
		err   error
		query = `UPDATE tournament_participants 
        SET status_winner = $1, status = $2, additional_info = $3, reward_point = $4, position = $5, transaction_code = $6, updated_date = $7
        WHERE tournament_id = $8 AND user_id = $9`
	)
	_, err = tx.Exec(ctx, query, statusWinner, status, additionalInfo, rewardPoint, position, transactionCode, time.Now().In(time.UTC), tournamentId, userId)
	if err != nil {
		return c.errHandler("model.UpdateTournamentParticipant", err, utils.ErrUpdatingTournamentParticipant)
	}
	return nil
}

func (c *Contract) DeleteTournamentParticipant(tx pgx.Tx, ctx context.Context, tournamentId, userId int64) error {
	var (
		err   error
		query = `DELETE FROM tournament_participants WHERE tournament_id=$1 AND user_id=$2`
	)
	_, err = tx.Exec(ctx, query, tournamentId, userId)
	if err != nil {
		return c.errHandler("model.DeleteTournamentParticipant", err, utils.ErrDeletingTournamentParticipant)
	}
	return nil
}
