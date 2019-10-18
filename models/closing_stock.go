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
	q := `
		SELECT year, month
		FROM saldo_stocks 
		WHERE company_id=? 
		ORDER BY CAST(CONCAT(year, month) AS UNSIGNED) DESC LIMIT 1 
	`
	err := db.QueryRowContext(ctx, q, ctx.Value(api.Ctx("auth")).(User).Company.ID).Scan(&u.Year, &u.Month)
	if err != nil {
		return err
	}

	q = `call closing_stocks(?, ?, ?)`
	stmt, err := db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, ctx.Value(api.Ctx("auth")).(User).Company.ID, u.Month, u.Year)
	return err
}
