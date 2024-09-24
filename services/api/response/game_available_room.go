package response

import "encoding/json"

type GameAvailableRoomRes struct {
	RoomID         int    `json:"room_id"`
	RoomCode       string `json:"room_code"`
	RoomImageUrl   string `json:"room_image_url"`
	CafeName       string `json:"cafe_name"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	GameMasterID   int    `json:"game_master_id"`
	GameMasterName string `json:"game_master_name"`
}

func BuildGameAvailableRoomResResp(data string) []GameAvailableRoomRes {
	var resp []GameAvailableRoomRes
	_ = json.Unmarshal([]byte(data), &resp)
	return resp
}
