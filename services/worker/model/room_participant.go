package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"

	"github.com/jackc/pgx/v4/pgxpool"
)

type RoomParticipantResp struct {
	UserCode        string         `db:"user_code"`
	UserName        string         `db:"user_name"`
	UserImgUrl      string         `db:"user_image_url"`
	UserXPlayer     string         `db:"user_x_player"`
	StatusWinner    bool           `db:"status_winner"`
	Status          string         `db:"status"`
	TransactionCode sql.NullString `db:"transaction_code"`
	AdditionalInfo  sql.NullString `db:"additional_info"`
	Position        int            `db:"position"`
	RewardPoint     sql.NullInt64  `db:"reward_point"`
}

func (c *Contract) GetAllParticipantByRoomCode(db *pgxpool.Pool, ctx context.Context, code string) ([]RoomParticipantResp, error) {
	var (
		err   error
		list  []RoomParticipantResp
		query = `
			SELECT 
				u.user_code AS user_code,
				COALESCE(u.username, '') AS user_name,
				u.image_url AS user_image_url,
				u.x_player AS user_x_player,
				rp.status_winner,
				rp.status,
				rp.position,
				rp.additional_info,
				rp.reward_point
			FROM rooms r
				JOIN rooms_participants rp ON rp.room_id = r.id
				JOIN users u ON rp.user_id = u.id
			WHERE room_code = $1 AND rp.status = 'active' AND r.deleted_date IS NULL`
	)

	rows, err := db.Query(ctx, query, code)
	if err != nil {
		return list, c.errHandler("model.GetAllParticipantByRoomCode", err, utils.ErrGettingAllParticipantByRoomCode)
	}

	defer rows.Close()
	for rows.Next() {
		var data RoomParticipantResp
		err = rows.Scan(
			&data.UserCode, &data.UserName, &data.UserImgUrl,
			&data.UserXPlayer, &data.StatusWinner, &data.Status,
			&data.Position, &data.AdditionalInfo, &data.RewardPoint,
		)
		if err != nil {
			return list, c.errHandler("model.GetAllParticipantByRoomCode", err, utils.ErrScanningListRoom)
		}
		list = append(list, data)
	}

	return list, nil
}
