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
	GameEnt struct {
		Id                 int64          `db:"id"`
		CafeId             string         `db:"cafe_id"`
		GameCode           string         `db:"game_code"`
		GameType           string         `db:"game_type"`
		Name               string         `db:"name"`
		ImageUrl           string         `db:"image_url"`
		CollectionUrl      string         `db:"collection_url"`
		Description        string         `db:"description"`
		Difficulty         sql.NullString `db:"difficulty"`
		Level              float64        `db:"level"`
		AdminCode          sql.NullString `db:"admin_code"`
		MinimalParticipant sql.NullInt64  `db:"minimal_participant"`
		MaximumParticipant int64          `db:"maximum_participant"`
		Status             string         `db:"status"`
		CreatedDate        time.Time      `db:"created_date"`
		UpdatedDate        sql.NullTime   `db:"updated_date"`
		DeletedDate        sql.NullTime   `db:"updated_date"`
	}

	GameResp struct {
		Id                 int64          `db:"id"`
		CafeCode           string         `db:"cafe_code"`
		CafeName           string         `db:"cafe_name"`
		CafeAddress        string         `db:"cafe_address"`
		Location           string         `db:"location"`
		GameCode           string         `db:"game_code"`
		GameType           string         `db:"game_type"`
		Name               string         `db:"name"`
		ImageUrl           string         `db:"image_url"`
		CollectionUrl      string         `db:"collection_url"`
		Description        string         `db:"description"`
		Duration           int64          `db:"duration"`
		Difficulty         sql.NullString `db:"difficulty"`
		Level              float64        `db:"level"`
		MinimalParticipant sql.NullInt64  `db:"minimal_participant"`
		MaximumParticipant int64          `db:"maximum_participant"`
		Status             string         `db:"status"`
		AdminCode          sql.NullString `db:"admin_code"`
		GameCategories     sql.NullString `db:"game_categories"`
		GameCharacteristic sql.NullString `db:"game_characteristics"`
		GameRelated        sql.NullString `db:"game_related_list"`
		GameRoomAvailables sql.NullString `db:"room_available_list"`
	}
)

func (c *Contract) GetGameList(db *pgxpool.Pool, ctx context.Context, param request.GameParam) ([]GameResp, request.GameParam, error) {
	var (
		err        error
		list       []GameResp
		where      []string
		paramQuery []interface{}
		totalData  int

		query = `SELECT 
		cafes.cafe_code, cafes.name as cafe_name, 
		games.game_code, games.game_type, games.name, games.image_url, 
		games.collection_url, games.description, games.status,
		games.duration, games.minimal_participant, games.maximum_participant,
		games.difficulty, games.level, admins.admin_code, 
		games_categories.categories,
		cafes.city AS location
		FROM games 
		LEFT JOIN cafes ON cafes.id = games.cafe_id
		LEFT JOIN admins ON admins.id = games.admin_id 
		LEFT JOIN (
			SELECT game_id, JSON_AGG(JSON_BUILD_OBJECT(
			'category_name', category_name 
			)) AS categories 
			FROM games_categories
			GROUP BY game_id
		) AS games_categories ON games_categories.game_id = games.id`
	)

	// Populate Search
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		// orWhere = append(orWhere, fmt.Sprintf("cafes.name iLIKE $%d", len(paramQuery)))
		orWhere = append(orWhere, fmt.Sprintf("games.name iLIKE $%d", len(paramQuery)))
		// orWhere = append(orWhere, fmt.Sprintf("games.game_type iLIKE $%d", len(paramQuery)))
		// orWhere = append(orWhere, fmt.Sprintf("games.description iLIKE $%d", len(paramQuery)))
		where = append(where, "("+strings.Join(orWhere, " OR ")+")")
	}
	if len(param.CafeCode) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.CafeCode)
		orWhere = append(orWhere, fmt.Sprintf("cafes.cafe_code = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}
	if len(param.Status) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Status)
		orWhere = append(orWhere, fmt.Sprintf("games.status = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}
	if len(param.GameType) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.GameType)
		orWhere = append(orWhere, fmt.Sprintf("games.game_type = ANY($%d)", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}
	if len(param.Difficulty) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Difficulty)
		orWhere = append(orWhere, fmt.Sprintf("games.difficulty = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}
	if len(param.GameCategoryName) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.GameCategoryName)
		orWhere = append(orWhere, fmt.Sprintf(`EXISTS (
			SELECT 1
			FROM json_array_elements(games_categories.categories) AS category
			WHERE lower(category->>'category_name') = ANY($%d)
		)`, len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}
	if len(param.Location) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Location)
		orWhere = append(orWhere, fmt.Sprintf("cafes.city = ANY($%d)", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	if param.Level > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Level)
		orWhere = append(orWhere, fmt.Sprintf("games.level = $%d ", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	if param.MinDuration > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.MinDuration)
		orWhere = append(orWhere, fmt.Sprintf("games.duration >= $%d ", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	if param.MaxDuration > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.MaxDuration)
		orWhere = append(orWhere, fmt.Sprintf("games.duration <= $%d ", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	if param.MinimalParticipant > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.MinimalParticipant)
		orWhere = append(orWhere, fmt.Sprintf("games.minimal_participant >= $%d ", len(paramQuery)))
		paramQuery = append(paramQuery, param.MinimalParticipant)
		orWhere = append(orWhere, fmt.Sprintf("games.maximum_participant >= $%d ", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	if param.MaximumParticipant > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.MaximumParticipant)
		orWhere = append(orWhere, fmt.Sprintf("games.maximum_participant <= $%d ", len(paramQuery)))
		paramQuery = append(paramQuery, param.MaximumParticipant)
		orWhere = append(orWhere, fmt.Sprintf("games.minimal_participant <= $%d ", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// Handling Soft Delete
	where = append(where, "games.deleted_date IS NULL")

	// Append All Where Conditions
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS data`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetGameList", err, utils.ErrCountingListCafe)
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
		return list, param, c.errHandler("model.GetGameList", err, utils.ErrGettingListCafe)
	}

	defer rows.Close()
	for rows.Next() {
		var data GameResp
		err = rows.Scan(&data.CafeCode, &data.CafeName, &data.GameCode, &data.GameType, &data.Name, &data.ImageUrl, &data.CollectionUrl, &data.Description, &data.Status, &data.Duration, &data.MinimalParticipant, &data.MaximumParticipant, &data.Difficulty, &data.Level, &data.AdminCode, &data.GameCategories, &data.Location)
		if err != nil {
			return list, param, c.errHandler("model.GetGameList", err, utils.ErrScanningListCafe)
		}
		list = append(list, data)
	}
	return list, param, nil
}

func (c *Contract) GetGameByCode(db *pgxpool.Pool, ctx context.Context, code string) (GameResp, error) {
	var (
		err  error
		data GameResp
		sql  = `SELECT 
		games.id, cafes.cafe_code, cafes.name as cafe_name, cafes.address as cafe_address,
		games.game_code, games.game_type, games.name, games.image_url, 
		games.collection_url, games.description, games.status,
		games.duration, games.minimal_participant, games.maximum_participant,
		games.difficulty, games.level, admins.admin_code,  
		games_categories.categories,
		games_related.game_related_list,
		game_room_available.room_available_list,
		cafes.city AS location
		FROM games 
		LEFT JOIN cafes ON cafes.id = games.cafe_id
		LEFT JOIN admins ON admins.id = games.admin_id 
		LEFT JOIN (
			SELECT game_id, JSON_AGG(JSON_BUILD_OBJECT(
			'category_name', category_name 
			)) AS categories 
			FROM games_categories
			GROUP BY game_id
		) AS games_categories ON games_categories.game_id = games.id
		LEFT JOIN (
			SELECT 
				g1.id,
				(SELECT JSON_AGG(JSON_BUILD_OBJECT(
				'game_id', g2.id,
				'name', g2.name,
				'game_code', g2.game_code,
				'game_type', g2.game_type,
				'difficulty', g2.difficulty,
				'image_url', g2.image_url,
				'minimal_participant', g2.minimal_participant,
				'maximum_participant', g2.maximum_participant,
				'duration', g2.duration,
				'location', c2.city
				)) AS game_related_list
				FROM games AS g2 
					JOIN cafes c2 ON c2.id = g2.cafe_id
				WHERE g2.game_type = g1.game_type AND g1.id <> g2.id) AS game_related_list
			FROM 
				games AS g1
		) AS games_related ON games_related.id = games.id
		left join (
			select r.game_id,  JSON_AGG(JSON_BUILD_OBJECT(
				'room_id', r.id,
				'room_code', r.room_code,
				'room_image_url', r.image_url,
				'cafe_name', c."name",
				'start_date', r.start_date,
				'end_date', r.end_date,
				'game_master_id', r.game_master_id,
				'game_master_name', a."name"
			)) AS room_available_list 
			from rooms r 
			left join games g on r.game_id = g.id 
			left join cafes c on g.cafe_id = c.id 
			left join admins a on r.game_master_id = a.id
			where r.start_date > NOW() and r.deleted_date is null
			group by r.game_id
		) as game_room_available on game_room_available.game_id=games.id
		WHERE games.game_code = $1`
	)

	err = db.QueryRow(ctx, sql, code).Scan(&data.Id, &data.CafeCode, &data.CafeName, &data.CafeAddress, &data.GameCode, &data.GameType, &data.Name, &data.ImageUrl, &data.CollectionUrl, &data.Description, &data.Status, &data.Duration, &data.MinimalParticipant, &data.MaximumParticipant, &data.Difficulty, &data.Level, &data.AdminCode, &data.GameCategories, &data.GameRelated, &data.GameRoomAvailables, &data.Location)
	if err != nil {
		return data, c.errHandler("model.GetGameByCode", err, utils.ErrGettingGameByCode)
	}

	return data, nil
}

func (c *Contract) AddGame(tx pgx.Tx, ctx context.Context, cafeId int64, code, gameType, name, imgUrl, collectionUrl, desc, difficulty, status string, level float64, minimalParticipant, maximumParticipant, duration, adminId int64) (int64, error) {
	var (
		err error
		id  int64
		sql = `INSERT INTO games(cafe_id, game_code, game_type, name, image_url, collection_url, difficulty, level, description, status, minimal_participant, maximum_participant, duration, admin_id, created_date)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`
	)

	err = tx.QueryRow(ctx, sql, cafeId, code, gameType, name, imgUrl, collectionUrl, difficulty, level, desc, status, minimalParticipant, maximumParticipant, duration, adminId, time.Now().In(time.UTC)).Scan(&id)
	if err != nil {
		return id, c.errHandler("model.AddGame", err, utils.ErrAddingCafe)
	}

	return id, nil
}

func (c *Contract) UpdateGameByCode(tx pgx.Tx, ctx context.Context, cafeId int64, code, gameType, name, imgUrl, collectionUrl, desc, difficulty, status string, level float64, minimalParticipant, maximumParticipant, adminId, duration int64) error {
	var (
		err error
		sql = `
		UPDATE games 
		SET cafe_id=$1, game_type=$2, name=$3, image_url=$4, collection_url=$5, description=$6, difficulty=$7, status=$8, level=$9, minimal_participant=$10, maximum_participant=$11, admin_id=$12, duration=$13, updated_date=$14
		WHERE game_code=$15`
	)

	_, err = tx.Exec(ctx, sql, cafeId, gameType, name, imgUrl, collectionUrl, desc, difficulty, status, level, minimalParticipant, maximumParticipant, adminId, duration, time.Now().In(time.UTC), code)
	if err != nil {
		return c.errHandler("model.UpdateGameByCode", err, utils.ErrUpdatingGame)
	}

	return nil
}

func (c *Contract) CheckExistGameUsed(db *pgxpool.Pool, ctx context.Context, id int64) (bool, error) {
	var exists bool
	sql := `
		SELECT EXISTS (
			SELECT 1 FROM rooms WHERE game_id = $1
			UNION
			SELECT 1 FROM tournaments WHERE game_id = $1
		)
	`

	err := db.QueryRow(ctx, sql, id).Scan(&exists)
	if err != nil {
		return false, c.errHandler("model.CheckExistGameUsed", err, utils.ErrCheckExistGameUsed)
	}

	return exists, nil
}

func (c *Contract) DeleteGameById(db *pgxpool.Pool, ctx context.Context, id int64) error {
	var (
		err error
		sql = `
		UPDATE games 
		SET updated_date=$1, deleted_date=$2 
		WHERE id=$3`
	)
	_, err = db.Exec(ctx, sql, time.Now().In(time.UTC), time.Now().In(time.UTC), id)
	if err != nil {
		return c.errHandler("model.DeleteGameByCode", err, utils.ErrUpdatingGame)
	}

	return nil
}

func (c *Contract) GetGameIdByCode(db *pgxpool.Pool, ctx context.Context, code string) (int64, error) {
	var (
		err   error
		id    int64
		query = `SELECT id FROM games WHERE game_code=$1`
	)
	err = db.QueryRow(ctx, query, code).Scan(&id)
	if err != nil {
		return id, c.errHandler("model.GetGameIdByCode", err, utils.ErrGettingGameByCode)
	}
	return id, nil
}
