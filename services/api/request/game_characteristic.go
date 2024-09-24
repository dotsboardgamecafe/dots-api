package request

type (
	GameCharacteristicReq struct {
		CharacteristicLeft  string `json:"characteristic_left"`
		CharacteristicRight string `json:"characteristic_right"`
		Value               int    `json:"value"`
	}
)
