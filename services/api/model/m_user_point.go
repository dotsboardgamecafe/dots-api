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
	UserCode         string    `db:"user_code"`
	UserImageUrl     string    `db:"user_image_url"`
	UserName         string    `db:"user_name"`
	UserStyle        UserStyle `db:"user_style"`
	DataSource       string    `db:"data_source"`
	SourceCode       string    `db:"source_code"`
	TitleDescription string    `db:"title_description"`
	ItemName         string    `db:"item_name"`
	ItemCode         string    `db:"item_code"`
	ItemImgUrl       string    `db:"item_img_url"`
	Point            int       `db:"point"`
	Description      string    `db:"description"`
	CreatedDate      time.Time `db:"created_date"`
}

type UserPointHistory struct {
	SourceId       int64     `db:"source_id"`
	SourceUserCode string    `db:"source_user_code"`
	SourceType     string    `db:"source_type"`
	SourceCode     string    `db:"source_code"`
	SourceName     string    `db:"source_name"`
	Point          int       `db:"point"`
	Description    string    `db:"description"`
	CreatedDate    time.Time `db:"created_date"`
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

	// Get latest point and tier
	userCode, currentUserPoint, currentUserTierId, _ := c.GetLatestPointAndTier(tx, ctx, userId)

	// Calculate total point and define latest tier
	finalTotalPoint := point + currentUserPoint
	if finalTotalPoint < 0 {
		finalTotalPoint = 0
	}
	finalTier, _ := c.GetTierByPoinCriteria(tx, ctx, finalTotalPoint)

	return c.handleTierTransition(tx, ctx, userId, userCode, currentUserTierId, finalTier, finalTotalPoint)
}

func (c *Contract) AddUserPointWithDescription(tx pgx.Tx, ctx context.Context, userId int64, dataSource string, sourceCode string, point int, description *string) error {
	sql := `INSERT INTO users_points(
		user_id,
		data_source,
		source_code,
		point,
		description,
		created_date
	)
	VALUES($1, $2, $3, $4, $5, $6);`

	_, err := tx.Exec(ctx, sql, userId, dataSource, sourceCode, point, description, time.Now().In(time.UTC))
	if err != nil {
		return c.errHandler("model.AddUserPoint", err, utils.ErrAddUserPoint)
	}

	// Get latest point and tier
	userCode, currentUserPoint, currentUserTierId, _ := c.GetLatestPointAndTier(tx, ctx, userId)

	// Calculate total point and define latest tier
	finalTotalPoint := point + currentUserPoint
	if finalTotalPoint < 0 {
		finalTotalPoint = 0
	}
	finalTier, _ := c.GetTierByPoinCriteria(tx, ctx, finalTotalPoint)

	return c.handleTierTransition(tx, ctx, userId, userCode, currentUserTierId, finalTier, finalTotalPoint)
}

func (c *Contract) handleTierTransition(tx pgx.Tx, ctx context.Context, userId int64, userCode string, currentUserTierId int64, finalTier TierEnt, finalTotalPoint int) error {
	if currentUserTierId == finalTier.Id {
		sqlUpdatePoint := `UPDATE users SET latest_point = $1 WHERE id = $2;`
		_, err := tx.Exec(ctx, sqlUpdatePoint, finalTotalPoint, userId)
		return err
	}

	// Update user's latest_point and latest_tier_id
	sqlUpdatePointAndTier := `UPDATE users SET latest_point = $1, latest_tier_id = $2 WHERE id = $3;`
	_, err := tx.Exec(ctx, sqlUpdatePointAndTier, finalTotalPoint, finalTier.Id, userId)
	if err != nil {
		return err
	}

	if finalTier.Id > currentUserTierId {
		// LEVEL UP
		// 1. Generate Notification code
		notifCode := utils.GeneratePrefixCode(utils.NotifPrefix)
		description := "Anda sekarang berada di tingkat/tier baru! Selamat datang di " + finalTier.Name + "yang lebih tinggi dengan akses lebih banyak fitur dan manfaat eksklusif."
		descriptionJSON, err := json.Marshal(description)
		if err != nil {
			return err
		}
		// Insert data into db
		err = c.AddNotificationWithTx(tx, ctx, notifCode, "user", userCode, userCode, utils.LevelUpType, utils.LevelUpTitle, descriptionJSON, "")
		if err != nil {
			return fmt.Errorf("error adding notification: %v", err)
		}

		// 2. Storing Tier Benefit Notification
		listBenefits, _ := c.GetBenefitsByTierId(tx, ctx, finalTier.Id)
		if len(listBenefits) > 0 {
			for _, benefit := range listBenefits {
				notifBenefitCode := utils.GeneratePrefixCode(utils.NotifPrefix)
				descBenefit := "Selamat! Anda telah beruntung dan mendapatkan " + benefit.Name + " dari kami. Kami berterima kasih atas partisipasi Anda!"
				descBenefitJSON, err := json.Marshal(descBenefit)
				if err != nil {
					return err
				}
				err = c.AddNotificationWithTx(tx, ctx, notifBenefitCode, "user", userCode, benefit.RewardCode, utils.RewardsType, utils.RewardsTitle, descBenefitJSON, benefit.ImageUrl)
				if err != nil {
					return fmt.Errorf("error adding notification: %v", err)
				}
			}
		}

		// 3. Add tier entry to users_points
		sqlInsertTier := `INSERT INTO users_points(user_id, data_source, source_code, point, created_date) VALUES($1, $2, $3, $4, $5);`
		_, err = tx.Exec(ctx, sqlInsertTier, userId, utils.UserPointType["TIER"], finalTier.TierCode, 0, time.Now().In(time.UTC))
		if err != nil {
			return c.errHandler("model.handleTierTransition.insertTier", err, utils.ErrAddUserPoint)
		}
	} else {
		// DEGRADED
		var currentTierName string
		queryCurrentTier := `SELECT COALESCE(name, '') FROM tiers WHERE id = $1;`
		_ = tx.QueryRow(ctx, queryCurrentTier, currentUserTierId).Scan(&currentTierName)
		if currentTierName == "" {
			currentTierName = "Previous Tier"
		}
		degradedDesc := "Degraded From " + currentTierName

		// 1. Remove current and higher tier entries from users_points
		sqlDeleteHigherTiers := `DELETE FROM users_points 
		WHERE user_id = $1 AND data_source = $2 AND source_code IN (
			SELECT tier_code FROM tiers WHERE id > $3
		);`
		_, err = tx.Exec(ctx, sqlDeleteHigherTiers, userId, utils.UserPointType["TIER"], finalTier.Id)
		if err != nil {
			return c.errHandler("model.handleTierTransition.deleteTier", err, utils.ErrRemovingUserPoint)
		}

		// 2. Update the previous/lower tier entry's created_date and give description
		sqlUpdatePrevTier := `UPDATE users_points 
			SET created_date = $1, description = $2 
			WHERE id = (
				SELECT id FROM users_points 
				WHERE user_id = $3 AND data_source = $4 AND source_code = $5 
				ORDER BY created_date DESC, id DESC 
				LIMIT 1
			);`
		cmdTag, err := tx.Exec(ctx, sqlUpdatePrevTier, time.Now().In(time.UTC), degradedDesc, userId, utils.UserPointType["TIER"], finalTier.TierCode)
		if err != nil {
			return c.errHandler("model.handleTierTransition.updatePrevTier", err, utils.ErrUpdatingTier)
		}

		// If no entry exists for the lower tier in users_points, insert entry with description
		if cmdTag.RowsAffected() == 0 {
			sqlInsertTier := `INSERT INTO users_points(user_id, data_source, source_code, point, description, created_date) VALUES($1, $2, $3, $4, $5, $6);`
			_, err = tx.Exec(ctx, sqlInsertTier, userId, utils.UserPointType["TIER"], finalTier.TierCode, 0, degradedDesc, time.Now().In(time.UTC))
			if err != nil {
				return c.errHandler("model.handleTierTransition.insertTier", err, utils.ErrAddUserPoint)
			}
		}
	}

	return nil
}

func (c *Contract) RemoveUserPoint(tx pgx.Tx, ctx context.Context, userId int64, dataSource string, sourceCode string, point int) error {
	var (
		err   error
		query = `DELETE FROM users_points WHERE user_id = $1 AND data_source = $2 AND source_code = $3`
	)

	cmdTag, err := tx.Exec(ctx, query, userId, dataSource, sourceCode)
	if err != nil {
		return c.errHandler("model.RemoveUserPoint", err, utils.ErrRemovingUserPoint)
	}

	// if point is greater than 0 and a row was deleted, it should
	// calculate user latest_point and define latest tier
	if point > 0 && cmdTag.RowsAffected() > 0 {
		// Get latest point and tier
		userCode, currentUserPoint, currentUserTierId, _ := c.GetLatestPointAndTier(tx, ctx, userId)

		// Calculate total point and define latest tier
		finalTotalPoint := currentUserPoint - point
		if finalTotalPoint < 0 {
			finalTotalPoint = 0
		}
		finalTier, _ := c.GetTierByPoinCriteria(tx, ctx, finalTotalPoint)

		return c.handleTierTransition(tx, ctx, userId, userCode, currentUserTierId, finalTier, finalTotalPoint)
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
