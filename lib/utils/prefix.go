package utils

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	CafePrefix        = "CAFE-"
	GamePrefix        = "GAME-"
	BannerPrefix      = "BANNER-"
	SettingPrefix     = "SET-"
	UserPrefix        = "USR-"
	AdminPrefix       = "ADM-"
	RoomPrefix        = "ROOM-"
	TournamentPrefix  = "TOUR-"
	TierPrefix        = "TIER-"
	BadgePrefix       = "BDG-"
	ParentBadgePrefix = "PAR-"
	BadgeRulePrefix   = "BDGRULE-"
	RewardPrefix      = "RWRD-"
	RedeemPrefix      = "REEDEM-"
	TransactionPrefix = "TRX-"
	RolePrefix        = "ROLE-"
	PermissionPrefix  = "PRMS-"
	NotifPrefix       = "NOTIF-"
	SeasonPrefix      = "SEA-"
)

// TODO: Make increment generated prefix based on database data
func GeneratePrefixCode(prefix string) string {
	var (
		code string
		now  = time.Now().In(time.UTC)
	)

	rand.New(rand.NewSource(now.UnixNano()))
	code, _ = Generate(`[A-Z]{10}`)
	return fmt.Sprintf("%s%s%s", prefix, now.Format("20060102"), code)
}
