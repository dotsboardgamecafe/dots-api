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

type TierEnt struct {
	Id          int64          `db:"id"`
	TierCode    string         `db:"tier_code"`
	Name        string         `db:"name"`
	MinPoint    int64          `db:"min_point"`
	MaxPoint    int64          `db:"max_point"`
	Description sql.NullString `db:"description"`
	Status      sql.NullString `db:"status"`
	CreatedDate time.Time      `db:"created_date"`
	UpdatedDate sql.NullTime   `db:"updated_date"`
	DeletedDate sql.NullTime   `db:"deleted_date"`
}

type TierWithRewardEnt struct {
	TierId            int64          `db:"tier_id"`
	TierName          string         `db:"tier_name"`
	RewardCode        string         `db:"reward_code"`
	RewardName        sql.NullString `db:"reward_name"`
	RewardImageUrl    sql.NullString `db:"reward_img_url"`
	RewardDescription sql.NullString `db:"reward_description"`
}

func (c *Contract) GetTiersList(db *pgxpool.Pool, ctx context.Context, param request.TierParam) ([]TierEnt, request.TierParam, error) {
	var (
		err        error
		list       []TierEnt
		paramQuery []interface{}
		totalData  int
		query      = `
		SELECT 
			id, tier_code, name, min_point, max_point, description, status, created_date, updated_date, deleted_date
		FROM tiers`
	)

	// Populate Search
	paramQuery, query = generateTierFilterByQuery(param, query)

	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS total`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetTiersList", err, utils.ErrCountingListBadge)
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
	query += fmt.Sprintf("OFFSET $%d ", len(paramQuery))

	paramQuery = append(paramQuery, param.Limit)
	query += fmt.Sprintf("LIMIT $%d ", len(paramQuery))

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, param, c.errHandler("model.GetTiersList", err, utils.ErrGettingListTier)
	}

	defer rows.Close()
	for rows.Next() {
		var data TierEnt
		err = rows.Scan(&data.Id, &data.TierCode, &data.Name, &data.MinPoint, &data.MaxPoint, &data.Description, &data.Status, &data.CreatedDate, &data.UpdatedDate, &data.DeletedDate)
		if err != nil {
			return list, param, c.errHandler("model.GetTiersList", err, utils.ErrScanningListTier)
		}
		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) GetTierByCode(db *pgxpool.Pool, ctx context.Context, code string) (TierEnt, error) {
	var (
		err   error
		data  TierEnt
		query = `
		SELECT 
			id, tier_code, name, min_point, max_point, description , status, created_date, updated_date, deleted_date
		FROM tiers
		WHERE tier_code = $1 AND deleted_date is null`
	)

	err = db.QueryRow(ctx, query, code).Scan(&data.Id, &data.TierCode, &data.Name, &data.MinPoint, &data.MaxPoint, &data.Description, &data.Status, &data.CreatedDate, &data.UpdatedDate, &data.DeletedDate)
	if err != nil {
		return data, c.errHandler("model.GetTierByCode", err, utils.ErrGettingTierByCode)
	}

	return data, nil
}

func (c *Contract) GetTierByPoinCriteria(tx pgx.Tx, ctx context.Context, point int) (TierEnt, error) {
	var (
		data  TierEnt
		query = `SELECT 
			COALESCE(tiers.id, 4) AS id, 
			COALESCE(tiers.name, 'Legend') AS name
		FROM (SELECT 1)
		LEFT JOIN tiers ON $1 BETWEEN tiers.min_point AND tiers.max_point;`
	)

	_ = tx.QueryRow(ctx, query, point).Scan(&data.Id, &data.Name)

	return data, nil
}

func (c *Contract) InsertTier(db *pgxpool.Pool, ctx context.Context, tierCode, name, description string, minPoint, maxPoint int, status string) error {
	query := `
		INSERT INTO tiers (tier_code, name, min_point, max_point, description, status, created_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.Exec(ctx, query, tierCode, name, minPoint, maxPoint, description, status, time.Now().UTC())
	if err != nil {
		return c.errHandler("model.InsertTier", err, utils.ErrAddingTier)
	}

	return nil
}

func (c *Contract) UpdateTier(db *pgxpool.Pool, ctx context.Context, tierCode, name, description string, minPoint, maxPoint int, status string) error {
	query := `
		UPDATE tiers 
		SET name = $2, min_point = $3, max_point = $4, description = $5, status = $6, updated_date = $7
		WHERE tier_code = $1`

	_, err := db.Exec(ctx, query, tierCode, name, minPoint, maxPoint, description, status, time.Now().UTC())
	if err != nil {
		return c.errHandler("model.UpdateTier", err, utils.ErrUpdatingTier)
	}

	return nil
}

func (c *Contract) DeleteTier(db *pgxpool.Pool, ctx context.Context, tierCode string) error {
	query := `
		UPDATE tiers 
		SET deleted_date = $2
		WHERE tier_code = $1`

	_, err := db.Exec(ctx, query, tierCode, time.Now().UTC())
	if err != nil {
		return c.errHandler("model.DeleteTier", err, utils.ErrDeletingTier)
	}

	return nil
}

func (c *Contract) GetTierIdByCode(db *pgxpool.Pool, ctx context.Context, code string) (int64, error) {
	var (
		err   error
		id    int64
		query = `SELECT id FROM tiers WHERE tier_code=$1`
	)
	err = db.QueryRow(ctx, query, code).Scan(&id)
	if err != nil {
		return id, c.errHandler("model.GetTierIdByCode", err, utils.ErrGettingTierByCode)
	}
	return id, nil
}

// Private Function
func generateTierFilterByQuery(param request.TierParam, query string) ([]interface{}, string) {
	var (
		where      []string
		paramQuery []interface{}
	)

	// STATUS
	if len(param.Status) > 0 {
		paramQuery = append(paramQuery, param.Status)
		where = append(where, "status = $"+strconv.Itoa(len(paramQuery)))
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

func (c *Contract) GetTierWithReward(db *pgxpool.Pool, ctx context.Context, tierId int) ([]TierWithRewardEnt, error) {
	var (
		err  error
		list []TierWithRewardEnt

		query = `SELECT 
			t.id AS tier_id,
			t.name AS tier_name,
			r.name AS reward_name,
			r.reward_code AS reward_code,
			r.image_url AS reward_img_url,
			r.description AS reward_description
    FROM tiers t JOIN rewards r ON t.id = r.tier_id AND r.status = 'active'
		WHERE t.id = $1 AND r.expired_date > NOW()`
	)

	rows, err := db.Query(ctx, query, tierId)
	if err != nil {
		return list, c.errHandler("model.GetTierWithReward", err, utils.ErrGetTierWithReward)
	}

	defer rows.Close()
	for rows.Next() {
		var data TierWithRewardEnt
		err = rows.Scan(
			&data.TierId, &data.TierName, &data.RewardName, &data.RewardCode, &data.RewardImageUrl, &data.RewardDescription,
		)

		if err != nil {
			return list, c.errHandler("model.GetTierWithReward", err, utils.ErrScanTierWithReward)
		}

		list = append(list, data)
	}

	return list, nil
}
