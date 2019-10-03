package models

import (
	"context"
	"database/sql"

	"github.com/jacky-htg/inventory/libraries/api"
)

// Customer : struct of customer
type Customer struct {
	ID      uint64
	Company Company
	Name    string
	Email   string
	Address string
	Hp      string
}

const qCustomers = `SELECT id, name, email, address, hp FROM customers`

// List of customers
func (u *Customer) List(ctx context.Context, tx *sql.Tx) ([]Customer, error) {
	var list []Customer

	rows, err := tx.QueryContext(ctx, qCustomers+" WHERE company_id=?", ctx.Value(api.Ctx("auth")).(User).Company.ID)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Customer
		c.Company = ctx.Value(api.Ctx("auth")).(User).Company
		err = rows.Scan(&c.ID, &c.Name, &c.Email, &c.Address, &c.Hp)
		if err != nil {
			return list, err
		}

		list = append(list, c)
	}

	return list, rows.Err()
}

// Create new customer
func (u *Customer) Create(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO customers (company_id, name, email, address, hp) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	userLogin := ctx.Value(api.Ctx("auth")).(User)
	res, err := stmt.ExecContext(ctx, userLogin.Company.ID, u.Name, u.Email, u.Address, u.Hp)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	u.ID = uint64(id)
	u.Company = userLogin.Company

	return err
}

// View customer by id
func (u *Customer) View(ctx context.Context, tx *sql.Tx) error {
	u.Company = ctx.Value(api.Ctx("auth")).(User).Company

	return tx.QueryRowContext(
		ctx,
		qCustomers+" WHERE id=? AND company_id=?",
		u.ID,
		ctx.Value(api.Ctx("auth")).(User).Company.ID,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Address, &u.Hp)
}

// Update customer by id
func (u *Customer) Update(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `
		UPDATE customers  
		SET name = ?, 
			email = ?, 
			address = ?,
			hp = ?
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

// Delete customer by id
func (u *Customer) Delete(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM customers WHERE id = ? AND company_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)

	return err
}
