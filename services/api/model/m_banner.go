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

	"github.com/jackc/pgx/v4/pgxpool"
)

type BannerEnt struct {
	Id          int64        `db:"id"`
	BannerCode  string       `db:"banner_code"`
	BannerType  string       `db:"banner_type"`
	Title       string       `db:"title"`
	Description string       `db:"description"`
	ImageURL    string       `db:"image_url"`
	Status      string       `db:"status"`
	CreatedDate time.Time    `db:"created_date"`
	UpdatedDate sql.NullTime `db:"updated_date"`
	DeletedDate sql.NullTime `db:"deleted_date"`
}

func (c *Contract) GetBannerList(db *pgxpool.Pool, ctx context.Context, param request.BannerParam) ([]BannerEnt, request.BannerParam, error) {
	var (
		err        error
		list       []BannerEnt
		where      []string
		paramQuery []interface{}
		totalData  int
	)

	query := `SELECT id, banner_code, banner_type, title, description, image_url, status, created_date, updated_date, deleted_date FROM banners`

	// Populate Search
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		orWhere = append(orWhere, fmt.Sprintf("banner_code iLIKE $%d", len(paramQuery)))
		orWhere = append(orWhere, fmt.Sprintf("banner_type iLIKE $%d", len(paramQuery)))
		orWhere = append(orWhere, fmt.Sprintf("title iLIKE $%d", len(paramQuery)))
		orWhere = append(orWhere, fmt.Sprintf("description iLIKE $%d", len(paramQuery)))
		where = append(where, "("+strings.Join(orWhere, " OR ")+")")
	}

	if len(param.Status) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Status)
		orWhere = append(orWhere, fmt.Sprintf("status = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// Handling Soft Delete
	where = append(where, "deleted_date IS NULL")

	// Append All Where Conditions
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	// Count Query
	newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS data`
	err = db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
	if err != nil {
		return list, param, c.errHandler("model.GetBannerList", err, utils.ErrCountingListBanner)
	}
	param.Count = totalData

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
		return list, param, c.errHandler("model.GetBannerList", err, utils.ErrGettingListBanner)
	}
	defer rows.Close()

	for rows.Next() {
		var data BannerEnt
		err = rows.Scan(&data.Id, &data.BannerCode, &data.BannerType, &data.Title, &data.Description, &data.ImageURL, &data.Status, &data.CreatedDate, &data.UpdatedDate, &data.DeletedDate)
		if err != nil {
			return list, param, c.errHandler("model.GetBannerList", err, utils.ErrScanningListBanner)
		}
		list = append(list, data)
	}

	return list, param, nil
}

func (c *Contract) GetBannerByCode(db *pgxpool.Pool, ctx context.Context, code string) (BannerEnt, error) {
	var (
		err  error
		data BannerEnt
		sql  = `SELECT id, banner_code, banner_type, title, description, image_url, status, created_date, updated_date, deleted_date
		FROM banners 
		WHERE banner_code = $1`
	)
	err = db.QueryRow(ctx, sql, code).Scan(&data.Id, &data.BannerCode, &data.BannerType, &data.Title, &data.Description, &data.ImageURL, &data.Status, &data.CreatedDate, &data.UpdatedDate, &data.DeletedDate)
	if err != nil {
		return data, c.errHandler("model.GetBannerByCode", err, utils.ErrGettingBannerByCode)
	}

	return data, nil
}

func (c *Contract) AddBanner(db *pgxpool.Pool, ctx context.Context, code, title, desc, bannerType, imageUrl, status string) error {
	sql := `INSERT INTO banners(banner_code, banner_type, title, description, image_url, status, created_date)
	VALUES($1,$2,$3,$4,$5,$6,$7)`

	_, err := db.Exec(ctx, sql, code, bannerType, title, desc, imageUrl, status, time.Now().In(time.UTC))
	if err != nil {
		return c.errHandler("model.AddBanner", err, utils.ErrAddingBanner)
	}

	return nil
}

func (c *Contract) UpdateBannerByCode(db *pgxpool.Pool, ctx context.Context, code, bannerType, title, desc, imageUrl, status string) error {
	var (
		err error
		sql = `
		UPDATE banners 
		SET banner_type=$1, title=$2, description=$3, image_url= $4, status=$5, updated_date=$6
		WHERE banner_code=$7`
	)

	_, err = db.Exec(ctx, sql, bannerType, title, desc, imageUrl, status, time.Now().In(time.UTC), code)
	if err != nil {
		return c.errHandler("model.UpdateBannerByCode", err, utils.ErrUpdatingBanner)
	}

	return nil
}

func (c *Contract) DeleteBannerByCode(db *pgxpool.Pool, ctx context.Context, code string) error {
	var (
		err error
		sql = `
		UPDATE banners 
		SET updated_date=$1, deleted_date=$2 
		WHERE banner_code=$3`
	)
	_, err = db.Exec(ctx, sql, time.Now().In(time.UTC), time.Now().In(time.UTC), code)
	if err != nil {
		return c.errHandler("model.DeleteBannerByCode", err, utils.ErrUpdatingBanner)
	}

	return nil
}
