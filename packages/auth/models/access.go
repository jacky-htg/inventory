package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/libraries/token"
)

//Access : struct of Access
type Access struct {
	ID       uint32
	ParentID sql.NullInt64
	Name     string
	Alias    string
}

const qAccess = `SELECT id, parent_id, name, alias FROM access`

//List of access
func (u *Access) List(ctx context.Context, tx *sql.Tx) ([]Access, error) {
	list := []Access{}

	rows, err := tx.QueryContext(ctx, qAccess)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var access Access
		err = rows.Scan(access.getArgs()...)
		if err != nil {
			return list, err
		}

		list = append(list, access)
	}

	if err := rows.Err(); err != nil {
		return list, err
	}

	if len(list) <= 0 {
		return list, errors.New("Access not found")
	}

	return list, nil
}

//GetByName : get access by name
func (u *Access) GetByName(ctx context.Context, tx *sql.Tx) error {
	return tx.QueryRowContext(ctx, qAccess+" WHERE name=?", u.Name).Scan(u.getArgs()...)
}

//GetByAlias : get access by alias
func (u *Access) GetByAlias(ctx context.Context, tx *sql.Tx) error {
	return tx.QueryRowContext(ctx, qAccess+" WHERE alias=?", u.Alias).Scan(u.getArgs()...)
}

//Get : get access by id
func (u *Access) Get(ctx context.Context, tx *sql.Tx) error {
	return tx.QueryRowContext(ctx, qAccess+" WHERE id=?", u.ID).Scan(u.getArgs()...)
}

//Create new Access
func (u *Access) Create(ctx context.Context, tx *sql.Tx) error {
	const query = `
		INSERT INTO access (parent_id, name, alias, created)
		VALUES (?, ?, ?, NOW())
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, u.ParentID, u.Name, u.Alias)
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

//Delete : delete access
func (u *Access) Delete(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM access WHERE id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID)
	return err
}

// GetIDs : get array of access id
func (u *Access) GetIDs(ctx context.Context, db *sql.DB) ([]uint32, error) {
	var access []uint32

	rows, err := db.QueryContext(ctx, "SELECT id FROM access WHERE name != 'root'")
	if err != nil {
		return access, err
	}

	defer rows.Close()

	for rows.Next() {
		var id uint32
		err = rows.Scan(&id)
		if err != nil {
			return access, err
		}
		access = append(access, id)
	}

	return access, rows.Err()
}

// IsAuth for check user authorization
func (u *Access) IsAuth(ctx context.Context, db *sql.DB, tokenparam interface{}, controller string, route string) (bool, User, error) {
	query := `
	SELECT true
	FROM users
	JOIN roles_users ON users.id = roles_users.user_id
	JOIN roles ON roles_users.role_id = roles.id
	JOIN access_roles ON roles.id = access_roles.role_id
	JOIN access ON access_roles.access_id = access.id
	WHERE (access.name = 'root' OR access.name = ? OR access.name = ?)
	AND users.id = ?
	`
	var isAuth bool
	var err error

	if tokenparam == nil {
		return isAuth, User{}, api.ErrBadRequest(errors.New("Bad request for token"), "")
	}

	isValid, username := token.ValidateToken(tokenparam.(string))
	if !isValid {
		return isAuth, User{}, api.ErrBadRequest(errors.New("Bad request for invalid token"), "")
	}

	user := User{Username: username}
	err = user.GetByUsername(ctx, db)
	if err != nil {
		return isAuth, User{}, err
	}

	err = db.QueryRowContext(ctx, query, controller, route, user.ID).Scan(&isAuth)

	return isAuth, user, err
}

func (u *Access) getArgs() []interface{} {
	var args []interface{}
	args = append(args, &u.ID)
	args = append(args, &u.ParentID)
	args = append(args, &u.Name)
	args = append(args, &u.Alias)
	return args
}
