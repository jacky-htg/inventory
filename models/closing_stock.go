package models

import (
	"context"
	"database/sql"

	"github.com/jacky-htg/inventory/libraries/api"
)

// ClosingStock : struct of ClosingStock
type ClosingStock struct {
	Month int
	Year  int
}

// Closing stock
func (u *ClosingStock) Closing(ctx context.Context, db *sql.DB) error {
	var closingResult int
	q := `
		SELECT closing_stocks(lastClosing.month, lastClosing.year)
		FROM (
			SELECT year, month
			FROM saldo_stocks 
			WHERE company_id=? 
			ORDER BY CAST(CONCAT(year, month) AS UNSIGNED) DESC LIMIT 1 ) lastClosing
	`
	return db.QueryRowContext(ctx, q, ctx.Value(api.Ctx("auth")).(User).Company.ID).Scan(&closingResult)
}
