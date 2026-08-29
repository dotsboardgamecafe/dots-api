package model

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"dots-api/lib/utils"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	ParticipantAdditionalInfo struct {
		RegistrationType string `db:"registration_type" json:"registration_type"`
	}

	RoomParticipantResp struct {
		UserCode        string                    `db:"user_code"`
		UserName        string                    `db:"user_name"`
		UserImgUrl      string                    `db:"user_image_url"`
		UserXPlayer     string                    `db:"user_x_player"`
		StatusWinner    bool                      `db:"status_winner"`
		Status          string                    `db:"status"`
		TransactionCode sql.NullString            `db:"transaction_code"`
		AdditionalInfo  ParticipantAdditionalInfo `db:"additional_info"`
		Position        int                       `db:"position"`
		RewardPoint     sql.NullInt64             `db:"reward_point"`
	}
)

func (p *ParticipantAdditionalInfo) Value() (driver.Value, error) {
	if p == nil || p.RegistrationType == "" {
		return json.Marshal(ParticipantAdditionalInfo{RegistrationType: "self_booking"})
	}
	return json.Marshal(p)
}

func (p *ParticipantAdditionalInfo) Scan(value interface{}) error {
	if value == nil {
		*p = ParticipantAdditionalInfo{RegistrationType: "self_booking"}
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if ok {
			if str == "" || str == "member" {
				*p = ParticipantAdditionalInfo{RegistrationType: "self_booking"}
				return nil
			}
			return json.Unmarshal([]byte(str), p)
		}

		return errors.New("worker.ParticipantAdditionalInfo.Scan: " + utils.ErrTypeAssertionFailed + " expected []byte, got " + fmt.Sprintf("%T", value))
	}

	if string(b) == "member" || string(b) == "" {
		*p = ParticipantAdditionalInfo{RegistrationType: "self_booking"}
		return nil
	}

	return json.Unmarshal(b, p)
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
