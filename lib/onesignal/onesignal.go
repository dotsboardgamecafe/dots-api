package onesignal

import (
	"bytes"
	"dots-api/bootstrap"
	"dots-api/lib/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

// device_type list reference https://documentation.onesignal.com/reference/add-a-device
const (
	URL_HOST                      = "https://onesignal.com/api"
	VERSION                       = "v1"
	API_PREFIX_DEVICE             = "players"
	API_PREFIX_NOTIFICATION       = "notifications"
	DEVICE_TYPE_IOS               = 0
	DEVICE_TYPE_ANDROID           = 1
	DEVICE_TYPE_BROWSER_CHROMEAPP = 4
	DEVICE_TYPE_BROWSER_CHROMEWEB = 5
	DEVICE_TYPE_BROWSER_SAFARI    = 7
	DEVICE_TYPE_BROWSER_FIREFOX   = 8
)

type service struct {
	app *bootstrap.App
}

func New(app *bootstrap.App) *service {
	return &service{app}
}

func getURLByPrefix(prefix string) string {
	url := fmt.Sprintf("%s/%s/%s", URL_HOST, VERSION, prefix)
	return url
}

func (s *service) getApiKey() string {
	return s.app.Config.GetString("onesignal.api_key")
}

func (s *service) getAppID() string {
	return s.app.Config.GetString("onesignal.app_id")
}

// Player represents a OneSignal player.
type Player struct {
	ID                string            `json:"id"`
	Playtime          int               `json:"playtime"`
	SDK               string            `json:"sdk"`
	Identifier        string            `json:"identifier"`
	SessionCount      int               `json:"session_count"`
	Language          string            `json:"language"`
	Timezone          int               `json:"timezone"`
	GameVersion       string            `json:"game_version"`
	DeviceOS          string            `json:"device_os"`
	DeviceType        int               `json:"device_type"`
	DeviceModel       string            `json:"device_model"`
	AdID              string            `json:"ad_id"`
	Tags              map[string]string `json:"tags"`
	LastActive        int               `json:"last_active"`
	AmountSpent       float32           `json:"amount_spent"`
	CreatedAt         int               `json:"created_at"`
	InvalidIdentifier bool              `json:"invalid_identifier"`
	BadgeCount        int               `json:"badge_count"`
}

func (s *service) AddDevice(deviceType int) (string, error) {
	player := Player{}
	url := getURLByPrefix(API_PREFIX_DEVICE)

	values := map[string]interface{}{
		"app_id":      s.getAppID(),
		"device_type": deviceType,
	}

	request, err := utils.RequestHandler(values, url, http.MethodPost)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err := utils.ResponseAsyncHandler(request)
	if err != nil {
		return "", err
	}

	json.Unmarshal([]byte(response[0]), &player)

	return player.ID, err
}

func (s *service) CreateOSNotifications(xPlayer, title, description, types string) (map[string]interface{}, error) {
	var (
		err         error
		method, url string
		res         map[string]interface{}
	)
	if xPlayer == "" {
		return res, nil
	}
	bodyData := map[string]interface{}{
		"app_id":             s.app.Config.GetString("onesignal.app_id"),
		"headings":           map[string]interface{}{"en": title},
		"contents":           map[string]interface{}{"en": description},
		"include_player_ids": []string{xPlayer},
		"data": map[string]interface{}{
			"type": types,
		},
		// "send_after":         "2021-10-25 16:31:58.000 +0700",
	}

	method = "POST"
	url = URL_HOST + "/v1/notifications"
	payload, err := json.Marshal(bodyData)
	if err != nil {
		return res, err
	}

	// Populate Http Request
	requestData, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return res, err
	}
	requestData.Header.Set("Authorization", "Basic "+s.app.Config.GetString("onesignal.app_key"))
	requestData.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}

	// Do Http Request
	resp, err := client.Do(requestData)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	// Read Response Body
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	// Encode Byte and Bind the Data
	err = json.Unmarshal([]byte(responseData), &res)
	if err != nil {
		return res, err
	}

	// If status is ok
	if s.checkHttpStatusOK(resp.StatusCode) {
		return res, nil
	}

	// Log Error
	s.app.Log.FromDefault().WithFields(logrus.Fields{
		"FunctionName": "[CreateOSNotifications]",
	}).Errorf("Error messages: %v", res["errors"].([]interface{})[0])
	return res, err
}

// checkHttpStatusOK ...
func (s *service) checkHttpStatusOK(status int) bool {
	if status == http.StatusOK || status == http.StatusCreated || status == http.StatusAccepted || status == 201 {
		return true
	}
	return false
}

func (s *service) GetPlayerDevice(playerID string) (Player, error) {
	player := Player{}
	url := fmt.Sprintf("%s/%s", getURLByPrefix(API_PREFIX_DEVICE), playerID)

	request, err := utils.RequestHandler(nil, url, http.MethodGet)
	if err != nil {
		return player, err
	}
	request.Header.Set("Authorization", "Basic "+s.getApiKey())
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err := utils.ResponseHandler(request)
	if err != nil {
		return player, err
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return player, err
	}

	json.Unmarshal(jsonResponse, &player)

	return player, err
}
