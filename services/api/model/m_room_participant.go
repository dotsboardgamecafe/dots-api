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
	RoomParticipantEnt struct {
		Id              int64          `db:"id"`
		RoomId          int64          `db:"room_id"`
		UserId          int64          `db:"user_id"`
		StatusWinner    bool           `db:"status_winner"`
		Status          string         `db:"status"`
		Position        int            `db:"position"`
		AdditionalInfo  sql.NullString `db:"additional_info"`
		RewardPoint     sql.NullInt64  `db:"reward_point"`
		TransactionCode sql.NullString `db:"transaction_code"`
	}

	RoomParticipantResp struct {
		UserCode           string         `db:"user_code"`
		UserName           string         `db:"user_name"`
		UserImgUrl         string         `db:"user_image_url"`
		UserXPlayer        string         `db:"user_x_player"`
		StatusWinner       bool           `db:"status_winner"`
		Status             string         `db:"status"`
		TransactionCode    sql.NullString `db:"transaction_code"`
		AdditionalInfo     sql.NullString `db:"additional_info"`
		Position           int            `db:"position"`
		RewardPoint        sql.NullInt64  `db:"reward_point"`
		RoomId             int64          `db:"room_id"`
		UserId             int64          `db:"user_id"`
		ParticipationPoint int            `db:"participation_point"`
		LatestTier         sql.NullString `db:"latest_tier"`
		RoomBannerUri      string         `db:"room_banner_uri"`
	}
)

func (c *Contract) GetAllParticipantByRoomCode(db *pgxpool.Pool, ctx context.Context, code string) ([]RoomParticipantResp, error) {
	var (
		err   error
		list  []RoomParticipantResp
		query = `SELECT 
			u.user_code AS user_code,
			COALESCE(u.username, '') AS user_name,
			u.image_url AS user_image_url,
			u.x_player AS user_x_player, 
			rp.status_winner,
			rp.status,
			rp.position,
			rp.additional_info,
			rp.reward_point,
			tr.name as latest_tier_name 
			FROM rooms r
				JOIN rooms_participants rp ON rp.room_id = r.id
				JOIN users u ON rp.user_id = u.id
				LEFT JOIN tiers tr ON tr.id = u.latest_tier_id 
			WHERE room_code = $1 AND rp.status = 'active' `
	)

	rows, err := db.Query(ctx, query, code)
	if err != nil {
		return list, c.errHandler("model.GetAllParticipantByRoomCode", err, utils.ErrGettingAllParticipantByRoomCode)
	}

	defer rows.Close()
	for rows.Next() {
		var data RoomParticipantResp
		err = rows.Scan(
			&data.UserCode, &data.UserName, &data.UserImgUrl, &data.UserXPlayer,
			&data.StatusWinner, &data.Status,
			&data.Position, &data.AdditionalInfo, &data.RewardPoint,
			&data.LatestTier,
		)
		if err != nil {
			return list, c.errHandler("model.GetAllParticipantByRoomCode", err, utils.ErrScanningAllParticipantByRoomCode)
		}
		list = append(list, data)
	}

	return list, nil
}

func (c *Contract) GetParticipantByRoomCodeAndUserCode(db *pgxpool.Pool, ctx context.Context, roomCode, userCode string) (RoomParticipantResp, error) {
	var (
		err  error
		data RoomParticipantResp

		queryGetRoomGameDetail = `SELECT 
		u.user_code AS user_code,
		u.username AS user_name,
		u.image_url AS user_image_url,
		u.x_player AS user_x_player, 
		rp.status_winner,
		rp.status,
		rp.transaction_code,
		rp.room_id,
		rp.user_id,
		r.reward_point AS participation_point,
		tr.name as latest_tier_name,
		r.image_url AS room_banner_uri
		FROM rooms_participants rp
			JOIN rooms r ON rp.room_id = r.id
			JOIN users u ON rp.user_id = u.id
			LEFT JOIN tiers tr ON tr.id = u.latest_tier_id 
		WHERE room_code = $1 and u.user_code =$2 `
	)

	err = db.QueryRow(ctx, queryGetRoomGameDetail, roomCode, userCode).Scan(
		&data.UserCode, &data.UserName, &data.UserImgUrl, &data.UserXPlayer,
		&data.StatusWinner, &data.Status, &data.TransactionCode,
		&data.RoomId, &data.UserId, &data.ParticipationPoint, &data.LatestTier, &data.RoomBannerUri,
	)
	if err != nil && err != pgx.ErrNoRows {
		return data, c.errHandler("model.GetParticipantByRoomCodeAndUserCode", err, utils.ErrGettingRoomByCodeAndUserCode)
	}

	return data, nil
}

func (c *Contract) GetOneRoomParticipant(db *pgxpool.Pool, ctx context.Context, roomId, userId int64) (*RoomParticipantEnt, error) {
	var (
		participant RoomParticipantEnt
		query       = `
			SELECT 
				id, room_id, user_id, status_winner, status, position, additional_info, reward_point, transaction_code
			FROM rooms_participants
			WHERE room_id = $1 AND user_id = $2
		`
	)

	err := db.QueryRow(ctx, query, roomId, userId).Scan(
		&participant.Id,
		&participant.RoomId,
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
			return nil, c.errHandler("model.GetOneRoomParticipant", err, utils.EmptyData)
		}
		return nil, c.errHandler("model.GetOneRoomParticipant", err, utils.ErrGettingRoomByCodeAndUserCode)
	}

	return &participant, nil
}

func (c *Contract) CountRoomParticipantByUserId(db *pgxpool.Pool, ctx context.Context, userId int64) (int64, error) {
	var count int64
	query := `
		SELECT 
			COUNT(*)
		FROM rooms_participants rp
		LEFT JOIN users u ON rp.user_id = u.id
		WHERE u.id = $1 AND rp.status = 'active AND game_master_id != NULL'
	`

	err := db.QueryRow(ctx, query, userId).Scan(&count)
	if err != nil {
		return 0, c.errHandler("model.CountRoomParticipantByUserId", err, utils.ErrGettingRoomByCodeAndUserCode)
	}

	return count, nil
}

func (c *Contract) CountRoomParticipantByUserIdAndGameIdAndIsGameMasterAndBookingPrice(db *pgxpool.Pool, ctx context.Context, userId, gameId int64, bookingPrice float64, isGameMaster bool) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
			FROM rooms_participants rp
		LEFT JOIN rooms r ON rp.room_id = r.id
		LEFT JOIN games g ON g.id = r.game_id
		LEFT JOIN users u ON rp.user_id = u.id
		WHERE u.id = $1 AND rp.status = 'active' AND r.game_id = $2 AND r.booking_price >= $3 
	`
	if isGameMaster {
		query += " AND r.game_master_id IS NOT NULL"
	}

	err := db.QueryRow(ctx, query, userId, gameId, bookingPrice).Scan(&count)
	if err != nil {
		return 0, c.errHandler("model.CountRoomParticipantByUserIdAndGameIdAndIsGameMaster", err, utils.ErrCountParticipantRoomByUserIdAndGameIdAndIsGameMaster)
	}

	return count, nil
}

func (c *Contract) CountRoomParticipantByUserIdAndStartDateAndEndDate(db *pgxpool.Pool, ctx context.Context, userId int64, startDate, endDate string) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM rooms_participants rp
		LEFT JOIN rooms r ON rp.room_id = r.id
		LEFT JOIN games g ON g.id = r.game_id
		LEFT JOIN users u ON rp.user_id = u.id
		WHERE u.id = $1 AND rp.status = 'active' AND ( r.start_date BETWEEN $2 AND $3 )
	`

	err := db.QueryRow(ctx, query, userId, startDate, endDate).Scan(&count)
	if err != nil {
		return 0, c.errHandler("model.CountRoomParticipantByUserIdAndStartDateAndEndDate", err, utils.ErrCountParticipantRoomByStartDateAndEndDate)
	}

	return count, nil
}

func (c *Contract) CountRoomParticipantByUserIdAndStartDateAndLifeTime(db *pgxpool.Pool, ctx context.Context, userId int64, startDate string, endDate string) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM rooms_participants rp
		LEFT JOIN rooms r ON rp.room_id = r.id
		LEFT JOIN games g ON g.id = r.game_id
		LEFT JOIN users u ON rp.user_id = u.id
		WHERE u.id = $1 AND rp.status = 'active' AND (r.start_date BETWEEN $2 AND $3) AND r.end_date = $4
	`

	err := db.QueryRow(ctx, query, userId, startDate, startDate, nil).Scan(&count)
	if err != nil {
		return 0, c.errHandler("model.CountRoomParticipantByUserIdAndStartDateAndEndDate", err, utils.ErrCountParticipantRoomByStartDateAndEndDate)
	}

	return count, nil
}

func (c *Contract) InsertOneRoomParticipant(tx pgx.Tx, ctx context.Context, roomId, userId int64, status string, rewardPoint int64, transactionCode string) error {
	var (
		err   error
		query = `INSERT INTO rooms_participants(room_id, user_id, status, reward_point, transaction_code) VALUES($1,$2,$3,$4,$5)`
	)
	_, err = tx.Exec(ctx, query, roomId, userId, status, rewardPoint, transactionCode)
	if err != nil {
		return c.errHandler("model.InsertOneRoomParticipant", err, utils.ErrAddingRoomParticipant)
	}
	return nil
}

func (c *Contract) UpdateRoomParticipant(tx pgx.Tx, ctx context.Context, roomId, userId int64, statusWinner bool, position int, status string, additionalInfo string, rewardPoint int64, transactionCode string) error {
	var (
		err   error
		query = `UPDATE rooms_participants 
		SET status_winner=$1, status=$2, additional_info=$3, reward_point=$4, position=$5, transaction_code=$6, updated_date=$7
        WHERE room_id=$8 AND user_id=$9`
	)
	_, err = tx.Exec(ctx, query, statusWinner, status, additionalInfo, rewardPoint, position, transactionCode, time.Now().UTC(), roomId, userId)
	if err != nil {
		return c.errHandler("model.UpdateRoomParticipant", err, utils.ErrUpdatingRoomParticipant)
	}
	return nil
}

func (c *Contract) DeleteRoomParticipant(tx pgx.Tx, ctx context.Context, roomId, userId int64) error {
	var (
		err   error
		query = `DELETE FROM rooms_participants WHERE room_id=$1 AND user_id=$2`
	)
	_, err = tx.Exec(ctx, query, roomId, userId)
	if err != nil {
		return c.errHandler("model.DeleteRoomParticipant", err, utils.ErrDeletingRoomParticipant)
	}
	return nil
}

func (c *Contract) CountParticipantRoomByRoomId(db *pgxpool.Pool, ctx context.Context, id int64) (int, error) {
	var (
		err              error
		totalParticipant int

		queryGetTotalParticipant = `select count(id)  from rooms_participants rp 
		where rp.status!= 'cancel' and room_id = $1
		group by room_id `
	)

	err = db.QueryRow(ctx, queryGetTotalParticipant, id).Scan(
		&totalParticipant,
	)
	if err != nil && err != pgx.ErrNoRows {
		return totalParticipant, c.errHandler("model.CountParticipantRoomByRoomId", err, utils.ErrCountParticipantRoomByRoomId)
	}

	return totalParticipant, nil
}
