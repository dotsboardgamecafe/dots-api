package model

import (
	"context"
	"dots-api/lib/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserPointEnt struct {
	Id               int64     `db:"id"`
	UserId           int64     `db:"user_id"`
	UserName         string    `db:"user_name"`
	DataSource       string    `db:"data_source"`
	SourceCode       string    `db:"source_code"`
	TitleDescription string    `db:"title_description"`
	GameName         string    `db:"game_name"`
	GameCode         string    `db:"game_code"`
	GameImgUrl       string    `db:"game_img_url"`
	Point            int       `db:"point"`
	CreatedDate      time.Time `db:"created_date"`
}

func (c *Contract) AddUserPoint(tx pgx.Tx, ctx context.Context, userId int64, dataSource string, sourceCode string, point int) error {
	sql := `INSERT INTO users_points(
		user_id,
		data_source,
		source_code,
		point,
		created_date
	)
	VALUES($1, $2, $3, $4, $5);`

	_, err := tx.Exec(ctx, sql, userId, dataSource, sourceCode, point, time.Now().In(time.UTC))
	if err != nil {
		return c.errHandler("model.AddUserPoint", err, utils.ErrAddUserPoint)
	}

	// Get latest point
	userCode, currentUserPoint, currentUserTierId, _ := c.GetLatestPointAndTier(tx, ctx, userId)

	// Calculate total point and define latest tier
	finalTotalPoint := point + currentUserPoint
	finalTier, _ := c.GetTierByPoinCriteria(tx, ctx, finalTotalPoint)

	if currentUserTierId != finalTier.Id {
		sqlUpdatePointAndTier := `UPDATE users
			SET latest_point = $1, latest_tier_id = $2
			WHERE id = $3;`

		_, err = tx.Exec(ctx, sqlUpdatePointAndTier, finalTotalPoint, finalTier.Id, userId)
		if err != nil {
			return err
		}

		// Generate Notification code
		notifCode := utils.GeneratePrefixCode(utils.NotifPrefix)

		description := "Anda sekarang berada di tingkat/tier baru! Selamat datang di" + finalTier.Name + "yang lebih tinggi dengan akses lebih banyak fitur dan manfaat eksklusif."

		descriptionJSON, err := json.Marshal(description)
		if err != nil {
			return err
		}
		// Insert data into db
		err = c.AddNotificationWithTx(tx, ctx, notifCode, "user", userCode, userCode, utils.LevelUpType, utils.LevelUpTitle, descriptionJSON, "")
		if err != nil {
			return fmt.Errorf("error adding notification: %v", err)
		}

		// Storing Tier Benefit Notification
		listBenefits, _ := c.GetBenefitsByTierId(tx, ctx, finalTier.Id)
		if len(listBenefits) > 0 {
			for _, benefit := range listBenefits {
				notifBenefitCode := utils.GeneratePrefixCode(utils.NotifPrefix)

				description = "Selamat! Anda telah beruntung dan mendapatkan " + benefit.Name + " dari kami. Kami berterima kasih atas partisipasi Anda!"

				descriptionJSON, err := json.Marshal(description)
				if err != nil {
					return err
				}

				// Insert data into db
				err = c.AddNotificationWithTx(tx, ctx, notifBenefitCode, "user", userCode, benefit.RewardCode, utils.RewardsType, utils.RewardsTitle, descriptionJSON, benefit.ImageUrl)
				if err != nil {
					return fmt.Errorf("error adding notification: %v", err)
				}
			}
		}
	} else {
		sqlUpdatePointAndTier := `UPDATE users SET latest_point = $1 WHERE id = $2;`

		_, err = tx.Exec(ctx, sqlUpdatePointAndTier, finalTotalPoint, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Contract) GetCurrentUserTotalPoint(db *pgxpool.Pool, ctx context.Context, userId int) (int, error) {
	var (
		TotalPoint int

		sql = `SELECT SUM(point) FROM users_points WHERE user_id = $1;`
	)

	err := db.QueryRow(ctx, sql, userId).Scan(&TotalPoint)
	if err != nil {
		return TotalPoint, c.errHandler("model.GetCurrentUserTotalPoint", err, utils.ErrGetCurrentUserTotalPoint)
	}

	return TotalPoint, err
}
