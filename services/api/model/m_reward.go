package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"
	"dots-api/services/api/request"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RewardEnt struct {
	Id           int64 `db:"id"`
	Tier         TierEnt
	Name         string         `db:"name"`
	ImageUrl     string         `db:"image_url"`
	CategoryType string         `db:"category_type"`
	RewardCode   string         `db:"reward_code"`
	Description  sql.NullString `db:"description"`
	Status       string         `db:"status"`
	ExpiredDate  sql.NullTime   `db:"expired_date"`
	CreatedDate  time.Time      `db:"created_date"`
	UpdatedDate  sql.NullTime   `db:"updated_date"`
	DeletedDate  sql.NullTime   `db:"deleted_date"`
}

type RewardResEnt struct {
	Id           int64          `db:"id"`
	TierId       int64          `db:"tier_id"`
	Name         string         `db:"name"`
	ImageUrl     string         `db:"image_url"`
	CategoryType string         `db:"category_type"`
	RewardCode   string         `db:"reward_code"`
	Description  sql.NullString `db:"description"`
	Status       string         `db:"status"`
	ExpiredDate  sql.NullTime   `db:"expired_date"`
	CreatedDate  time.Time      `db:"created_date"`
	UpdatedDate  sql.NullTime   `db:"updated_date"`
	DeletedDate  sql.NullTime   `db:"deleted_date"`
}

func (c *Contract) GetRewardList(db *pgxpool.Pool, ctx context.Context, param request.RewardParam) ([]RewardEnt, request.RewardParam, error) {
	var (
		err        error
		list       []RewardEnt
		paramQuery []interface{}
		totalData  int
		query      = `
        SELECT 
            r.id, 
            t.id AS tier_id,
            t.tier_code, 
            t.name AS tier_name, 
            t.min_point, 
            t.max_point, 
            t.description AS tier_description,
            r.name AS reward_name,
            r.image_url,
            r.category_type,
            r.reward_code,
            r.status,
			r.description,
            r.expired_date,
            r.created_date,
            r.updated_date,
			r.deleted_date
        FROM rewards r
        JOIN tiers t ON r.tier_id = t.id
        `
	)

	// Populate Search
	paramQuery, query = generateRewardFilterByQuery(param, query)

	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS total`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetBadgeList", err, utils.ErrCountingListReward)
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
		return list, param, c.errHandler("model.GetRewardList", err, utils.ErrGettingListReward)
	}

	defer rows.Close()
	for rows.Next() {
		var data RewardEnt
		err := rows.Scan(&data.Id, &data.Tier.Id, &data.Tier.TierCode, &data.Tier.Name, &data.Tier.MinPoint, &data.Tier.MaxPoint, &data.Tier.Description, &data.Name, &data.ImageUrl, &data.CategoryType, &data.RewardCode, &data.Status, &data.Description, &data.ExpiredDate, &data.CreatedDate, &data.UpdatedDate, &data.DeletedDate)
		if err != nil {
			return list, param, c.errHandler("model.GetRewardList", err, utils.ErrScanningListReward)
		}
		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) GetRewardByCode(db *pgxpool.Pool, ctx context.Context, rewardCode string) (RewardEnt, error) {
	var (
		data  RewardEnt
		query = `
        SELECT 
            r.id, 
            t.id AS tier_id,
            t.tier_code, 
            t.name AS tier_name, 
            t.min_point, 
            t.max_point, 
            t.description AS tier_description,
            r.name AS reward_name,
            r.image_url,
            r.category_type,
            r.reward_code,
            r.status,
			r.description,
            r.expired_date,
            r.created_date,
            r.updated_date,
			r.deleted_date 
        FROM rewards r
        JOIN tiers t ON r.tier_id = t.id
        WHERE r.reward_code = $1
        `
	)

	err := db.QueryRow(ctx, query, rewardCode).Scan(
		&data.Id, &data.Tier.Id, &data.Tier.TierCode, &data.Tier.Name, &data.Tier.MinPoint, &data.Tier.MaxPoint, &data.Tier.Description, &data.Name, &data.ImageUrl, &data.CategoryType, &data.RewardCode, &data.Status, &data.Description, &data.ExpiredDate, &data.CreatedDate, &data.UpdatedDate, &data.DeletedDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return data, c.errHandler("model.GetRewardByCode", errors.New("reward not found"), utils.ErrRewardNotFound)
		}
		return data, c.errHandler("model.GetRewardByCode", err, utils.ErrGettingRewardByCode)
	}

	return data, nil
}

func (c *Contract) AddReward(db *pgxpool.Pool, ctx context.Context, tierID int64, name, imageUrl, categoryType, rewardCode, description, status string, expiredDate interface{}) error {
	query := `
        INSERT INTO rewards (tier_id, name, image_url, category_type, reward_code, status, description, expired_date, created_date)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := db.Exec(ctx, query, tierID, name, imageUrl, categoryType, rewardCode, status, description, expiredDate, time.Now().In(time.UTC))
	if err != nil {
		return c.errHandler("model.AddReward", err, utils.ErrAddingReward)
	}

	return nil
}

func (c *Contract) UpdateReward(db *pgxpool.Pool, ctx context.Context, rewardCode string, tierID int64, name, imageUrl, categoryType, description, status string, expiredDate interface{}) error {
	query := `
        UPDATE rewards 
        SET tier_id = $2, name = $3, image_url = $4, category_type = $5, status = $6, description = $7, expired_date = $8, updated_date = $9
        WHERE reward_code = $1`

	_, err := db.Exec(ctx, query, rewardCode, tierID, name, imageUrl, categoryType, status, description, expiredDate, time.Now().In(time.UTC))
	if err != nil {
		return c.errHandler("model.UpdateReward", err, utils.ErrUpdatingReward)
	}

	return nil
}

func (c *Contract) DeleteReward(db *pgxpool.Pool, ctx context.Context, code string) error {
	query := `
       UPDATE rewards 
	   SET updated_date = $1
        WHERE reward_code = $2`

	_, err := db.Exec(ctx, query, time.Now().In(time.UTC), code)
	if err != nil {
		return c.errHandler("model.DeleteReward", err, utils.ErrDeletingReward)
	}

	return nil
}

func (c *Contract) GetBenefitsByTierId(tx pgx.Tx, ctx context.Context, tierId int64) ([]RewardEnt, error) {
	var (
		list  []RewardEnt
		query = `
			SELECT 
				r.id, 
				r.name AS reward_name,
				r.image_url,
				r.category_type,
				r.reward_code
			FROM rewards r
			JOIN tiers t ON r.tier_id = t.id
			WHERE t.id = $1 AND r.status = 'active'
		`
	)

	rows, err := tx.Query(ctx, query, tierId)
	if err != nil {
		return list, c.errHandler("model.GetBenefitsByTierId", err, utils.ErrGettingListReward)
	}

	defer rows.Close()
	for rows.Next() {
		var data RewardEnt
		err := rows.Scan(&data.Id, &data.Name, &data.ImageUrl, &data.CategoryType, &data.RewardCode)
		if err != nil {
			return list, c.errHandler("model.GetBenefitsByTierId", err, utils.ErrScanningListReward)
		}
		list = append(list, data)
	}

	return list, nil
}

// Private Function
func generateRewardFilterByQuery(param request.RewardParam, query string) ([]interface{}, string) {
	var (
		where      []string
		paramQuery []interface{}
	)

	// STATUS
	if len(param.Status) > 0 {
		paramQuery = append(paramQuery, param.Status)
		where = append(where, "r.status = $"+strconv.Itoa(len(paramQuery)))
	}

	// NAME
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		orWhere = append(orWhere, "r.name iLIKE $"+strconv.Itoa(len(paramQuery)))
		orWhere = append(orWhere, "t.name iLIKE $"+strconv.Itoa(len(paramQuery)))
		where = append(where, "("+strings.Join(orWhere, " OR ")+")")
	}

	// CATEGORY TYPE
	if len(param.CategoryType) > 0 {
		paramQuery = append(paramQuery, param.CategoryType)
		where = append(where, "r.category_type = $"+strconv.Itoa(len(paramQuery)))
	}

	// EXPIRED DATE
	if len(param.ExpiredDate) > 0 {
		paramQuery = append(paramQuery, param.ExpiredDate)
		where = append(where, "r.expired_date > $"+strconv.Itoa(len(paramQuery)))
	}

	// Append All Where Conditions
	if len(where) > 0 {
		query += " WHERE r.deleted_date IS NULL AND " + strings.Join(where, " AND ")
	} else {
		query += " WHERE r.deleted_date IS NULL "
	}

	return paramQuery, query
}
