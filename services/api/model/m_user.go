package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"dots-api/services/api/request"
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
	UserName           sql.NullString `db:"username"`
	PhoneNumber        string         `db:"phone_number"`
	FullName           string         `db:"fullname"`
	ImageURL           sql.NullString `db:"image_url"`
	LatestPoint        int            `db:"latest_point"`
	LatestTierId       int            `db:"latest_tier_id"`
	LatestTierName     string         `db:"latest_tier_name"`
	TierMinRangePoint  int            `db:"tier_min_range_point"`
	TierMaxRangePoint  int            `db:"tier_max_range_point"`
	Password           string         `db:"password"`
	XPlayer            string         `db:"x_player"`
	RoleId             int            `db:"role_id"`
	StatusVerification bool           `db:"status_verification"`
	Status             string         `db:"status"`
	TotalSpent         int            `db:"total_spent"`
	CreatedDate        time.Time      `db:"created_date"`
	UpdatedDate        sql.NullTime   `db:"updated_date"`
	DeletedDate        sql.NullTime   `db:"deleted_date"`
}

func (c *Contract) GetAllUsers(db *pgxpool.Pool, ctx context.Context) ([]UserEnt, error) {
	var users []UserEnt

	query := `
		SELECT id, user_code, email, username, phone_number, fullname, image_url, 
		       latest_point, latest_tier_id, password, x_player, status_verification, 
		       status, created_date, updated_date, deleted_date 
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
						users.created_date AS created_date,
						users.updated_date AS updated_date,
						users.deleted_date AS deleted_date,
						users.role_id
					FROM users
					JOIN tiers ON tiers.id = users.latest_tier_id
					JOIN roles r ON r.id = users.role_id 
					WHERE email = $1 AND users.deleted_date IS NULL AND r.deleted_date IS NULL`

	err := db.QueryRow(ctx, sql, email).Scan(
		&res.ID,
		&res.UserCode,
		&res.Email,
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

		query = `SELECT users.id AS id,
						user_code,
						email,
						username,
						phone_number,
						fullname,
						image_url,
						latest_point,
						tiers.id AS latest_tier_id,
						tiers.name AS latest_tier_name,
						password,
						x_player,
						status_verification,
						users.status,
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
					) AS booking_transactions ON users.id = booking_transactions.user_id`
	)

	// Populate Search
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		orWhere = append(orWhere, fmt.Sprintf("fullname iLIKE $%d", len(paramQuery)))
		orWhere = append(orWhere, fmt.Sprintf("email iLIKE $%d", len(paramQuery)))
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
			&res.UserName,
			&res.PhoneNumber,
			&res.FullName,
			&res.ImageURL,
			&res.LatestPoint,
			&res.LatestTierId,
			&res.LatestTierName,
			&res.Password,
			&res.XPlayer,
			&res.StatusVerification,
			&res.Status,
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

func (c *Contract) UpdateUserProfile(db *pgxpool.Pool, ctx context.Context, userCode, fullName, imageURL string, phoneNumber string) error {
	sql := `
		UPDATE users
		SET image_url = $1, fullname = $2, phone_number = $3, updated_date = $4
		WHERE user_code = $5
	`

	currentTime := time.Now().In(time.UTC)
	_, err := db.Exec(ctx, sql, imageURL, fullName, phoneNumber, currentTime, userCode)
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

	_, err := db.Exec(ctx, sql, email, userIdentifier)
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

func (c *Contract) UpdateUser(db *pgxpool.Pool, ctx context.Context, userCode, fullName, imageURL, phoneNumber, email, userName, status string) error {
	sql := `
		UPDATE users
		SET image_url = $1, fullname = $2, phone_number = $3, email = $4, username = $5, status = $6, updated_date = $7
		WHERE user_code = $8
	`

	currentTime := time.Now().In(time.UTC)
	_, err := db.Exec(ctx, sql, imageURL, fullName, phoneNumber, email, userName, status, currentTime, userCode)
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

		query = `SELECT
    		u.id, 
				COALESCE(u.username, '') AS username,
				data_source, 
				source_code,
				g.name AS game_name,
				g.game_code,
    		g.image_url AS game_url,
				point, 
				up.created_date
    	FROM users_points up JOIN users u ON up.user_id = u.id
				LEFT JOIN tournaments t ON t.tournament_code = up.source_code
				LEFT JOIN rooms r ON r.room_code = up.source_code 
				JOIN games g ON g.id = r.game_id OR g.id = t.game_id
    	WHERE up.data_source != 'redeem'
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
		err = rows.Scan(&data.Id,
			&data.UserName,
			&data.DataSource,
			&data.SourceCode,
			&data.GameName,
			&data.GameCode,
			&data.GameImgUrl,
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

		query = `SELECT
    		u.id, 
				COALESCE(u.username, '') AS username,
				CASE
					-- Room type (normal and special_event)
					WHEN (up.data_source = 'room') THEN (
						SELECT CONCAT('Joined: ', rooms."name") AS info
						FROM rooms
						WHERE room_code = up.source_code
					)
					-- Tournament type
					WHEN (up.data_source = 'tournament') THEN (
						SELECT CONCAT('Joined: ', tournaments."name") AS info
						FROM tournaments
						WHERE tournament_code = up.source_code
					)
					-- Redeem Invoice
					WHEN (up.data_source = 'redeem') THEN (
						SELECT CONCAT('Purchased: ', description) AS info
						FROM user_redeem_histories
						WHERE custom_id = up.source_code
					)
					-- Badges
					WHEN (up.data_source = 'badge') THEN (
						SELECT CONCAT('Claimed: ', badges."name") AS info
						FROM badges
						WHERE badge_code = up.source_code
					)
				END AS title_description,
				data_source, 
				source_code,
				point, 
				up.created_date
    	FROM users_points up JOIN users u ON up.user_id = u.id
    	WHERE u.user_code = $1
			ORDER BY up.id DESC
			LIMIT 5;`
	)

	rows, err := db.Query(ctx, query, UserCode)
	if err != nil {
		return list, c.errHandler("model.GetUserPointActivities", err, utils.ErrGetUsersPointActivity)
	}

	defer rows.Close()
	for rows.Next() {
		var data UserPointEnt
		err = rows.Scan(&data.Id,
			&data.UserName,
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
		sql = `
		UPDATE users 
		SET updated_date=$1, deleted_date=$2 
		WHERE user_code=$3`
	)
	_, err = db.Exec(ctx, sql, time.Now().In(time.UTC), time.Now().In(time.UTC), userCode)
	if err != nil {
		return c.errHandler("model.DeleteUserByCode", err, utils.ErrDeletingAdmin)
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
