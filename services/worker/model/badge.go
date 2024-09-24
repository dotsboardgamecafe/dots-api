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

	for _, userId := range userIdList {
		badgeRuleList, err := m.GetBadgeRuleByBadgeCode(h.DB, ctx, badgeCode)
		if err != nil {
			return h.errHandler("model.CheckBadge", err, utils.ErrGettingBadgeRule)
		}

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

				var totalGamesPlayed int64
				for _, gameCode := range specificBoardGameCategory.GameCode {

					gameId, err := m.GetGameIdByCode(h.DB, ctx, gameCode)
					if err != nil {
						return h.errHandler("model.CheckBadge", err, utils.ErrGettingGameByCode)
					}

					roomGameCount, err := m.CountRoomParticipantByUserIdAndGameIdAndIsGameMaster(h.DB, ctx, userId, gameId, isGameMaster)
					if err != nil {
						return h.errHandler("model.CheckBadge", err, utils.ErrCountingRoomParticipants)
					}

					tournamentGameCount, err := m.CountTournamentParticipantByUserIdAndGameIdAndIsGameMaster(h.DB, ctx, userId, gameId)
					if err != nil {
						return h.errHandler("model.CheckBadge", err, utils.ErrCountingTournamentParticipants)
					}

					totalGamesPlayed += roomGameCount + tournamentGameCount
				}

				if totalGamesPlayed >= specificBoardGameCategory.TotalPlayed {
					badgeRules = append(badgeRules, true)
				} else {
					badgeRules = append(badgeRules, false)
				}
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

					if roomCount > 1 && tournamentCount > 1 {
						badgeRules = append(badgeRules, true)
					} else {
						badgeRules = append(badgeRules, false)
					}

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

					if roomCount > 1 && tournamentCount > 1 {
						badgeRules = append(badgeRules, true)
					} else {
						badgeRules = append(badgeRules, false)
					}
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
				if requiredSpendAmount <= totalSpentAmount {
					badgeRules = append(badgeRules, true)
				} else {
					badgeRules = append(badgeRules, false)
				}
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

				if requiredTournamentWon <= totalTournamenWon {
					badgeRules = append(badgeRules, true)
				} else {
					badgeRules = append(badgeRules, false)
				}

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

				totalGames, err := m.CountDifferentGamesByUserID(h.DB, ctx, userId)
				if err != nil {
					return h.errHandler("model.CheckBadge", err, utils.ErrGettingTotalInvoiceAmount)
				}

				if requiredDifferentGames <= totalGames {
					badgeRules = append(badgeRules, true)
				} else {
					badgeRules = append(badgeRules, false)
				}
			}
		}
		if utils.ContainsFalse(badgeRules) {
			badgeId, err := m.GetBadgeIdByCode(h.DB, ctx, badgeCode)
			if err != nil {
				return h.errHandler("model.CheckBadge", err, utils.ErrGettingBadgeID)
			}

			err = m.AddUserBadge(h.DB, ctx, userId, badgeId)
			if err != nil {
				return h.errHandler("model.CheckBadge", err, utils.ErrAddingUserBadge)
			}
		}
	}
	return nil
}
