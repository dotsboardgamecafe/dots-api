package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type RoomEnt struct {
	CafeCode           string         `db:"cafe_code"`
	CafeName           string         `db:"cafe_name"`
	CafeAddress        string         `db:"cafe_address"`
	RoomCode           string         `db:"room_code"`
	RoomType           string         `db:"room_type"`
	RoomImgUrl         string         `db:"image_url"`
	Name               string         `db:"name"`
	Difficulty         string         `db:"difficulty"`
	Description        string         `db:"description"`
	Instruction        string         `db:"instruction"`
	StartDate          time.Time      `db:"start_date"`
	EndDate            time.Time      `db:"end_date"`
	MaximumParticipant int            `db:"maximum_participant"`
	CurrentUsedSlot    int            `db:"current_used_slot"`
	InstagramLink      string         `db:"instagram_link"`
	Status             string         `db:"status"`
	BookingPrice       int64          `db:"booking_price"`
	GameMasterName     sql.NullString `db:"game_master_name"`
	GameMasterImageUrl sql.NullString `db:"game_master_image_url"`
	GameCode           string         `db:"game_code"`
	GameName           string         `db:"game_name"`
	GameImgUrl         string         `db:"game_img_url"`
}

func (c *Contract) GetRoomList(db *pgxpool.Pool, ctx context.Context, startDate string) ([]RoomEnt, error) {
	var (
		err  error
		list []RoomEnt

		query = `SELECT 
				c.cafe_code AS cafe_code,
				c.name AS cafe_name,
				c.address AS cafe_address,
				rooms.room_code,
				rooms.room_type,
				rooms.image_url,
				rooms.name,
				rooms.difficulty,
				rooms.description,
				COALESCE(rooms.instruction, '') AS instruction,
				rooms.start_date,
				rooms.end_date,
				rooms.maximum_participant,
				rooms.instagram_link,
				rooms.status,
				rooms.booking_price,
				COALESCE(cp.count_participants, 0) AS current_used_slot,
				a."name" AS game_master_name,
				a.image_url AS game_master_image_url,
				g.game_code,
				g.name AS game_name,
				g.image_url AS game_img_url
			FROM rooms
				JOIN games g ON rooms.game_id = g.id 
				JOIN cafes c ON c.id = g.cafe_id
				LEFT JOIN admins a ON rooms.game_master_id = a.id
				LEFT JOIN (
					SELECT count(rp.user_id) AS count_participants, rp.room_id
					FROM rooms_participants rp
					WHERE rp.status != 'cancel'
					GROUP BY rp.room_id
				) AS cp ON cp.room_id = rooms.id
			WHERE rooms.start_date = $1 AND rooms.deleted_date IS NULL`
	)

	rows, err := db.Query(ctx, query, startDate)

	if err != nil {
		return nil, c.errHandler("model.GetRoomList", err, utils.ErrGettingListRoom)
	}

	defer rows.Close()

	for rows.Next() {
		var data RoomEnt
		err = rows.Scan(
			&data.CafeCode, &data.CafeName, &data.CafeAddress,
			&data.RoomCode, &data.RoomType, &data.RoomImgUrl,
			&data.Name, &data.Difficulty, &data.Description, &data.Instruction, &data.StartDate,
			&data.EndDate, &data.MaximumParticipant, &data.InstagramLink, &data.Status, &data.BookingPrice,
			&data.CurrentUsedSlot,
			&data.GameMasterName, &data.GameMasterImageUrl,
			&data.GameCode, &data.GameName, &data.GameImgUrl,
		)
		if err != nil {
			return nil, c.errHandler("model.GetRoomList", err, utils.ErrScanningListRoom)
		}
		list = append(list, data)
	}
	if err := rows.Err(); err != nil {
		return nil, c.errHandler("model.GetRoomList", err, utils.ErrGettingListRoom)
	}

	return list, nil
}

func (c *Contract) UpdateRoomStatus(db *pgxpool.Pool, ctx context.Context, roomCode, status string) error {
	sql := `
		UPDATE rooms
		SET updated_date = NOW(), status = $1
		WHERE room_code = $2 AND deleted_date IS NULL
	`

	_, err := db.Exec(ctx, sql, status, roomCode)
	if err != nil {
		return c.errHandler("model.UpdateRoomStatus", err, utils.ErrUpdatingRoomStatus)
	}

	return nil
}

func (c *Contract) GetListRoomCodes(db *pgxpool.Pool, ctx context.Context) ([]string, error) {
	var (
		err       error
		roomCodes []string
		query     = `SELECT room_code FROM rooms WHERE end_date < NOW() AND deleted_date IS NULL`
	)

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, c.errHandler("model.GetAllRoomCodes", err, utils.ErrGettingListRoomCode)
	}
	defer rows.Close()

	for rows.Next() {
		var roomCode string
		if err = rows.Scan(&roomCode); err != nil {
			return nil, c.errHandler("model.GetAllRoomCodes", err, utils.ErrGettingListRoomCode)
		}
		roomCodes = append(roomCodes, roomCode)
	}

	if rows.Err() != nil {
		return nil, c.errHandler("model.GetAllRoomCodes", rows.Err(), utils.ErrGettingListRoomCode)
	}

	return roomCodes, nil
}
