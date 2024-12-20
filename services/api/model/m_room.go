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
	RoomEnt struct {
		RoomId             int64           `db:"room_id"`
		GameMasterCode     sql.NullString  `db:"game_master_code"`
		GameMasterName     sql.NullString  `db:"game_master_name"`
		GameMasterImgUrl   sql.NullString  `db:"game_master_img_url"`
		GameId             int64           `db:"game_id"`
		GameCode           string          `db:"game_code"`
		GameName           string          `db:"game_name"`
		GameImgUrl         string          `db:"game_img_url"`
		CafeCode           string          `db:"cafe_code"`
		CafeName           string          `db:"cafe_name"`
		CafeAddress        string          `db:"cafe_address"`
		RoomCode           string          `db:"room_code"`
		RoomType           string          `db:"room_type"`
		BannerRoomUrl      string          `db:"room_banner_url"`
		Name               string          `db:"name"`
		Description        string          `db:"description"`
		SpecialInstruction string          `db:"special_instruction"`
		Difficulty         string          `db:"difficulty"`
		StartDate          sql.NullTime    `db:"start_date"`
		EndDate            sql.NullTime    `db:"end_date"`
		StartTime          time.Time       `db:"start_time"`
		EndTime            time.Time       `db:"end_time"`
		MaximumParticipant int             `db:"maximum_participant"`
		BookingPrice       float64         `db:"booking_price"`
		RewardPoint        int             `db:"reward_point"`
		InstagramLink      string          `db:"instagram_link"`
		Status             string          `db:"status"`
		DayPastEndDate     sql.NullFloat64 `db:"day_past_end_date"`
		CurrentUsedSlot    int             `db:"current_used_slot"`
	}

	RoomListEnt struct {
		CafeId             int64           `db:"cafe_id"`
		CafeCode           string          `db:"cafe_code"`
		CafeName           string          `db:"cafe_name"`
		CafeAddress        string          `db:"cafe_address"`
		RoomCode           string          `db:"room_code"`
		RoomType           string          `db:"room_type"`
		RoomImgUrl         string          `db:"image_url"`
		Name               string          `db:"name"`
		Difficulty         string          `db:"difficulty"`
		Description        string          `db:"description"`
		Instruction        string          `db:"instruction"`
		StartDate          sql.NullTime    `db:"start_date"`
		EndDate            sql.NullTime    `db:"end_date"`
		StartTime          time.Time       `db:"start_time"`
		EndTime            time.Time       `db:"end_time"`
		MaximumParticipant int             `db:"maximum_participant"`
		CurrentUsedSlot    int             `db:"current_used_slot"`
		InstagramLink      string          `db:"instagram_link"`
		DayPastEndDate     sql.NullFloat64 `db:"day_past_end_date"`
		Status             string          `db:"status"`
		BookingPrice       float64         `db:"booking_price"`
		GameMasterName     sql.NullString  `db:"game_master_name"`
		GameMasterImageUrl sql.NullString  `db:"game_master_image_url"`
		GameCode           string          `db:"game_code"`
		GameName           string          `db:"game_name"`
		GameImgUrl         string          `db:"game_img_url"`
	}
)

func (c *Contract) GetRoomList(db *pgxpool.Pool, ctx context.Context, param request.RoomParam) ([]RoomListEnt, request.RoomParam, error) {
	var (
		err        error
		list       []RoomListEnt
		paramQuery []interface{}
		totalData  int
		where      []string

		query = `
			SELECT 
				c.id AS cafe_id,
				c.cafe_code AS cafe_code,
				c.name AS cafe_name,
				c.address AS cafe_address,
				rooms.room_code,
				rooms.room_type,
				COALESCE(rooms.image_url, '') AS image_url,
				rooms.name,
				rooms.difficulty,
				rooms.description,
				COALESCE(rooms.instruction, '') AS instruction,
				rooms.start_date,
				rooms.end_date,
				rooms.start_time,
				rooms.end_time,
				rooms.maximum_participant,
				rooms.instagram_link,
				rooms.status,
				DATE_PART('day', NOW() - rooms.end_date) AS days_past_end_date, 
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
				left join admins a on rooms.game_master_id = a.id
				left join (
					select count(rp.user_id) count_participants, rp.room_id
					from rooms_participants rp
					WHERE rp.status != 'cancel'
					group by rp.room_id
				) as cp on cp.room_id = rooms.id
			`
	)
	// ROOM NAME (KEYWORD)
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		orWhere = append(orWhere, fmt.Sprintf("rooms.name iLIKE $%d", len(paramQuery)))
		where = append(where, "("+strings.Join(orWhere, " AND ")+")")
	}

	// ROOM TYPE
	if len(param.RoomType) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.RoomType)
		orWhere = append(orWhere, fmt.Sprintf("rooms.room_type = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// CAFE CITY
	if len(param.CafeCity) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, strings.ToLower(param.CafeCity))
		orWhere = append(orWhere, fmt.Sprintf("LOWER(rooms.location_city) = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// ROOM STATUS
	if len(param.Status) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Status)
		orWhere = append(orWhere, fmt.Sprintf("rooms.status = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// Append All Where Conditions
	if len(where) > 0 {
		query += " WHERE rooms.deleted_date is null AND " + strings.Join(where, " AND ")
	} else {
		query += " WHERE rooms.deleted_date is null "
	}

	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS total`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetRoomList", err, utils.ErrCountingListRoom)
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
		return list, param, c.errHandler("model.GetRoomList", err, utils.ErrGettingListRoom)
	}

	defer rows.Close()
	for rows.Next() {
		var data RoomListEnt
		err = rows.Scan(
			&data.CafeId, &data.CafeCode, &data.CafeName, &data.CafeAddress,
			&data.RoomCode, &data.RoomType, &data.RoomImgUrl,
			&data.Name, &data.Difficulty, &data.Description, &data.Instruction, &data.StartDate, &data.EndDate,
			&data.StartTime, &data.EndTime,
			&data.MaximumParticipant, &data.InstagramLink, &data.Status, &data.DayPastEndDate, &data.BookingPrice,
			&data.CurrentUsedSlot,
			&data.GameMasterName, &data.GameMasterImageUrl,
			&data.GameCode, &data.GameName, &data.GameImgUrl,
		)

		if err != nil {
			return list, param, c.errHandler("model.GetRoomList", err, utils.ErrScanningListRoom)
		}
		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) GetRoomByCode(db *pgxpool.Pool, ctx context.Context, code string) (RoomEnt, error) {
	var (
		err  error
		data RoomEnt

		queryGetRoomGameDetail = `SELECT 
			a.admin_code AS game_master_code,
			a.name AS game_master_name,
			a.image_url AS game_master_img_url,
			r.game_id,
			g.game_code,
			g.name AS game_name,
			g.image_url AS game_img_url,
			c.cafe_code AS cafe_code,
			c.name AS cafe_name,
			c.address AS cafe_address,
			r.id as room_id,
			r.room_code,
			r.room_type,
			r.name,
			r.description,
			COALESCE(r.instruction, '') AS special_instruction,
			r.difficulty,
			r.start_date,
			r.end_date,
			r.start_time,
			r.end_time,
			r.maximum_participant,
			r.booking_price,
			r.reward_point,
			r.instagram_link,
			r.status,
			DATE_PART('day', NOW() - r.end_date) AS days_past_end_date, 
			COALESCE(r.image_url, '') AS room_banner_url,
			COALESCE(cp.count_participants, 0) as current_used_slot
		FROM rooms r 
			JOIN admins a ON r.game_master_id = a.id
			JOIN games g ON r.game_id = g.id 
			JOIN cafes c ON c.id = g.cafe_id  
			left join (
				select count(rp.user_id) count_participants, rp.room_id
				from rooms_participants rp
				WHERE rp.status != 'cancel'
				group by rp.room_id
			) as cp on cp.room_id = r.id
		WHERE room_code = $1 AND r.deleted_date IS NULL`
	)

	err = db.QueryRow(ctx, queryGetRoomGameDetail, code).Scan(
		&data.GameMasterCode, &data.GameMasterName, &data.GameMasterImgUrl,
		&data.GameId, &data.GameCode, &data.GameName, &data.GameImgUrl,
		&data.CafeCode, &data.CafeName, &data.CafeAddress,
		&data.RoomId, &data.RoomCode, &data.RoomType, &data.Name, &data.Description, &data.SpecialInstruction, &data.Difficulty,
		&data.StartDate, &data.EndDate, &data.StartTime, &data.EndTime,
		&data.MaximumParticipant, &data.BookingPrice, &data.RewardPoint, &data.InstagramLink, &data.Status, &data.DayPastEndDate, &data.BannerRoomUrl, &data.CurrentUsedSlot,
	)

	if err != nil {
		return data, c.errHandler("model.GetRoomByCode", err, utils.ErrGettingRoomByCode)
	}

	return data, nil
}

func (c *Contract) AddRoom(db *pgxpool.Pool, ctx context.Context, gameMasterId int64, gameId int64, roomCode, roomType, roomName, description string, startDate, endDate, startTime, endTime interface{}, bookingPrice float64, rewardPoint int, intagramLink, status, difficulty, instruction string, maximumParticipant int, imageUrl string, locationCity string) error {
	sql := `INSERT INTO rooms(
		game_master_id, game_id, room_code, room_type, "name", description, start_date, end_date, start_time, end_time, booking_price, reward_point, instagram_link, status, difficulty, instruction, maximum_participant, image_url, created_date, updated_date, location_city
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`

	_, err := db.Exec(ctx, sql, gameMasterId, gameId, roomCode, roomType, roomName, description, startDate, endDate, startTime, endTime, bookingPrice, rewardPoint, intagramLink, status, difficulty, instruction, maximumParticipant, imageUrl, time.Now().In(time.UTC), time.Now().In(time.UTC), locationCity)
	if err != nil {
		return c.errHandler("model.AddRoom", err, utils.ErrAddingRoom)
	}

	return nil
}

func (c *Contract) UpdateRoom(db *pgxpool.Pool, ctx context.Context, code string, gameMasterId int64, gameId int64, roomCode, roomType, roomName, description string, startDate, endDate, startTime, endTime interface{}, bookingPrice float64, rewardPoint int, instagramLink, status, difficulty, instruction string, maximumParticipant int, imageUrl string, locationCity string) error {
	var (
		err error
		sql = `
		UPDATE rooms 
		SET game_master_id = $1,
		    game_id = $2,
		    room_code = $3,
		    room_type = $4,
		    name = $5,
		    description = $6,
		    start_date = $7,
		    end_date = $8,
		    start_time = $9,
		    end_time = $10,
		    booking_price = $11,
		    reward_point = $12,
		    instagram_link = $13,
		    status = $14,
		    difficulty = $15,
		    instruction = $16,
		    maximum_participant = $17,
		    image_url = $18,
		    updated_date = $19,
		    location_city = $20
		WHERE room_code = $21`
	)
	_, err = db.Exec(ctx, sql, gameMasterId, gameId, roomCode, roomType, roomName, description, startDate, endDate, startTime, endTime, bookingPrice, rewardPoint, instagramLink, status, difficulty, instruction, maximumParticipant, imageUrl, time.Now().In(time.UTC), locationCity, code)
	if err != nil {
		return c.errHandler("model.UpdateRoom", err, utils.ErrUpdatingRoom)
	}

	return nil
}

func (c *Contract) SetWinnerRoom(tx pgx.Tx, ctx context.Context, roomCode, userCode string) error {
	var (
		err error
		sql = `
		update rooms_participants 
		set status_winner = true,
		reward_point = r.reward_point, 
		updated_date = $1
		from rooms_participants rp
		join rooms r on rp.room_id = r.id 
		join users u on rp.user_id = u.id 
		where r.room_code =$2 and u.user_code=$3`
	)
	_, err = tx.Exec(ctx, sql, time.Now().In(time.UTC), roomCode, userCode)
	if err != nil {
		return c.errHandler("model.SetWinnerRoom", err, utils.ErrUpdatingRoom)
	}

	return nil
}

func (c *Contract) ResetWinnerRoom(tx pgx.Tx, ctx context.Context, roomCode string) error {
	var (
		err error
		sql = `
		update rooms_participants 
		set status_winner = false,
		reward_point = 0, 
		updated_date = $1
		from rooms_participants rp
		join rooms r on rp.room_id = r.id 
		join users u on rp.user_id = u.id 
		where r.room_code =$2 `
	)
	_, err = tx.Exec(ctx, sql, time.Now().In(time.UTC), roomCode)
	if err != nil {
		return c.errHandler("model.SetWinnerRoom", err, utils.ErrUpdatingRoom)
	}

	return nil
}

func (c *Contract) UpdateRoomStatus(db *pgxpool.Pool, ctx context.Context, roomCode, status string) error {
	sql := `
		UPDATE rooms
		SET status = $1, updated_date = NOW()
		WHERE room_code = $2 AND deleted_date IS NULL
	`

	_, err := db.Exec(ctx, sql, status, roomCode)
	if err != nil {
		return c.errHandler("model.UpdateRoomStatus", err, utils.ErrUpdatingRoomStatus)
	}

	return nil
}

func (c *Contract) UpdateRoomStatusTrx(tx pgx.Tx, ctx context.Context, roomCode, status string) error {
	sql := `
		UPDATE rooms
		SET status = $1, updated_date = NOW()
		WHERE room_code = $2 AND deleted_date IS NULL
	`

	_, err := tx.Exec(ctx, sql, status, roomCode)
	if err != nil {
		return c.errHandler("model.UpdateRoomStatusTrx", err, utils.ErrUpdatingRoomStatus)
	}

	return nil
}

func (c *Contract) DeleteRoomByCode(db *pgxpool.Pool, ctx context.Context, roomCode string) error {
	sql := `
		UPDATE rooms
		SET updated_date = NOW(), deleted_date = NOW()
		WHERE room_code = $1
	`

	_, err := db.Exec(ctx, sql, roomCode)
	if err != nil {
		return c.errHandler("model.DeleteRoomByCode", err, utils.ErrDeletingRoom)
	}

	return nil
}

func (c *Contract) GetRoomIdByCode(db *pgxpool.Pool, ctx context.Context, code string) (int64, error) {
	var (
		err   error
		id    int64
		query = `SELECT id FROM rooms WHERE room_code=$1 AND deleted_date IS NULL`
	)
	err = db.QueryRow(ctx, query, code).Scan(&id)
	if err != nil {
		return id, c.errHandler("model.GetRoomIdByCode", err, utils.ErrGettingRoomByCode)
	}
	return id, nil
}
