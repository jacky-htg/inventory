package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jacky-htg/inventory/libraries/api"
)

// Product : struct of Product
type Product struct {
	ID            uint64
	Code          string
	Name          string
	PurchasePrice float64
	SalePrice     float64
	Company       Company
}

const qProducts = `
SELECT 	products.id, 
		products.code, 
		products.name,
		products.sale_price, 
		companies.id, 
		companies.code, 
		companies.name,
		companies.address  
FROM products
JOIN companies ON products.company_id = companies.id
`

// List of products
func (u *Product) List(ctx context.Context, db *sql.DB) ([]Product, error) {
	list := []Product{}

	rows, err := db.QueryContext(ctx, qProducts+" WHERE companies.id=?", ctx.Value(api.Ctx("auth")).(User).Company.ID)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var r Product
		err = rows.Scan(r.getArgs()...)
		if err != nil {
			return list, err
		}

		list = append(list, r)
	}

	if err := rows.Err(); err != nil {
		return list, err
	}

	if len(list) <= 0 {
		return list, errors.New("Product not found")
	}

	return list, nil
}

// Get product by id
func (u *Product) Get(ctx context.Context, db *sql.DB) error {
	return db.QueryRowContext(ctx, qProducts+" WHERE products.id=? AND companies.id=?", u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID).Scan(u.getArgs()...)
}

// Create new product
func (u *Product) Create(ctx context.Context, db *sql.DB) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	const query = `
		INSERT INTO products (company_id, code, name, sale_price, created)
		VALUES (?, ?, ?, ?, NOW())
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, userLogin.Company.ID, u.Code, u.Name, u.SalePrice)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = uint64(id)
	u.Company = userLogin.Company

	return nil
}

// Update product
func (u *Product) Update(ctx context.Context, db *sql.DB) error {

	stmt, err := db.PrepareContext(ctx, `
		UPDATE products 
		SET name = ?,
			sale_price = ?,
			updated = NOW()
		WHERE id = ?
		AND company_id = ?
	`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Name, u.SalePrice, u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)
	return err
}

// Delete product
func (u *Product) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM products WHERE id = ? AND company_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)
	return err
}

func (u *Product) getArgs() []interface{} {
	var args []interface{}
	args = append(args, &u.ID)
	args = append(args, &u.Code)
	args = append(args, &u.Name)
	args = append(args, &u.SalePrice)
	args = append(args, &u.Company.ID)
	args = append(args, &u.Company.Code)
	args = append(args, &u.Company.Name)
	args = append(args, &u.Company.Address)

	return args
}
