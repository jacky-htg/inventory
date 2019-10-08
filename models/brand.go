package models

import (
	"context"
	"database/sql"

	"github.com/jacky-htg/inventory/libraries/api"
)

// Brand : struct of brand
type Brand struct {
	ID      uint64
	Company Company
	Code    string
	Name    string
}

const qBrands = `SELECT id, code, name FROM brands`

// List of brands
func (u *Brand) List(ctx context.Context, tx *sql.Tx) ([]Brand, error) {
	var list []Brand

	rows, err := tx.QueryContext(ctx, qBrands+" WHERE company_id=?", ctx.Value(api.Ctx("auth")).(User).Company.ID)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Brand
		c.Company = ctx.Value(api.Ctx("auth")).(User).Company
		err = rows.Scan(&c.ID, &c.Code, &c.Name)
		if err != nil {
			return list, err
		}

		list = append(list, c)
	}

	return list, rows.Err()
}

// Create new Brand
func (u *Brand) Create(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO brands (company_id, code, name) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	userLogin := ctx.Value(api.Ctx("auth")).(User)
	res, err := stmt.ExecContext(ctx, userLogin.Company.ID, u.Code, u.Name)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	u.ID = uint64(id)
	u.Company = userLogin.Company

	return err
}

// View Brand by id
func (u *Brand) View(ctx context.Context, tx *sql.Tx) error {
	u.Company = ctx.Value(api.Ctx("auth")).(User).Company

	return tx.QueryRowContext(
		ctx,
		qBrands+" WHERE id=? AND company_id=?",
		u.ID,
		ctx.Value(api.Ctx("auth")).(User).Company.ID,
	).Scan(&u.ID, &u.Code, &u.Name)
}

// Update Brand by id
func (u *Brand) Update(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `
		UPDATE brands  
		SET code = ?,
			name = ?
		WHERE id = ? AND company_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	userLogin := ctx.Value(api.Ctx("auth")).(User)
	u.Company = userLogin.Company
	_, err = stmt.ExecContext(ctx, u.Code, u.Name, u.ID, userLogin.Company.ID)

	return err
}

// Delete Brand by id
func (u *Brand) Delete(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM brands WHERE id = ? AND company_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)

	return err
}
