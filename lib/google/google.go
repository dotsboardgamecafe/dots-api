package google

import (
	"bytes"
	"dots-api/bootstrap"
	"dots-api/lib/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

const (
	GOOGLE_VERIFY_URL = "https://oauth2.googleapis.com/tokeninfo"
)

type service struct {
	app *bootstrap.App
}

type GoogleTokenResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func New(app *bootstrap.App) *service {
	return &service{app}
}

func (s *service) Verify(token string) (GoogleTokenResponse, error) {
	var (
		err    error
		client = &http.Client{}
		res    = GoogleTokenResponse{}
		param  = url.Values{}
	)
	param.Set("id_token", token)

	var payload = bytes.NewBufferString(param.Encode())
	req, err := http.NewRequest(http.MethodPost, GOOGLE_VERIFY_URL, payload)
	if err != nil {
		return res, errors.New(utils.ErrInvalidToken)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	if err != nil {
		return res, errors.New(utils.ErrInvalidToken)
	}

	defer response.Body.Close()

	respData, err := io.ReadAll(response.Body)
	json.Unmarshal(respData, &res)

	if len(res.Email) < 1 {
		return res, errors.New(utils.ErrInvalidToken)
	}

	return res, err
}
