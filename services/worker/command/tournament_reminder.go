package command

import (
	"context"
	"dots-api/lib/onesignal"
	"dots-api/lib/utils"
	"dots-api/services/api/response"
	"dots-api/services/worker/model"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/urfave/cli/v2"
)

// UserRoomReminder ...
func (app Contract) UserTournamentReminder(c *cli.Context) error {
	var dataListTournament []model.TournamentsEnt
	var err error
	now := time.Now().UTC()

	// Calculate h-3 reminder
	h3Reminder := now.AddDate(0, 0, 3).Format(utils.DATE_FORMAT)

	// Calculate h-1 reminder
	h1Reminder := now.AddDate(0, 0, 1).Format(utils.DATE_FORMAT)

	// Begin Context
	ctx := context.Background()
	m := model.Contract{App: app.App}

	dataListTournament, err = m.GetTournamentList(m.DB, ctx, h1Reminder)
	if err != nil {
		return err
	}

	for _, tournament := range dataListTournament {

		dataListUser, err := m.GetAllParticipantByTournamentCode(m.DB, ctx, tournament.TournamentCode)
		if err != nil {
			return err
		}
		// Populate response
		for _, user := range dataListUser {

			// Generate Notification code
			notifCode := utils.GeneratePrefixCode(utils.NotifPrefix)

			description := response.NotificationTournamentResp{
				StartDate:   tournament.StartDate.Time.Format("2006-01-02"),
				StartTime:   tournament.StartTime.Format("15:04:05"),
				EndTime:     tournament.EndTime.Format("15:04:05"),
				CafeName:    tournament.CafeName,
				GameName:    tournament.Name.String,
				CafeAddress: tournament.CafeAddress,
				Level:       tournament.Level,
			}

			descriptionJSON, err := json.Marshal(description)
			if err != nil {
				return err
			}

			// Insert data into db
			err = m.AddNotification(m.DB, ctx, notifCode, "user", user.UserCode, tournament.TournamentCode, utils.TournamentBookingType, tournament.Name.String, descriptionJSON, tournament.ImageUrl.String)
			if err != nil {
				return err
			}

			onesignal := onesignal.New(m.App)
			OSDescription := utils.TournamentReminderPushNotificationDescription + "\n\n" + "Tournament name: " + tournament.Name.String + "\n" + "Date: " + tournament.StartDate.Time.Format("2006-01-02") + "\n" + "Location: " + tournament.CafeName
			_, err = onesignal.CreateOSNotifications(user.UserXPlayer, utils.TournamentReminderPushNotificationTitle, OSDescription, utils.Tournament)
			if err != nil {
				log.Printf("Error : %s", err)
			}

		}
	}

	dataListTournament, err = m.GetTournamentList(m.DB, ctx, h3Reminder)
	if err != nil {
		return err
	}

	for _, tournament := range dataListTournament {

		dataListUser, err := m.GetAllParticipantByTournamentCode(m.DB, ctx, tournament.TournamentCode)
		if err != nil {
			return err
		}
		// Populate response
		for _, user := range dataListUser {

			// Generate Notification code
			notifCode := utils.GeneratePrefixCode(utils.NotifPrefix)

			description := response.NotificationTournamentResp{
				StartDate:   tournament.StartDate.Time.Format("2006-01-02"),
				StartTime:   tournament.StartTime.Format("15:04:05"),
				EndTime:     tournament.EndTime.Format("15:04:05"),
				CafeName:    tournament.CafeName,
				GameName:    tournament.Name.String,
				CafeAddress: tournament.CafeAddress,
				Level:       tournament.Level,
			}

			descriptionJSON, err := json.Marshal(description)
			if err != nil {
				return err
			}

			// Insert data into db
			err = m.AddNotification(m.DB, ctx, notifCode, "user", user.UserCode, tournament.TournamentCode, utils.TournamentBookingType, tournament.Name.String, descriptionJSON, tournament.ImageUrl.String)
			if err != nil {
				return err
			}

			onesignal := onesignal.New(m.App)
			OSDescription := utils.TournamentReminderPushNotificationDescription + "\n\n" + "Tournament name: " + tournament.Name.String + "\n" + "Date: " + tournament.StartDate.Time.Format("2006-01-02") + "\n" + "Location: " + tournament.CafeName
			_, err = onesignal.CreateOSNotifications(user.UserXPlayer, utils.TournamentReminderPushNotificationTitle, OSDescription, utils.Tournament)
			if err != nil {
				log.Printf("Error : %s", err)
			}

		}
	}

	fmt.Printf("Send Tournament Reminder Notification success at %v", now.Format("Monday 2006-01-02 15:04:05"))
	return nil
}
