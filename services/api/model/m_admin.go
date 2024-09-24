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
	"golang.org/x/crypto/bcrypt"
)

type AdminEnt struct {
	ID          int          `db:"id"`
	AdminCode   string       `db:"admin_code"`
	Email       string       `db:"email"`
	Name        string       `db:"name"`
	UserName    string       `db:"username"`
	PhoneNumber string       `db:"phone_number"`
	Password    string       `db:"password"`
	Status      string       `db:"status"`
	ImageURL    string       `db:"image_url"`
	RoleId      int          `db:"role_id"`
	CreatedDate time.Time    `db:"created_date"`
	UpdatedDate sql.NullTime `db:"updated_date"`
	DeletedDate sql.NullTime `db:"deleted_date"`
}

func (c *Contract) GetAdminList(db *pgxpool.Pool, ctx context.Context, param request.AdminParam) ([]AdminEnt, request.AdminParam, error) {
	var (
		err        error
		list       []AdminEnt
		where      []string
		paramQuery []interface{}
		totalData  int

		query = `SELECT 
		admin_code, email, name, username, status, image_url, phone_number
		FROM admins`
	)

	// Populate Search
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		orWhere = append(orWhere, fmt.Sprintf("name iLIKE $%d", len(paramQuery)))
		orWhere = append(orWhere, fmt.Sprintf("email iLIKE $%d", len(paramQuery)))
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
			return list, param, c.errHandler("model.GetAdminList", err, utils.ErrCountingListAdmin)
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
		return list, param, c.errHandler("model.GetAdminList", err, utils.ErrGettingListAdmin)
	}

	defer rows.Close()
	for rows.Next() {
		var data AdminEnt
		err = rows.Scan(&data.AdminCode, &data.Email, &data.Name, &data.UserName, &data.Status, &data.ImageURL, &data.PhoneNumber)
		if err != nil {
			return list, param, c.errHandler("model.GetAdminList", err, utils.ErrScanningListAdmin)
		}
		list = append(list, data)
	}
	return list, param, nil
}

func (c *Contract) GetAdminIdByCode(db *pgxpool.Pool, ctx context.Context, adminCode string) (int64, error) {
	var (
		err   error
		id    int64
		query = `SELECT id FROM admins WHERE admin_code = $1`
	)
	err = db.QueryRow(ctx, query, adminCode).Scan(&id)
	if err != nil {
		return id, c.errHandler("model.GetAdminIdByCode", err, utils.ErrGettingAdminByCode)
	}

	return id, nil
}

func (c *Contract) GetAdminByCode(db *pgxpool.Pool, ctx context.Context, adminCode string) (AdminEnt, error) {
	var (
		err  error
		data AdminEnt
		sql  = `SELECT admin_code, email, name, username, password, status, phone_number, image_url
		FROM admins 
		WHERE admin_code = $1`
	)

	err = db.QueryRow(ctx, sql, adminCode).Scan(&data.AdminCode, &data.Email, &data.Name, &data.UserName, &data.Password, &data.Status, &data.PhoneNumber, &data.ImageURL)
	if err != nil {
		if err != pgx.ErrNoRows {
			return data, nil
		}
		return data, c.errHandler("model.GetAdminByCode", err, utils.ErrGettingAdminByCode)
	}

	return data, nil
}

func (c *Contract) GetAdminByEmail(db *pgxpool.Pool, ctx context.Context, email string) (AdminEnt, error) {
	var (
		err      error
		data     AdminEnt
		sqlQuery = `SELECT admin_code, email, name, username, password, status, phone_number, image_url, role_id
		FROM admins 
		WHERE email = $1`
	)

	err = db.QueryRow(ctx, sqlQuery, email).Scan(&data.AdminCode, &data.Email, &data.Name, &data.UserName, &data.Password, &data.Status, &data.PhoneNumber, &data.ImageURL, &data.RoleId)
	if err != nil && err != pgx.ErrNoRows {
		return data, c.errHandler("model.GetAdminByEmail", err, utils.ErrGettingAdminByEmail)
	}

	return data, nil
}

func (c *Contract) GetAdminByPhoneNumber(db *pgxpool.Pool, ctx context.Context, phoneNumber string) (AdminEnt, error) {
	var (
		err      error
		data     AdminEnt
		sqlQuery = `SELECT admin_code, email, name, username, password, status, phone_number, image_url
		FROM admins 
		WHERE phone_number = $1`
	)

	err = db.QueryRow(ctx, sqlQuery, phoneNumber).Scan(&data.AdminCode, &data.Email, &data.Name, &data.UserName, &data.Password, &data.Status, &data.PhoneNumber, &data.ImageURL)
	if err != nil {
		if err != sql.ErrNoRows && err != pgx.ErrNoRows {
			return data, c.errHandler("model.GetAdminByPhoneNumber", err, utils.ErrGettingAdminByPhone)
		}
	}

	return data, nil
}

func (c *Contract) AddAdmin(db *pgxpool.Pool, ctx context.Context, adminCode, email, name, userName, password, status, phoneNumber, imageURL string) error {
	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.errHandler("model.UpdatePasswordUser", err, utils.ErrHashingPassword)
	}

	// Insert data to database
	sql := `INSERT INTO admins(admin_code, email, name, username, password, status, phone_number, image_url, role_id, created_date)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	_, err = db.Exec(ctx, sql, adminCode, email, name, userName, hashedPassword, status, phoneNumber, imageURL, utils.RoleAdminId, time.Now().In(time.UTC))
	if err != nil {
		return c.errHandler("model.AddAdmin", err, utils.ErrAddingAdmin)
	}

	return nil
}

func (c *Contract) UpdateAdminByCode(db *pgxpool.Pool, ctx context.Context, adminCode, email, name, userName, password, status, phoneNumber, imageURL string) error {
	var (
		err error
		sql = `
		UPDATE admins 
		SET email=$1,name=$2,username=$3,password=$4,status=$5,phone_number=$6,image_url=$7,updated_date=$8
		WHERE admin_code = $9`
	)

	// Validate Code
	adminExist, err := c.GetAdminByCode(db, ctx, adminCode)
	if err != nil {
		return err
	}

	//ignore update
	if len(email) < 1 {
		email = adminExist.Email
	}

	if len(phoneNumber) < 1 {
		phoneNumber = adminExist.PhoneNumber
	}

	if len(name) < 1 {
		name = adminExist.Name
	}

	if len(userName) < 1 {
		userName = adminExist.UserName
	}

	if len(imageURL) < 1 {
		imageURL = adminExist.ImageURL
	}

	if len(password) < 1 {
		password = adminExist.Password
	} else {
		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return c.errHandler("model.UpdateAdminByCode", err, utils.ErrHashingPassword)
		}

		password = string(hashedPassword)
	}

	if len(status) < 1 {
		status = adminExist.Status
	}

	_, err = db.Exec(ctx, sql, email, name, userName, password, status, phoneNumber, imageURL, time.Now().In(time.UTC), adminCode)
	if err != nil {
		return c.errHandler("model.UpdateAdminByCode", err, utils.ErrUpdatingAdmin)
	}

	return nil
}

func (c *Contract) DeleteAdminByCode(db *pgxpool.Pool, ctx context.Context, adminCode string) error {
	var (
		err error
		sql = `
		UPDATE admins 
		SET updated_date=$1, deleted_date=$2 
		WHERE admin_code=$3`
	)
	_, err = db.Exec(ctx, sql, time.Now().In(time.UTC), time.Now().In(time.UTC), adminCode)
	if err != nil {
		return c.errHandler("model.DeleteAdminByCode", err, utils.ErrUpdatingAdmin)
	}

	return nil
}

func (c *Contract) UpdateAdminStatus(db *pgxpool.Pool, ctx context.Context, userCode, status string) error {
	sql := `
		UPDATE admins
		SET status = $1
		WHERE admin_code = $2
	`

	_, err := db.Exec(ctx, sql, status, userCode)
	if err != nil {
		return c.errHandler("model.UpdateAdminStatus", err, utils.ErrUpdatingAdminStatus)
	}

	return nil
}

func (c *Contract) IsAdminUsernameExist(db *pgxpool.Pool, ctx context.Context, userName string) (bool, error) {
	var (
		isExist = false

		query = `SELECT id FROM admins WHERE username = $1`
	)

	userId := ""
	_ = db.QueryRow(ctx, query, userName).Scan(&userId)
	if userId != "" {
		return true, nil
	}

	return isExist, nil
}
