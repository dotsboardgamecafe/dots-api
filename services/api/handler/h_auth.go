package handler

import (
	"context"
	"dots-api/bootstrap"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"dots-api/services/api/request"
	"dots-api/services/api/response"
	"net/http"
	"time"
)

func (h *Contract) LoginUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		ctx = context.TODO()
		m   = model.Contract{App: h.App}
		req = request.LoginUserReq{}
		res = response.LoginUserRes{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	switch r.Header.Get("X-Actor-Type") {
	case utils.User:

		// check x player token
		xPlayer := r.Header.Get(bootstrap.XPlayer)

		dataUser, role, permissions, expAtUnix, jwtToken, err := m.UserLogin(h.DB, ctx, req.Email, req.Password, xPlayer)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		// convert unix timestamp
		expAt := time.Unix(expAtUnix, 0)
		res = response.LoginUserRes{
			Token:       jwtToken,
			ImageURL:    dataUser.ImageURL.String,
			DateOfBirth: dataUser.DateOfBirth.String,
			Gender:      dataUser.Gender.String,
			UserCode:    dataUser.UserCode,
			FullName:    dataUser.FullName,
			PhoneNumber: dataUser.PhoneNumber,
			Email:       dataUser.Email.String,
			ExpiredAt:   expAt.Format(utils.DATE_TIME_FORMAT),
			ActorType:   utils.User,
			CreatedDate: time.Now().In(time.UTC).Format(utils.DATE_TIME_FORMAT),
			RoleId:      dataUser.RoleId,
			RoleCode:    role.RoleCode,
			Permissions: ToResponsePermissionList(permissions),
		}

		// Populate Response
		h.SendSuccess(w, res, nil)

	case utils.Admin:
		dataAdmin, role, permissions, expAtUnix, jwtToken, err := m.AdminLogin(h.DB, ctx, req.Email, req.Password)
		if err != nil {
			h.SendBadRequest(w, err.Error())
			return
		}

		// convert unix timestamp
		expAt := time.Unix(expAtUnix, 0)
		res = response.LoginUserRes{
			Token:       jwtToken,
			UserCode:    dataAdmin.AdminCode,
			FullName:    dataAdmin.Name,
			Email:       dataAdmin.Email,
			ExpiredAt:   expAt.Format(utils.DATE_TIME_FORMAT),
			ActorType:   utils.Admin,
			CreatedDate: time.Now().In(time.UTC).Format(utils.DATE_TIME_FORMAT),
			RoleId:      dataAdmin.RoleId,
			RoleCode:    role.RoleCode,
			Permissions: ToResponsePermissionList(permissions),
		}

		// Populate Response
		h.SendSuccess(w, res, nil)
	default:
		h.SendBadRequest(w, err.Error())
		return
	}
}

func (h *Contract) RegisterUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		ctx = context.TODO()
		m   = model.Contract{App: h.App}
		req = request.RegisterUserReq{}
		res = response.RegisterUserRes{}
	)

	xPlayer := r.Header.Get(bootstrap.XPlayer)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	// Compare Password
	if req.Password != req.ConfirmPassword {
		h.SendBadRequest(w, utils.ErrPasswordMismatch)
		return
	}

	userCode, email, err := m.RegisterUser(h.DB, ctx, req.Fullname, req.DateOfBirth, req.Gender, req.Email, req.PhoneNumber, req.ConfirmPassword, req.Username, xPlayer)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.RequestVerifyEmailUser(h.DB, ctx, req.Email, utils.VerifyRegistration)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	res = response.RegisterUserRes{
		UserCode:    userCode,
		Email:       email,
		CreatedDate: time.Now().In(time.UTC).Format(utils.DATE_TIME_FORMAT),
	}

	h.SendSuccess(w, res, nil)
}

func (h *Contract) RequestVerifyEmailUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		check bool
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		req   = request.RequestVerifyEmailReq{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if check = utils.Contains(utils.VerificationType, req.Type); !check {
		h.SendBadRequest(w, utils.ErrInvalidSendingEmailType)
		return
	}

	err = m.RequestVerifyEmailUser(h.DB, ctx, req.Email, req.Type)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

func (h *Contract) RequestVerifyUpdateEmailUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		check          bool
		ctx            = context.TODO()
		m              = model.Contract{App: h.App}
		req            = request.RequestVerifyEmailReq{}
		userIdentifier = bootstrap.GetIdentifierCodeFromToken(ctx, r)
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if check = (req.Type == utils.UpdateEmail); !check {
		h.SendBadRequest(w, utils.ErrInvalidUpdateEmailType)
		return
	}

	err = m.RequestVerifyUpdateEmailUser(h.DB, ctx, req.Email, req.Type, userIdentifier)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

func (h *Contract) VerifyTokenUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		ctx       = context.TODO()
		m         = model.Contract{App: h.App}
		jwtToken  string
		expAt     time.Time
		expAtUnix int64
		dataUser  model.UserEnt
		res       = response.LoginUserRes{}

		// Initiate Query Param
		param = map[string]interface{}{
			"type":  "",
			"token": "",
		}
	)

	if token, ok := r.URL.Query()["token"]; ok && len(token[0]) > 0 {
		param["token"] = token[0]
	}

	if types, ok := r.URL.Query()["type"]; ok && len(types) > 0 {
		paramType := types[0]
		if !utils.Contains(utils.VerificationType, paramType) {
			h.SendBadRequest(w, utils.ErrInvalidSendingEmailType)
			return
		}
		param["type"] = paramType
	} else {
		h.SendBadRequest(w, utils.ErrInvalidTypeQueryParameter)
		return
	}

	dataUser, role, permissions, expAtUnix, jwtToken, err := m.CheckTokenAndExpiration(h.DB, ctx, param["type"].(string), utils.User, param["token"].(string))
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	//convert unix timestamp
	expAt = time.Unix(expAtUnix, 0)

	// Populate response
	res = response.LoginUserRes{
		Token:       jwtToken,
		ImageURL:    dataUser.ImageURL.String,
		UserCode:    dataUser.UserCode,
		FullName:    dataUser.FullName,
		PhoneNumber: dataUser.PhoneNumber,
		Email:       dataUser.Email.String,
		ExpiredAt:   expAt.Format(utils.DATE_TIME_FORMAT),
		ActorType:   utils.User,
		CreatedDate: time.Now().In(time.UTC).Format(utils.DATE_TIME_FORMAT),
		RoleId:      dataUser.RoleId,
		RoleCode:    role.RoleCode,
		Permissions: ToResponsePermissionList(permissions),
	}

	h.SendSuccess(w, res, nil)
}

func (h *Contract) VerifyTokenUpdateEmailUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		ctx       = context.TODO()
		m         = model.Contract{App: h.App}
		jwtToken  string
		expAt     time.Time
		expAtUnix int64
		check     bool
		dataUser  model.UserEnt
		res       = response.LoginUserRes{}

		// Initiate Query Param
		param = map[string]interface{}{
			"usercode": "",
			"type":     "",
			"token":    "",
		}
	)

	if token, ok := r.URL.Query()["token"]; ok && len(token[0]) > 0 {
		param["token"] = token[0]
	}

	if userCode, ok := r.URL.Query()["usercode"]; ok && len(userCode[0]) > 0 {
		param["usercode"] = userCode[0]
	}

	if types, ok := r.URL.Query()["type"]; ok && len(types) > 0 {
		paramType := types[0]
		if check = (paramType == utils.UpdateEmail); !check {
			h.SendBadRequest(w, utils.ErrInvalidUpdateEmailType)
			return
		}
		param["type"] = paramType
	} else {
		h.SendBadRequest(w, utils.ErrInvalidTypeQueryParameter)
		return
	}

	dataUser, role, permissions, expAtUnix, jwtToken, err := m.CheckTokenAndExpirationUpdateEmail(h.DB, ctx, param["type"].(string), utils.User, param["token"].(string), param["usercode"].(string))
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	//convert unix timestamp
	expAt = time.Unix(expAtUnix, 0)

	// Populate response
	res = response.LoginUserRes{
		Token:       jwtToken,
		ImageURL:    dataUser.ImageURL.String,
		UserCode:    dataUser.UserCode,
		FullName:    dataUser.FullName,
		PhoneNumber: dataUser.PhoneNumber,
		Email:       dataUser.Email.String,
		ExpiredAt:   expAt.Format(utils.DATE_TIME_FORMAT),
		ActorType:   utils.User,
		CreatedDate: time.Now().In(time.UTC).Format(utils.DATE_TIME_FORMAT),
		RoleId:      dataUser.RoleId,
		RoleCode:    role.RoleCode,
		Permissions: ToResponsePermissionList(permissions),
	}

	h.SendSuccess(w, res, nil)
}

func (h *Contract) ResetPasswordUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		userIdentifier = bootstrap.GetIdentifierCodeFromToken(ctx, r)
		m              = model.Contract{App: h.App}
		req            = request.ResetPasswordReq{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	err = m.ResetPassword(h.DB, ctx, userIdentifier, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

func (h *Contract) VerifyPasswordAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		userIdentifier = bootstrap.GetIdentifierCodeFromToken(ctx, r)
		m              = model.Contract{App: h.App}
		req            = request.VerifyPasswordReq{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	err = m.VerifyUserPassword(h.DB, ctx, userIdentifier, req.Password)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	res := response.VerifyPasswordRes{
		Message: "Verify Password Success",
		Status:  true,
	}

	// Populate response
	h.SendSuccess(w, res, nil)
}
