package model

import (
	"context"
	"dots-api/lib/utils"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RolePermissionEnt struct {
	ID           int       `db:"id"`
	RoleId       int       `db:"role_id"`
	PermissionId int       `db:"permission_id"`
	CreatedDate  time.Time `db:"created_date"`
	RoleEnt
	PermissionEnt
}

func (c *Contract) GetRolePermissionByRoleId(db *pgxpool.Pool, ctx context.Context, roleId int) ([]PermissionEnt, error) {
	var (
		err  error
		list []PermissionEnt

		query = `SELECT p.id, p.permission_code,
		p."name" as permission_name, p.route_pattern ,p.route_method , p.description,
		p.status, p.created_date , p.updated_date , p.deleted_date 
		FROM role_permissions rp 
		inner join roles r on r.id = rp.role_id 
		inner join permissions p ON p.id = rp.permission_id and p.deleted_date is null and p.status = 'active'
		where rp.role_id = $1`
	)

	rows, err := db.Query(ctx, query, roleId)
	if err != nil {
		return list, c.errHandler("model.GetRoleList", err, utils.ErrGettingListRolePermission)
	}

	defer rows.Close()
	for rows.Next() {
		var data PermissionEnt
		err = rows.Scan(
			&data.ID,
			&data.PermissionCode,
			&data.Name,
			&data.RoutePattern,
			&data.RouteMethod,
			&data.Description,
			&data.Status,
			&data.CreatedDate,
			&data.UpdatedDate,
			&data.DeletedDate,
		)
		if err != nil {
			return list, c.errHandler("model.GetRoleList", err, utils.ErrScanningListRolePermission)
		}
		list = append(list, data)
	}
	return list, nil
}

func (c *Contract) AddRolePermission(db pgx.Tx, ctx context.Context, roleId, permissionId int) error {

	// Insert data to database
	sql := `INSERT INTO role_permissions(role_id,permission_id, created_date)
	VALUES($1,$2,$3)`

	_, err := db.Exec(ctx, sql, roleId, permissionId, time.Now().In(time.UTC))
	if err != nil {
		return c.errHandler("model.AddRolePermission", err, utils.ErrAddingRolePermission)
	}

	return nil
}

func (c *Contract) DeleteRolePermissionByRoleId(db pgx.Tx, ctx context.Context, roleId int) error {
	var (
		err error
		sql = `
		DELETE FROM role_permissions
		WHERE role_id=$1`
	)
	_, err = db.Exec(ctx, sql, roleId)
	if err != nil {
		return c.errHandler("model.DeleteRolePermissionByRoleId", err, utils.ErrDeletingRolePermissionByRoleID)
	}

	return nil
}
