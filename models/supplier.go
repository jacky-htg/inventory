package models

import (
	"context"
	"database/sql"

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

// Get supplier by id
func (u *Supplier) Get(ctx context.Context, tx *sql.Tx) error {
	return tx.QueryRowContext(ctx, qSuppliers+" WHERE suppliers.id=? AND companies.id=?", u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID).Scan(u.getArgs()...)
}
