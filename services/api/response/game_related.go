package response

import "encoding/json"

type GameRelatedRes struct {
	GameId             int     `json:"game_id"`
	GameName           string  `json:"name"`
	GameCode           string  `json:"game_code"`
	GameType           string  `json:"game_type"`
	Location           string  `json:"location"`
	Difficulty         string  `json:"difficulty"`
	Level              float64 `json:"level"`
	ImageUrl           string  `json:"image_url"`
	Duration           int     `json:"duration"`
	MinimalParticipant int     `json:"minimal_participant"`
	MaximumParticipant int     `json:"maximum_participant"`
}

func BuildGameRelatedResp(data string) []GameRelatedRes {
	var resp []GameRelatedRes
	_ = json.Unmarshal([]byte(data), &resp)
	return resp
}
