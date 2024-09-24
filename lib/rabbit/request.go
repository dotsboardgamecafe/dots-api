package rabbit

import (
	"encoding/json"
	"log"
)

type QueueUserBadgeData struct {
	BadgeType string `json:"badge_type"`
	UserId    int64  `json:"user_id"`
}

type QueueBadgeData struct {
	BadgeCode string `json:"badge_code"`
}

// constructor for checek user after badge created
func QueueUserBadgeReq(badgeType string, userId int64) interface{} {
	var (
		err           error
		interfaceData interface{}
		data          = QueueUserBadgeData{
			BadgeType: badgeType,
			UserId:    userId,
		}
	)
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}
	err = json.Unmarshal(dataBytes, &interfaceData)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	return interfaceData
}

// constructor for check user before badge created
func QueueBadgeReq(badgeCode string) interface{} {
	var (
		err           error
		interfaceData interface{}
		data          = QueueBadgeData{
			BadgeCode: badgeCode,
		}
	)
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}
	err = json.Unmarshal(dataBytes, &interfaceData)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	return interfaceData
}
