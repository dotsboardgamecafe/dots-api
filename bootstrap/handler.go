package bootstrap

import (
	"context"
	"dots-api/lib/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	validator "github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
)

const (
	// Custom headers
	XSignature     = "X-SIGNATURE" // Custom header to hold signature string
	XTimestamp     = "X-TIMESTAMP" // Custom header to hold timestamp for the signature
	XPlayer        = "X-PLAYER"    // Token obtained from OneSignal Push notification
	XChannelHeader = "X-CHANNEL"   // Custom header to determine the channel

	// Success and error messages
	MsgSuccess          = "APP:SUCCESS"              // Success message
	MsgErrValidation    = "ERR:VALIDATION"           // Error due to validation
	MsgEmptyData        = "ERR:EMPTY_DATA"           // Data not found error
	MsgErrParam         = "ERR:INVALID_PARAM"        // Error due to invalid parameter in the query string
	MsgBadReq           = "ERR:BAD_REQUEST"          // General bad request error
	MsgUnprocessable    = "ERR:UNPROCESSABLE_ENTITY" // 422 Error
	MsgNotFound         = "ERR:NOT_FOUND"            // 404 Not Found error
	MsgAuthErr          = "ERR:AUTHENTICATION"       // Authentication error
	MsgAuthorizedErr    = "ERR:AUTHORIZED"           // Authorization error
	MsgForbiddenErr     = "ERR:FORBIDDEN"            // Forbidden error
	MsgEmailNotFoundErr = "ERR:EMAIL_NOT_FOUND"      // Error indicating email not found

	Success          = "Success"
	Forbidden        = "Forbidden"
	AuthHeader       = "Authorization" // Authorization header
	UnauthorizedData = "Unauthorized data"
	ValidationError  = "Validation error"
	AuthBase64Error  = "[base64:Invalid]" // Error flag for invalid base64
)

// ErrorBase64 give error string of invalid base64
func (h *App) ErrorBase64() error {
	return fmt.Errorf(AuthBase64Error)
}

// Bind bind the API request payload (body) into request struct.
func (h *App) Bind(r *http.Request, input interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&input)

	return err
}

func (h *App) BindAndValidate(r *http.Request, input interface{}) error {
	err := h.Bind(r, &input)
	if err != nil {
		return err
	}

	return h.Validator.Driver.Struct(input)
}

func GetIdentifierCodeFromToken(ctx context.Context, r *http.Request) string {
	claims := &CustomClaims{}
	tokenAuth := r.Header.Get("Authorization")
	_, _ = jwt.ParseWithClaims(tokenAuth, claims, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	return claims.Code
}

func GetIdentifierChannelFromToken(ctx context.Context, r *http.Request) string {
	claims := &CustomClaims{}
	tokenAuth := r.Header.Get("Authorization")
	_, _ = jwt.ParseWithClaims(tokenAuth, claims, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	return claims.Channel
}

func GetIdentifierEmailFromToken(ctx context.Context, r *http.Request) string {
	claims := &CustomClaims{}
	tokenAuth := r.Header.Get("Authorization")
	_, _ = jwt.ParseWithClaims(tokenAuth, claims, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	return claims.Email
}

func GetIdentifierRoleFromToken(ctx context.Context, r *http.Request) string {
	claims := &CustomClaims{}
	tokenAuth := r.Header.Get("Authorization")
	_, _ = jwt.ParseWithClaims(tokenAuth, claims, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	return claims.Issuer
}

func (h *App) GetChannel(r *http.Request) string {
	return r.Header.Get(XChannelHeader)
}

func (h *App) GetXPlayer(r *http.Request) string {
	return r.Header.Get(XPlayer)
}

func (h *App) GetToken(r *http.Request) string {
	return r.Header.Get(AuthHeader)
}

func (h *App) EmptyJSONArr() []map[string]interface{} {
	return []map[string]interface{}{}
}

// SendSuccess send success into response with 200 http code.
func (h *App) SendSuccess(w http.ResponseWriter, payload interface{}, pagination interface{}) {
	if pagination == nil {
		pagination = h.EmptyJSONArr()
	}
	h.RespondWithJSON(w, 200, MsgSuccess, Success, payload, pagination)
}

func (h *App) SendEmptyDataSuccess(w http.ResponseWriter, payload interface{}, pagination interface{}) {
	if pagination == nil {
		pagination = h.EmptyJSONArr()
	}

	h.RespondWithCustomJSON(w, 200, MsgSuccess, utils.EmptyData, payload, pagination)
}

func (h *App) SendSuccessCustomMsg(w http.ResponseWriter, payload interface{}, pagination interface{}, message string) {
	if pagination == nil {
		pagination = h.EmptyJSONArr()
	}
	if len(message) <= 0 {
		message = Success
	}
	h.RespondWithJSON(w, 200, MsgSuccess, message, payload, pagination)
}

// SendBadRequest send bad request into response with 400 http code.
func (h *App) SendBadRequest(w http.ResponseWriter, message string) {
	h.RespondWithJSON(w, 400, MsgBadReq, message, h.EmptyJSONArr(), h.EmptyJSONArr())
}

// SendUnprocessableEntity send 422 error response
func (h *App) SendUnprocessableEntity(w http.ResponseWriter, message string) {
	h.RespondWithJSON(w, 422, MsgUnprocessable, message, nil, nil)
}

// SendBadWithNilDataRequest send bad request into response with 400 http code.
func (h *App) SendBadWithNilDataRequest(w http.ResponseWriter, message string) {
	h.RespondWithJSON(w, 400, MsgBadReq, message, nil, h.EmptyJSONArr())
}

// SendNotfound send bad request into response with 400 http code.
func (h *App) SendNotfound(w http.ResponseWriter, message string) {
	h.RespondWithJSON(w, 404, MsgNotFound, message, h.EmptyJSONArr(), h.EmptyJSONArr())
}

// SendAuthError send bad request into response with 400 http code.
func (h *App) SendAuthError(w http.ResponseWriter, message string) {
	h.RespondWithJSON(w, 401, MsgAuthErr, message, h.EmptyJSONArr(), h.EmptyJSONArr())
}

// SendUnAuthorizedData send bad request into response with 401 http code.
func (h *App) SendUnAuthorizedData(w http.ResponseWriter) {
	h.RespondWithJSON(w, 401, MsgAuthorizedErr, UnauthorizedData, h.EmptyJSONArr(), h.EmptyJSONArr())
}

// // SendForbidden send forbidden into response with 403 http code.
func (h *App) SendForbidden(w http.ResponseWriter, msg string) {
	if msg == "" {
		msg = Forbidden
	}
	h.RespondWithJSON(w, 403, MsgForbiddenErr, msg, h.EmptyJSONArr(), h.EmptyJSONArr())
}

// SendAuthError send bad request into response with 400 http code.
func (h *App) SendInternalServerErr(w http.ResponseWriter, message string) {
	h.RespondWithJSON(w, 502, MsgAuthErr, message, h.EmptyJSONArr(), h.EmptyJSONArr())
}

// SendRequestValidationError Send validation error response to consumers.
func (h *App) SendRequestValidationError(w http.ResponseWriter, validationErrors validator.ValidationErrors) {
	errorResponse := map[string][]string{}
	errorTranslation := validationErrors.Translate(h.Validator.Translator)
	for _, err := range validationErrors {
		errKey := utils.Underscore(err.StructField())
		errorResponse[errKey] = append(
			errorResponse[errKey],
			strings.Replace(errorTranslation[err.Namespace()], err.StructField(), "[]", -1),
		)
	}

	h.RespondWithJSON(w, 400, MsgErrValidation, ValidationError, errorResponse, h.EmptyJSONArr())
}

// SendBindAndValidateError handles errors related to binding and validation.
func (h *App) SendBindAndValidateError(w http.ResponseWriter, err interface{}) {
	switch v := err.(type) {
	case validator.ValidationErrors:
		h.SendRequestValidationError(w, v)
	case error:
		h.SendBadRequest(w, v.Error())
	default:
		log.Fatal(utils.ErrInvalidTypeError)
	}
}

// RespondWithJSON write json response format
func (h *App) RespondWithJSON(
	w http.ResponseWriter,
	httpCode int,
	statCode string,
	message string,
	payload interface{},
	pagination interface{},
) {
	respPayload := map[string]interface{}{
		"stat_code":  statCode,
		"stat_msg":   message,
		"pagination": pagination,
		"data":       payload,
	}

	response, _ := json.Marshal(respPayload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	_, _ = w.Write(response)
}

// RespondWithJSON write json response format
func (h *App) RespondWithCustomJSON(
	w http.ResponseWriter,
	httpCode int,
	statCode string,
	message string,
	payload interface{},
	pagination interface{},
) {
	respPayload := map[string]interface{}{
		"stat_code":  statCode,
		"stat_msg":   message,
		"pagination": pagination,
		"data":       payload,
	}

	nJson := jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 "rw",
	}.Froze()

	response, _ := nJson.Marshal(&respPayload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	_, _ = w.Write(response)
}

// GetUserCode ...
func (h *App) GetUserCode(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value("identifier").(map[string]string)["mcode"])
}

// GetUserRole ...
func (h *App) GetUserRole(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value("identifier").(map[string]string)["role"])
}

// ParamOrder ...
type ParamOrder struct {
	Field string
	By    string
}

/*
GetIntParam Parse the url param to get value as integer.
for example, we need to get limit and offset param
*/
func (h *App) GetIntParam(r *http.Request, name string) (int, error) {
	param := r.URL.Query().Get(name)
	if len(param) == 0 {
		return 0, nil
	}

	return strconv.Atoi(param)
}

// GetStringParam Parse the url param to get value as string.
func (h *App) GetStringParam(r *http.Request, name string) (string, error) {
	param := r.URL.Query().Get(name)
	if len(param) == 0 {
		return "", nil
	}

	return param, nil
}

// GetBoolParam Parse the url param to get value as boolean
func (h *App) GetBoolParam(r *http.Request, name string) (bool, error) {
	param := r.URL.Query().Get(name)
	if len(param) == 0 {
		return false, nil
	}

	result, err := strconv.ParseBool(param)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (h *App) PingAction(w http.ResponseWriter, r *http.Request) {
	h.SendSuccess(w, h.EmptyJSONArr(), nil)
}

func (h *App) VersionAction(w http.ResponseWriter, r *http.Request) {
	type Version struct {
		Build string `json:"build"`
	}
	res := Version{Build: "1.0.2"}
	h.SendSuccess(w, res, nil)
}
