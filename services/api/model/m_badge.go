package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"dots-api/services/api/request"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BadgeEnt struct {
	Id            int64          `db:"id"`
	BadgeCode     string         `db:"badge_code"`
	BadgeCategory string         `db:"badge_category"`
	Name          string         `db:"name"`
	ImageURL      string         `db:"image_url"`
	VPPoint       int64          `db:"vp_point"`
	Status        string         `db:"status"`
	Description   sql.NullString `db:"description"`
	ParentCode    sql.NullString `db:"parent_code"`
	CreatedDate   time.Time      `db:"created_date"`
	UpdatedDate   sql.NullTime   `db:"updated_date"`
	DeletedDate   sql.NullTime   `db:"deleted_date"`
}

// GetBadgeList retrieves a list of all badges from the database.
func (c *Contract) GetBadgeList(db *pgxpool.Pool, ctx context.Context, param request.BadgeParam) ([]BadgeEnt, request.BadgeParam, error) {
	var (
		err        error
		list       []BadgeEnt
		paramQuery []interface{}
		totalData  int
		query      = `SELECT id, badge_code, badge_category, description, vp_point, name, image_url, status, parent_code, created_date, updated_date, deleted_date FROM badges`
	)

	// Populate Search
	paramQuery, query = generateBadgeFilterByQuery(param, query)

	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS total`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetBadgeList", err, utils.ErrCountingListBadge)
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
	query += fmt.Sprintf("OFFSET $%d ", len(paramQuery))

	paramQuery = append(paramQuery, param.Limit)
	query += fmt.Sprintf("LIMIT $%d ", len(paramQuery))

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, param, c.errHandler("model.GetBadgeList", err, utils.ErrGettingListBadge)
	}

	defer rows.Close()

	for rows.Next() {
		var badge BadgeEnt
		err = rows.Scan(
			&badge.Id, &badge.BadgeCode, &badge.BadgeCategory, &badge.Description, &badge.VPPoint, &badge.Name,
			&badge.ImageURL, &badge.Status, &badge.ParentCode, &badge.CreatedDate,
			&badge.UpdatedDate, &badge.DeletedDate,
		)
		if err != nil {
			return nil, param, c.errHandler("model.GetBadgeList", err, utils.ErrScanningListBadge)
		}
		list = append(list, badge)
	}

	return list, param, nil
}

// GetBadgeDetail retrieves details of a badge by its ID.
func (c *Contract) GetBadgeDetailByCode(db *pgxpool.Pool, ctx context.Context, code string) (BadgeEnt, error) {
	var badge BadgeEnt

	query := `SELECT id, badge_code, badge_category,  vp_point, description, name, image_url, status, created_date, updated_date, deleted_date FROM badges WHERE badge_code = $1 AND deleted_date is null`
	err := db.QueryRow(ctx, query, code).Scan(
		&badge.Id, &badge.BadgeCode, &badge.BadgeCategory, &badge.VPPoint, &badge.Description, &badge.Name,
		&badge.ImageURL, &badge.Status, &badge.CreatedDate,
		&badge.UpdatedDate, &badge.DeletedDate,
	)
	if err != nil {
		return badge, c.errHandler("model.GetBadgeDetailByCode", err, utils.ErrGettingBadgeByCode)
	}

	return badge, nil
}

// GetBadgeDetailByParentCode retrieves details of a badge by its code.
func (c *Contract) GetBadgeDetailByParentCode(db *pgxpool.Pool, ctx context.Context, code string) ([]BadgeEnt, error) {
	var badges []BadgeEnt

	query := `SELECT id, badge_code, badge_category, vp_point, description, name, image_url, status, created_date, updated_date, deleted_date 
	          FROM badges 
	          WHERE parent_code = $1 AND deleted_date IS NULL`

	rows, err := db.Query(ctx, query, code)
	if err != nil {
		return nil, c.errHandler("model.GetBadgeDetailByParentCode", err, utils.ErrGettingListBadgeByParentCode)
	}
	defer rows.Close()

	for rows.Next() {
		var badge BadgeEnt
		err := rows.Scan(
			&badge.Id, &badge.BadgeCode, &badge.BadgeCategory, &badge.VPPoint, &badge.Description, &badge.Name,
			&badge.ImageURL, &badge.Status, &badge.CreatedDate,
			&badge.UpdatedDate, &badge.DeletedDate,
		)
		if err != nil {
			return nil, c.errHandler("model.GetBadgeDetailByParentCode", err, utils.ErrGettingBadgeByCodeByParentCode)
		}
		badges = append(badges, badge)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, c.errHandler("model.GetBadgeDetailByParentCode", err, utils.ErrGettingListBadgeByParentCode)
	}

	return badges, nil
}

// GetBadgeRuleList retrieves a list of all badge rules from the database.
func (c *Contract) GetBadgeListByKeyCondition(db *pgxpool.Pool, ctx context.Context, keyCondition string) ([]string, error) {
	var (
		err  error
		list []string
	)

	query := `
	SELECT 
		b.badge_code 
	FROM badges_rules br 
	LEFT JOIN badges b ON b.id = br.badge_id 
	WHERE br.key_condition = $1`
	rows, err := db.Query(ctx, query, keyCondition)
	if err != nil {
		return list, c.errHandler("model.GetBadgeListByKeyCondition", err, utils.ErrGettingBadgeRuleList)
	}
	defer rows.Close()

	for rows.Next() {
		var ent string
		err = rows.Scan(
			&ent,
		)
		if err != nil {
			return list, c.errHandler("model.GetBadgeListByKeyCondition", err, utils.ErrScanningBadgeRule)
		}
		list = append(list, ent)
	}

	return list, nil
}

// AddBadge adds a new badge to the database within a transaction.
func (c *Contract) AddBadge(tx pgx.Tx, ctx context.Context, badgeCode, badgeCategory, name, imageURL string, vpPoint int64, status, description, parentCode string) (int64, error) {
	var id int64

	query := `INSERT INTO badges(badge_code, badge_category, vp_point, name, image_url, status, description, parent_code, created_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	err := tx.QueryRow(ctx, query, badgeCode, badgeCategory, vpPoint, name, imageURL, status, description, parentCode, time.Now().UTC()).Scan(&id)
	if err != nil {
		return 0, c.errHandler("model.AddBadge", err, utils.ErrAddingBadge)
	}

	return id, nil
}

// UpdateBadge updates an existing badge in the database.
func (c *Contract) UpdateBadgeByCode(tx pgx.Tx, ctx context.Context, badgeCategory, name, description, imageURL, status string, vpPoint int64, badgeCode string) error {
	query := `UPDATE badges SET badge_category = $1, name = $2, description = $3, image_url = $4, status = $5, vp_point = $6, updated_date = $7 WHERE badge_code = $8`

	_, err := tx.Exec(ctx, query, badgeCategory, name, description, imageURL, status, vpPoint, time.Now().UTC(), badgeCode)
	if err != nil {
		return c.errHandler("model.UpdateBadgeByCode", err, utils.ErrUpdatingBadge)
	}

	return nil
}

// DeleteBadge marks a badge as deleted in the database.
func (c *Contract) DeleteBadge(tx pgx.Tx, ctx context.Context, id int64) error {
	query := `UPDATE badges SET deleted_date = $1 WHERE id = $2`

	_, err := tx.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return c.errHandler("model.DeleteBadge", err, utils.ErrDeletingBadge)
	}

	return nil
}

func (c *Contract) GetBadgeIdByCode(db *pgxpool.Pool, ctx context.Context, code string) (int64, error) {
	var (
		err   error
		id    int64
		query = `SELECT id FROM badges WHERE badge_code=$1`
	)
	err = db.QueryRow(ctx, query, code).Scan(&id)
	if err != nil {
		return id, c.errHandler("model.GetBadgeIdByCode", err, utils.ErrGettingBadgeByCode)
	}
	return id, nil
}

// Private Function
func generateBadgeFilterByQuery(param request.BadgeParam, query string) ([]interface{}, string) {
	var (
		where      []string
		paramQuery []interface{}
	)

	// STATUS
	if len(param.Status) > 0 {
		paramQuery = append(paramQuery, param.Status)
		where = append(where, "status = $"+strconv.Itoa(len(paramQuery)))
	}

	// BADGE CATEGORY
	if len(param.BadgeCategory) > 0 {
		paramQuery = append(paramQuery, param.BadgeCategory)
		where = append(where, "badge_category = $"+strconv.Itoa(len(paramQuery)))
	}

	// NAME (KEYWORD)
	if len(param.Keyword) > 0 {
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		where = append(where, "name iLIKE $"+strconv.Itoa(len(paramQuery)))
	}

	// Append All Where Conditions
	if len(where) > 0 {
		query += " WHERE deleted_date IS NULL AND " + strings.Join(where, " AND ")
	} else {
		query += " WHERE deleted_date IS NULL "
	}

	return paramQuery, query
}
