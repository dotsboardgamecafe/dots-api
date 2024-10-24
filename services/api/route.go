package api

import (
	"dots-api/bootstrap"
	"dots-api/services/api/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// RegisterRoutes all routes for the apps
func RegisterRoutes(r *chi.Mux, app *bootstrap.App) {
	r.Route("/v1", func(r chi.Router) {
		r.Get("/ping", app.PingAction)
		r.Get("/version", app.VersionAction)

		AppSubsRoute(r, app)
	})
}

func AppSubsRoute(r chi.Router, app *bootstrap.App) {
	h := handler.Contract{App: app}

	// Auth
	r.Route("/auths", func(r chi.Router) {
		r.Post("/login", nrWrap(h.LoginUserAct, app.NewRelic))
		r.Post("/register", nrWrap(h.RegisterUserAct, app.NewRelic))

		// Request Token for Registration and ForgotPassword
		r.Post("/resend-verification", nrWrap(h.RequestVerifyEmailUserAct, app.NewRelic))
		r.Post("/verify-token", nrWrap(h.VerifyTokenUserAct, app.NewRelic))
		r.Post("/verify-token-email", nrWrap(h.VerifyTokenUpdateEmailUserAct, app.NewRelic))

		// Verifying User Password (required by UI team before changing email or delete user)
		r.With(app.VerifyJwtToken).Post("/verify-password", nrWrap(h.VerifyPasswordAct, app.NewRelic))
	})

	// forgot password
	r.Route("/forgot-password", func(r chi.Router) {
		r.Post("/", nrWrap(h.RequestVerifyEmailUserAct, app.NewRelic))
		r.Post("/validate-token", nrWrap(h.VerifyTokenUserAct, app.NewRelic))
		r.With(app.VerifyJwtToken).Put("/new-password", nrWrap(h.ResetPasswordUserAct, app.NewRelic))
	})

	// User
	r.Route("/users", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetUserListAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateUserAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetUserDetailAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}/status", nrWrap(h.UpdateUserStatusAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeleteUserAct, app.NewRelic))

		// Other User's Activities
		r.With(app.VerifyAccessRoute).Get("/{code}/activity", nrWrap(h.GetAllPlayerActivities, app.NewRelic))

		// User's Point Activities
		r.With(app.VerifyAccessRoute).Get("/{code}/point-activity", nrWrap(h.GetUserPointActivities, app.NewRelic))

		// User's Transactions
		r.With(app.VerifyAccessRoute).Get("/{code}/transactions", nrWrap(h.GetUserTransactions, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}/transactions/{trx_code}", nrWrap(h.GetUserTransactionDetail, app.NewRelic))

		//User Badges
		r.With(app.VerifyAccessRoute).Get("/{code}/badges", nrWrap(h.GetUserBadgeListAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}/badges/{badge-code}", nrWrap(h.GetUserBadgeByBadgeCodeAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}/badges/{badge-code}", nrWrap(h.UpdateUserBadgeByBadgeCodeAct, app.NewRelic))

		//User favourite game
		r.With(app.VerifyAccessRoute).Get("/{code}/favourite-games", nrWrap(h.GetUserFavouriteGameAct, app.NewRelic))

		//User game collection
		r.With(app.VerifyAccessRoute).Get("/{code}/game-collection", nrWrap(h.GetUserGameCollectionAct, app.NewRelic))

		//User history game
		r.With(app.VerifyAccessRoute).Get("/{code}/history-games", nrWrap(h.GetUserGameHistoryAct, app.NewRelic))

		r.With(app.VerifyAccessRoute).Put("/update-password", nrWrap(h.UpdatePasswordUserAct, app.NewRelic))

		// User Profile
		r.Route("/profile", func(r chi.Router) {
			r.Use(app.VerifyJwtToken)
			r.Post("/new-email", nrWrap(h.RequestVerifyUpdateEmailUserAct, app.NewRelic))
			r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetUserProfileAct, app.NewRelic))
			r.With(app.VerifyAccessRoute).Put("/", nrWrap(h.UpdateUserProfileAct, app.NewRelic))
		})

	})

	// Master Cafe
	r.Route("/cafes", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetCafeListAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddCafeAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetCafeDetailAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateCafeAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeleteCafeAct, app.NewRelic))
	})

	// Master Cafe
	r.Route("/banners", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetBannerListAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetBannerDetailAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddBannerAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateBannerAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeleteBannerAct, app.NewRelic))
	})

	// Master Game
	r.Route("/games", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetGameListAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddGameAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetGameDetailAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateGameAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeleteGameAct, app.NewRelic))
	})

	// Master Admin
	r.Route("/admins", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetAdminListAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddAdminAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetAdminDetailAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateAdminAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeleteAdminAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}/status", nrWrap(h.UpdateAdminStatus, app.NewRelic))
	})

	// Master Setting
	r.Route("/settings", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetSettingListAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetSettingDetailAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddSettingAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateSettingAct, app.NewRelic))
	})

	r.Route("/game-mechanics", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)

		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetGameMechanicList, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetDetailGameMechanic, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddGameMechanic, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateGameMechanic, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeleteGameMechanic, app.NewRelic))
	})

	r.Route("/game-types", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)

		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetGameTypeList, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetDetailGameType, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddGameType, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateGameType, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeleteGameType, app.NewRelic))
	})

	// Room
	r.Route("/rooms", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetRoomList, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddRoom, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateRoom, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}/close", nrWrap(h.SetWinnerRoomAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetRoomByCode, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/{code}/book", nrWrap(h.BookingRoom, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}/status", nrWrap(h.UpdateRoomStatus, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeleteRoom, app.NewRelic))
	})

	// Badges
	r.Route("/badges", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetBadgeListAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetBadgeDetailAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddBadgeAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateBadgeAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeleteBadgeAct, app.NewRelic))
	})

	// Badges
	r.Route("/tournament-badges", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetTournamentBadgeDetailAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddTournamentBadgeAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateTournamentBadgeAct, app.NewRelic))
	})

	// Tournament
	r.Route("/tournaments", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetTournamentList, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetTournamentDetailAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddTournamentAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateTournamentAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}/close", nrWrap(h.SetWinnerTournamentAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}/status", nrWrap(h.UpdateTournamentStatus, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeleteTournamentAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/{code}/book", nrWrap(h.BookingTournamentAct, app.NewRelic))
	})

	// Upload
	r.Route("/uploads", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.UploadFileAct, app.NewRelic))
	})

	// Tiers
	r.Route("/tiers", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetTiersListAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetTiersDetailAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddTierAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateTierAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeleteTierAct, app.NewRelic))
	})

	// Rewards
	r.Route("/rewards", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetRewardsListAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetRewardDetailAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddRewardAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateRewardAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeleteRewardAct, app.NewRelic))
	})

	// Hall of Fame, Most VP & Most Unique Games
	r.Route("/players", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/hall-of-fame", nrWrap(h.GetHallOfFame, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/monthly-top-achiever", nrWrap(h.GetMonthlyTopAchiever, app.NewRelic))
	})

	// User Invoice Redemption (Histories & Redeem)
	r.Route("/redeems", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetRedeemHistory, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.Redeem, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{invoice_code}", nrWrap(h.GetRedeemDetail, app.NewRelic))
	})

	// CMS Management for Invoice Redemption (Histories & Redeem)
	r.Route("/invoices", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/{user_code}/history", nrWrap(h.GetInvoicesClaimedHistory, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/{user_code}/claim", nrWrap(h.Claim, app.NewRelic))
	})

	// Transaction Callback
	r.Route("/transaction", func(r chi.Router) {
		r.With(app.VerifyXenditCallbackToken).Post("/callback", nrWrap(h.TransactionCallback, app.NewRelic))
	})

	r.Route("/payment", func(r chi.Router) {
		r.Get("/success-callback", nrWrap(h.SuccessCallback, app.NewRelic))
		r.Get("/failure-callback", nrWrap(h.FailureCallback, app.NewRelic))
	})

	// Master Permission
	r.Route("/permissions", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetPermissionListAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddPermissionAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetPermissionDetailAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdatePermissionAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeletePermissionAct, app.NewRelic))
	})

	// Master Role
	r.Route("/roles", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.With(app.VerifyAccessRoute).Get("/", nrWrap(h.GetRoleListAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Post("/", nrWrap(h.AddRoleAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Get("/{code}", nrWrap(h.GetRoleDetailAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Put("/{code}", nrWrap(h.UpdateRoleAct, app.NewRelic))
		r.With(app.VerifyAccessRoute).Delete("/{code}", nrWrap(h.DeleteRoleAct, app.NewRelic))
	})

	// User Notification
	r.Route("/notifications", func(r chi.Router) {
		r.Use(app.VerifyJwtToken)
		r.Get("/", nrWrap(h.GetNotificationList, app.NewRelic))
		r.Put("/{code}", nrWrap(h.UpdateNotificationIsSeenAct, app.NewRelic))
	})
}

// nrWrap wraps an HTTP handler with New Relic instrumentation
func nrWrap(handler http.HandlerFunc, nrapp *newrelic.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txn := nrapp.StartTransaction(r.URL.Path)
		defer txn.End()

		handler(w, r)
	}
}
