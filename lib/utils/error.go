package utils

var (
	// General error
	EmptyData                          = "Data not found"
	ErrNotFoundPage                    = "Sorry. We couldn't find that page"
	ErrSystemError                     = "Something error with our system. Please contact our administrator"
	ErrInvalidTokenChannel             = "Invalid token channel"
	ErrConfigKeyNotFound               = "Config ['app.key'] doesn't exists"
	ErrEmailAlreadyRegistered          = "Your email has been registered. Please change to a new email"
	ErrPhoneAlreadyRegistered          = "Your phone number has been registered. Please change to a new phone number"
	ErrUsernameAlreadyRegistered       = "Username already taken. Please use another one."
	ErrEmailNotVerified                = "Email not verified"
	ErrInvalidEmailPassword            = "Email/password is incorrect"
	ErrGeneratingJWT                   = "Error generating JWT token"
	ErrGettingVerificationsData        = "Error getting data from verifications table"
	ErrBeginningTransaction            = "Error beginning transaction"
	ErrInvalidToken                    = "Token is invalid"
	ErrTokenExpired                    = "Token is expired"
	ErrMarshalData                     = "Error Mashall data"
	ErrTokenUsed                       = "Token has been used"
	ErrMarkingToken                    = "Error marking token as used"
	ErrCommittingTransaction           = "Error committing transaction"
	ErrSendingResetPasswordEmail       = "Error sending email for reset password"
	ErrAddingResetPasswordVerification = "Error adding verification for reset password"
	ErrSendingVerifyEmail              = "Error sending email for verify email"
	ErrSendingForgotPasswordEmail      = "Error sending email for forgot password"
	ErrSendingUpdateEmail              = "Error sending email for update email"
	ErrInvalidSendingEmailType         = "Type must be one of the following: (verify_registration | forgot_password | update_email)"
	ErrInvalidUpdateEmailType          = "Type must be update_email"
	ErrInvalidTypeQueryParameter       = "Type query parameter is missing"
	ErrPasswordMismatch                = "Password does not match"
	ErrHashingPassword                 = "Error hashing the new password"
	ErrInvalidTypeError                = "Incorrect error type provided for 'err' parameter. It must be an instance of 'validator.ValidationErrors' or 'error'"
	ErrInvalidStartDateFormat          = "Invalid start date format"
	ErrInvalidEndDateFormat            = "Invalid end date format"
	ErrInvalidStartTimeFormat          = "Invalid start time format"
	ErrInvalidEndTimeFormat            = "Invalid end time format"
	ErrXPlayerTokenCantBeEmpty         = "X-Player can't be empty"
	ErrForbiddenAuth                   = "Forbidden, you're not authorized"

	ErrInvalidExpiredDate = "Invalid expired date format"
	// Error for module AMQP
	ErrConnectAMQP           = "Can't connect to AMQP"
	ErrCreateChannelAMQP     = "Can't create a amqpChannel"
	ErrContentTypeNotAllowed = "Content type is not allowed"

	// Error for module user
	ErrGettingUserData                = "Error getting user data"
	ErrInsertingUser                  = "Error inserting user"
	ErrUpdatingUserPassword           = "Error updating user's password"
	ErrFetchingUserPassword           = "Error fetching user's current password"
	ErrUpdatingUserEmail              = "Error updating user email"
	ErrUpdatingUserEmailStatus        = "Error updating user email status"
	ErrUpdatingUserProfile            = "Error updating user profile"
	ErrUpdatingUserXPlayer            = "Error updating user x player"
	ErrRetrievingUserByUserIdentifier = "Error retrieving user by user identifier"
	ErrGettingUserByEmail             = "Error getting user by email"
	ErrGettingListUser                = "Error getting list user"
	ErrGettingListUserId              = "Error getting list user id"
	ErrCountingListUser               = "Error counting list user"
	ErrScanningListUser               = "Error scanning list user"
	ErrScanningListUserId             = "Error scanning list user id"
	ErrGettingUserLatestPoint         = "Error getting user latest point"
	ErrUpdatingUserStatus             = "Error updating user status"
	ErrDeletingAdmin                  = "Error deleting admin by code"

	// Error for module user address
	ErrGettingUserAddresses      = "Error getting user addresses by user ID"
	ErrScanningUserAddresses     = "Error scanning user addresses"
	ErrIteratingUserAddresses    = "Error iterating over user addresses"
	ErrInsertingUserAddress      = "Error inserting user address"
	ErrCheckingAddressIdentifier = "Error checking address identifier"
	ErrUpdatingUserAddress       = "Error updating user address"
	ErrDeletingUserAddress       = "Error deleting user address"

	// Error for module setting
	ErrCountingListSetting       = "Error counting list setting"
	ErrGettingListSetting        = "Error getting list setting"
	ErrScanningListSetting       = "Error scanning list setting"
	ErrAddingSetting             = "Error adding setting"
	ErrGettingSettingByCode      = "Error getting setting by code"
	ErrUpdatingSetting           = "Error updating setting by code"
	ErrGettingSettingByKey       = "Error getting setting by key"
	ErrGettingSettingListByGroup = "Error getting settings list by group"

	// Error for module cafe
	ErrCountingListCafe      = "Error counting list cafe"
	ErrGettingListCafe       = "Error getting list cafe"
	ErrScanningListCafe      = "Error scanning list cafe"
	ErrAddingCafe            = "Error adding cafe"
	ErrGettingCafeByCode     = "Error getting cafe by code"
	ErrGettingCafeCityByCode = "Error getting cafe city by code"
	ErrUpdatingCafe          = "Error updating cafe by code"

	// Error for module notification
	ErrCountingListNotification  = "Error counting list notification"
	ErrGettingListNotification   = "Error getting list notification"
	ErrScanningListNotification  = "Error scanning list notification"
	ErrAddingNotification        = "Error adding notification"
	ErrGettingNotificationByCode = "Error getting notification by code"
	ErrUpdatingNotification      = "Error updating notification by code"

	// Error for module game
	ErrCountingListGame    = "Error counting list game"
	ErrGettingListGame     = "Error getting list game"
	ErrScanningListGame    = "Error scanning list game"
	ErrAddingGame          = "Error adding game"
	ErrGettingGameByCode   = "Error getting game by code"
	ErrUpdatingGame        = "Error updating game by code"
	ErrCheckExistGameUsed  = "Error on checking game is used by other room or tournament"
	ErrForbiddenDeleteGame = "Game cannot be deleted, because is used in another room or tournament"

	// Error for module game category
	ErrAddingGameCategory   = "Error adding game category"
	ErrDeletingGameCategory = "Error deleting game category"

	// Error for module game characteristic
	ErrAddingGameCharacteristic   = "Error adding game characteristic"
	ErrDeletingGameCharacteristic = "Error deleting game characteristic"

	// Error for module banner
	ErrCountingListBanner  = "Error counting list banner"
	ErrGettingListBanner   = "Error getting list banner"
	ErrScanningListBanner  = "Error scanning list banner"
	ErrAddingBanner        = "Error adding banner"
	ErrGettingBannerByCode = "Error getting banner by code"
	ErrUpdatingBanner      = "Error updating banner by code"

	// Error for admin setting
	ErrCountingListAdmin   = "Error counting list admin"
	ErrGettingListAdmin    = "Error getting list admin"
	ErrScanningListAdmin   = "Error scanning list admin"
	ErrAddingAdmin         = "Error adding admin"
	ErrUpdatingAdmin       = "Error updating admin by code"
	ErrGettingAdminByCode  = "Error getting admin by code"
	ErrGettingAdminByEmail = "Error getting admin by email"
	ErrGettingAdminByPhone = "Error getting admin by phone number"
	ErrUpdatingAdminStatus = "Error updating admin status by code"

	// Error for tournament
	ErrCountingListTournament                   = "Error counting list tournament"
	ErrGettingListTournament                    = "Error getting list tournament"
	ErrGettingListTournamentCode                = "Error getting list tournament code"
	ErrGetRemainingOfNonWinnerPlayers           = "Error getting list of non winner players"
	ErrScanningListTournament                   = "Error scanning list tournament"
	ErrAddingTournament                         = "Error adding tournament"
	ErrGettingTournamentByCode                  = "Error getting tournament by code"
	ErrUpdatingTournament                       = "Error updating tournament by code"
	ErrDeletingTournament                       = "Error deleting tournament"
	ErrUpdatingTournamentStatus                 = "Error updating status tournament by code"
	ErrFetchingTournamentByStartDate            = "Error fetching tournament by start date"
	ErrCountParticipantTournamentByTournamentId = "Error counting total participant by tournament id"
	ErrGettingTournamentByCodeAndUserCode       = "Error getting tournament by code and usercode"

	// Error for module tournament participant
	ErrAddingTournamentParticipant                     = "Error adding tournament participant"
	ErrUpdatingTournamentParticipant                   = "Error deleting tournament participant"
	ErrDeletingTournamentParticipant                   = "Error deleting tournament participant"
	ErrGettingAllParticipantByTournamentCode           = "Error getting all participant by tournament code"
	ErrFetchingTournamentParticipant                   = "Error fething tournament participant"
	ErrCountParticipantTournamentByStartDateAndEndDate = "Error counting total participant by start date and end date"
	ErrCountParticipantWonTournament                   = "Error counting total participant won tournament"

	// Error for module room
	ErrCountingListRoom                                     = "Error counting list room"
	ErrGettingListRoom                                      = "Error getting list room"
	ErrGettingListRoomCode                                  = "Error getting list room code"
	ErrScanningListRoom                                     = "Error scanning list room"
	ErrAddingRoom                                           = "Error adding room"
	ErrGettingRoomByCode                                    = "Error getting room by code"
	ErrGettingRoomByCodeAndUserCode                         = "Error getting room by code and usercode"
	ErrUpdatingRoom                                         = "Error updating room by code"
	ErrDeletingRoom                                         = "Error deleting room by code"
	ErrUpdatingRoomStatus                                   = "Error updating status room by code"
	ErrCountParticipantRoomByRoomId                         = "Error counting total participant by room id"
	ErrCountParticipantRoomByUserIdAndGameIdAndIsGameMaster = "Error counting total participant by user id and game id and have game master"
	ErrCountParticipantRoomByStartDateAndEndDate            = "Error counting total participant by start date and end date"

	// Error for module room participant
	ErrAddingRoomParticipant            = "Error adding room participant"
	ErrUpdatingRoomParticipant          = "Error updating room participant"
	ErrUpdatingRoomParticipantStatus    = "Error updating status room participant"
	ErrDeletingRoomParticipant          = "Error deleting room participant"
	ErrGettingAllParticipantByRoomCode  = "Error getting all participant by room code"
	ErrScanningAllParticipantByRoomCode = "Error scanning all participant by room code"

	// Error for module tier
	ErrGettingListTier    = "Error getting list tier"
	ErrGettingTierByCode  = "Error getting tier by code"
	ErrScanningListTier   = "Error scanning list tier"
	ErrGetTierWithReward  = "Error getting list tier with reward"
	ErrScanTierWithReward = "Error scanning list tier with reward"
	ErrAddingTier         = "Error adding tier"
	ErrUpdatingTier       = "Error updating tier"
	ErrDeletingTier       = "Error deleting tier"

	// Error for module user point
	ErrGetPlayerAndOtherActivities  = "Error getting other user activites"
	ErrScanPlayerAndOtherActivities = "Error scanning other user activites"
	ErrGetUsersPointActivity        = "Error getting user point activites"
	ErrScanUsersPointActivity       = "Error scanning user point activites"
	ErrAddUserPoint                 = "Error adding user point"
	ErrGetCurrentUserTotalPoint     = "Error counting user point"

	// Error for module reward
	ErrGettingListReward   = "Error getting list reward"
	ErrGettingRewardByCode = "Error getting reward by code"
	ErrCountingListReward  = "Error counting list reward"
	ErrScanningListReward  = "Error scanning list reward"
	ErrAddingReward        = "Error adding reward"
	ErrUpdatingReward      = "Error updating reward"
	ErrDeletingReward      = "Error deleting reward"
	ErrRewardNotFound      = "Error reward not found"

	// Error for module badge
	ErrGettingListBadge               = "Error getting list badge"
	ErrGettingListBadgeByParentCode   = "Error getting list badge by parent code"
	ErrGettingBadgeByCode             = "Error getting badge by code"
	ErrGettingBadgeByCodeByParentCode = "Error getting badge by code by parent code"
	ErrScanningListBadge              = "Error scanning list badge"
	ErrAddingBadge                    = "Error adding badge"
	ErrUpdatingBadge                  = "Error updating badge"
	ErrDeletingBadge                  = "Error deleting badge"
	ErrBadgeNotFound                  = "Badge not found"
	ErrCountingListBadge              = "Error counting list badge"

	ErrGettingBadgeRuleList   = "Error getting list of badge rules"
	ErrScanningBadgeRule      = "Error scanning badge rule"
	ErrGettingBadgeRuleDetail = "Error getting badge rule detail"
	ErrAddingBadgeRule        = "Error adding badge rule"
	ErrUpdatingBadgeRule      = "Error updating badge rule"
	ErrDeletingBadgeRule      = "Error deleting badge rule"

	ErrGettingTournamentBadge   = "Error getting tournament badge"
	ErrScanningTournamentBadge  = "Error scanning tournament badge"
	ErrInsertingTournamentBadge = "Error inserting tournament badge"
	ErrDeletingTournamentBadge  = "Error deleting tournament badge"

	// Error for User Badge
	ErrGettingListUserBadge         = "Error getting list user badge"
	ErrScanningListUserBadge        = "Error scanning list user badge"
	ErrCountingListUserBadge        = "Error counting list user badge"
	ErrGettingtUserBadgeByBadgeCode = "Error getting  user badge by badge code"
	ErrorUpdatingUserBadge          = "Error updating user badge"
	ErrorAddingUserBadge            = "Error adding user badge"
	ErrorDeletingUserBadge          = "Error deleting user badge"
	ErrorCheckingUserBadge          = "Error checking user badge"

	//Error for User Favourite Game
	ErrGettingListUserFavouriteGame  = "Error getting list user favourite game"
	ErrScanningListUserFavouriteGame = "Error scanning list user favourite game"
	ErrCountingListUserFavouriteGame = "Error counting list user favourite game"

	//Error for User Game Collection
	ErrGettingListUserGameCollection  = "Error getting list user game collection"
	ErrScanningListUserGameCollection = "Error scanning list user game collection"
	ErrCountingListUserGameCollection = "Error counting list user game collection"

	//Error for User Game History
	ErrGettingListUserGameHistory  = "Error getting list user game history"
	ErrScanningListUserGameHistory = "Error scanning list user game history"
	ErrCountingListUserGameHistory = "Error counting list user game history"
	ErrCountingDifferentGames      = "Error counting different game"
	// Error for User Redeem History
	ErrGettingListUserRedeemHistory        = "Error getting list user redeem history"
	ErrScanningListUserRedeemHistory       = "Error scanning list user redeem history"
	ErrCountingListUserRedeemHistory       = "Error counting list user redeem history"
	ErrGettingtUserRedeemHistoryDetailCode = "Error getting user redeem history detail"
	ErrIsInvoiceCodeExist                  = "Error getting is invoice code exist"
	ErrFailedUpdateRedeemInfo              = "Error updating invoice"

	// Error for User Transaction Module
	ErrGettingListUserTransaction        = "Error getting list user transaction"
	ErrCountingTotalUserTransaction      = "Error counting total user transaction"
	ErrScanningListUserTransaction       = "Error scanning list user transaction"
	ErrCountingListUserTransaction       = "Error counting list user transaction"
	ErrGettingtUserTransactionDetailCode = "Error getting user transaction detail"

	// Error for Top Player (most VP, most unique games & hall of fame)
	ErrGetHallOfFame             = "Error getting hall of fame"
	ErrScanHallOfFame            = "Error scanning hall of fame"
	ErrGetMonthlyTopAchiever     = "Error getting monthly top achiever"
	ErrScanGetMonthlyTopAchiever = "Error scanning monthly top achiever"

	// Error Payment (UsersTransactions)
	ErrCreatingOneInvoice            = "Error creating invoice from Xendit"
	ErrGetInvoiceByAggregatorCode    = "Error fetching selected users transaction"
	ErrUpdateInvoiceByAggregatorCode = "Error update selected users transaction"

	// Error for module RBAC
	// Error Permission
	ErrCheckPermission         = "Error checking permission rbac"
	ErrCountingListPermission  = "Error counting list permission"
	ErrGettingListPermission   = "Error getting list permission"
	ErrScanningListPermission  = "Error scanning permission"
	ErrGettingPermissionByCode = "Error getting permission by code"
	ErrAddingPermission        = "Error adding permission"
	ErrUpdatingPermission      = "Error updating permission"
	ErrDeletingPermission      = "Error deleting permission"

	// Error Role
	ErrCountingListRole  = "Error counting list role"
	ErrGettingListRole   = "Error getting list role"
	ErrScanningListRole  = "Error scanning role"
	ErrGettingRoleByCode = "Error getting role by code"
	ErrGettingRoleById   = "Error getting role by role_id"
	ErrAddingRole        = "Error adding role"
	ErrUpdatingRole      = "Error updating role"
	ErrDeletingRole      = "Error deleting role"

	// Error RolePermission
	ErrGettingListRolePermission      = "Error getting list role permisson"
	ErrScanningListRolePermission     = "Error scanning role permission"
	ErrAddingRolePermission           = "Error adding role permisson"
	ErrDeletingRolePermissionByRoleID = "Error deleting role permission by role id"

	// Error Season
	ErrAddingSeason       = "Error adding season"
	ErrUpdatingSeason     = "Error updating season"
	ErrDeletingSeason     = "Error deleting season"
	ErrGettingAllSeasons  = "Error getting all seasons"
	ErrScanningAllSeasons = "Error scanning all seasons"
	ErrFetchingSeason     = "Error fetching season"

	// Check Badge
	ErrGettingBadgeList               = "error getting badge list"
	ErrGettingBadgeRule               = "error getting badge rule"
	ErrUnmarshallingBadgeRule         = "error unmarshalling badge rule"
	ErrCountingRoomParticipants       = "error counting room participants"
	ErrCountingTournamentParticipants = "error counting tournament participants"
	ErrGettingBadgeID                 = "error getting badge ID by code"
	ErrAddingUserBadge                = "error adding user badge"
	ErrGettingTotalInvoiceAmount      = "error getting total invoice amount"
	ErrGettingTotalBookingAmount      = "error getting total booking amount"
)
