package response

import "encoding/json"

type GameCategoryRes struct {
	CategoryName string `json:"category_name"`
}

func BuildGameCategoryResp(data string) []GameCategoryRes {
	var resp []GameCategoryRes
	_ = json.Unmarshal([]byte(data), &resp)
	return resp
}
