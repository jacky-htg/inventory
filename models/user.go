package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/libraries/array"
)

//User : struct of User
type User struct {
	ID       uint64
	Username string
	Password string
	Email    string
	IsActive bool
	Roles    []Role
	Company  Company
	Region   Region
	Branch   Branch
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
	var err error
	var where []string
	var query string
	var branches []uint32
	var params []interface{}

	userLogin := ctx.Value(api.Ctx("auth")).(User)

	where = append(where, "users.company_id=?")
	params = append(params, userLogin.Company.ID)

	query = qUsers
	if userLogin.Region.ID > 0 {
		tx, err := db.Begin()
		if err != nil {
			return list, err
		}

		branches, err = u.Region.GetIDBranches(ctx, tx)
		if err != nil {
			tx.Rollback()
			return list, err
		}
		tx.Commit()
	} else {
		if userLogin.Branch.ID > 0 {
			branches = append(branches, userLogin.Branch.ID)
		}
	}

	if len(branches) > 0 {
		var oRBranchConditions []string
		for _, b := range branches {
			oRBranchConditions = append(oRBranchConditions, "users.branch_id=?")
			params = append(params, b)
		}
		if len(oRBranchConditions) > 0 {
			where = append(where, "("+strings.Join(oRBranchConditions, " OR ")+")")
		}
	}

	query += " WHERE " + strings.Join(where, " AND ") + " GROUP BY users.id"

	rows, err := db.QueryContext(ctx, query, params...)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		var roleIDs, roleNames string
		var regionID, branchID sql.NullInt64
		var regionCode, regionName, branchCode, branchName, branchType sql.NullString
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
			&regionID,
			&regionCode,
			&regionName,
			&branchID,
			&branchCode,
			&branchName,
			&branchType,
			&user.Branch.Address,
		)
		if err != nil {
			return list, err
		}

		if regionID.Int64 > 0 {
			user.Region = Region{ID: uint32(regionID.Int64), Code: regionCode.String, Name: regionName.String}
		}

		if branchID.Int64 > 0 {
			user.Branch = Branch{
				ID:      uint32(branchID.Int64),
				Code:    branchCode.String,
				Name:    branchName.String,
				Type:    branchType.String,
				Address: user.Branch.Address,
			}
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
	var err error
	var where []string
	var query string
	var branches []uint32
	var params []interface{}

	userLogin := ctx.Value(api.Ctx("auth")).(User)

	where = append(where, "users.id=?")
	params = append(params, u.ID)

	where = append(where, "users.company_id=?")
	params = append(params, userLogin.Company.ID)

	query = qUsers
	if userLogin.Region.ID > 0 {
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		branches, err = u.Region.GetIDBranches(ctx, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	} else {
		if userLogin.Branch.ID > 0 {
			branches = append(branches, userLogin.Branch.ID)
		}
	}

	if len(branches) > 0 {
		var oRBranchConditions []string
		for _, b := range branches {
			oRBranchConditions = append(oRBranchConditions, "users.branch_id=?")
			params = append(params, b)
		}
		if len(oRBranchConditions) > 0 {
			where = append(where, "("+strings.Join(oRBranchConditions, " OR ")+")")
		}
	}

	query += " WHERE " + strings.Join(where, " AND ") + " GROUP BY users.id"

	var roleIDs, roleNames string
	var regionID, branchID sql.NullInt64
	var regionCode, regionName, branchCode, branchName, branchType sql.NullString
	err = db.QueryRowContext(ctx, query, params...).Scan(
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
		&regionID,
		&regionCode,
		&regionName,
		&branchID,
		&branchCode,
		&branchName,
		&branchType,
		&u.Branch.Address,
	)
	if err != nil {
		return err
	}

	if regionID.Int64 > 0 {
		u.Region = Region{ID: uint32(regionID.Int64), Code: regionCode.String, Name: regionName.String}
	}

	if branchID.Int64 > 0 {
		u.Branch = Branch{
			ID:      uint32(branchID.Int64),
			Code:    branchCode.String,
			Name:    branchName.String,
			Type:    branchType.String,
			Address: u.Branch.Address,
		}
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

// GetByUsername : get user by username
// BEWARE : DONT CALL THIS FUNCTION
// this function just call in login only
func (u *User) GetByUsername(ctx context.Context, db *sql.DB) error {
	var roleIDs, roleNames string
	var regionID, branchID sql.NullInt64
	var regionCode, regionName, branchCode, branchName, branchType sql.NullString
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
		&regionID,
		&regionCode,
		&regionName,
		&branchID,
		&branchCode,
		&branchName,
		&branchType,
		&u.Branch.Address,
	)
	if err != nil {
		return err
	}

	if regionID.Int64 > 0 {
		u.Region = Region{ID: uint32(regionID.Int64), Code: regionCode.String, Name: regionName.String}
	}

	if branchID.Int64 > 0 {
		u.Branch = Branch{
			ID:      uint32(branchID.Int64),
			Code:    branchCode.String,
			Name:    branchName.String,
			Type:    branchType.String,
			Address: u.Branch.Address,
		}
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
	var regionID, branchID sql.NullInt64
	var branches, regions []uint32
	var err error
	userLogin := ctx.Value(api.Ctx("auth")).(User)

	switch {
	case userLogin.Branch.ID <= 0 && userLogin.Region.ID <= 0:
		branches, err = u.Company.GetIDBranches(ctx, tx)
		if err != nil {
			return err
		}

		regions, err = u.Company.GetIDRegions(ctx, tx)
		if err != nil {
			return err
		}
	case userLogin.Region.ID > 0:
		branches, err = u.Region.GetIDBranches(ctx, tx)
		if err != nil {
			return err
		}

		// user region can not create user head office
		if u.Branch.ID <= 0 {
			return api.ErrForbidden(errors.New("Forbidden data owner"), "")
		}

		regions = []uint32{}
	case userLogin.Branch.ID > 0:
		branches = []uint32{userLogin.Branch.ID}
		regions = []uint32{}
	}

	if u.Branch.ID > 0 {
		var arr array.ArrUint32
		isExist, _ := arr.InArray(u.Branch.ID, branches)

		if !isExist {
			return api.ErrForbidden(errors.New("Forbidden data owner"), "")
		}
		branchID.Valid = true
		branchID.Int64 = int64(u.Branch.ID)
	}

	if u.Region.ID > 0 {
		var arr array.ArrUint32
		isExist, _ := arr.InArray(u.Region.ID, regions)

		if !isExist {
			return api.ErrForbidden(errors.New("Forbidden data owner"), "")
		}
		regionID.Valid = true
		regionID.Int64 = int64(u.Region.ID)
	}

	const query = `
		INSERT INTO users (username, password, email, is_active, company_id, region_id, branch_id, created, updated)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, u.Username, u.Password, u.Email, u.IsActive, userLogin.Company.ID, regionID, branchID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = uint64(id)
	u.Company = userLogin.Company
	if u.Branch.ID > 0 {
		err = u.Branch.Get(ctx, tx)
		if err != nil {
			return err
		}
	}

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
	var regionID, branchID sql.NullInt64
	var branches, regions []uint32
	var err error
	userLogin := ctx.Value(api.Ctx("auth")).(User)

	switch {
	case userLogin.Branch.ID <= 0 && userLogin.Region.ID <= 0:
		branches, err = u.Company.GetIDBranches(ctx, tx)
		if err != nil {
			return err
		}

		regions, err = u.Company.GetIDRegions(ctx, tx)
		if err != nil {
			return err
		}
	case userLogin.Region.ID > 0:
		branches, err = u.Region.GetIDBranches(ctx, tx)
		if err != nil {
			return err
		}

		// user region can not create user head office
		if u.Branch.ID <= 0 {
			return api.ErrForbidden(errors.New("Forbidden data owner"), "")
		}

		regions = []uint32{}
	case userLogin.Branch.ID > 0:
		branches = []uint32{userLogin.Branch.ID}
		regions = []uint32{}
	}

	if u.Branch.ID > 0 {
		var arr array.ArrUint32
		isExist, _ := arr.InArray(u.Branch.ID, branches)

		if !isExist {
			return api.ErrForbidden(errors.New("Forbidden data owner"), "")
		}
		branchID.Valid = true
		branchID.Int64 = int64(u.Branch.ID)
	}

	if u.Region.ID > 0 {
		var arr array.ArrUint32
		isExist, _ := arr.InArray(u.Region.ID, regions)

		if !isExist {
			return api.ErrForbidden(errors.New("Forbidden data owner"), "")
		}
		regionID.Valid = true
		regionID.Int64 = int64(u.Region.ID)
	}

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE users 
		SET username = ?,
			password = ?,
			is_active = ?,
			region_id = ?,
			branch_id = ?,
			updated = NOW()
		WHERE id = ?
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Username, u.Password, u.IsActive, regionID, branchID, u.ID)
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
