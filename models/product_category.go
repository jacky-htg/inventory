package models

import (
	"context"
	"database/sql"

	"github.com/jacky-htg/inventory/libraries/api"
)

// ProductCategory : struct of ProductCategory
type ProductCategory struct {
	ID       uint64
	Company  Company
	Name     string
	Category Category
}

// Category : struct of Category
type Category struct {
	ID   uint
	Name string
}

const qProductCategories = `
	SELECT product_categories.id, categories.id category_id, categories.name parent_category, product_categories.name 
	FROM product_categories
	JOIN categories ON product_categories.category_id = categories.id
`

// List of ProductCategories
func (u *ProductCategory) List(ctx context.Context, tx *sql.Tx) ([]ProductCategory, error) {
	var list []ProductCategory

	rows, err := tx.QueryContext(ctx, qProductCategories+" WHERE company_id=?", ctx.Value(api.Ctx("auth")).(User).Company.ID)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var c ProductCategory
		c.Company = ctx.Value(api.Ctx("auth")).(User).Company
		err = rows.Scan(&c.ID, &c.Category.ID, &c.Category.Name, &c.Name)
		if err != nil {
			return list, err
		}

		list = append(list, c)
	}

	return list, rows.Err()
}

// Create new ProductCategory
func (u *ProductCategory) Create(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO product_categories (category_id, company_id, name) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	userLogin := ctx.Value(api.Ctx("auth")).(User)
	res, err := stmt.ExecContext(ctx, u.Category.ID, userLogin.Company.ID, u.Name)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	u.ID = uint64(id)
	u.Company = userLogin.Company
	u.Category.Get(ctx, tx)

	return err
}

// View ProductCategory by id
func (u *ProductCategory) View(ctx context.Context, tx *sql.Tx) error {
	u.Company = ctx.Value(api.Ctx("auth")).(User).Company

	return tx.QueryRowContext(
		ctx,
		qProductCategories+" WHERE product_categories.id=? AND product_categories.company_id=?",
		u.ID,
		ctx.Value(api.Ctx("auth")).(User).Company.ID,
	).Scan(&u.ID, &u.Category.ID, &u.Category.Name, &u.Name)
}

// Update ProductCategory by id
func (u *ProductCategory) Update(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `
		UPDATE product_categories  
		SET category_id = ?,
			name = ?
		WHERE id = ? AND company_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	userLogin := ctx.Value(api.Ctx("auth")).(User)
	u.Company = userLogin.Company
	_, err = stmt.ExecContext(ctx, u.Category.ID, u.Name, u.ID, userLogin.Company.ID)

	return err
}

// Delete ProductCategory by id
func (u *ProductCategory) Delete(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM product_categories WHERE id = ? AND company_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)

	return err
}

// Get Category by id
func (u *Category) Get(ctx context.Context, tx *sql.Tx) error {
	return tx.QueryRowContext(
		ctx,
		"SELECT id, name FROM categories WHERE id=?",
		u.ID,
	).Scan(&u.ID, &u.Name)
}
