package model

import (
	"context"
	"dots-api/lib/utils"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BadgeRuleEnt struct {
	Id            int64       `db:"id"`
	BadgeRuleCode string      `db:"badge_rule_code"`
	BadgeId       int64       `db:"badge_id"`
	KeyCondition  string      `db:"key_condition"`
	ValueType     string      `db:"value_type"`
	Value         interface{} `db:"value"`
}

// GetBadgeRuleList retrieves a list of all badge rules from the database.
func (c *Contract) GetBadgeRuleList(db *pgxpool.Pool, ctx context.Context, badgeId int64) ([]BadgeRuleEnt, error) {
	var (
		err  error
		list []BadgeRuleEnt
	)

	query := `
	SELECT 
		br.id, br.badge_rule_code,  br.badge_id, br.key_condition, br.value_type, br.value 
	FROM badges_rules br 
	LEFT JOIN badges b ON b.id = br.badge_id WHERE b.id = $1`
	rows, err := db.Query(ctx, query, badgeId)
	if err != nil {
		return list, c.errHandler("model.GetBadgeRuleList", err, utils.ErrGettingBadgeRuleList)
	}
	defer rows.Close()

	for rows.Next() {
		var badgeRule BadgeRuleEnt
		err = rows.Scan(
			&badgeRule.Id, &badgeRule.BadgeRuleCode, &badgeRule.BadgeId,
			&badgeRule.KeyCondition, &badgeRule.ValueType, &badgeRule.Value,
		)
		if err != nil {
			return list, c.errHandler("model.GetBadgeRuleDetail", err, utils.ErrScanningBadgeRule)
		}
		list = append(list, badgeRule)
	}

	return list, nil
}

// GetBadgeRuleByBadgeCode retrieves a list of all badge rules from the database.
func (c *Contract) GetBadgeRuleByBadgeCode(db *pgxpool.Pool, ctx context.Context, badgeCode string) ([]BadgeRuleEnt, error) {
	var (
		err  error
		list []BadgeRuleEnt
	)

	query := `SELECT br.id, br.badge_rule_code, br.badge_id, br.key_condition, br.value_type, br.value 
	          FROM badges_rules br 
	          LEFT JOIN badges b ON br.badge_id = b.id 
	          WHERE b.badge_code = $1`
	rows, err := db.Query(ctx, query, badgeCode)
	if err != nil {
		return list, c.errHandler("model.GetBadgeRuleByBadgeCode", err, utils.ErrGettingBadgeRuleList)
	}
	defer rows.Close()

	for rows.Next() {
		var badgeRule BadgeRuleEnt
		err := rows.Scan(
			&badgeRule.Id, &badgeRule.BadgeRuleCode, &badgeRule.BadgeId,
			&badgeRule.KeyCondition, &badgeRule.ValueType, &badgeRule.Value,
		)
		if err != nil {
			return list, c.errHandler("model.GetBadgeRuleByBadgeCode", err, utils.ErrScanningBadgeRule)
		}
		list = append(list, badgeRule)
	}

	if rows.Err() != nil {
		return list, c.errHandler("model.GetBadgeRuleByBadgeCode", rows.Err(), utils.ErrGettingBadgeRuleList)
	}

	return list, nil
}

// AddBadgeRule adds a new badge rule to the database.
func (c *Contract) AddBadgeRule(tx pgx.Tx, ctx context.Context, badgeRuleCode string, badgeId int64, keyCondition, valueType string, value interface{}) error {
	query := `INSERT INTO badges_rules(badge_rule_code, badge_id,  key_condition, value_type, value) VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(ctx, query, badgeRuleCode, badgeId, keyCondition, valueType, value)
	if err != nil {
		return c.errHandler("model.AddBadgeRule", err, utils.ErrAddingBadgeRule)
	}

	return nil
}

// DeleteBadgeRule deletes a badge rule from the database.
func (c *Contract) DeleteBadgeRule(tx pgx.Tx, ctx context.Context, badgeId int64) error {
	query := `DELETE FROM badges_rules WHERE badge_id = $1`

	_, err := tx.Exec(ctx, query, badgeId)
	if err != nil {
		return c.errHandler("model.DeleteBadgeRule", err, utils.ErrDeletingBadgeRule)
	}

	return nil
}
