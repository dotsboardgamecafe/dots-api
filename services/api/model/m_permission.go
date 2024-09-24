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

type PermissionEnt struct {
	ID             int64        `db:"id"`
	PermissionCode string       `db:"admin_code"`
	Name           string       `db:"name"`
	RoutePattern   string       `db:"route_pattern"`
	RouteMethod    string       `db:"route_method"`
	Description    string       `db:"description"`
	Status         string       `db:"status"`
	CreatedDate    time.Time    `db:"created_date"`
	UpdatedDate    sql.NullTime `db:"updated_date"`
	DeletedDate    sql.NullTime `db:"deleted_date"`
}

func (c *Contract) GetPermissionList(db *pgxpool.Pool, ctx context.Context, param request.PermissionParam) ([]PermissionEnt, request.PermissionParam, error) {
	var (
		err        error
		list       []PermissionEnt
		where      []string
		paramQuery []interface{}
		totalData  int

		query = `SELECT 
		id,permission_code, "name", route_pattern, route_method, description, status, created_date, updated_date, deleted_date
		FROM permissions`
	)

	// Populate Search
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		orWhere = append(orWhere, fmt.Sprintf("name iLIKE $%d", len(paramQuery)))
		orWhere = append(orWhere, fmt.Sprintf("route_pattern iLIKE $%d", len(paramQuery)))
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

	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS data`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, param, c.errHandler("model.GetPermissionList", err, utils.ErrCountingListPermission)
		}
		param.Count = totalData
	}

	if param.Limit > 0 && param.Page > 0 {
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
	}

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, param, c.errHandler("model.GetPermissionList", err, utils.ErrGettingListPermission)
	}

	defer rows.Close()
	for rows.Next() {
		var data PermissionEnt
		err = rows.Scan(&data.ID, &data.PermissionCode, &data.Name, &data.RoutePattern, &data.RouteMethod, &data.Description, &data.Status, &data.CreatedDate, &data.UpdatedDate, &data.DeletedDate)
		if err != nil {
			return list, param, c.errHandler("model.GetPermissionList", err, utils.ErrScanningListPermission)
		}
		list = append(list, data)
	}
	return list, param, nil
}

func (c *Contract) GetPermissionByCode(db *pgxpool.Pool, ctx context.Context, permissionCode string) (PermissionEnt, error) {
	var (
		err  error
		data PermissionEnt
		sql  = `SELECT id, permission_code, "name", route_pattern, route_method, description, status, created_date, updated_date, deleted_date
		FROM permissions 
		WHERE permission_code = $1`
	)

	err = db.QueryRow(ctx, sql, permissionCode).Scan(&data.ID, &data.PermissionCode, &data.Name, &data.RoutePattern, &data.RouteMethod, &data.Description, &data.Status, &data.CreatedDate, &data.UpdatedDate, &data.DeletedDate)
	if err != nil {
		return data, c.errHandler("model.GetPermissionByCode", err, utils.ErrGettingPermissionByCode)
	}

	return data, nil
}

func (c *Contract) AddPermission(db *pgxpool.Pool, ctx context.Context, permissionCode, name, routePattern, routeMethod, description, status string) error {

	// Insert data to database
	sql := `INSERT INTO permissions(permission_code, "name", route_pattern, route_method, description, status, created_date)
	VALUES($1,$2,$3,$4,$5,$6,$7)`

	_, err := db.Exec(ctx, sql, permissionCode, name, routePattern, routeMethod, description, status, time.Now().In(time.UTC))
	if err != nil {
		return c.errHandler("model.AddPermission", err, utils.ErrAddingPermission)
	}

	return nil
}

func (c *Contract) UpdatePermissionByCode(db *pgxpool.Pool, ctx context.Context, permissionCode, name, routePattern, routeMethod, description, status string) error {
	var (
		err error
		sql = `
		UPDATE permissions 
		SET name=$1,route_pattern=$2,route_method=$3,description=$4,status=$5,updated_date=$6
		WHERE permission_code=$7`
	)
	_, err = db.Exec(ctx, sql, name, routePattern, routeMethod, description, status, time.Now().In(time.UTC), permissionCode)
	if err != nil {
		return c.errHandler("model.UpdatePermissionByCode", err, utils.ErrUpdatingPermission)
	}

	return nil
}

func (c *Contract) DeletePermissionByCode(db *pgxpool.Pool, ctx context.Context, permissionCode string) error {
	var (
		err error
		sql = `
		UPDATE permissions 
		SET updated_date=$1, deleted_date=$2 
		WHERE permission_code=$3`
	)
	_, err = db.Exec(ctx, sql, time.Now().In(time.UTC), time.Now().In(time.UTC), permissionCode)
	if err != nil {
		return c.errHandler("model.DeletePermissionByCode", err, utils.ErrUpdatingPermission)
	}

	return nil
}
