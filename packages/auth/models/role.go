package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jacky-htg/inventory/libraries/api"
)

//Role : struct of Role
type Role struct {
	ID   uint32
	Name string
}

const qRoles = `SELECT id, name FROM roles`

//List of roles
func (u *Role) List(ctx context.Context, db *sql.DB) ([]Role, error) {
	list := []Role{}

	rows, err := db.QueryContext(ctx, qRoles+" WHERE company_id=?", ctx.Value(api.Ctx("auth")).(User).Company.ID)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var role Role
		err = rows.Scan(role.getArgs()...)
		if err != nil {
			return list, err
		}

		list = append(list, role)
	}

	if err := rows.Err(); err != nil {
		return list, err
	}

	if len(list) <= 0 {
		return list, errors.New("Role not found")
	}

	return list, nil
}

//Get role by id
func (u *Role) Get(ctx context.Context, db *sql.DB) error {
	return db.QueryRowContext(ctx, qRoles+" WHERE id=? AND company_id=?", u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID).Scan(u.getArgs()...)
}

//Create new role
func (u *Role) Create(ctx context.Context, db *sql.DB) error {
	const query = `
		INSERT INTO roles (name, company_id, created)
		VALUES (?, ?, NOW())
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, u.Name, ctx.Value(api.Ctx("auth")).(User).Company.ID)
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

//Update role
func (u *Role) Update(ctx context.Context, db *sql.DB) error {

	stmt, err := db.PrepareContext(ctx, `
		UPDATE roles 
		SET name = ?
		WHERE id = ?
		AND company_id = ?
	`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Name, u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)
	return err
}

//Delete role
func (u *Role) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM roles WHERE id = ? AND company_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)
	return err
}

//Grant access to role
func (u *Role) Grant(ctx context.Context, db *sql.DB, accessID uint32) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO access_roles (access_id, role_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, accessID, u.ID)
	return err
}

//Revoke access from role
func (u *Role) Revoke(ctx context.Context, db *sql.DB, accessID uint32) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM access_roles WHERE access_id= ? AND role_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, accessID, u.ID)
	return err
}

func (u *Role) getArgs() []interface{} {
	var args []interface{}
	args = append(args, &u.ID)
	args = append(args, &u.Name)
	return args
}
