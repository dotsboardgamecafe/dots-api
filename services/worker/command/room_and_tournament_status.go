package command

import (
	"context"
	"dots-api/lib/utils"
	"dots-api/services/worker/model"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

// UpdateStatusRoomAndTournament ...
func (app Contract) UpdateStatusRoomAndTournament(c *cli.Context) error {
	var (
		dataListTournamentCode []string
		err                    error
		// Begin Context
		ctx = context.Background()
		m   = model.Contract{App: app.App}
		now = time.Now().UTC()
	)

	rooms, err := m.GetRoomCodeAndTypeExpiredRoomLists(m.DB, ctx)
	if err != nil {
		return err
	}

	for _, room := range rooms {
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

	dataListTournamentCode, err = m.GetListTournamentCodes(m.DB, ctx)
	if err != nil {
		return err
	}

	for _, tournamentCode := range dataListTournamentCode {
		err = m.UpdateTournamentStatus(app.DB, ctx, tournamentCode, "inactive")
		if err != nil {
			return err
		}

		removeUpcomingRoomNotification(&m, ctx, tournamentCode, utils.TournamentReminder)
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
