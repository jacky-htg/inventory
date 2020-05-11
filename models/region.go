package models

import (
	"context"
	"database/sql"

	"github.com/jacky-htg/inventory/libraries/api"
)

//Region : struct of Region
type Region struct {
	ID      uint32
	Code    string
	Name    string
	Company Company
}

const qRegions = `
SELECT 	regions.id, 
		regions.code, 
		regions.name, 
		companies.id, 
		companies.code, 
		companies.name,
		companies.address  
FROM regions
JOIN companies ON regions.company_id = companies.id
`

//List of regions
func (u *Region) List(ctx context.Context, db *sql.DB) ([]Region, error) {
	list := []Region{}

	rows, err := db.QueryContext(ctx, qRegions+" WHERE companies.id=?", ctx.Value(api.Ctx("auth")).(User).Company.ID)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var r Region
		err = rows.Scan(r.getArgs()...)
		if err != nil {
			return list, err
		}

		list = append(list, r)
	}

	if err := rows.Err(); err != nil {
		return list, err
	}

	/*if len(list) <= 0 {
		return list, errors.New("Region not found")
	}*/

	return list, nil
}

//Get region by id
func (u *Region) Get(ctx context.Context, db *sql.DB) error {
	return db.QueryRowContext(ctx, qRegions+" WHERE regions.id=? AND companies.id=?", u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID).Scan(u.getArgs()...)
}

//Create new region
func (u *Region) Create(ctx context.Context, db *sql.DB) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	const query = `
		INSERT INTO regions (company_id, code, name, created)
		VALUES (?, ?, ?, NOW())
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, userLogin.Company.ID, u.Code, u.Name)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = uint32(id)
	u.Company = userLogin.Company

	return nil
}

//Update region
func (u *Region) Update(ctx context.Context, db *sql.DB) error {

	stmt, err := db.PrepareContext(ctx, `
		UPDATE regions 
		SET name = ?,
			updated = NOW()
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

//Delete region
func (u *Region) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM regions WHERE id = ? AND company_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)
	return err
}

//AddBranch to region
func (u *Region) AddBranch(ctx context.Context, tx *sql.Tx, branchID uint32) error {
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO branches_regions (branch_id, region_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, branchID, u.ID)
	return err
}

//DeleteBranch from region
func (u *Region) DeleteBranch(ctx context.Context, tx *sql.Tx, branchID uint32) error {
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM branches_regions WHERE branch_id= ? AND region_id = ?`)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, branchID, u.ID)
	return err
}

// GetIDBranches by region id
func (u *Region) GetIDBranches(ctx context.Context, tx *sql.Tx) ([]uint32, error) {
	var list []uint32
	var err error

	rows, err := tx.QueryContext(ctx, "SELECT branch_id FROM branches_regions WHERE region_id=?", u.ID)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var temp uint32
		err = rows.Scan(&temp)
		if err != nil {
			return list, err
		}

		list = append(list, temp)
	}

	return list, rows.Err()
}

func (u *Region) getArgs() []interface{} {
	var args []interface{}
	args = append(args, &u.ID)
	args = append(args, &u.Code)
	args = append(args, &u.Name)
	args = append(args, &u.Company.ID)
	args = append(args, &u.Company.Code)
	args = append(args, &u.Company.Name)
	args = append(args, &u.Company.Address)

	return args
}
