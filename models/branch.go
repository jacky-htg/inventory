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
func (u *Branch) Get(ctx context.Context, tx *sql.Tx) error {
	return tx.QueryRowContext(ctx, qBranches+" WHERE branches.id=? AND companies.id=?", u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID).Scan(u.getArgs()...)
}
