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

type UserBadgeEnt struct {
	Id          int64     `db:"id"`
	UserId      int64     `db:"user_id"`
	BadgeId     int64     `db:"badge_id"`
	IsClaim     bool      `db:"is_claim"`
	CreatedDate time.Time `db:"created_date"`
}

type UserBadgeResp struct {
	BadgeId       int64          `db:"badge_id"`
	UserId        sql.NullInt64  `db:"user_id"`
	BadgeName     string         `db:"badge_name"`
	BadgeImageURL string         `db:"badge_image_url"`
	BadgeCode     string         `db:"badge_code"`
	VPPoint       int            `db:"vp_point"`
	BadgeCategory string         `db:"badge_category"`
	Description   sql.NullString `db:"description"`
	IsClaim       sql.NullBool   `db:"is_claim"`
	CreatedDate   time.Time      `db:"created_date"`
	IsBadgeOwned  sql.NullBool   `db:"is_badge_owned"`
	NeedToClaim   sql.NullBool   `db:"need_to_claim"`
}

func (c *Contract) GetUserBadgeList(db *pgxpool.Pool, ctx context.Context, userCode string, param request.UserBadgeParam) ([]UserBadgeResp, request.UserBadgeParam, error) {
	var (
		err        error
		list       []UserBadgeResp
		paramQuery []interface{}
		totalData  int
		where      []string

		query = `SELECT 
					b.id as badge_id,
					ub.user_id ,
					b."name" as badge_name,
					b.image_url as badge_image_url,
					b.badge_code ,
					b.badge_category,
					b.vp_point,
					b.description,
					ub.is_claim,
					b.created_date,
					CASE
						WHEN user_id IS NULL THEN false
						ELSE true
					END AS is_badge_owned,
					CASE
							WHEN ub.user_id is not NULL and ub.is_claim = false THEN true
							ELSE false
					END AS need_to_claim
					FROM
							badges b
					LEFT JOIN 
						(
							SELECT ub.*
							FROM public.users_badges ub
							INNER JOIN public.users u ON ub.user_id = u.id
							WHERE u.user_code = $1
						) AS ub 
					ON 
						b.id = ub.badge_id`
	)

	paramQuery = append(paramQuery, userCode)

	if len(param.IsClaim) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.IsClaim)
		orWhere = append(orWhere, fmt.Sprintf(" ub.is_claim=$%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// Append All Where Conditions
	if len(where) > 0 {
		query += " WHERE b.status = 'active' AND " + strings.Join(where, " AND ")
	} else {
		query += " WHERE b.status = 'active' "
	}

	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS data`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetUserBadgeList", err, utils.ErrCountingListUserBadge)
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
		return list, param, c.errHandler("model.GetUserBadgeList", err, utils.ErrGettingListUserBadge)
	}

	defer rows.Close()
	for rows.Next() {
		var data UserBadgeResp
		err = rows.Scan(&data.BadgeId, &data.UserId, &data.BadgeName, &data.BadgeImageURL, &data.BadgeCode, &data.BadgeCategory, &data.VPPoint, &data.Description, &data.IsClaim, &data.CreatedDate, &data.IsBadgeOwned, &data.NeedToClaim)
		if err != nil {
			return list, param, c.errHandler("model.GetUserBadgeList", err, utils.ErrScanningListUserBadge)
		}
		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) GetUserBadgeByBadgeCode(db *pgxpool.Pool, ctx context.Context, userCode, badgeCode string) (UserBadgeResp, error) {
	var (
		err  error
		data UserBadgeResp

		query = `SELECT 
					b.id as badge_id,
					ub.user_id ,
					b."name" as badge_name,
					b.image_url as badge_image_url,
					b.badge_code ,
					b.badge_category,
					b.vp_point,
					b.description,
					ub.is_claim,
					b.created_date,
					CASE
						WHEN user_id IS NULL THEN false
						ELSE true
					END AS is_badge_owned,
					CASE
							WHEN ub.user_id is not NULL and ub.is_claim = false THEN true
							ELSE false
					END AS need_to_claim
				FROM
						badges b
				LEFT JOIN 
					(
						SELECT ub.*
						FROM public.users_badges ub
						INNER JOIN public.users u ON ub.user_id = u.id
						WHERE u.user_code = $1
					) AS ub 
				ON 
					b.id = ub.badge_id
				WHERE   b.badge_code=$2 `
	)

	err = db.QueryRow(ctx, query, userCode, badgeCode).Scan(&data.BadgeId, &data.UserId, &data.BadgeName, &data.BadgeImageURL, &data.BadgeCode, &data.BadgeCategory, &data.VPPoint, &data.Description, &data.IsClaim, &data.CreatedDate, &data.IsBadgeOwned, &data.NeedToClaim)
	if err != nil {
		return data, c.errHandler("model.GetUserBadgeByBadgeCode", err, utils.ErrGettingtUserBadgeByBadgeCode)
	}

	return data, nil
}

func (c *Contract) AddUserBadgeTx(tx pgx.Tx, ctx context.Context, userId, badgeId int64) error {
	// Query to check if the user badge already exists
	checkQuery := `
        SELECT 1 FROM users_badges WHERE user_id = $1 AND badge_id = $2
    `
	var exists int
	err := tx.QueryRow(ctx, checkQuery, userId, badgeId).Scan(&exists)
	if err != nil && err != pgx.ErrNoRows {
		return c.errHandler("model.AddUserBadge", err, utils.ErrorCheckingUserBadge)
	}

	// If the badge already exists, skip adding it
	if err == nil {
		return nil
	}

	// Query to add the user badge if it does not exist
	insertQuery := `
        INSERT INTO users_badges (user_id, badge_id, is_claim, created_date) 
        VALUES ($1, $2, $3, $4)
    `
	_, err = tx.Exec(ctx, insertQuery, userId, badgeId, false, time.Now())
	if err != nil {
		return c.errHandler("model.AddUserBadge", err, utils.ErrorAddingUserBadge)
	}

	return nil
}

func (c *Contract) AddUserBadge(db *pgxpool.Pool, ctx context.Context, userId, badgeId int64) error {
	// Query to check if the user badge already exists
	checkQuery := `
        SELECT 1 FROM users_badges WHERE user_id = $1 AND badge_id = $2
    `
	var exists int
	err := db.QueryRow(ctx, checkQuery, userId, badgeId).Scan(&exists)
	if err != nil && err != pgx.ErrNoRows {
		return c.errHandler("model.AddUserBadge", err, utils.ErrorCheckingUserBadge)
	}

	// If the badge already exists, skip adding it
	if err == nil {
		return nil
	}

	// Query to add the user badge if it does not exist
	insertQuery := `
        INSERT INTO users_badges (user_id, badge_id, is_claim, created_date) 
        VALUES ($1, $2, $3, $4)
    `
	_, err = db.Exec(ctx, insertQuery, userId, badgeId, false, time.Now())
	if err != nil {
		return c.errHandler("model.AddUserBadge", err, utils.ErrorAddingUserBadge)
	}

	return nil
}

func (c *Contract) UpdateUserBadge(tx pgx.Tx, ctx context.Context, userId, badgeId int64, isClaim bool) error {
	query := `
        UPDATE users_badges 
        SET is_claim = $1, updated_date = $2 
        WHERE user_id = $3 AND badge_id = $4
    `
	_, err := tx.Exec(ctx, query, isClaim, time.Now(), userId, badgeId)
	if err != nil {
		return c.errHandler("model.UpdateUserBadge", err, utils.ErrorUpdatingUserBadge)
	}
	return nil
}

func (c *Contract) DeleteUserBadge(tx pgx.Tx, ctx context.Context, userId, badgeId int64) error {
	// Check if the user badge exists
	exists, err := c.IsUserBadgeExists(tx, ctx, userId, badgeId)
	if err != nil {
		return err
	}
	if !exists {
		// User badge doesn't exist, no need to delete
		return nil
	}

	// User badge exists, proceed with deletion
	query := `
        DELETE FROM users_badges 
        WHERE user_id = $1 AND badge_id = $2
    `
	_, err = tx.Exec(ctx, query, userId, badgeId)
	if err != nil {
		return c.errHandler("model.DeleteUserBadge", err, utils.ErrorDeletingUserBadge)
	}
	return nil
}

func (c *Contract) IsUserBadgeExists(tx pgx.Tx, ctx context.Context, userId, badgeId int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users_badges WHERE user_id = $1 AND badge_id = $2)`
	err := tx.QueryRow(ctx, query, userId, badgeId).Scan(&exists)
	if err != nil {
		return false, c.errHandler("model.IsUserBadgeExists", err, "Error checking existence of user badge")
	}
	return exists, nil
}
