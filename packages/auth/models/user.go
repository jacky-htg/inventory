package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/jacky-htg/inventory/libraries/array"
	master "github.com/jacky-htg/inventory/packages/master/models"
)

//User : struct of User
type User struct {
	ID       uint64
	Username string
	Password string
	Email    string
	IsActive bool
	Roles    []Role
	Company  master.Company
	Region   master.Region
	Branch   master.Branch
}

const qUsers = `
SELECT users.id, users.username, users.password, users.email, users.is_active, 
	JSON_ARRAYAGG(roles.id) as roles_id, JSON_ARRAYAGG(roles.name) as roles_name,
	companies.id, companies.code, companies.name, companies.address,
	regions.id, regions.code, regions.name,
	branches.id, branches.code, branches.name, branches.type, branches.address
FROM users
JOIN companies ON users.company_id = companies.id
LEFT JOIN regions ON users.region_id = regions.id
LEFT JOIN branches ON users.branch_id = branches.id
LEFT JOIN roles_users ON users.id=roles_users.user_id
LEFT JOIN roles ON roles_users.role_id=roles.id
`

//List : List of users
func (u *User) List(ctx context.Context, db *sql.DB) ([]User, error) {
	list := []User{}

	rows, err := db.QueryContext(ctx, qUsers+" GROUP BY users.id")
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		var roleIDs, roleNames string
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.IsActive,
			&roleIDs,
			&roleNames,
			&user.Company.ID,
			&user.Company.Code,
			&user.Company.Name,
			&user.Company.Address,
			&user.Region.ID,
			&user.Region.Code,
			&user.Region.Name,
			&user.Branch.ID,
			&user.Branch.Code,
			&user.Branch.Name,
			&user.Branch.Type,
			&user.Branch.Address,
		)
		if err != nil {
			return list, err
		}

		if len(roleIDs) > 0 {
			var ids []int32
			err = json.Unmarshal([]byte(roleIDs), &ids)
			if err != nil {
				return list, err
			}
			var names []string
			err = json.Unmarshal([]byte(roleNames), &names)
			if err != nil {
				return list, err
			}

			for i, v := range ids {
				user.Roles = append(user.Roles, Role{ID: uint32(v), Name: names[i]})
			}
		}

		list = append(list, user)
	}

	if err := rows.Err(); err != nil {
		return list, err
	}

	if len(list) <= 0 {
		return list, errors.New("Users not found")
	}

	return list, nil
}

//Get : get user by id
func (u *User) Get(ctx context.Context, db *sql.DB) error {
	var roleIDs, roleNames string
	err := db.QueryRowContext(ctx, qUsers+" WHERE users.id=? GROUP BY users.id", u.ID).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Email,
		&u.IsActive,
		&roleIDs,
		&roleNames,
		&u.Company.ID,
		&u.Company.Code,
		&u.Company.Name,
		&u.Company.Address,
		&u.Region.ID,
		&u.Region.Code,
		&u.Region.Name,
		&u.Branch.ID,
		&u.Branch.Code,
		&u.Branch.Name,
		&u.Branch.Type,
		&u.Branch.Address,
	)
	if err != nil {
		return err
	}

	if len(roleIDs) > 0 {
		var ids []int32
		err = json.Unmarshal([]byte(roleIDs), &ids)
		if err != nil {
			return err
		}
		var names []string
		err = json.Unmarshal([]byte(roleNames), &names)
		if err != nil {
			return err
		}

		for i, v := range ids {
			u.Roles = append(u.Roles, Role{ID: uint32(v), Name: names[i]})
		}
	}

	return nil
}

//GetByUsername : get user by username
func (u *User) GetByUsername(ctx context.Context, db *sql.DB) error {
	var roleIDs, roleNames string
	err := db.QueryRowContext(ctx, qUsers+" WHERE users.username=? GROUP BY users.id", u.Username).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Email,
		&u.IsActive,
		&roleIDs,
		&roleNames,
		&u.Company.ID,
		&u.Company.Code,
		&u.Company.Name,
		&u.Company.Address,
		&u.Region.ID,
		&u.Region.Code,
		&u.Region.Name,
		&u.Branch.ID,
		&u.Branch.Code,
		&u.Branch.Name,
		&u.Branch.Type,
		&u.Branch.Address,
	)
	if err != nil {
		return err
	}

	if len(roleIDs) > 0 {
		var ids []int32
		err = json.Unmarshal([]byte(roleIDs), &ids)
		if err != nil {
			return err
		}
		var names []string
		err = json.Unmarshal([]byte(roleNames), &names)
		if err != nil {
			return err
		}

		for i, v := range ids {
			u.Roles = append(u.Roles, Role{ID: uint32(v), Name: names[i]})
		}
	}

	return nil
}

//Create new user
func (u *User) Create(ctx context.Context, tx *sql.Tx) error {
	const query = `
		INSERT INTO users (username, password, email, is_active, created, updated)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, u.Username, u.Password, u.Email, u.IsActive)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = uint64(id)

	if len(u.Roles) > 0 {
		for _, r := range u.Roles {
			err = u.AddRole(ctx, tx, r)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//Update : update user
func (u *User) Update(ctx context.Context, tx *sql.Tx) error {

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE users 
		SET username = ?,
			password = ?,
			is_active = ?,
			updated = NOW()
		WHERE id = ?
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Username, u.Password, u.IsActive, u.ID)
	if err != nil {
		return err
	}

	existingRoles, err := u.GetRoleIDs(ctx, tx)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	for _, r := range u.Roles {
		if r.ID > 0 {
			var array array.ArrUint32
			isExist, index := array.InArray(r.ID, existingRoles)
			if !isExist {
				err = u.AddRole(ctx, tx, r)
				if err != nil {
					return err
				}
			} else {
				existingRoles = array.RemoveByIndex(existingRoles, index)
			}
		}
	}

	for _, r := range existingRoles {
		err = u.DeleteRole(ctx, tx, r)
		if err != nil {
			return err
		}
	}

	return nil
}

//Delete : delete user
func (u *User) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM users WHERE id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID)
	return err
}

// GetRoleIDs : get array of role id from user
func (u *User) GetRoleIDs(ctx context.Context, tx *sql.Tx) ([]uint32, error) {
	var roles []uint32

	rows, err := tx.QueryContext(ctx, "SELECT role_id FROM roles_users WHERE user_id = ?", u.ID)
	if err != nil {
		return roles, err
	}

	defer rows.Close()

	for rows.Next() {
		var id uint32
		err = rows.Scan(&id)
		if err != nil {
			return roles, err
		}
		roles = append(roles, id)
	}

	return roles, rows.Err()
}

//AddRole to user
func (u *User) AddRole(ctx context.Context, tx *sql.Tx, r Role) error {
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO roles_users (role_id, user_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, r.ID, u.ID)
	return err
}

//DeleteRole from user
func (u *User) DeleteRole(ctx context.Context, tx *sql.Tx, roleID uint32) error {
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM roles_users WHERE role_id=? AND user_id=?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, roleID, u.ID)
	return err
}
