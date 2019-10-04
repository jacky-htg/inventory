package models

import (
	"context"
	"database/sql"

	"github.com/jacky-htg/inventory/libraries/api"
)

//Branch : struct of Branch
type Branch struct {
	ID      uint32
	Code    string
	Name    string
	Address sql.NullString
	Type    string
	Company Company
}

const qBranches = `
SELECT 	branches.id, 
	branches.code, 
	branches.name,
	branches.address,
	branches.type, 
	companies.id, 
	companies.code, 
	companies.name,
	companies.address  
FROM branches
JOIN companies ON branches.company_id = companies.id
`

const qOnlyBranches = `
SELECT 	id, 
	code, 
	name,
	address,
	type
FROM branches
`

func (u *Branch) getArgs() []interface{} {
	var args []interface{}
	args = append(args, &u.ID)
	args = append(args, &u.Code)
	args = append(args, &u.Name)
	args = append(args, &u.Address)
	args = append(args, &u.Type)
	args = append(args, &u.Company.ID)
	args = append(args, &u.Company.Code)
	args = append(args, &u.Company.Name)
	args = append(args, &u.Company.Address)

	return args
}

// Get branch by id
func (b *Branch) Get(ctx context.Context, tx *sql.Tx) error {
	return tx.QueryRowContext(ctx, qBranches+" WHERE branches.id=? AND companies.id=?", b.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID).Scan(b.getArgs()...)
}

// List all branches
func (b *Branch) List(ctx context.Context, tx *sql.Tx) ([]Branch, error) {
	var list []Branch

	rows, err := tx.QueryContext(ctx, qOnlyBranches+"WHERE company_id=?", ctx.Value(api.Ctx("auth")).(User).Company.ID)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var r Branch
		r.Company = ctx.Value(api.Ctx("auth")).(User).Company
		err = rows.Scan(&r.ID, &r.Code, &r.Name, &r.Address, &r.Type)
		if err != nil {
			return list, err
		}

		list = append(list, r)
	}

	return list, rows.Err()
}

// Create new branch
func (b *Branch) Create(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO branches (code, name, address, type, company_id) VALUES (?,?,?,?,?)
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	userLogin := ctx.Value(api.Ctx("auth")).(User)
	res, err := stmt.ExecContext(ctx, b.Code, b.Name, b.Address, b.Type, userLogin.Company.ID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	b.ID = uint32(id)
	b.Company = userLogin.Company

	return err
}

// Update branch by id
func (b *Branch) Update(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `
		UPDATE branches
		SET 
			code = ?,
			name = ?,
			address = ?,
			type = ?
		WHERE id = ? AND company_id = ?
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	company := ctx.Value(api.Ctx("auth")).(User).Company
	b.Company = company
	_, err = stmt.ExecContext(ctx, b.Code, b.Name, b.Address, b.Type, b.ID, company.ID)

	return err
}

// Delete branch by id
func (b *Branch) Delete(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM branches WHERE id = ? AND company_id = ?`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, b.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)

	return nil
}
