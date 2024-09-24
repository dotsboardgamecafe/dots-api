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

type CafeEnt struct {
	Id          int64        `db:"id"`
	CafeCode    string       `db:"cafe_code"`
	Name        string       `db:"name"`
	Address     string       `db:"address"`
	Description string       `db:"description"`
	Status      string       `db:"status"`
	Province    string       `db:"province"`
	City        string       `db:"city"`
	CreatedDate time.Time    `db:"created_date"`
	UpdatedDate sql.NullTime `db:"updated_date"`
	DeletedDate sql.NullTime `db:"updated_date"`
}

func (c *Contract) GetCafeList(db *pgxpool.Pool, ctx context.Context, param request.CafeParam) ([]CafeEnt, request.CafeParam, error) {
	var (
		err        error
		list       []CafeEnt
		where      []string
		paramQuery []interface{}
		totalData  int

		query = `SELECT 
		cafe_code, name, address, description, status, province, city
		FROM cafes`
	)

	// Populate Search
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		orWhere = append(orWhere, fmt.Sprintf("name iLIKE $%d", len(paramQuery)))
		orWhere = append(orWhere, fmt.Sprintf("address iLIKE $%d", len(paramQuery)))
		orWhere = append(orWhere, fmt.Sprintf("description iLIKE $%d", len(paramQuery)))
		where = append(where, "("+strings.Join(orWhere, " OR ")+")")
	}
	if len(param.Status) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Status)
		orWhere = append(orWhere, fmt.Sprintf("status = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}
	if len(param.Location) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.Location)
		orWhere = append(orWhere, fmt.Sprintf("city = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// Handling Soft Delete
	where = append(where, "deleted_date IS NULL")

	// Append All Where Conditions
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS data`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetCafeList", err, utils.ErrCountingListCafe)
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
		return list, param, c.errHandler("model.GetCafeList", err, utils.ErrGettingListCafe)
	}

	defer rows.Close()
	for rows.Next() {
		var data CafeEnt
		err = rows.Scan(&data.CafeCode, &data.Name, &data.Address, &data.Description, &data.Status, &data.Province, &data.City)
		if err != nil {
			return list, param, c.errHandler("model.GetCafeList", err, utils.ErrScanningListCafe)
		}
		list = append(list, data)
	}
	return list, param, nil
}

func (c *Contract) GetCafeByCode(db *pgxpool.Pool, ctx context.Context, code string) (CafeEnt, error) {
	var (
		err  error
		data CafeEnt
		sql  = `SELECT cafe_code, name, address, description, status, province, city
		FROM cafes 
		WHERE cafe_code = $1`
	)
	err = db.QueryRow(ctx, sql, code).Scan(&data.CafeCode, &data.Name, &data.Address, &data.Description, &data.Status, &data.Province, &data.City)
	if err != nil {
		return data, c.errHandler("model.GetCafeByCode", err, utils.ErrGettingCafeByCode)
	}

	return data, nil
}

func (c *Contract) AddCafe(db *pgxpool.Pool, ctx context.Context, code, name, address, desc, status, province, city string) error {
	var (
		err error
		id  int64
		sql = `INSERT INTO cafes(cafe_code, name, address, description, status, province, city, created_date)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`
	)

	err = db.QueryRow(ctx, sql, code, name, address, desc, status, province, city, time.Now().In(time.UTC)).Scan(&id)
	if err != nil {
		return c.errHandler("model.AddCafe", err, utils.ErrAddingCafe)
	}

	return nil
}

func (c *Contract) UpdateCafeByCode(db *pgxpool.Pool, ctx context.Context, code, name, address, desc, status, province, city string) error {
	var (
		err error
		sql = `
		UPDATE cafes 
		SET name=$1,address=$2,description=$3,status=$4,province=$5,city=$6,updated_date=$7
		WHERE cafe_code=$8`
	)

	_, err = db.Exec(ctx, sql, name, address, desc, status, province, city, time.Now().In(time.UTC), code)
	if err != nil {
		return c.errHandler("model.UpdateCafeByCode", err, utils.ErrUpdatingCafe)
	}

	return nil
}

func (c *Contract) DeleteCafeByCode(db *pgxpool.Pool, ctx context.Context, code string) error {
	var (
		err error
		sql = `
		UPDATE cafes 
		SET updated_date=$1, deleted_date=$2 
		WHERE cafe_code=$3`
	)
	_, err = db.Exec(ctx, sql, time.Now().In(time.UTC), time.Now().In(time.UTC), code)
	if err != nil {
		return c.errHandler("model.DeleteCafeByCode", err, utils.ErrUpdatingCafe)
	}

	return nil
}

func (c *Contract) GetCafeIdByCode(db *pgxpool.Pool, ctx context.Context, code string) (int64, error) {
	var (
		err   error
		id    int64
		query = `SELECT id FROM cafes WHERE cafe_code = $1`
	)
	err = db.QueryRow(ctx, query, code).Scan(&id)
	if err != nil {
		return id, c.errHandler("model.GetCafeIdByCode", err, utils.ErrGettingCafeByCode)
	}
	return id, nil
}

func (c *Contract) GetCafeLocationCityByCode(db *pgxpool.Pool, ctx context.Context, code string) (string, error) {
	var (
		err   error
		city  string
		query = `SELECT city FROM cafes WHERE cafe_code = $1`
	)
	err = db.QueryRow(ctx, query, code).Scan(&city)
	if err != nil {
		return city, c.errHandler("model.GetCafeLocationCityByCode", err, utils.ErrGettingCafeCityByCode)
	}
	return city, nil
}
