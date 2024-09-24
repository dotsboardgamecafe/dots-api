package model

import (
	"context"
	"dots-api/lib/utils"

	"github.com/jackc/pgx/v4"
)

type (
	GameCharacteristicEnt struct {
		Id                  int64  `db:"id"`
		GameId              int64  `db:"game_id"`
		CharacteristicLeft  string `db:"characteristic_left"`
		CharacteristicRight string `db:"characteristic_right"`
		Value               string `db:"value"`
	}
)

func (c *Contract) InsertOneGameCharacteristic(tx pgx.Tx, ctx context.Context, gameId int64, charLeft, charRight string, value int) error {
	var (
		err   error
		query = `INSERT INTO games_characteristics(game_id, characteristic_left, characteristic_right, value) VALUES($1,$2,$3,$4)`
	)
	_, err = tx.Exec(ctx, query, gameId, charLeft, charRight, value)
	if err != nil {
		return c.errHandler("model.InsertOneGameCharacteristic", err, utils.ErrAddingGameCharacteristic)
	}
	return nil
}

func (c *Contract) DeleteGameCharacteristic(tx pgx.Tx, ctx context.Context, gameId int64) error {
	var (
		err   error
		query = `DELETE FROM games_characteristics WHERE game_id=$1`
	)
	_, err = tx.Exec(ctx, query, gameId)
	if err != nil {
		return c.errHandler("model.DeleteGameCharacteristic", err, utils.ErrDeletingGameCharacteristic)
	}
	return nil
}
