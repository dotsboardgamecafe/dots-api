package model

import (
	"context"
	"database/sql"
	"dots-api/lib/qr"
	"dots-api/lib/utils"
	"dots-api/services/api/request"
	"fmt"
	"math"
	"os"
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
		NumberOfPopularity int64          `db:"number_of_popularity"`
	}
)

func (c *Contract) GetGameList(db *pgxpool.Pool, ctx context.Context, param request.GameParam) ([]GameResp, request.GameParam, error) {
	var (
		err        error
		list       []GameResp
		where      []string
		paramQuery []interface{}
		totalData  int

		query = `WITH 
		unique_participants AS (
				SELECT DISTINCT r.game_id, rp.user_id
				FROM rooms r 
				INNER JOIN rooms_participants rp ON r.id = rp.room_id AND rp.status = 'active'
				UNION
				SELECT DISTINCT t.game_id, tp.user_id
				FROM tournaments t
				INNER JOIN tournament_participants tp ON t.id = tp.tournament_id AND tp.status = 'active'
				UNION
				SELECT DISTINCT game_id, user_id FROM users_game_collections
		),		
		games_popularity AS (
			SELECT game_id, COUNT(DISTINCT user_id) as number_of_popularity
			FROM unique_participants
			GROUP BY game_id
		),
		game_collections AS (
			SELECT game_id, COUNT(distinct user_id) as number_of_collection
			FROM (SELECT DISTINCT ugc.game_id, ugc.user_id FROM users_game_collections ugc) unique_collections 
			GROUP BY game_id
		)
		SELECT 
			g.id,
			c.cafe_code, 
			c.name as cafe_name, 
			g.game_code, 
			g.game_type, 
			g.name, 
			g.image_url, 
			g.collection_url, 
			g.description, 
			g.status, 
			g.duration, 
			g.minimal_participant, 
			g.maximum_participant, 
			g.difficulty, 
			g.level, 
			a.admin_code, 
			COALESCE(gc.categories, '[]'::json) as categories,
			c.city AS location, 
			COALESCE(gp.number_of_popularity, 0) + coalesce(gcols.number_of_collection, 0) AS number_of_popularity
		FROM games g
		LEFT JOIN games_popularity gp ON gp.game_id = g.id
		LEFT JOIN game_collections gcols on gcols.game_id = g.id
		LEFT JOIN cafes c ON c.id = g.cafe_id
		LEFT JOIN admins a ON a.id = g.admin_id 
		LEFT JOIN LATERAL (
			SELECT 
				game_id, 
				JSON_AGG(JSON_BUILD_OBJECT('category_name', category_name)) AS categories 
			FROM games_categories 
			WHERE game_id = g.id 
			GROUP BY game_id
		) gc ON true
		`
	)

	// Populate Search
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		// orWhere = append(orWhere, fmt.Sprintf("cafes.name iLIKE $%d", len(paramQuery)))
		orWhere = append(orWhere, fmt.Sprintf("g.name iLIKE $%d", len(paramQuery)))
		// orWhere = append(orWhere, fmt.Sprintf("games.game_type iLIKE $%d", len(paramQuery)))
		// orWhere = append(orWhere, fmt.Sprintf("games.description iLIKE $%d", len(paramQuery)))
		where = append(where, "("+strings.Join(orWhere, " OR ")+")")
	}
	if len(param.CafeCode) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.CafeCode)
		orWhere = append(orWhere, fmt.Sprintf("c.cafe_code = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}
	if len(param.Status) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Status)
		orWhere = append(orWhere, fmt.Sprintf("g.status = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}
	if len(param.GameType) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.GameType)
		orWhere = append(orWhere, fmt.Sprintf("g.game_type = ANY($%d)", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}
	if len(param.Difficulty) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Difficulty)
		orWhere = append(orWhere, fmt.Sprintf("g.difficulty = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}
	if len(param.GameCategoryName) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.GameCategoryName)
		orWhere = append(orWhere, fmt.Sprintf(`EXISTS (
			SELECT 1
			FROM json_array_elements(gc.categories) AS category
			WHERE lower(category->>'category_name') = ANY($%d)
		)`, len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}
	if len(param.Location) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Location)
		orWhere = append(orWhere, fmt.Sprintf("c.city = ANY($%d)", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	if param.Level > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Level)
		orWhere = append(orWhere, fmt.Sprintf("g.level = $%d ", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	if param.MinDuration > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.MinDuration)
		orWhere = append(orWhere, fmt.Sprintf("g.duration >= $%d ", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	if param.MaxDuration > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.MaxDuration)
		orWhere = append(orWhere, fmt.Sprintf("g.duration <= $%d ", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	if param.NumberOfPlayers > 0 {
		paramQuery = append(paramQuery, param.NumberOfPlayers)
		where = append(where, fmt.Sprintf("$%[1]d >= g.minimal_participant AND g.maximum_participant >= $%[1]d", len(paramQuery)))
	}

	// Handling Soft Delete
	where = append(where, "g.deleted_date IS NULL")

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
	query += " ORDER BY " + param.SortKey + " " + param.Sort + " "

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
		err = rows.Scan(&data.Id, &data.CafeCode, &data.CafeName, &data.GameCode, &data.GameType, &data.Name, &data.ImageUrl, &data.CollectionUrl, &data.Description, &data.Status, &data.Duration, &data.MinimalParticipant, &data.MaximumParticipant, &data.Difficulty, &data.Level, &data.AdminCode, &data.GameCategories, &data.Location, &data.NumberOfPopularity)
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
		sql  = `WITH 
		unique_participants AS (
				SELECT DISTINCT r.game_id, rp.user_id
				FROM rooms r 
				INNER JOIN rooms_participants rp ON r.id = rp.room_id AND rp.status = 'active'
				UNION
				SELECT DISTINCT t.game_id, tp.user_id
				FROM tournaments t
				INNER JOIN tournament_participants tp ON t.id = tp.tournament_id AND tp.status = 'active'
				UNION
				SELECT DISTINCT game_id, user_id
				FROM users_game_collections
		),		
		games_popularity AS (
			SELECT game_id, COUNT(DISTINCT user_id) as number_of_popularity
			FROM unique_participants
			GROUP BY game_id
		),
		game_collections AS (
			SELECT game_id, COUNT(distinct user_id) as number_of_collection
			FROM (SELECT DISTINCT ugc.game_id, ugc.user_id FROM users_game_collections ugc) unique_collections 
			GROUP BY game_id
		)
		SELECT 
			games.id, cafes.cafe_code, cafes.name as cafe_name, cafes.address as cafe_address,
			games.game_code, games.game_type, games.name, games.image_url, 
			games.collection_url, games.description, games.status, 
			COALESCE(gp.number_of_popularity, 0) + coalesce(gcols.number_of_collection, 0) AS number_of_popularity,
			games.duration, games.minimal_participant, games.maximum_participant,
			games.difficulty, games.level, admins.admin_code,  
			games_categories.categories,
			games_related.game_related_list,
			game_room_available.room_available_list,
			cafes.city AS location
		FROM games 
		LEFT JOIN games_popularity gp ON gp.game_id = games.id
		LEFT JOIN game_collections gcols on gcols.game_id = games.id
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
			SELECT g1.id,
			(
				SELECT JSON_AGG(JSON_BUILD_OBJECT(
					'game_id', g2.id,
					'name', g2.name,
					'game_code', g2.game_code,
					'game_type', g2.game_type,
					'level', g2.level,
					'difficulty', g2.difficulty,
					'image_url', g2.image_url,
					'minimal_participant', g2.minimal_participant,
					'maximum_participant', g2.maximum_participant,
					'duration', g2.duration,
					'location', c2.city,
					'categories', (
						SELECT JSON_AGG(JSON_BUILD_OBJECT('category_name', gc2.category_name))
						FROM games_categories gc2
						WHERE gc2.game_id = g2.id
					)
				))
				FROM games g2
					JOIN cafes c2 ON c2.id = g2.cafe_id
				WHERE g2.id <> g1.id
				AND EXISTS (
					SELECT 1 
					FROM games_categories gc1
					WHERE gc1.game_id = g1.id
					AND EXISTS (
						SELECT 1 
						FROM games_categories gc2
						WHERE gc2.game_id = g2.id
						AND gc2.category_name = gc1.category_name
					)
				)
				AND g2.deleted_date IS NULL
			) AS game_related_list
			FROM games g1
		) AS games_related ON games_related.id = games.id
		LEFT JOIN (
			SELECT r.game_id,
			JSON_AGG(JSON_BUILD_OBJECT(
				'room_id', r.id,
				'room_code', r.room_code,
				'room_name', r.name,
				'status', r.status,
				'start_date', r.start_date,
				'end_date', r.end_date,
				'maximum_participant', r.maximum_participant,
				'current_participant', (
					SELECT COUNT(*) 
					FROM rooms_participants rp 
					WHERE rp.room_id = r.id AND rp.status = 'active'
				)
			)) AS room_available_list
			FROM rooms r
			WHERE r.status = 'open'
			GROUP BY r.game_id
		) AS game_room_available ON game_room_available.game_id = games.id
		WHERE games.game_code = $1 AND games.deleted_date IS NULL`
	)

	err = db.QueryRow(ctx, sql, code).Scan(&data.Id, &data.CafeCode, &data.CafeName, &data.CafeAddress, &data.GameCode, &data.GameType, &data.Name, &data.ImageUrl, &data.CollectionUrl, &data.Description, &data.Status, &data.NumberOfPopularity, &data.Duration, &data.MinimalParticipant, &data.MaximumParticipant, &data.Difficulty, &data.Level, &data.AdminCode, &data.GameCategories, &data.GameRelated, &data.GameRoomAvailables, &data.Location)
	if err != nil {
		return data, c.errHandler("model.GetGameByCode", err, utils.ErrGettingGameByCode)
	}

	return data, nil
}

func (c *Contract) AddGame(tx pgx.Tx, ctx context.Context, cafeId int64, code, gameType, name, imgUrl, collectionUrl, desc, difficulty, status string, level float64, minimalParticipant, maximumParticipant, duration int64) (int64, error) {
	var (
		err error
		id  int64
		sql = `INSERT INTO games(cafe_id, game_code, game_type, name, image_url, collection_url, difficulty, level, description, status, minimal_participant, maximum_participant, duration, created_date)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`
	)

	err = tx.QueryRow(ctx, sql, cafeId, code, gameType, name, imgUrl, collectionUrl, difficulty, level, desc, status, minimalParticipant, maximumParticipant, duration, time.Now().In(time.UTC)).Scan(&id)
	if err != nil {
		return id, c.errHandler("model.AddGame", err, utils.ErrAddingCafe)
	}

	return id, nil
}

func (c *Contract) UpdateGameByCode(tx pgx.Tx, ctx context.Context, cafeId int64, code, gameType, name, imgUrl, collectionUrl, desc, difficulty, status string, level float64, minimalParticipant, maximumParticipant, duration int64) error {
	var (
		err error
		sql = `
		UPDATE games 
		SET cafe_id=$1, game_type=$2, name=$3, image_url=$4, collection_url=$5, description=$6, difficulty=$7, status=$8, level=$9, minimal_participant=$10, maximum_participant=$11, duration=$12, updated_date=$13
		WHERE game_code=$14`
	)

	_, err = tx.Exec(ctx, sql, cafeId, gameType, name, imgUrl, collectionUrl, desc, difficulty, status, level, minimalParticipant, maximumParticipant, duration, time.Now().In(time.UTC), code)
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

func (c *Contract) GetGameQRCodeByCode(ctx context.Context, code string) (string, error) {
	fileName := code + "_QR"
	qrString, err := utils.ImageToBase64(c.Config.GetString("upload_path") + "/" + fileName)

	if err != nil {
		if !os.IsNotExist(err) {
			return "", c.errHandler("model.AddGameQRCode", err, utils.ErrGettingGameQrCode)
		}
		fileName, err = qr.GenerateQRCode(code, fileName, c.Config.GetString("upload_path"))
		if err != nil {
			return "", c.errHandler("model.AddGameQRCode", err, utils.ErrGettingGameQrCode)
		}

		qrString, err = utils.ImageToBase64(c.Config.GetString("upload_path") + "/" + fileName)
		if err != nil {
			return "", c.errHandler("model.AddGameQRCode", err, utils.ErrGettingGameQrCode)
		}
	}

	return qrString, nil
}
