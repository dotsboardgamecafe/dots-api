package model

import (
	"context"
	"database/sql"
	"dots-api/lib/utils"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TournamentParticipantRespEnt struct {
	Id             int64          `db:"id"`
	TournamentId   int64          `db:"tournament_id"`
	UserCode       string         `db:"user_code"`
	UserName       string         `db:"user_name"`
	UserImgUrl     string         `db:"user_image_url"`
	UserXPlayer    string         `db:"user_x_player"`
	StatusWinner   bool           `db:"status_winner"`
	Status         string         `db:"status"`
	AdditionalInfo sql.NullString `db:"additional_info"`
	Position       int            `db:"position"`
	RewardPoint    sql.NullInt64  `db:"reward_point"`
}

func (c *Contract) GetAllParticipantByTournamentCode(db *pgxpool.Pool, ctx context.Context, code string) ([]TournamentParticipantRespEnt, error) {
	var (
		err   error
		list  []TournamentParticipantRespEnt
		query = `
		SELECT 
			tp.id,
			tp.tournament_id,
			u.user_code,
			COALESCE(u.username, '') AS username,
			u.image_url,
			u.x_player AS user_x_player,
			tp.status_winner,
			tp.status,
			tp.position,
			tp.additional_info,
			tp.reward_point
		FROM tournaments t
			LEFT JOIN tournament_participants tp ON tp.tournament_id = t.id
			LEFT JOIN users u ON tp.user_id = u.id
		WHERE tournament_code = $1 AND tp.status = 'active' AND t.deleted_date IS NULL `
	)

	rows, err := db.Query(ctx, query, code)
	if err != nil {
		return list, c.errHandler("model.GetAllParticipantByTournamentCode", err, utils.ErrGettingAllParticipantByTournamentCode)
	}

	defer rows.Close()
	for rows.Next() {
		var data TournamentParticipantRespEnt
		err = rows.Scan(
			&data.Id, &data.TournamentId, &data.UserCode,
			&data.UserName, &data.UserImgUrl, &data.UserXPlayer,
			&data.StatusWinner, &data.Status, &data.Position,
			&data.AdditionalInfo, &data.RewardPoint,
		)
		if err != nil {
			return list, c.errHandler("model.GetAllParticipantByTournamentCode", err, utils.ErrScanningListTournament)
		}
		list = append(list, data)
	}

	return list, nil
}
