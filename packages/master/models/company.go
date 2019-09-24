package models

import (
	"context"
	"database/sql"
	"errors"
)

//Company : struct of Company
type Company struct {
	ID      uint32
	Code    string
	Name    string
	Address string
}

const qCompanies = `SELECT id, code, name, address FROM companies`

//List of companies
func (u *Company) List(ctx context.Context, db *sql.DB) ([]Company, error) {
	list := []Company{}

	rows, err := db.QueryContext(ctx, qCompanies)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Company
		err = rows.Scan(c.getArgs()...)
		if err != nil {
			return list, err
		}

		list = append(list, c)
	}

	if err := rows.Err(); err != nil {
		return list, err
	}

	if len(list) <= 0 {
		return list, errors.New("Company not found")
	}

	return list, nil
}

//Get company by id
func (u *Company) Get(ctx context.Context, db *sql.DB) error {
	return db.QueryRowContext(ctx, qCompanies+" WHERE id=?", u.ID).Scan(u.getArgs()...)
}

//Create new company
func (u *Company) Create(ctx context.Context, db *sql.DB) error {
	const query = `
		INSERT INTO companies (code, name, address, created)
		VALUES (?, ?, ?, NOW())
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, u.Code, u.Name, u.Address)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = uint32(id)

	return nil
}

//Update company
func (u *Company) Update(ctx context.Context, db *sql.DB) error {

	stmt, err := db.PrepareContext(ctx, `
		UPDATE companies 
		SET name = ?,
			address = ?,
			updated = NOW()
		WHERE id = ?
	`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Name, u.Address, u.ID)
	return err
}

//Delete company
func (u *Company) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM companies WHERE id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID)
	return err
}

func (u *Company) getArgs() []interface{} {
	var args []interface{}
	args = append(args, &u.ID)
	args = append(args, &u.Code)
	args = append(args, &u.Name)
	args = append(args, &u.Address)

	return args
}
