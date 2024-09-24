package response

import "encoding/json"

type GameCharacteristicRes struct {
	CharacteristicLeft  string `json:"characteristic_left"`
	CharacteristicRight string `json:"characteristic_right"`
	Value               int    `json:"value"`
}

func BuildGameCharacteristicResp(data string) []GameCharacteristicRes {
	var resp []GameCharacteristicRes
	_ = json.Unmarshal([]byte(data), &resp)
	return resp
}
