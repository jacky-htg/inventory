package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jacky-htg/inventory/libraries/api"
)

// Product : struct of Product
type Product struct {
	ID              uint64
	Code            string
	Name            string
	PurchasePrice   float64
	SalePrice       float64
	MinimumStock    uint
	Company         Company
	Brand           Brand
	ProductCategory ProductCategory
}

const qProducts = `
SELECT 	products.id, 
		products.code, 
		products.name,
		products.sale_price,
		products.minimum_stock, 
		companies.id as company_id, 
		companies.code as company_code, 
		companies.name as company_name,
		companies.address as company_address,
		brands.id as brand_id,
		brands.code as brand_code,
		brands.name as brand_name,
		product_categories.id as product_category_id,
		product_categories.name as product_category_name  
FROM products
JOIN companies ON products.company_id = companies.id
JOIN brands ON products.brand_id = brands.id
JOIN product_categories ON products.product_category_id = product_categories.id
`

// List of products
func (u *Product) List(ctx context.Context, tx *sql.Tx) ([]Product, error) {
	list := []Product{}

	rows, err := tx.QueryContext(ctx, qProducts+" WHERE companies.id=?", ctx.Value(api.Ctx("auth")).(User).Company.ID)
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
func (u *Product) Get(ctx context.Context, tx *sql.Tx) error {
	return tx.QueryRowContext(ctx, qProducts+" WHERE products.id=? AND companies.id=?", u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID).Scan(u.getArgs()...)
}

// Create new product
func (u *Product) Create(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	const query = `
		INSERT INTO products (company_id, brand_id, product_category_id, code, name, sale_price, minimum_stock, created)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, userLogin.Company.ID, u.Brand.ID, u.ProductCategory.ID, u.Code, u.Name, u.SalePrice, u.MinimumStock)
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
func (u *Product) Update(ctx context.Context, tx *sql.Tx) error {

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE products 
		SET name = ?,
			sale_price = ?,
			brand_id = ?,
			product_category_id = ?,
			minimum_stock = ?,
			updated = NOW()
		WHERE id = ?
		AND company_id = ?
	`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Name, u.SalePrice, u.Brand.ID, u.ProductCategory.ID, u.MinimumStock, u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)
	return err
}

// Delete product
func (u *Product) Delete(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM products WHERE id = ? AND company_id = ?`)
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
	args = append(args, &u.MinimumStock)
	args = append(args, &u.Company.ID)
	args = append(args, &u.Company.Code)
	args = append(args, &u.Company.Name)
	args = append(args, &u.Company.Address)
	args = append(args, &u.Brand.ID)
	args = append(args, &u.Brand.Code)
	args = append(args, &u.Brand.Name)
	args = append(args, &u.ProductCategory.ID)
	args = append(args, &u.ProductCategory.Name)

	return args
}
