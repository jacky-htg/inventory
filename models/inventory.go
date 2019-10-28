package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jacky-htg/inventory/libraries/api"
)

// Inventory : struct of Receive
type Inventory struct {
	ID              uint64
	CompanyID       uint32
	BranchID        uint32
	ShelveID        uint64
	ProductID       uint64
	ProductCode     string
	TransactionID   uint64
	Code            string
	TransactionDate time.Time
	Type            string
	InOut           bool
	Qty             uint
	Created         time.Time
	Updated         time.Time
}

// Create new inventory
func (u *Inventory) Create(ctx context.Context, tx *sql.Tx) error {
	var err error
	userLogin := ctx.Value(api.Ctx("auth")).(User)

	const queryDetail = `
		INSERT INTO inventories (company_id, branch_id, product_id, product_code, transaction_id, code, transaction_date, type, in_out, qty, shelve_id, created, updated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1, ?, NOW(), NOW())
	`
	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		userLogin.Company.ID,
		userLogin.Branch.ID,
		u.ProductID,
		u.ProductCode,
		u.TransactionID,
		u.Code,
		u.TransactionDate,
		u.Type,
		u.InOut,
		u.ShelveID,
	)
	return err
}

// Update inventory
func (u *Inventory) Update(ctx context.Context, tx *sql.Tx) error {
	var err error
	userLogin := ctx.Value(api.Ctx("auth")).(User)

	const queryUpdate = `
		UPDATE inventories
		SET shelve_id = ?,
			product_id = ?, 
			product_code = ?, 
			transaction_id = ?, 
			code = ?, 
			transaction_date = ?, 
			updated= NOW()
		WHERE id = ? AND company_id = ? AND branch_id = ?
	`
	stmt, err := tx.PrepareContext(ctx, queryUpdate)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.ShelveID,
		u.ProductID,
		u.ProductCode,
		u.TransactionID,
		u.Code,
		u.TransactionDate,
		u.ID,
		userLogin.Company.ID,
		userLogin.Branch.ID,
	)
	return err
}

// Delete Inventory
func (u *Inventory) Delete(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Company.ID != u.CompanyID || u.BranchID <= 0 || userLogin.Branch.ID != u.BranchID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `DELETE FROM inventories WHERE id = ? AND company_id = ? AND branch_id = ?`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID, userLogin.Company.ID, userLogin.Branch.ID)
	return err
}

// DeleteByComposit Inventory
func (u *Inventory) DeleteByComposit(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Company.ID != u.CompanyID || u.BranchID <= 0 || userLogin.Branch.ID != u.BranchID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `DELETE FROM inventories WHERE product_id = ? AND product_code = ? AND transaction_id = ? AND type = ? AND company_id = ? AND branch_id = ?`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ProductID, u.ProductCode, u.TransactionID, u.Type, userLogin.Company.ID, userLogin.Branch.ID)
	return err
}

// GetByComposit Inventory
func (u *Inventory) GetByComposit(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	const query = `
		SELECT id, company_id, branch_id, shelve_id, code, transaction_date, in_out, qty  
		FROM inventories 
		WHERE product_id = ? AND product_code = ? AND transaction_id = ? AND type = ? AND company_id = ? AND branch_id = ?`
	err := tx.QueryRowContext(ctx, query,
		u.ProductID, u.ProductCode, u.TransactionID, u.Type, userLogin.Company.ID, userLogin.Branch.ID).Scan(
		&u.ID, &u.CompanyID, &u.BranchID, &u.ShelveID, &u.Code, &u.TransactionDate, &u.InOut, &u.Qty)

	return err
}
