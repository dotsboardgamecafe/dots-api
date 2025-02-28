package utils

type UserType string

type BadgeCategory string

func (b BadgeCategory) String() string {
	return string(b)
}

func (b BadgeCategory) Valid() bool {
	switch b {
	case "play":
		return true
	case "tournament":
		return true
	case "spent":
		return true
	case "gift":
		return true
	default:
		return false
	}
}

var (
	DATE_TIME_FORMAT = "2006-01-02 15:04:05"
	DATE_FORMAT      = "2006-01-02"
	DATE_DAY_FORMAT  = "01 Jan 2006"
	YEAR_FORMAT      = "2006"
	TIME_FORMAT      = "15:04:05"

	// actor type for register and login
	User  = "user"
	Admin = "admin"

	// role id
	RoleSuperAdminId = 1
	RoleAdminId      = 2
	RoleMemberId     = 3
	RoleCashierId    = 4

	// verification type
	VerifyRegistration = "verify_registration"
	ForgotPassword     = "forgot_password"
	UpdateEmail        = "update_email"

	// reset password route
	ResetPassRoute        = "reset-password?token="
	VerifyTokenRoute      = "auths/verify-token?token="
	VerifyTokenEmailRoute = "auths/verify-token-email?token="
	ForgotPasswordRoute   = "forgot-password?token="
	TypeRoute             = "&type="
	UserCodeRoute         = "&usercode="

	VerificationType           = []string{VerifyRegistration, ForgotPassword}
	StatusBanner               = []string{"publish", "unpublish"}
	StatusAdmin                = []string{"active", "inactive"}
	StatusUser                 = []string{"active", "inactive"}
	StatusCafe                 = []string{"active", "inactive"}
	StatusGame                 = []string{"active", "inactive"}
	StatusReward               = []string{"active", "inactive"}
	StatusTournament           = []string{"active", "inactive"}
	StatusTier                 = []string{"active", "inactive"}
	StatusBadges               = []string{"active", "inactive"}
	StatusRole                 = []string{"active", "inactive"}
	StatusPermission           = []string{"active", "inactive"}
	StatusRoom                 = []string{"active", "inactive"}
	StatusRoomParticipant      = []string{"active", "pending", "cancel"}
	RoomType                   = []string{"normal", "special_event"}
	RoomDifficulty             = []string{"easy", "medium", "hard"}
	RewardType                 = []string{"fnb", "game", "tournament"}
	RewardUsed                 = []string{"1", "0"}
	MonthlyTopAchieverCategory = []string{"vp", "unique_game"}
	HTTPMethodList             = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	XenditTransactionStatus    = []string{"PENDING", "PAID", "SETTLED", "EXPIRED"}

	// Notification title
	UpcomingTournament  = "upcoming_tournament"
	BookingConfirmation = "booking_confirmation"
	CanceledRoom        = "canceled_room"
	LevelUp             = "level_up"
	RoomReminder        = "room_reminder"
	TournamentReminder  = "tournament_reminder"
	Reward              = "reward"

	// Badge Category
	BadgeCategoryPlay       BadgeCategory = "play"
	BadgeCategoryTournament BadgeCategory = "tournament"
	BadgeCategorySpent      BadgeCategory = "spent"
	BadgeCategoryGift       BadgeCategory = "gift"

	// Badge Rule
	TotalSpend                = "total_spend"
	TimeLimit                 = "time_limit"
	LifeTime                  = "life_time"
	SpesificBoardGameCategory = "spesific_board_game_category"
	TournamentWon             = "tournament_won"
	PlayingGames              = "playing_games"
	Quantity                  = "quantity"
	Tournament                = "tournament"
	Room                      = "room"
	Transaction               = "transaction"
	Badge                     = "badge"

	// Notification Title
	FailPaymentType        = "payment_failed"
	FailPaymentTitle       = "Pembayaran Gagal!"
	FailPaymentDescription = "Pembayaran Anda tidak dapat diproses pada saat ini. Pastikan informasi pembayaran Anda benar dan coba kembali dalam beberapa saat."

	ExpiredPaymentType        = "payment_expired"
	ExpiredPaymentTitle       = "Pembayaran Kedaluwarsa!"
	ExpiredPaymentDescription = "Pembayaran Anda telah melewati batas waktu pembayaran."

	SuccessPaymentType                            = "payment_success"
	SuccessPaymentTitle                           = "Pembayaran Berhasil!"
	RoomReminderPushNotificationTitle             = "Roll for Initiative! Epic D&D Special Event Awaits"
	TournamentReminderPushNotificationTitle       = "Heads Up! Don't Miss the Upcoming Tournament"
	RoomReminderPushNotificationDescription       = "Calling all adventurers! Dust off your dice and sharpen your spells for a legendary Dungeons & Dragons. Glory, loot, and bragging rights await!"
	TournamentReminderPushNotificationDescription = "Calling all thrill-seekers and competitive spirits! Get ready for an epic tournament that'll have you on the edge of your seat. Mark your calendars and join the fun!"

	LevelUpType  = "level_up"
	LevelUpTitle = "Selamat! Anda Naik Level!"

	RewardsType  = "got_reward"
	RewardsTitle = "Selamat! Anda Mendapatkan Reward!"

	RoomBookingType       = "room_booking_confirmation"
	TournamentBookingType = "tournament_booking_confirmation"

	// Mapping UserPointType, RedeemPlatform, PaymentStatus, RoomStatus & TournamentStatus
	UserPointType = map[string]string{
		"TOURNAMENT_TYPE":      "tournament",
		"TOURNAMENT_TYPE_PLAY": "tournament_play",
		"TOURNAMENT_TYPE_PAID": "tournament_paid",
		"ROOM_TYPE":            "room",
		"ROOM_TYPE_PLAY":       "room_play",
		"ROOM_TYPE_PAID":       "room_paid",
		"BADGE_TYPE":           "badge",
		"REDEEM_TYPE":          "redeem",
		"GAME_COLLECTION":      "game",
		"TIER":                 "tier",
	}

	// UserPointTypeText
	UserPointTypeText = map[string]string{
		"TOURNAMENT_TYPE":      "[username] has joined [name] and gained [vp]",
		"TOURNAMENT_TYPE_PLAY": "[username] played in the [name] tournament and earned [vp]",
		"ROOM_TYPE":            "[username] has joined [name] and gained [vp]",
		"ROOM_TYPE_PLAY":       "[username] played in the [name] room and earned [vp]",
		"BADGE_TYPE":           "[username] just claimed the [name] badge",
		"REDEEM_TYPE":          "[username] has redeemed [name]",
		"GAME_COLLECTION":      "[username] played [name] for the first time",
		"TIER":                 "[username] has advanced to [name] tier",
	}

	RedeemPlatform = map[string]string{
		"APP": "app",
		"CMS": "cms",
	}

	PaymentStatus = map[string]string{
		"PENDING": "PENDING",
		"PAID":    "PAID",
		"EXPIRED": "EXPIRED",
	}

	RoomStatus = map[string]string{
		"ACTIVE":   "active",
		"INACTIVE": "inactive",
		"CLOSED":   "closed",
	}

	TournamentStatus = map[string]string{
		"ACTIVE":   "active",
		"INACTIVE": "inactive",
		"CLOSED":   "closed",
	}
)

func GetUserPointTypeByValue(value string) (string, bool) {
	for k, v := range UserPointType {
		if v == value {
			return k, true
		}
	}

	return "", false
}
