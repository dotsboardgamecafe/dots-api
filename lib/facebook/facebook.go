package facebook

import (
	"dots-api/bootstrap"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	FACEBOOK_VERIFY_URL = "https://graph.facebook.com/v2.5/me"
)

type service struct {
	app *bootstrap.App
}

type FacebookTokenResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func New(app *bootstrap.App) *service {
	return &service{app}
}

func (s *service) Verify(token string) (FacebookTokenResponse, error) {
	var (
		err    error
		client = &http.Client{}
		res    = FacebookTokenResponse{}
	)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?access_token=%s&fields=%s", FACEBOOK_VERIFY_URL, token, "id,name,email"), nil)
	if err != nil {
		return res, errors.New("invalid token")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	if err != nil {
		return res, errors.New("invalid token")
	}
	defer response.Body.Close()

	respData, err := io.ReadAll(response.Body)
	if err != nil {
		return res, errors.New("invalid token")
	}
	json.Unmarshal(respData, &res)

	if len(res.Email) < 1 {
		return res, errors.New("invalid token")
	}

	return res, err
}
