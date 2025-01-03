package model

import (
	"context"
	"dots-api/bootstrap"
	"dots-api/lib/mail"
	"dots-api/lib/utils"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func (c *Contract) GenerateTokenJWT(userIdentifier, actorType, email string) (string, int64, error) {
	var (
		token, channel string
		expAt          int64
		err            error
	)

	key := c.Config.GetString("app.key")
	if len(key) == 0 {
		return token, expAt, errors.New(utils.ErrConfigKeyNotFound)
	}

	if actorType == utils.Admin {
		channel = ChannelCMS
	} else {
		channel = ChannelApp
	}

	expAt = time.Now().UTC().AddDate(7, 0, 0).Unix()
	claims := &bootstrap.CustomClaims{
		Code:    userIdentifier,
		Email:   email,
		Channel: channel,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAt,
			Issuer:    actorType,
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = rawToken.SignedString([]byte(key))
	if err != nil {
		return token, expAt, err
	}
	return token, expAt, nil
}

func (c *Contract) RegisterUser(db *pgxpool.Pool, ctx context.Context, fullName, dateOfBirth, gender, email, phoneNumber, password, userName, xPlayer string) (string, string, error) {
	var (
		err           error
		id            int64
		tierId        int64
		userInsertSQL string

		// Generate User Identifier
		userCode = utils.GeneratePrefixCode(utils.UserPrefix)

		// Replace hashPassword with your actual password hashing function
		passwordHash, _ = bcrypt.GenerateFromPassword([]byte(password), 14)

		existingUserId int64
	)

	// Check for existing unverified email and delete if found
	err = db.QueryRow(ctx, "SELECT id FROM users WHERE email = $1 AND status_verification = false;", email).Scan(&existingUserId)
	if err != nil && err != pgx.ErrNoRows {
		return userCode, email, c.errHandler("model.CheckExistingUser", err, "Error checking existing unverified user")
	}

	if existingUserId != 0 {
		_, err = db.Exec(ctx, "DELETE FROM users WHERE id = $1;", existingUserId)
		if err != nil {
			return userCode, email, c.errHandler("model.DeleteExistingUser", err, "Error deleting existing unverified user")
		}
	}

	// Check for existing verified email
	err = db.QueryRow(ctx, "SELECT id FROM users WHERE email = $1 AND status_verification = true;", email).Scan(&existingUserId)
	if err != nil && err != pgx.ErrNoRows {
		return userCode, email, c.errHandler("model.CheckExistingUser", err, "Error checking existing verified user")
	}

	if existingUserId != 0 {
		return userCode, email, errors.New(utils.ErrEmailAlreadyRegistered)
	}

	// Get the first Novice Tier to store as default value of registered users
	db.QueryRow(ctx, "SELECT id FROM tiers WHERE tier_code = 'TIER-001';").Scan(&tierId)

	// Insert user data into 'users' table
	userInsertSQL = `INSERT INTO users (user_code, fullname, date_of_birth, gender, email, phone_number, password, status_verification, status, created_date, latest_tier_id, username, x_player, role_id) 
        VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`
	err = db.QueryRow(ctx, userInsertSQL, userCode, fullName, dateOfBirth, gender, email, phoneNumber, passwordHash, false, "active", time.Now().In(time.UTC), tierId, userName, xPlayer, utils.RoleMemberId).Scan(&id)

	if err != nil {
		// Handle specific error cases
		switch {
		case strings.Contains(err.Error(), "users_username_key"):
			return userCode, email, errors.New(utils.ErrUsernameAlreadyRegistered)
		case strings.Contains(err.Error(), "users_email_key"):
			return userCode, email, errors.New(utils.ErrEmailAlreadyRegistered)
		case strings.Contains(err.Error(), "users_phone_number_key"):
			return userCode, email, errors.New(utils.ErrPhoneAlreadyRegistered)
		// Add other specific error cases here if needed
		default:
			return userCode, email, c.errHandler("model.RegisterUser", err, utils.ErrInsertingUser)
		}
	}

	return userCode, email, nil
}

func (c *Contract) UserLogin(db *pgxpool.Pool, ctx context.Context, email, password, xPlayer string) (UserEnt, RoleEnt, []PermissionEnt, int64, string, error) {
	var (
		err            error
		userData       UserEnt
		roleData       RoleEnt
		permissionList []PermissionEnt
		expAt          int64
		jwtToken       string
	)

	// Get user data by email
	userData, err = c.GetUserByEmail(db, ctx, email)
	if err != nil {
		if err.Error() == utils.EmptyData {
			return userData, roleData, permissionList, expAt, jwtToken, errors.New(utils.ErrInvalidEmailPassword)
		}
		return userData, roleData, permissionList, expAt, jwtToken, c.errHandler("model.UserLogin", err, err.Error())
	}

	// If email is not verified
	if !userData.StatusVerification {
		return userData, roleData, permissionList, expAt, jwtToken, errors.New(utils.ErrEmailNotVerified)
	}

	dataPassword := []byte(userData.Password)
	err = bcrypt.CompareHashAndPassword(dataPassword, []byte(password))
	if err != nil {
		return userData, roleData, permissionList, expAt, jwtToken, errors.New(utils.ErrInvalidEmailPassword)
	}

	// Generate JWT Token
	jwtToken, expAt, err = c.GenerateTokenJWT(userData.UserCode, utils.User, email)
	if err != nil {
		return userData, roleData, permissionList, expAt, jwtToken, c.errHandler("model.UserLogin", err, utils.ErrGeneratingJWT)
	}

	//get role data
	roleData, err = c.GetRoleById(db, ctx, userData.RoleId)
	if err != nil {
		return userData, roleData, permissionList, expAt, jwtToken, c.errHandler("model.UserLogin", err, utils.ErrGettingRoleById)
	}

	//get permission
	permissionList, err = c.GetRolePermissionByRoleId(db, ctx, userData.RoleId)
	if err != nil {
		return userData, roleData, permissionList, expAt, jwtToken, c.errHandler("model.UserLogin", err, utils.ErrGettingListRolePermission)
	}

	err = c.UpdateUserXPlayer(db, ctx, userData.UserCode, xPlayer)
	if err != nil {
		return userData, roleData, permissionList, expAt, jwtToken, c.errHandler("model.UserLogin", err, err.Error())
	}

	return userData, roleData, permissionList, expAt, jwtToken, nil
}

func (c *Contract) AdminLogin(db *pgxpool.Pool, ctx context.Context, email, password string) (AdminEnt, RoleEnt, []PermissionEnt, int64, string, error) {
	var (
		err            error
		userData       AdminEnt
		roleData       RoleEnt
		permissionList []PermissionEnt
		expAt          int64
		jwtToken       string
	)

	// Get user data by email
	userData, err = c.GetAdminByEmail(db, ctx, email)
	if err != nil {
		if err.Error() == utils.EmptyData {
			return userData, roleData, permissionList, expAt, jwtToken, errors.New(utils.ErrInvalidEmailPassword)
		}
		return userData, roleData, permissionList, expAt, jwtToken, c.errHandler("model.AdminLogin", err, err.Error())
	}

	dataPassword := []byte(userData.Password)
	err = bcrypt.CompareHashAndPassword(dataPassword, []byte(password))
	if err != nil {
		return userData, roleData, permissionList, expAt, jwtToken, errors.New(utils.ErrInvalidEmailPassword)
	}

	// Generate JWT Token
	jwtToken, expAt, err = c.GenerateTokenJWT(userData.AdminCode, utils.Admin, email)
	if err != nil {
		return userData, roleData, permissionList, expAt, jwtToken, c.errHandler("model.AdminLogin", err, utils.ErrGeneratingJWT)
	}

	//get role data
	roleData, err = c.GetRoleById(db, ctx, userData.RoleId)
	if err != nil {
		return userData, roleData, permissionList, expAt, jwtToken, c.errHandler("model.AdminLogin", err, utils.ErrGettingRoleById)
	}

	//get permission
	permissionList, err = c.GetRolePermissionByRoleId(db, ctx, userData.RoleId)
	if err != nil {
		return userData, roleData, permissionList, expAt, jwtToken, c.errHandler("model.AdminLogin", err, utils.ErrGettingListRolePermission)
	}

	return userData, roleData, permissionList, expAt, jwtToken, nil
}

func (c *Contract) CheckTokenAndExpiration(db *pgxpool.Pool, ctx context.Context, verificationType, actorType, token string) (UserEnt, RoleEnt, []PermissionEnt, int64, string, error) {
	var (
		jwtToken                string
		expAt                   int64
		checkVerificationSQL    string
		updateVerificationSQL   string
		email                   string
		expiredDateVerification time.Time
		isUsedData              bool
		dataUser                UserEnt
		roleData                RoleEnt
		permissionList          []PermissionEnt
		err                     error
	)

	// Check if token exists and is not used
	checkVerificationSQL = `SELECT email, expired_date, is_used FROM verifications WHERE token = $1 AND verification_type = $2`
	err = db.QueryRow(ctx, checkVerificationSQL, token, verificationType).Scan(&email, &expiredDateVerification, &isUsedData)
	if err != nil {
		return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiredForgotPassword", err, utils.ErrGettingVerificationsData)
	}

	// Check if token is used
	if isUsedData {
		return dataUser, roleData, permissionList, expAt, jwtToken, errors.New(utils.ErrTokenUsed)
	}

	// Check if token is expired
	if expiredDateVerification.Before(time.Now().UTC()) {
		return dataUser, roleData, permissionList, expAt, jwtToken, errors.New(utils.ErrTokenExpired)
	}

	if actorType == utils.User {
		// Get data user for create jwt token
		dataUser, err = c.GetUserByEmail(db, ctx, email)
		if err != nil {
			return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiration", err, utils.ErrGettingUserData)
		}

		//get role data
		roleData, err = c.GetRoleById(db, ctx, dataUser.RoleId)
		if err != nil {
			return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.UserLogin", err, utils.ErrGettingRoleById)
		}

		//get permission
		permissionList, err = c.GetRolePermissionByRoleId(db, ctx, dataUser.RoleId)
		if err != nil {
			return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.UserLogin", err, utils.ErrGettingListRolePermission)
		}
	}
	// Generate JWT Token
	jwtToken, expAt, err = c.GenerateTokenJWT(dataUser.UserCode, utils.User, email)
	if err != nil {
		return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiration", err, utils.ErrGeneratingJWT)
	}
	// Start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiration", err, utils.ErrBeginningTransaction)
	}

	// Update the verification record to mark it as used
	updateVerificationSQL = "UPDATE verifications SET is_used = $1 WHERE token = $2"
	_, err = tx.Exec(ctx, updateVerificationSQL, true, token)
	if err != nil {
		tx.Rollback(ctx)
		return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiredForgotPassword", err, utils.ErrMarkingToken)
	}

	switch verificationType {
	// Check is type of verify email
	case utils.VerifyRegistration:
		sql := `
			UPDATE users
			SET status_verification = $1
			WHERE user_code = $2
		`

		_, err := tx.Exec(ctx, sql, true, dataUser.UserCode)
		if err != nil {
			tx.Rollback(ctx)
			return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiredForgotPassword", err, utils.ErrUpdatingUserEmailStatus)
		}
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiredForgotPassword", err, utils.ErrCommittingTransaction)
	}

	return dataUser, roleData, permissionList, expAt, jwtToken, nil
}

func (c *Contract) CheckTokenAndExpirationUpdateEmail(db *pgxpool.Pool, ctx context.Context, verificationType, actorType, token, userCode string) (UserEnt, RoleEnt, []PermissionEnt, int64, string, error) {
	var (
		jwtToken                string
		expAt                   int64
		checkVerificationSQL    string
		updateVerificationSQL   string
		email                   string
		expiredDateVerification time.Time
		isUsedData              bool
		dataUser                UserEnt
		roleData                RoleEnt
		permissionList          []PermissionEnt
		err                     error
	)

	// Check if token exists and is not used
	checkVerificationSQL = `SELECT email, expired_date, is_used FROM verifications WHERE token = $1 AND verification_type = $2`
	err = db.QueryRow(ctx, checkVerificationSQL, token, verificationType).Scan(&email, &expiredDateVerification, &isUsedData)
	if err != nil {
		return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiredForgotPassword", err, utils.ErrGettingVerificationsData)
	}

	// Check if token is used
	if isUsedData {
		return dataUser, roleData, permissionList, expAt, jwtToken, errors.New(utils.ErrTokenUsed)
	}

	// Check if token is expired
	if expiredDateVerification.Before(time.Now().UTC()) {
		return dataUser, roleData, permissionList, expAt, jwtToken, errors.New(utils.ErrTokenExpired)
	}

	if actorType == utils.User {
		// Get data user for create jwt token
		dataUser, err = c.GetUserByUserCode(db, ctx, userCode)
		if err != nil {
			return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiration", err, utils.ErrGettingUserData)
		}
	}
	// Generate JWT Token
	jwtToken, expAt, err = c.GenerateTokenJWT(dataUser.UserCode, utils.User, email)
	if err != nil {
		return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiration", err, utils.ErrGeneratingJWT)
	}
	// Start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiration", err, utils.ErrBeginningTransaction)
	}

	// Update the verification record to mark it as used
	updateVerificationSQL = "UPDATE verifications SET is_used = $1 WHERE token = $2"
	_, err = tx.Exec(ctx, updateVerificationSQL, true, token)
	if err != nil {
		tx.Rollback(ctx)
		return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiredForgotPassword", err, utils.ErrMarkingToken)
	}

	sql := `
			UPDATE users
			SET email = $1 , status_verification = $2
			WHERE user_code = $3
		`

	_, err = tx.Exec(ctx, sql, email, true, userCode)
	if err != nil {
		tx.Rollback(ctx)
		return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiredForgotPassword", err, utils.ErrUpdatingUserEmail)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return dataUser, roleData, permissionList, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiredForgotPassword", err, utils.ErrCommittingTransaction)
	}

	// Set New Email
	dataUser.Email.String = email
	return dataUser, roleData, permissionList, expAt, jwtToken, nil
}

func (c *Contract) RequestForgotPassword(db *pgxpool.Pool, ctx context.Context, email string) error {
	var (
		err error
		// Generate Token 50 digits
		token, _ = utils.Generate(`[a-zA-Z0-9]{50}`)

		// Forgot password route
		linkNewPass = c.Config.GetString("web_url") + utils.ResetPassRoute + token + utils.TypeRoute + utils.ForgotPassword

		// Determine expired at
		expAt = time.Now().UTC().Add(time.Minute * 5)

		// Import contract send email
		mailContract = mail.New(c.App)
	)

	// Check email and get user data
	userData, err := c.GetUserByEmail(db, ctx, email)
	if err != nil {
		return c.errHandler("model.RequestForgotPassword", err, utils.ErrInvalidEmailPassword)
	}

	// Sending Forgot Password Mail
	err = mailContract.SendMail(mail.UserForgotPassword, mail.MailSubj[mail.UserForgotPassword], email, mail.EmailData{Name: userData.FullName, Email: email, Link: linkNewPass})
	if err != nil {
		return c.errHandler("model.RequestForgotPassword", err, utils.ErrSendingResetPasswordEmail)
	}

	// Insert verification data into 'verifications' table
	err = c.insertVerificationData(db, ctx, utils.User, utils.ForgotPassword, email, token, false, expAt)
	if err != nil {
		return c.errHandler("model.RequestForgotPassword", err, utils.ErrAddingResetPasswordVerification)
	}

	return nil
}

func (c *Contract) RequestVerifyEmailUser(db *pgxpool.Pool, ctx context.Context, email, types string) error {
	var (
		err error
		// Generate Token 50 digits
		token, _ = utils.Generate(`[a-zA-Z0-9]{50}`)

		// Determine expired at
		expAt = time.Now().UTC().Add(time.Minute * 5)

		// Import contract send email
		mailContract = mail.New(c.App)
	)

	// Check email and get user data
	userData, err := c.GetUserByEmail(db, ctx, email)
	if err != nil {
		return c.errHandler("model.RequestVerifyEmailUser", err, utils.ErrInvalidEmailPassword)
	}

	switch types {
	case utils.VerifyRegistration:
		link := c.Config.GetString("web_url") + utils.VerifyTokenRoute + token + utils.TypeRoute + types

		err = mailContract.SendMail(mail.UserVerifyEmail, mail.MailSubj[mail.UserVerifyEmail], email, mail.EmailData{Name: userData.FullName, Email: email, Link: link})
		if err != nil {
			return c.errHandler("model.RequestVerifyEmailUser", err, utils.ErrSendingVerifyEmail)
		}
	case utils.ForgotPassword:
		link := c.Config.GetString("web_url") + utils.ForgotPasswordRoute + token + utils.TypeRoute + types

		err = mailContract.SendMail(mail.UserForgotPassword, mail.MailSubj[mail.UserForgotPassword], email, mail.EmailData{Name: userData.FullName, Email: email, Link: link})
		if err != nil {
			return c.errHandler("model.RequestVerifyEmailUser", err, utils.ErrSendingForgotPasswordEmail)
		}
	default:
		return errors.New(utils.ErrInvalidSendingEmailType)
	}

	err = c.insertVerificationData(db, ctx, utils.User, types, email, token, false, expAt)
	if err != nil {
		return c.errHandler("model.RequestVerifyEmailUser", err, utils.ErrAddingResetPasswordVerification)
	}

	return nil
}

func (c *Contract) RequestVerifyUpdateEmailUser(db *pgxpool.Pool, ctx context.Context, email, types, userCode string) error {
	var (
		err error
		// Generate Token 50 digits
		token, _ = utils.Generate(`[a-zA-Z0-9]{50}`)

		// Determine expired at
		expAt = time.Now().UTC().Add(time.Minute * 5)

		// Import contract send email
		mailContract = mail.New(c.App)
	)

	// Check email and get user data
	userData, err := c.GetUserByUserCode(db, ctx, userCode)
	if err != nil {
		return c.errHandler("model.RequestVerifyEmailUser", err, utils.ErrInvalidEmailPassword)
	}

	link := c.Config.GetString("web_url") + utils.VerifyTokenEmailRoute + token + utils.TypeRoute + types + utils.UserCodeRoute + userCode

	err = mailContract.SendMail(mail.UserUpdateEmail, mail.MailSubj[mail.UserUpdateEmail], email, mail.EmailData{Name: userData.FullName, Email: email, Link: link})
	if err != nil {
		return c.errHandler("model.RequestVerifyEmailUser", err, utils.ErrSendingUpdateEmail)
	}

	err = c.insertVerificationData(db, ctx, utils.User, types, email, token, false, expAt)
	if err != nil {
		return c.errHandler("model.RequestVerifyEmailUser", err, utils.ErrAddingResetPasswordVerification)
	}

	return nil
}

func (c *Contract) ResetPassword(db *pgxpool.Pool, ctx context.Context, userCode, NewPassword, ConfirmPassword string) error {
	var (
		err      error
		dataUser UserEnt
	)
	// Check if new password matches the confirmation
	if NewPassword != ConfirmPassword {
		return errors.New(utils.ErrPasswordMismatch)
	}

	dataUser, err = c.GetUserByUserCode(db, ctx, userCode)
	if err != nil {
		return c.errHandler("model.ResetPassword", err, utils.ErrFetchingUserPassword)
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(NewPassword), 14)
	if err != nil {
		return c.errHandler("model.ResetPassword", err, utils.ErrHashingPassword)
	}

	// Update the user's password in the database
	sql := "UPDATE users SET password = $1, updated_date = $3 WHERE id = $2"
	_, err = db.Exec(ctx, sql, string(hashedPassword), dataUser.ID, time.Now().UTC())
	if err != nil {
		return c.errHandler("model.ResetPassword", err, utils.ErrUpdatingUserPassword)
	}

	return nil
}

func (c *Contract) insertVerificationData(db *pgxpool.Pool, ctx context.Context, actorType, verificationType, email, token string, isUsed bool, expiredDate time.Time) error {
	sql := `INSERT INTO verifications(actor_type, verification_type, email, token, is_used, expired_date, created_date)
        VALUES($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.Exec(ctx, sql, actorType, verificationType, email, token, isUsed, expiredDate, time.Now().In(time.UTC))
	if err != nil {
		return err
	}
	return nil
}

func (c *Contract) VerifyUserPassword(db *pgxpool.Pool, ctx context.Context, userCode string, inputtedPassword string) error {
	var (
		err      error
		dataUser UserEnt
	)

	dataUser, err = c.GetUserByUserCode(db, ctx, userCode)
	if err != nil {
		return c.errHandler("model.verifyUserPassword", err, utils.ErrFetchingUserPassword)
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte(inputtedPassword))
	if err != nil {
		return errors.New("Password is incorrect")
	}

	return nil
}
