package command

import (
	"context"
	"dots-api/services/worker/model"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

// UpdateStatusRoomAndTournament ...
func (app Contract) UpdateStatusRoomAndTournament(c *cli.Context) error {
	var (
		dataListRoomCode       []string
		dataListTournamentCode []string
		err                    error
		// Begin Context
		ctx = context.Background()
		m   = model.Contract{App: app.App}
		now = time.Now().UTC()
	)

	dataListRoomCode, err = m.GetListRoomCodes(m.DB, ctx)
	if err != nil {
		return err
	}

	for _, roomCode := range dataListRoomCode {
		err = m.UpdateRoomStatus(app.DB, ctx, roomCode, "inactive")
		if err != nil {
			return err
		}
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
	}

	fmt.Printf("Set inactive tournament and room success at %v", now.Format("Monday 2006-01-02 15:04:05"))
	return nil
}
