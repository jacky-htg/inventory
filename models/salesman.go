package models

import (
	"context"
	"database/sql"

	"github.com/jacky-htg/inventory/libraries/api"
)

// Salesman : struct of salesman
type Salesman struct {
	ID      uint64
	Company Company
	Code    string
	Name    string
	Email   string
	Address string
	Hp      string
}

const qSalesmen = `SELECT id, code, name, email, address, hp FROM salesmen`

// List of salesmen
func (u *Salesman) List(ctx context.Context, db *sql.DB) ([]Salesman, error) {
	var list []Salesman

	rows, err := db.QueryContext(ctx, qSalesmen+" WHERE company_id=?", ctx.Value(api.Ctx("auth")).(User).Company.ID)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Salesman
		c.Company = ctx.Value(api.Ctx("auth")).(User).Company
		err = rows.Scan(&c.ID, &c.Code, &c.Name, &c.Email, &c.Address, &c.Hp)
		if err != nil {
			return list, err
		}

		list = append(list, c)
	}

	return list, rows.Err()
}

// Create new salesman
func (u *Salesman) Create(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO salesmen (company_id, code, name, email, address, hp, created) VALUES (?, ?, ?, ?, ?, ?, NOW())`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	userLogin := ctx.Value(api.Ctx("auth")).(User)
	res, err := stmt.ExecContext(ctx, userLogin.Company.ID, u.Code, u.Name, u.Email, u.Address, u.Hp)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	u.ID = uint64(id)
	u.Company = userLogin.Company

	return err
}

// Get salesman by id
func (u *Salesman) Get(ctx context.Context, tx *sql.Tx) error {
	u.Company = ctx.Value(api.Ctx("auth")).(User).Company

	return tx.QueryRowContext(
		ctx,
		qSalesmen+" WHERE id=? AND company_id=?",
		u.ID,
		ctx.Value(api.Ctx("auth")).(User).Company.ID,
	).Scan(&u.ID, &u.Code, &u.Name, &u.Email, &u.Address, &u.Hp)
}

// Update salesman by id
func (u *Salesman) Update(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `
		UPDATE salesmen  
		SET name = ?, 
			email = ?, 
			address = ?,
			hp = ?,
			updated = NOW()
		WHERE id = ? AND company_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	userLogin := ctx.Value(api.Ctx("auth")).(User)
	u.Company = userLogin.Company
	_, err = stmt.ExecContext(ctx, u.Name, u.Email, u.Address, u.Hp, u.ID, userLogin.Company.ID)

	return err
}

// Delete salesman by id
func (u *Salesman) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM salesmen WHERE id = ? AND company_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)

	return err
}
