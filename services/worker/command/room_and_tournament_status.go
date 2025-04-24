package command

import (
	"context"
	"dots-api/lib/utils"
	"dots-api/services/worker/model"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// UpdateStatusRoomAndTournament ...
func (app Contract) UpdateStatusRoomAndTournament(c *cli.Context) error {
	var (
		err error
		// Begin Context
		ctx = context.Background()
		m   = model.Contract{App: app.App}
		now = time.Now().UTC()
	)

	now, err = utils.FromUTCLocationToGMT7(now)
	if err != nil {
		app.Log.File().Log(logrus.ErrorLevel, "set-inactive-room-and-tournament: ", err)
		return err
	}

	rooms, err := m.GetExpiredRoomLists(m.DB, ctx)
	if err != nil {
		return err
	}

	for _, room := range rooms {
		endDate := time.Date(
			room.EndDate.Time.Year(),
			room.EndDate.Time.Month(),
			room.EndDate.Time.Day(),
			room.EndTime.Hour(),
			room.EndTime.Minute(),
			room.EndTime.Second(),
			0,
			now.Location(),
		)

		if room.Status == "inactive" || now.Before(endDate) {
			continue
		}

		err = m.UpdateRoomStatus(app.DB, ctx, room.RoomCode, "inactive")
		if err != nil {
			return err
		}

		notificationType := utils.UpcomingSession

		if room.RoomType == utils.RoomType[1] {
			notificationType = utils.UpcomingEvent
		}

		removeUpcomingRoomNotification(&m, ctx, room.RoomCode, notificationType)
	}

	tournaments, err := m.GetExpiredTournamentLists(m.DB, ctx)
	if err != nil {
		return err
	}

	for _, tournament := range tournaments {
		endDate := time.Date(
			tournament.EndDate.Time.Year(),
			tournament.EndDate.Time.Month(),
			tournament.EndDate.Time.Day(),
			tournament.EndTime.Hour(),
			tournament.EndTime.Minute(),
			tournament.EndTime.Second(),
			0,
			now.Location(),
		)

		if tournament.Status == "inactive" || now.Before(endDate) {
			continue
		}

		err = m.UpdateTournamentStatus(app.DB, ctx, tournament.TournamentCode, "inactive")
		if err != nil {
			return err
		}

		removeUpcomingRoomNotification(&m, ctx, tournament.TournamentCode, utils.TournamentReminder)
	}

	fmt.Printf("Set inactive tournament and room success at %v", now.Format("Monday 2006-01-02 15:04:05"))
	return nil
}

func removeUpcomingRoomNotification(m *model.Contract, ctx context.Context, roomCode, notificationType string) {
	err := m.DeleteNotification(m.DB, ctx, roomCode, notificationType)
	if err != nil {
		fmt.Println("h.removeUpcomingRoomNotification: ", err.Error())
	}
}
