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

type RoleEnt struct {
	ID          int          `db:"id"`
	RoleCode    string       `db:"admin_code"`
	Name        string       `db:"name"`
	Description string       `db:"description"`
	Status      string       `db:"status"`
	CreatedDate time.Time    `db:"created_date"`
	UpdatedDate sql.NullTime `db:"updated_date"`
	DeletedDate sql.NullTime `db:"deleted_date"`
}

func (c *Contract) GetRoleList(db *pgxpool.Pool, ctx context.Context, param request.RoleParam) ([]RoleEnt, request.RoleParam, error) {
	var (
		err        error
		list       []RoleEnt
		where      []string
		paramQuery []interface{}
		totalData  int

		query = `SELECT 
		id,role_code, "name",  description, status, created_date, updated_date, deleted_date
		FROM roles`
	)

	// Populate Search
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		orWhere = append(orWhere, fmt.Sprintf("name iLIKE $%d", len(paramQuery)))
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
			return list, param, c.errHandler("model.GetRoleList", err, utils.ErrCountingListRole)
		}
		param.Count = totalData
	}

	if param.Limit > 0 && param.Page > 0 {
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
	}

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, param, c.errHandler("model.GetRoleList", err, utils.ErrGettingListRole)
	}

	defer rows.Close()
	for rows.Next() {
		var data RoleEnt
		err = rows.Scan(&data.ID, &data.RoleCode, &data.Name, &data.Description, &data.Status, &data.CreatedDate, &data.UpdatedDate, &data.DeletedDate)
		if err != nil {
			return list, param, c.errHandler("model.GetRoleList", err, utils.ErrScanningListRole)
		}
		list = append(list, data)
	}
	return list, param, nil
}

func (c *Contract) GetRoleByCode(db *pgxpool.Pool, ctx context.Context, roleCode string) (RoleEnt, error) {
	var (
		err  error
		data RoleEnt
		sql  = `SELECT id, role_code, "name",  description, status, created_date, updated_date, deleted_date
		FROM roles 
		WHERE role_code = $1`
	)

	err = db.QueryRow(ctx, sql, roleCode).Scan(&data.ID, &data.RoleCode, &data.Name, &data.Description, &data.Status, &data.CreatedDate, &data.UpdatedDate, &data.DeletedDate)
	if err != nil {
		return data, c.errHandler("model.GetRoleByCode", err, utils.ErrGettingRoleByCode)
	}

	return data, nil
}

func (c *Contract) GetRoleById(db *pgxpool.Pool, ctx context.Context, roleId int) (RoleEnt, error) {
	var (
		err  error
		data RoleEnt
		sql  = `SELECT id, role_code, "name",  description, status, created_date, updated_date, deleted_date
		FROM roles 
		WHERE id = $1`
	)

	err = db.QueryRow(ctx, sql, roleId).Scan(&data.ID, &data.RoleCode, &data.Name, &data.Description, &data.Status, &data.CreatedDate, &data.UpdatedDate, &data.DeletedDate)
	if err != nil {
		return data, c.errHandler("model.GetRoleById", err, utils.ErrGettingRoleById)
	}

	return data, nil
}

func (c *Contract) AddRole(db pgx.Tx, ctx context.Context, roleCode, name, description, status string) (int, error) {
	var id int

	// Insert data to database
	sql := `INSERT INTO roles(role_code, "name",  description, status, created_date)
	VALUES($1,$2,$3,$4,$5)`

	err := db.QueryRow(ctx, sql, roleCode, name, description, status, time.Now().In(time.UTC)).Scan(id)
	if err != nil {
		return id, c.errHandler("model.AddRole", err, utils.ErrAddingRole)
	}

	return id, nil
}

func (c *Contract) UpdateRoleByCode(db pgx.Tx, ctx context.Context, roleCode, name, description, status string) error {
	var (
		err error
		sql = `
		UPDATE roles 
		SET name=$1,description=$2,status=$3,updated_date=$4
		WHERE role_code=$5`
	)
	_, err = db.Exec(ctx, sql, name, description, status, time.Now().In(time.UTC), roleCode)
	if err != nil {
		return c.errHandler("model.UpdateRoleByCode", err, utils.ErrUpdatingRole)
	}

	return nil
}

func (c *Contract) DeleteRoleByCode(db *pgxpool.Pool, ctx context.Context, roleCode string) error {
	var (
		err error
		sql = `
		UPDATE roles 
		SET updated_date=$1, deleted_date=$2 
		WHERE role_code=$3`
	)
	_, err = db.Exec(ctx, sql, time.Now().In(time.UTC), time.Now().In(time.UTC), roleCode)
	if err != nil {
		return c.errHandler("model.DeleteRoleByCode", err, utils.ErrUpdatingRole)
	}

	return nil
}
