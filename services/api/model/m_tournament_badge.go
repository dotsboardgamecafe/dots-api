package model

import (
	"context"
	"dots-api/lib/utils"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TournamentBadge represents the association between tournaments and badges.
type TournamentBadge struct {
	ID           int64 `json:"id"`
	TournamentID int64 `json:"tournament_id"`
	BadgeID      int64 `json:"badge_id"`
}

func (c *Contract) GetTournamentBadgeByTournamentCode(db *pgxpool.Pool, ctx context.Context, tournamentCode string) ([]BadgeEnt, error) {
	var (
		err  error
		list []BadgeEnt
	)

	query := `
		SELECT
			b.id,
			b.badge_code,
			b.badge_category,
			b.name,
			b.image_url,
			b.status,
			b.vp_point,
			b.created_date,
			b.updated_date,
			b.deleted_date
		FROM
			tournament_badges tb
		JOIN
			badges b ON tb.badge_id = b.id
		JOIN
			tournaments t ON t.id = tb.tournament_id
		WHERE
			t.tournament_code = $1
	`

	rows, err := db.Query(ctx, query, tournamentCode)
	if err != nil {
		return list, c.errHandler("model.GetTournamentBadgeByTournamentID", err, utils.ErrGettingTournamentBadge)
	}
	defer rows.Close()

	for rows.Next() {
		var data BadgeEnt
		err = rows.Scan(
			&data.Id,
			&data.BadgeCode,
			&data.BadgeCategory,
			&data.Name,
			&data.ImageURL,
			&data.Status,
			&data.VPPoint,
			&data.CreatedDate,
			&data.UpdatedDate,
			&data.DeletedDate,
		)
		if err != nil {
			if err != pgx.ErrNoRows {
				return list, nil
			}
			return list, c.errHandler("model.GetTournamentBadgeByTournamentID", err, utils.ErrScanningTournamentBadge)
		}
		list = append(list, data)
	}

	return list, nil
}

// InsertTournamentBadge inserts a new tournament-badge association into the database.
func (c *Contract) InsertTournamentBadge(tx pgx.Tx, ctx context.Context, tournamentID, badgeID int64) (int64, error) {
	query := "INSERT INTO tournament_badges (tournament_id, badge_id) VALUES ($1, $2) RETURNING id"
	var id int64
	err := tx.QueryRow(ctx, query, tournamentID, badgeID).Scan(&id)
	if err != nil {
		return 0, c.errHandler("model.InsertTournamentBadge", err, utils.ErrInsertingTournamentBadge)
	}
	return id, nil
}

// DeleteTournamentBadge deletes a tournament-badge association from the database.
func (c *Contract) DeleteTournamentBadge(tx pgx.Tx, ctx context.Context, tournamentId int64) error {
	query := "DELETE FROM tournament_badges WHERE tournament_id = $1"
	_, err := tx.Exec(ctx, query, tournamentId)
	if err != nil {
		return c.errHandler("model.DeleteTournamentBadge", err, utils.ErrDeletingTournamentBadge)
	}
	return nil
}
