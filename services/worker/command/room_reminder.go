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
func (app Contract) UserRoomReminder(c *cli.Context) error {
	var dataListRoom []model.RoomEnt

	var err error
	now := time.Now().UTC()

	// Calculate h-3 reminder
	h3Reminder := now.AddDate(0, 0, 3).Format(utils.DATE_FORMAT)

	// Calculate h-1 reminder
	h1Reminder := now.AddDate(0, 0, 1).Format(utils.DATE_FORMAT)

	// Begin Context
	ctx := context.Background()
	m := model.Contract{App: app.App}

	dataListRoom, err = m.GetRoomList(m.DB, ctx, h1Reminder)
	if err != nil {
		return err
	}

	for _, room := range dataListRoom {

		dataListUser, err := m.GetAllParticipantByRoomCode(m.DB, ctx, room.RoomCode)
		if err != nil {
			return err
		}

		// Populate response
		for _, user := range dataListUser {

			// Generate Notification code
			notifCode := utils.GeneratePrefixCode(utils.NotifPrefix)

			description := response.NotificationRoomResp{
				StartDate:   room.StartDate.Format("2006-01-02"),
				CafeName:    room.CafeName,
				CafeAddress: room.CafeAddress,
				GameName:    room.CafeName,
				Level:       room.Difficulty,
			}

			descriptionJSON, err := json.Marshal(description)
			if err != nil {
				return err
			}

			// Insert data into db
			err = m.AddNotification(m.DB, ctx, notifCode, "user", user.UserCode, room.RoomCode, utils.RoomBookingType, room.Name, descriptionJSON, room.RoomImgUrl)
			if err != nil {
				return err
			}

			onesignal := onesignal.New(m.App)
			OSDescription := utils.RoomReminderPushNotificationDescription + "\n\n" + "Room name: " + room.Name + "\n" + "Date: " + room.StartDate.Format("2006-01-02") + "\n" + "Location: " + room.CafeName
			_, err = onesignal.CreateOSNotifications(user.UserXPlayer, utils.RoomReminderPushNotificationTitle, OSDescription, utils.Room)
			if err != nil {
				log.Printf("Error : %s", err)
			}

		}
	}

	dataListRoom, err = m.GetRoomList(m.DB, ctx, h3Reminder)
	if err != nil {
		return err
	}

	for _, room := range dataListRoom {

		dataListUser, err := m.GetAllParticipantByRoomCode(m.DB, ctx, room.RoomCode)
		if err != nil {
			return err
		}
		// Populate response
		for _, user := range dataListUser {

			// Generate Notification code
			notifCode := utils.GeneratePrefixCode(utils.NotifPrefix)

			description := response.NotificationRoomResp{
				StartDate:   room.StartDate.Format("2006-01-02"),
				CafeName:    room.CafeName,
				CafeAddress: room.CafeAddress,
				GameName:    room.CafeName,
				Level:       room.Difficulty,
			}

			descriptionJSON, err := json.Marshal(description)
			if err != nil {
				return err
			}

			// Insert data into db
			err = m.AddNotification(m.DB, ctx, notifCode, "user", user.UserCode, room.RoomCode, utils.RoomBookingType, room.Name, descriptionJSON, room.RoomImgUrl)
			if err != nil {
				return err
			}

			onesignal := onesignal.New(m.App)
			OSDescription := utils.RoomReminderPushNotificationDescription + "\n\n" + "Room name: " + room.Name + "\n" + "Date: " + room.StartDate.Format("2006-01-02") + "\n" + "Location: " + room.CafeName
			_, err = onesignal.CreateOSNotifications(user.UserXPlayer, utils.RoomReminderPushNotificationTitle, OSDescription, utils.Room)
			if err != nil {
				return err
			}

		}
	}

	fmt.Printf("Send Room Reminder Notification success at %v", now.Format("Monday 2006-01-02 15:04:05"))
	return nil
}
