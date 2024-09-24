package model

import (
	"context"
	"dots-api/lib/utils"

	"github.com/jackc/pgx/v4"
)

type (
	GameCategoryEnt struct {
		Id           int64  `db:"id"`
		GameId       int64  `db:"game_id"`
		CategoryName string `db:"category_name"`
	}
)

func (c *Contract) InsertOneGameCategory(tx pgx.Tx, ctx context.Context, gameId int64, catName string) error {
	var (
		err   error
		query = `INSERT INTO games_categories(game_id, category_name) VALUES($1,$2)`
	)
	_, err = tx.Exec(ctx, query, gameId, catName)
	if err != nil {
		return c.errHandler("model.InsertOneGameCategory", err, utils.ErrAddingGameCategory)
	}
	return nil
}

func (c *Contract) DeleteGameCategory(tx pgx.Tx, ctx context.Context, gameId int64) error {
	var (
		err   error
		query = `DELETE FROM games_categories WHERE game_id=$1`
	)
	_, err = tx.Exec(ctx, query, gameId)
	if err != nil {
		return c.errHandler("model.DeleteGameCategory", err, utils.ErrDeletingGameCategory)
	}
	return nil
}
