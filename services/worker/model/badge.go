package model

import (
	"context"
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"encoding/json"
)

func (h *Contract) CheckBadges(ctx context.Context, badgeCode string) error {
	var (
		m = model.Contract{App: h.App}
	)
	userIdList, err := m.GetListUsersByUserId(h.DB, ctx)
	if err != nil {
		return h.errHandler("model.CheckBadge", err, utils.ErrGettingBadgeRule)
	}

	badgeRuleList, err := m.GetBadgeRuleByBadgeCode(h.DB, ctx, badgeCode)
	if err != nil {
		return h.errHandler("model.CheckBadge", err, utils.ErrGettingBadgeRule)
	}

	for _, userId := range userIdList {

		var badgeRules []bool

		for _, badgeRule := range badgeRuleList {
			// check type is spesific board game category
			if badgeRule.KeyCondition == utils.SpesificBoardGameCategory {
				var specificBoardGameCategory SpesificBoardGameCategory
				valueJSON, err := json.Marshal(badgeRule.Value)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrUnmarshallingBadgeRule)
				}

				err = json.Unmarshal(valueJSON, &specificBoardGameCategory)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrUnmarshallingBadgeRule)
				}

				isGameMaster := specificBoardGameCategory.NeedGM

				var qualified []bool
				for _, gameCode := range specificBoardGameCategory.GameCode {

					gameId, err := m.GetGameIdByCode(h.DB, ctx, gameCode)
					if err != nil {
						return h.errHandler("model.CheckBadge", err, utils.ErrGettingGameByCode)
					}

					roomGameCount, err := m.CountRoomParticipantByUserIdAndGameIdAndIsGameMasterAndBookingPrice(h.DB, ctx, userId, gameId, 0, isGameMaster)
					if err != nil {
						return h.errHandler("model.CheckBadge", err, utils.ErrCountingRoomParticipants)
					}

					tournamentGameCount, err := m.CountTournamentParticipantByUserIdAndGameIdAndIsGameMasterAndBookingPrice(h.DB, ctx, userId, gameId, 0)
					if err != nil {
						return h.errHandler("model.CheckBadge", err, utils.ErrCountingTournamentParticipants)
					}

					hasCollection, err := m.CheckUserGameCollectionsByUserID(h.DB, ctx, userId, gameId)
					if err != nil {
						return h.errHandler("model.CheckBadge", err, utils.ErrCheckingUserGameCollection)
					}

					qualified = append(qualified, roomGameCount > 0 || tournamentGameCount > 0 || hasCollection)
				}

				badgeRules = append(badgeRules, utils.ContainsFalse(qualified))
				// check if condition is time limit
			} else if badgeRule.KeyCondition == utils.TimeLimit {
				var timeLimitCategory TimeLimitCategory
				valueJSON, err := json.Marshal(badgeRule.Value)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrUnmarshallingBadgeRule)
				}

				err = json.Unmarshal(valueJSON, &timeLimitCategory)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrUnmarshallingBadgeRule)
				}

				// check if type time limit
				if timeLimitCategory.Category == utils.TimeLimit {
					roomCount, err := m.CountRoomParticipantByUserIdAndStartDateAndEndDate(h.DB, ctx, userId, timeLimitCategory.StartDate, timeLimitCategory.EndDate)
					if err != nil {
						return h.errHandler("model.CheckBadge", err, utils.ErrCountingRoomParticipants)
					}

					tournamentCount, err := m.CountTournamentParticipantByUserIdAndStartDateAndEndDate(h.DB, ctx, userId, timeLimitCategory.StartDate, timeLimitCategory.EndDate)
					if err != nil {
						return h.errHandler("model.CheckBadge", err, utils.ErrCountingTournamentParticipants)
					}

					totalPlayedGames, err := m.CountUserGameCollectionsByUserID(h.DB, ctx, userId, timeLimitCategory.StartDate, timeLimitCategory.EndDate)
					if err != nil {
						return h.errHandler("model.CheckBadge", err, utils.ErrCountingRoomParticipants)
					}

					totalPlayed := roomCount + tournamentCount + int64(totalPlayedGames)

					badgeRules = append(badgeRules, totalPlayed >= 1)

					// check if type life time
				} else if timeLimitCategory.Category == utils.LifeTime {
					roomCount, err := m.CountRoomParticipantByUserIdAndStartDateAndLifeTime(h.DB, ctx, userId, timeLimitCategory.StartDate, timeLimitCategory.EndDate)
					if err != nil {
						return h.errHandler("model.CheckBadge", err, utils.ErrCountingRoomParticipants)
					}
					tournamentCount, err := m.CountTournamentParticipantByUserIdAndStartDateAndLifeTime(h.DB, ctx, userId, timeLimitCategory.StartDate, timeLimitCategory.EndDate)
					if err != nil {
						return h.errHandler("model.CheckBadge", err, utils.ErrCountingTournamentParticipants)
					}

					totalPlayedGames, err := m.CountUserGameCollectionsByUserID(h.DB, ctx, userId, timeLimitCategory.StartDate, timeLimitCategory.EndDate)
					if err != nil {
						return h.errHandler("model.CheckBadge", err, utils.ErrCountingRoomParticipants)
					}

					totalPlayed := roomCount + tournamentCount + int64(totalPlayedGames)

					badgeRules = append(badgeRules, totalPlayed >= 1)
				}

				// check if condition is total spend
			} else if badgeRule.KeyCondition == utils.TotalSpend {
				var requiredSpendAmount int
				valueJSON, err := json.Marshal(badgeRule.Value)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrUnmarshallingBadgeRule)
				}

				err = json.Unmarshal(valueJSON, &requiredSpendAmount)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrUnmarshallingBadgeRule)
				}

				totalClaimedInvoiceAmount, err := m.GetTotalInvoiceAmountByUserID(h.DB, ctx, userId)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrGettingTotalInvoiceAmount)
				}

				totalBookingAmount, err := m.GetTotalBookingAmountByUserID(h.DB, ctx, userId)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrGettingTotalBookingAmount)
				}

				totalSpentAmount := totalBookingAmount + totalClaimedInvoiceAmount
				badgeRules = append(badgeRules, requiredSpendAmount <= totalSpentAmount)
			} else if badgeRule.KeyCondition == utils.TournamentWon {
				var requiredTournamentWon int
				valueJSON, err := json.Marshal(badgeRule.Value)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrUnmarshallingBadgeRule)
				}

				err = json.Unmarshal(valueJSON, &requiredTournamentWon)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrUnmarshallingBadgeRule)
				}

				totalTournamenWon, err := m.CountTournamentWinnerByUserId(h.DB, ctx, userId)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrGettingTotalInvoiceAmount)
				}

				badgeRules = append(badgeRules, requiredTournamentWon <= totalTournamenWon)

			} else if badgeRule.KeyCondition == utils.PlayingGames {
				var requiredDifferentGames int
				valueJSON, err := json.Marshal(badgeRule.Value)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrUnmarshallingBadgeRule)
				}

				err = json.Unmarshal(valueJSON, &requiredDifferentGames)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrUnmarshallingBadgeRule)
				}
				totalGames, err := m.CountUserGameCollectionsByUserID(h.DB, ctx, userId, "", "")
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrGettingTotalInvoiceAmount)
				}

				badgeRules = append(badgeRules, requiredDifferentGames <= totalGames)
			}
		}

		if utils.ContainsFalse(badgeRules) {
			badgeId, err := m.GetBadgeIdByCode(h.DB, ctx, badgeCode)
			if err != nil {
				return h.errHandler("model.CheckBadge", err, utils.ErrGettingBadgeID)
			}

			if err := m.AddUserBadge(h.DB, ctx, userId, badgeId); err != nil {
				return h.errHandler("model.CheckBadge", err, utils.ErrAddingUserBadge)
			}
		}
	}

	return nil
}
