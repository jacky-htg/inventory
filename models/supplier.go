package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jacky-htg/inventory/libraries/api"
)

// Supplier : struct of Supplier
type Supplier struct {
	ID      uint64
	Code    string
	Name    string
	Address sql.NullString
	Company Company
}

const qSuppliers = `
SELECT 	suppliers.id, 
	suppliers.code, 
	suppliers.name,
	suppliers.address,
	companies.id, 
	companies.code, 
	companies.name,
	companies.address  
FROM suppliers
JOIN companies ON suppliers.company_id = companies.id
`

func (u *Supplier) getArgs() []interface{} {
	var args []interface{}
	args = append(args, &u.ID)
	args = append(args, &u.Code)
	args = append(args, &u.Name)
	args = append(args, &u.Address)
	args = append(args, &u.Company.ID)
	args = append(args, &u.Company.Code)
	args = append(args, &u.Company.Name)
	args = append(args, &u.Company.Address)

	return args
}

// List of suppliers
func (u *Supplier) List(ctx context.Context, db *sql.DB) ([]Supplier, error) {
	list := []Supplier{}

	rows, err := db.QueryContext(ctx, qSuppliers+" WHERE companies.id=?", ctx.Value(api.Ctx("auth")).(User).Company.ID)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var s Supplier
		err = rows.Scan(s.getArgs()...)
		if err != nil {
			return list, err
		}

		list = append(list, s)
	}

	if err := rows.Err(); err != nil {
		return list, err
	}

	if len(list) <= 0 {
		return list, errors.New("Supplier not found")
	}

	return list, nil
}

// Get supplier by id
func (u *Supplier) Get(ctx context.Context, tx *sql.Tx) error {
	return tx.QueryRowContext(ctx, qSuppliers+" WHERE suppliers.id=? AND companies.id=?", u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID).Scan(u.getArgs()...)
}

// Create new supplier
func (u *Supplier) Create(ctx context.Context, db *sql.DB) error {
	const query = `
		INSERT INTO suppliers (company_id, code, name, address, created)
		VALUES (?, ?, ?, ?, NOW())
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, ctx.Value(api.Ctx("auth")).(User).Company.ID, u.Code, u.Name, u.Address)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = uint64(id)

	return nil
}

// Update supplier
func (u *Supplier) Update(ctx context.Context, db *sql.DB) error {

	stmt, err := db.PrepareContext(ctx, `
		UPDATE suppliers 
		SET name = ?,
			address = ?,
			updated = NOW()
		WHERE id = ? AND company_id = ?
	`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Name, u.Address, u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)
	return err
}

// Delete supplier
func (u *Supplier) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM suppliers WHERE id = ? AND company_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)
	return err
}
