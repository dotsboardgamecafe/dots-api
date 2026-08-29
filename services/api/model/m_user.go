package model

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"dots-api/lib/utils"
	"dots-api/services/api/request"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type OriginUserEnt struct {
	ID                 int            `db:"id"`
	UserCode           string         `db:"user_code"`
	Email              string         `db:"email"`
	DateOfBirth        sql.NullString `db:"date_of_birth"`
	Gender             sql.NullString `db:"gender"`
	UserName           sql.NullString `db:"username"`
	PhoneNumber        string         `db:"phone_number"`
	FullName           string         `db:"fullname"`
	ImageURL           sql.NullString `db:"image_url"`
	LatestPoint        int            `db:"latest_point"`
	LatestTierId       int            `db:"latest_tier_id"`
	Password           string         `db:"password"`
	XPlayer            string         `db:"x_player"`
	StatusVerification bool           `db:"status_verification"`
	Status             string         `db:"status"`
	CreatedDate        time.Time      `db:"created_date"`
	UpdatedDate        sql.NullTime   `db:"updated_date"`
	DeletedDate        sql.NullTime   `db:"deleted_date"`
}

type UserEnt struct {
	ID                 int            `db:"id"`
	UserCode           string         `db:"user_code"`
	Email              sql.NullString `db:"email"`
	DateOfBirth        sql.NullString `db:"date_of_birth"`
	Gender             sql.NullString `db:"gender"`
	UserName           sql.NullString `db:"username"`
	PhoneNumber        string         `db:"phone_number"`
	FullName           string         `db:"fullname"`
	ImageURL           sql.NullString `db:"image_url"`
	LatestPoint        int            `db:"latest_point"`
	RoomNormalCount    sql.NullInt64  `db:"room_normal_count"`
	RoomEventCount     sql.NullInt64  `db:"room_event_count"`
	TournamentCount    sql.NullInt64  `db:"tournament_count"`
	BadgeCount         sql.NullInt64  `db:"badge_count"`
	GameCount          sql.NullInt64  `db:"game_count"`
	LatestTierId       int            `db:"latest_tier_id"`
	LatestTierName     string         `db:"latest_tier_name"`
	TierMinRangePoint  int            `db:"tier_min_range_point"`
	TierMaxRangePoint  int            `db:"tier_max_range_point"`
	Password           string         `db:"password"`
	XPlayer            string         `db:"x_player"`
	RoleId             int            `db:"role_id"`
	StatusVerification bool           `db:"status_verification"`
	Status             string         `db:"status"`
	UserStyle          UserStyle      `db:"styles"`
	TotalSpent         int            `db:"total_spent"`
	CreatedDate        time.Time      `db:"created_date"`
	UpdatedDate        sql.NullTime   `db:"updated_date"`
	DeletedDate        sql.NullTime   `db:"deleted_date"`
}

type UserStyle struct {
	Color string `db:"color" json:"color"`
}

func (s *UserStyle) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *UserStyle) Scan(value interface{}) error {
	if value == nil {
		*s = UserStyle{}
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if ok {
			return json.Unmarshal([]byte(str), s)
		}

		return errors.New("m.UserStyle.Scan: " + utils.ErrTypeAssertionFailed + " expected []byte, got " + fmt.Sprintf("%T", value))
	}

	return json.Unmarshal(b, s)
}

func (c *Contract) GetAllUsers(db *pgxpool.Pool, ctx context.Context) ([]UserEnt, error) {
	var users []UserEnt

	query := `
		SELECT id, user_code, email, date_of_birth, gender, username, phone_number, fullname, image_url, 
		       latest_point, latest_tier_id, password, x_player, status_verification, 
		       status, styles, created_date, updated_date, deleted_date
		FROM users
	`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, c.errHandler("model.GetAllUsers", err, utils.ErrGettingListUser)
	}
	defer rows.Close()

	for rows.Next() {
		var user UserEnt
		err := rows.Scan(
			&user.ID,
			&user.UserCode,
			&user.Email,
			&user.DateOfBirth,
			&user.Gender,
			&user.UserName,
			&user.PhoneNumber,
			&user.FullName,
			&user.ImageURL,
			&user.LatestPoint,
			&user.LatestTierId,
			&user.Password,
			&user.XPlayer,
			&user.StatusVerification,
			&user.Status,
			&user.UserStyle,
			&user.CreatedDate,
			&user.UpdatedDate,
			&user.DeletedDate,
		)
		if err != nil {
			return nil, c.errHandler("model.GetAllUsers", err, utils.ErrScanningListUser)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, c.errHandler("model.GetAllUsers", err, utils.ErrGettingListUser)
	}

	return users, nil
}

func (c *Contract) GetUserByEmail(db *pgxpool.Pool, ctx context.Context, email string) (UserEnt, error) {
	var res UserEnt

	sql := `SELECT users.id AS id,
						user_code,
						email,
						date_of_birth,
						gender,
						username,
						phone_number,
						fullname,
						image_url,
						latest_point,
						tiers.name AS latest_tier_name,
						password,
						x_player,
						role_id,
						status_verification,
						users.status,
						users.styles,
						users.created_date AS created_date,
						users.updated_date AS updated_date,
						users.deleted_date AS deleted_date,
						users.role_id
					FROM users
					JOIN tiers ON tiers.id = users.latest_tier_id
					JOIN roles r ON r.id = users.role_id 
					WHERE email = $1 AND users.deleted_date IS NULL AND r.deleted_date IS NULL`

	err := db.QueryRow(ctx, sql, strings.ToLower(strings.TrimSpace(email))).Scan(
		&res.ID,
		&res.UserCode,
		&res.Email,
		&res.DateOfBirth,
		&res.Gender,
		&res.UserName,
		&res.PhoneNumber,
		&res.FullName,
		&res.ImageURL,
		&res.LatestPoint,
		&res.LatestTierName,
		&res.Password,
		&res.XPlayer,
		&res.RoleId,
		&res.StatusVerification,
		&res.Status,
		&res.UserStyle,
		&res.CreatedDate,
		&res.UpdatedDate,
		&res.DeletedDate,
		&res.RoleId,
	)

	if err != nil {
		return res, c.errHandler("model.GetUserByEmail", err, utils.ErrGettingUserByEmail)
	}

	return res, nil
}

func (c *Contract) GetUserByUserCode(db *pgxpool.Pool, ctx context.Context, userCode string) (UserEnt, error) {
	var res UserEnt

	// Scan userId
	var userId int
	db.QueryRow(ctx, "SELECT id FROM users WHERE user_code = $1", userCode).Scan(&userId)

	sql := `SELECT users.id AS id,
						user_code,
						email,
						date_of_birth,
						gender,
						username,
						phone_number,
						fullname,
						image_url,
						latest_point,
						tiers.id AS latest_tier_id,
						tiers.name AS latest_tier_name,
						tiers.min_point AS tier_min_range_point,
						tiers.max_point AS tier_max_range_point,
						password,
						x_player,
						status_verification,
						users.status,
						users.styles,
						users.created_date AS created_date,
						users.updated_date AS updated_date,
						users.deleted_date AS deleted_date,
						(COALESCE(redeem_histories.total_redeemed_amount, 0) + COALESCE(booking_transactions.total_booking_price, 0)) AS total_spent
					FROM users
					JOIN tiers ON tiers.id = users.latest_tier_id
					LEFT JOIN (
						SELECT user_id, SUM(invoice_amount) AS total_redeemed_amount 
						FROM user_redeem_histories 
						WHERE user_id = $1
						GROUP BY user_id
					) AS redeem_histories ON users.id = redeem_histories.user_id
					LEFT JOIN (
						SELECT user_id, SUM(price) AS total_booking_price
						FROM users_transactions 
						WHERE user_id = $1 AND status = 'PAID'
						GROUP BY user_id
					) AS booking_transactions ON users.id = booking_transactions.user_id
					WHERE user_code = $2 AND users.deleted_date IS NULL`

	err := db.QueryRow(ctx, sql, userId, userCode).Scan(
		&res.ID,
		&res.UserCode,
		&res.Email,
		&res.DateOfBirth,
		&res.Gender,
		&res.UserName,
		&res.PhoneNumber,
		&res.FullName,
		&res.ImageURL,
		&res.LatestPoint,
		&res.LatestTierId,
		&res.LatestTierName,
		&res.TierMinRangePoint,
		&res.TierMaxRangePoint,
		&res.Password,
		&res.XPlayer,
		&res.StatusVerification,
		&res.Status,
		&res.UserStyle,
		&res.CreatedDate,
		&res.UpdatedDate,
		&res.DeletedDate,
		&res.TotalSpent,
	)

	if err != nil && err != pgx.ErrNoRows {
		return res, c.errHandler("model.GetUserByUserCode", err, utils.ErrRetrievingUserByUserIdentifier)
	}

	return res, nil
}

// CMS member page
func (c *Contract) GetUserList(db *pgxpool.Pool, ctx context.Context, param request.UserParam) ([]UserEnt, request.UserParam, error) {
	var (
		err        error
		list       []UserEnt
		where      []string
		paramQuery []interface{}
		totalData  int

		query = `
		with users_rooms_participation as (
		select user_id, count(r.room_type = 'normal' or NULL) as normal_count, count(r.room_type = 'special_event' or NULL) event_count  from rooms r
			join rooms_participants rp on r.id = rp.room_id 
			where
				r.status = 'inactive' and
				r.deleted_date is null and 
				rp.status = 'active'
			group by user_id
		),
		users_tournaments_participation as (
			select user_id, count(*) as tournament_count  from tournaments t
				join tournament_participants tp on t.id = tp.tournament_id
				where
					t.status = 'inactive' and
					t.deleted_date is null and 
					tp.status = 'active'
				group by user_id
		),
		users_badges_owned as (
			select ub.user_id, count(*) badge_count from badges b
			join users_badges ub on b.id = ub.badge_id and b.status = 'active' and b.deleted_date is null
			group by ub.user_id
		),
		users_game_boards_collection as (
			select user_id, count(*) game_count from users_game_collections
			group by user_id
		)
		SELECT users.id AS id,
						user_code,
						email,
						date_of_birth,
						gender,
						username,
						phone_number,
						fullname,
						image_url,
						latest_point,
						coalesce(urp.normal_count, 0) as room_normal_count,
						coalesce(urp.event_count, 0) as room_event_count,
						coalesce(utp.tournament_count, 0) as tournament_count,
						coalesce(ubo.badge_count, 0) as badge_count,
						coalesce(ugc.game_count, 0) as game_count,
						tiers.id AS latest_tier_id,
						tiers.name AS latest_tier_name,
						password,
						x_player,
						status_verification,
						users.status,
						users.styles,
						users.created_date AS created_date,
						users.updated_date AS updated_date,
						users.deleted_date AS deleted_date,
						(COALESCE(redeem_histories.total_redeemed_amount, 0) + COALESCE(booking_transactions.total_booking_price, 0)) AS total_spent
					FROM users
					JOIN tiers ON tiers.id = users.latest_tier_id
					LEFT JOIN (
						SELECT user_id, SUM(invoice_amount) AS total_redeemed_amount FROM user_redeem_histories GROUP BY user_id
					) AS redeem_histories ON users.id = redeem_histories.user_id
					LEFT JOIN (
						SELECT user_id, SUM(price) AS total_booking_price FROM users_transactions WHERE status = 'PAID' GROUP BY user_id
					) AS booking_transactions ON users.id = booking_transactions.user_id
					LEFT JOIN users_rooms_participation as urp ON users.id = urp.user_id
					LEFT JOIN users_tournaments_participation as utp ON users.id = utp.user_id
					LEFT JOIN users_badges_owned as ubo ON users.id = ubo.user_id
					LEFT JOIN users_game_boards_collection as ugc ON users.id = ugc.user_id`
	)

	// Populate Search
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		orWhere = append(orWhere, fmt.Sprintf("fullname iLIKE $%d", len(paramQuery)))
		orWhere = append(orWhere, fmt.Sprintf("username iLIKE $%d", len(paramQuery)))
		where = append(where, "("+strings.Join(orWhere, " OR ")+")")
	}

	if len(param.LatestTier) > 0 {
		var orWhere []string
		lowerString := make([]string, len(param.LatestTier))

		for i, s := range param.LatestTier {
			lowerString[i] = strings.ToLower(s)
		}

		paramQuery = append(paramQuery, lowerString)
		orWhere = append(orWhere, fmt.Sprintf(`lower(tiers.name) = ANY($%d::text[])`, len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	if len(param.Status) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Status)
		orWhere = append(orWhere, fmt.Sprintf("users.status = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// Handling Soft Delete
	where = append(where, "users.deleted_date IS NULL")

	// Append All Where Conditions
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS data`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetUserList", err, utils.ErrCountingListUser)
		}
		param.Count = totalData
	}

	// Select Max Page
	if param.Count > param.Limit && param.Page > int(param.Count/param.Limit) {
		param.Page = int(math.Ceil(float64(param.Count) / float64(param.Limit)))
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
		return list, param, c.errHandler("model.GetUserList", err, utils.ErrGettingListUser)
	}

	defer rows.Close()
	for rows.Next() {
		var res UserEnt
		err = rows.Scan(
			&res.ID,
			&res.UserCode,
			&res.Email,
			&res.DateOfBirth,
			&res.Gender,
			&res.UserName,
			&res.PhoneNumber,
			&res.FullName,
			&res.ImageURL,
			&res.LatestPoint,
			&res.RoomNormalCount,
			&res.RoomEventCount,
			&res.TournamentCount,
			&res.BadgeCount,
			&res.GameCount,
			&res.LatestTierId,
			&res.LatestTierName,
			&res.Password,
			&res.XPlayer,
			&res.StatusVerification,
			&res.Status,
			&res.UserStyle,
			&res.CreatedDate,
			&res.UpdatedDate,
			&res.DeletedDate,
			&res.TotalSpent,
		)

		if err != nil {
			return list, param, c.errHandler("model.GetUserList", err, utils.ErrScanningListAdmin)
		}
		list = append(list, res)
	}

	return list, param, nil
}

func (c *Contract) UpdateUserProfile(db *pgxpool.Pool, ctx context.Context, userCode, fullName, userName, gender, dateOfBirth, imageURL string, phoneNumber string) error {
	sql := `
		UPDATE users
		SET image_url = $1, fullname = $2, username = $3, gender = $4, date_of_birth = $5, phone_number = $6, updated_date = $7
		WHERE user_code = $8
	`

	currentTime := time.Now().In(time.UTC)

	_, err := db.Exec(ctx, sql, imageURL, fullName, userName, gender, dateOfBirth, phoneNumber, currentTime, userCode)
	if err != nil {
		return c.errHandler("model.UpdateUserProfile", err, utils.ErrUpdatingUserProfile)
	}

	return nil
}

func (c *Contract) UpdateUserEmail(db *pgxpool.Pool, ctx context.Context, userIdentifier, email string) error {
	sql := `
		UPDATE users
		SET email = $1
		WHERE user_code = $2
	`

	_, err := db.Exec(ctx, sql, strings.ToLower(strings.TrimSpace(email)), userIdentifier)
	if err != nil {
		return c.errHandler("model.UpdateUserEmail", err, utils.ErrUpdatingUserEmail)
	}

	return nil
}

func (c *Contract) UpdateUserXPlayer(db *pgxpool.Pool, ctx context.Context, userIdentifier, xPlayer string) error {
	sql := `
		UPDATE users
		SET x_player = $1
		WHERE user_code = $2
	`

	_, err := db.Exec(ctx, sql, xPlayer, userIdentifier)
	if err != nil {
		return c.errHandler("model.UpdateUserXPlayer", err, utils.ErrUpdatingUserXPlayer)
	}

	return nil
}

func (c *Contract) UpdatePasswordUser(db *pgxpool.Pool, ctx context.Context, userIdentifier, OldPassword, NewPassword, ConfirmPassword string) error {
	var (
		err      error
		dataUser UserEnt
	)
	// Check if new password matches the confirmation
	if NewPassword != ConfirmPassword {
		return errors.New(utils.ErrPasswordMismatch)
	}

	dataUser, err = c.GetUserByUserCode(db, ctx, userIdentifier)
	if err != nil {
		return c.errHandler("model.UpdatePasswordUser", err, utils.ErrFetchingUserPassword)
	}

	// Validate old password
	err = bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte(OldPassword))
	if err != nil {
		return errors.New("old password is incorrect")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.errHandler("model.UpdatePasswordUser", err, utils.ErrHashingPassword)
	}

	// Update the user's password in the database
	sql := "UPDATE users SET password = $1, updated_date = $3 WHERE id = $2"
	_, err = db.Exec(ctx, sql, string(hashedPassword), dataUser.ID, time.Now().UTC())
	if err != nil {
		return c.errHandler("model.UpdatePasswordUser", err, utils.ErrUpdatingUserPassword)
	}

	return nil
}

func (c *Contract) UpdateUserStatus(db *pgxpool.Pool, ctx context.Context, userCode, status string) error {
	sql := `
		UPDATE users
		SET status = $1
		WHERE user_code = $2
	`

	_, err := db.Exec(ctx, sql, status, userCode)
	if err != nil {
		return c.errHandler("model.UpdateUserStatus", err, utils.ErrUpdatingUserStatus)
	}

	return nil
}

func (c *Contract) UpdateUser(db *pgxpool.Pool, ctx context.Context, userCode, fullName, dateOfBirth, gender, imageURL, phoneNumber, email, userName, status string) error {
	sql := `
		UPDATE users
		SET image_url = $1, fullname = $2, gender = $3, date_of_birth = $4, phone_number = $5, email = $6, username = $7, status = $8, updated_date = $9
		WHERE user_code = $10
	`

	currentTime := time.Now().In(time.UTC)
	_, err := db.Exec(ctx, sql, imageURL, fullName, gender, dateOfBirth, phoneNumber, strings.ToLower(strings.TrimSpace(email)), userName, status, currentTime, userCode)
	if err != nil {
		return c.errHandler("model.UpdateUserProfile", err, utils.ErrUpdatingUserProfile)
	}

	return nil
}

func (c *Contract) GetLatestPointAndTier(tx pgx.Tx, ctx context.Context, userId int64) (string, int, int64, error) {
	var (
		UserCode     string
		LatestPoint  int
		LatestTierId int64

		sql = `SELECT user_code, latest_point, latest_tier_id FROM users WHERE id = $1;`
	)

	_ = tx.QueryRow(ctx, sql, userId).Scan(&UserCode, &LatestPoint, &LatestTierId)

	return UserCode, LatestPoint, LatestTierId, nil
}

func (c *Contract) GetPlayerAndOtherActivities(db *pgxpool.Pool, ctx context.Context, UserCode string) ([]UserPointEnt, error) {
	var (
		err  error
		list []UserPointEnt

		query = `
		SELECT
    		u.id, 
			u.user_code,
			u.image_url,
			u.username,
			u.styles,
			data_source, 
			source_code,
			COALESCE(t.name, r.name, b.name, tiers.name, g.name, '') AS item_name,
			COALESCE(t.tournament_code, r.room_code, b.badge_code, tiers.tier_code, g.game_code, '') AS item_code,
    		COALESCE(t.image_url, r.image_url, b.image_url, g.image_url, '') AS item_url,
			point, 
			up.created_date
    	FROM users_points up 
			JOIN users u ON up.user_id = u.id AND deleted_date IS NULL
			LEFT JOIN tournaments t ON t.tournament_code = up.source_code
			LEFT JOIN rooms r ON r.room_code = up.source_code
			LEFT JOIN badges b ON b.badge_code = up.source_code
			LEFT JOIN tiers ON tiers.tier_code = up.source_code
			LEFT JOIN users_game_collections ugc ON ugc.user_id = up.user_id 
				AND ugc.game_id = (SELECT id FROM games WHERE game_code = up.source_code)
			LEFT JOIN games g ON g.id = r.game_id OR g.id = t.game_id OR g.id = ugc.game_id 
    	WHERE up.data_source IN ('room', 'room_play', 'tournament', 'tournament_play', 'badge', 'tier', 'game')
			ORDER BY up.id DESC
			LIMIT 5;`
	)

	rows, err := db.Query(ctx, query)
	if err != nil {
		return list, c.errHandler("model.GetPlayerAndOtherActivities", err, utils.ErrGetPlayerAndOtherActivities)
	}

	defer rows.Close()
	for rows.Next() {
		var data UserPointEnt
		err = rows.Scan(
			&data.Id,
			&data.UserCode,
			&data.UserImageUrl,
			&data.UserName,
			&data.UserStyle,
			&data.DataSource,
			&data.SourceCode,
			&data.ItemName,
			&data.ItemCode,
			&data.ItemImgUrl,
			&data.Point,
			&data.CreatedDate,
		)
		if err != nil {
			return list, c.errHandler("model.GetPlayerAndOtherActivities", err, utils.ErrScanPlayerAndOtherActivities)
		}
		list = append(list, data)
	}

	return list, nil
}

func (c *Contract) GetUserPointActivities(db *pgxpool.Pool, ctx context.Context, UserCode string) ([]UserPointEnt, error) {
	var (
		err  error
		list []UserPointEnt

		query = `
			SELECT
				id,
				username,
				user_style,
				title_description,
				data_source, 
				source_code,
				point, 
				created_date
			FROM (
			SELECT
				u.id,
				u.username,
				COALESCE(u.styles, '{}') AS user_style,
				CASE
				WHEN up.data_source IN ('room','room_play')
					THEN 'Joined: ' || r.name
				WHEN up.data_source = 'room_paid'
					THEN 'Paid: ' || rr.name
				WHEN up.data_source = 'tournament'
					THEN 'Joined: ' || t.name
				WHEN up.data_source = 'tournament_play'
					THEN 'Participated: ' || t.name
				WHEN up.data_source = 'tournament_paid'
					THEN 'Paid: ' || tt.name
				WHEN up.data_source = 'redeem'
					THEN 'Purchased: ' || urh.description
				WHEN up.data_source = 'badge'
					THEN 'Claimed: ' || b.name
				WHEN up.data_source = 'game'
					THEN 'Collected: ' || g.name
				WHEN up.data_source = 'profile'
					THEN 'Profile: Updated Profile Image for the first time'
				WHEN up.data_source = 'point_subtract'
					THEN 'Point: ' || COALESCE(up.description, 'Point Deducted by System')
				WHEN up.data_source = 'point_add'
					THEN 'Point: ' || COALESCE(up.description, 'Point Added by System')
				END AS title_description,
				up.data_source, 
				up.source_code,
				up.point, 
				up.created_date
			FROM users_points up
			JOIN users u
				ON u.id = up.user_id

			-- rooms (room, room_play)
			LEFT JOIN rooms r
					ON up.data_source IN ('room','room_play')
					AND r.room_code = up.source_code

			-- room_paid -> users_transactions -> rooms
			LEFT JOIN users_transactions ut_room_paid
					ON up.data_source = 'room_paid'
					AND ut_room_paid.user_id = up.user_id
					AND ut_room_paid.transaction_code = up.source_code
			LEFT JOIN rooms rr
					ON rr.room_code = ut_room_paid.source_code

			-- tournament / tournament_play
			LEFT JOIN tournaments t
					ON up.data_source IN ('tournament','tournament_play')
					AND t.tournament_code = up.source_code

			-- tournament_paid -> users_transactions -> tournaments
			LEFT JOIN users_transactions ut_tour_paid
					ON up.data_source = 'tournament_paid'
					AND ut_tour_paid.user_id = up.user_id
					AND ut_tour_paid.transaction_code = up.source_code
			LEFT JOIN tournaments tt
					ON tt.tournament_code = ut_tour_paid.source_code

			-- redeem
			LEFT JOIN user_redeem_histories urh
					ON up.data_source = 'redeem'
					AND urh.custom_id = up.source_code

			-- badge
			LEFT JOIN badges b
					ON up.data_source = 'badge'
					AND b.badge_code = up.source_code

			-- game
			LEFT JOIN games g
					ON up.data_source = 'game'
					AND g.game_code = up.source_code

			WHERE u.user_code = $1
				AND up.data_source IN (
				'room','room_play',
				'room_paid',
				'tournament','tournament_play',
				'tournament_paid',
				'redeem',
				'badge',
				'game',
				'profile',
				'point_subtract',
				'point_add'
				)
			) AS t
			WHERE t.title_description IS NOT NULL
			ORDER BY t.created_date DESC
			LIMIT 5;`
	)

	rows, err := db.Query(ctx, query, UserCode)
	if err != nil {
		return list, c.errHandler("model.GetUserPointActivities", err, utils.ErrGetUsersPointActivity)
	}

	defer rows.Close()
	for rows.Next() {
		var data UserPointEnt
		err = rows.Scan(
			&data.Id,
			&data.UserName,
			&data.UserStyle,
			&data.TitleDescription,
			&data.DataSource,
			&data.SourceCode,
			&data.Point,
			&data.CreatedDate,
		)
		if err != nil {
			return list, c.errHandler("model.GetUserPointActivities", err, utils.ErrScanUsersPointActivity)
		}
		list = append(list, data)
	}

	return list, nil
}

func (c *Contract) GetUserPointHistories(db *pgxpool.Pool, ctx context.Context, UserCode string, param request.UserPointHistoryParam) ([]UserPointHistory, request.UserPointHistoryParam, error) {
	var (
		err  error
		list []UserPointHistory

		query = `
			SELECT
				t.source_id,
				t.source_user_code,
				t.source_type,
				t.source_code,
				t.source_name,
				t.point,
				t.description,
				t.created_date
			FROM (
			SELECT
				up.id AS source_id,
				u.user_code as source_user_code,
				up.data_source AS source_type,
				up.source_code,
				CASE
				WHEN up.data_source IN ('room','room_play')
					THEN 'Joined: ' || r.name
				WHEN up.data_source = 'room_paid'
					THEN 'Paid: ' || rr.name
				WHEN up.data_source = 'tournament'
					THEN 'Joined: ' || t.name
				WHEN up.data_source = 'tournament_play'
					THEN 'Participated: ' || t.name
				WHEN up.data_source = 'tournament_paid'
					THEN 'Paid: ' || tt.name
				WHEN up.data_source = 'redeem'
					THEN 'Purchased: ' || urh.description
				WHEN up.data_source = 'badge'
					THEN 'Claimed: ' || b.name
				WHEN up.data_source = 'game'
					THEN 'Collected: ' || g.name
				WHEN up.data_source = 'profile'
					THEN 'Profile: Updated Profile Image for the first time'
				WHEN up.data_source = 'point_subtract'
					THEN 'Point: Your point was decreased by ' || SUBSTR(up.point::text, 2)
				WHEN up.data_source = 'point_add'
					THEN 'Point: Your point was increased by ' || up.point
				WHEN up.data_source = 'tier'
					THEN 'Tier: You have reached tier ' || tiers.name
				END AS source_name,
				up.point,
				COALESCE(up.description, '') as description,
				up.created_date
			FROM users_points up
			JOIN users u
				ON u.id = up.user_id

			-- rooms (room, room_play)
			LEFT JOIN rooms r
					ON up.data_source IN ('room','room_play')
					AND r.room_code = up.source_code

			-- room_paid -> users_transactions -> rooms
			LEFT JOIN users_transactions ut_room_paid
					ON up.data_source = 'room_paid'
					AND ut_room_paid.user_id = up.user_id
					AND ut_room_paid.transaction_code = up.source_code
			LEFT JOIN rooms rr
					ON rr.room_code = ut_room_paid.source_code

			-- tournament / tournament_play
			LEFT JOIN tournaments t
					ON up.data_source IN ('tournament','tournament_play')
					AND t.tournament_code = up.source_code

			-- tournament_paid -> users_transactions -> tournaments
			LEFT JOIN users_transactions ut_tour_paid
					ON up.data_source = 'tournament_paid'
					AND ut_tour_paid.user_id = up.user_id
					AND ut_tour_paid.transaction_code = up.source_code
			LEFT JOIN tournaments tt
					ON tt.tournament_code = ut_tour_paid.source_code

			-- redeem
			LEFT JOIN user_redeem_histories urh
					ON up.data_source = 'redeem'
					AND urh.custom_id = up.source_code

			-- badge
			LEFT JOIN badges b
					ON up.data_source = 'badge'
					AND b.badge_code = up.source_code

			-- game
			LEFT JOIN games g
					ON up.data_source = 'game'
					AND g.game_code = up.source_code

			-- tier
			LEFT JOIN tiers
					ON up.data_source = 'tier'
					AND tiers.tier_code = up.source_code

			WHERE u.user_code = $1
				AND up.data_source IN (
				'room_play',
				'room_paid',
				'tournament_play',
				'tournament_paid',
				'redeem',
				'badge',
				'game',
				'profile',
				'point_subtract',
				'point_add',
				'tier'
				)
			) AS t
			WHERE t.source_name IS NOT NULL`
	)

	{
		totalData := 0
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS total`
		err := db.QueryRow(ctx, newQcount, UserCode).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetUserPointHistories", err, utils.ErrCountingListUserPointHistory)
		}
		param.Count = totalData
	}

	var (
		conditions []string
		args       []interface{}
	)

	args = append(args, UserCode)
	if param.SourceType != "" {
		if param.SourceType == "point" {
			conditions = append(conditions, fmt.Sprintf("t.source_type IN ($%d, $%d)", len(args)+1, len(args)+2))
			args = append(args, "point_subtract", "point_add")
		} else {
			args = append(args, param.SourceType)
			conditions = append(conditions, fmt.Sprintf("t.source_type = $%d", len(args)))
		}
	}

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s AND %s", query, strings.Join(conditions, " AND "))
	}

	if param.Order != "" && param.Sort != "" {
		query = fmt.Sprintf("%s ORDER BY t.%s %s", query, param.Order, param.Sort)
	}

	if param.Limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, param.Limit)
	}

	if param.Offset > 0 {
		query = fmt.Sprintf("%s OFFSET %d", query, param.Offset)
	}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return list, param, c.errHandler("model.GetUserPointHistories", err, utils.ErrGetUsersPointHistories)
	}

	defer rows.Close()
	for rows.Next() {
		var data UserPointHistory
		err = rows.Scan(
			&data.SourceId,
			&data.SourceUserCode,
			&data.SourceType,
			&data.SourceCode,
			&data.SourceName,
			&data.Point,
			&data.Description,
			&data.CreatedDate,
		)
		if err != nil {
			return list, param, c.errHandler("model.GetUserPointHistories", err, utils.ErrScanUsersPointHistories)
		}
		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) GetUserIdByUserCode(db *pgxpool.Pool, ctx context.Context, UserCode string) (int64, error) {
	var userId int64

	sql := `SELECT id FROM users WHERE user_code = $1;`

	err := db.QueryRow(ctx, sql, UserCode).Scan(&userId)
	if err != nil {
		return userId, c.errHandler("model.GetUserIdByUserCode", err, utils.ErrRetrievingUserByUserIdentifier)
	}

	return userId, nil
}

func (c *Contract) DeleteUserByCode(db *pgxpool.Pool, ctx context.Context, userCode string) error {
	var (
		err error
		sql = `DELETE FROM users WHERE user_code=$1`
	)
	_, err = db.Exec(ctx, sql, userCode)
	if err != nil {
		return c.errHandler("model.DeleteUserByCode", err, utils.ErrDeletingUser)
	}

	return nil
}

func (c *Contract) GetListUsersByUserId(db *pgxpool.Pool, ctx context.Context) ([]int64, error) {
	sql := "SELECT id  FROM users"
	rows, err := db.Query(ctx, sql)
	if err != nil {
		return nil, c.errHandler("model.GetListUsersByUserId", err, utils.ErrGettingListUserId)
	}
	defer rows.Close()

	var list []int64

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, c.errHandler("model.GetListUsersByUserId", err, utils.ErrScanningListUserId)
		}
		list = append(list, id)
	}

	if err := rows.Err(); err != nil {
		return nil, c.errHandler("model.GetListUsersByUserId", err, utils.ErrScanningListUserId)
	}

	return list, nil
}

func (c *Contract) CheckIfUsernameExists(db *pgxpool.Pool, ctx context.Context, username string) error {
	var exists bool

	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		return c.errHandler("model.CheckIfUsernameExists", err, utils.ErrCheckingIfUsernameExists)
	}

	if exists {
		return c.errHandler("model.CheckIfUsernameExists", errors.New(utils.ErrUsernameHasAlreadyTaken), utils.ErrUsernameHasAlreadyTaken)
	}

	return nil
}

func (c *Contract) StoreUserStyle(db *pgxpool.Pool, ctx context.Context, userID int64, styles UserStyle) error {
	var (
		err error
		sql = `UPDATE users SET styles = $2 WHERE id = $1`
	)

	_, err = db.Exec(ctx, sql, userID, styles)
	if err != nil {
		return c.errHandler("model.StoreUserStyle", err, utils.ErrStoringUserStyle)
	}

	return nil
}
