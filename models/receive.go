package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/libraries/array"
)

// Receive : struct of Receive
type Receive struct {
	ID             uint64
	Code           string
	Date           time.Time
	Remark         string
	Purchase       Purchase
	Company        Company
	Branch         Branch
	ReceiveDetails []ReceiveDetail
}

// ReceiveDetail struct
type ReceiveDetail struct {
	ID      uint64
	Product Product
	Qty     uint
	Code    string
	Shelve  Shelve
}

// List Receives
func (u *Receive) List(ctx context.Context, tx *sql.Tx) ([]Receive, error) {
	var list []Receive
	var err error

	query := `
	SELECT 	good_receivings.id, 
		good_receivings.code, 
		good_receivings.date,
		purchases.id,
		purchases.code,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type
	FROM good_receivings
	JOIN companies ON good_receivings.company_id = companies.id
	JOIN purchases ON good_receivings.purchase_id = purchases.id AND good_receivings.company_id = purchases.company_id AND good_receivings.branch_id = purchases.branch_id 
	JOIN branches ON good_receivings.branch_id = branches.id
	JOIN good_receiving_details ON good_receivings.id = good_receiving_details.good_receiving_id
	WHERE companies.id=?
	`
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	params := []interface{}{userLogin.Company.ID}

	switch {
	case userLogin.Region.ID > 0:
		branches, err := userLogin.Region.GetIDBranches(ctx, tx)
		if err != nil {
			return list, err
		}

		var orWhere []string
		for _, b := range branches {
			orWhere = append(orWhere, "branches.id=?")
			params = append(params, b)
		}

		query += " AND (" + strings.Join(orWhere, " OR ") + ")"

	case userLogin.Branch.ID > 0:
		query += " AND branches.id=?"
		params = append(params, userLogin.Branch.ID)
	}

	query += " GROUP BY good_receivings.id"

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var Receive Receive
		err = rows.Scan(
			&Receive.ID,
			&Receive.Code,
			&Receive.Date,
			&Receive.Purchase.ID,
			&Receive.Purchase.Code,
			&Receive.Company.ID,
			&Receive.Company.Code,
			&Receive.Company.Name,
			&Receive.Company.Address,
			&Receive.Branch.ID,
			&Receive.Branch.Code,
			&Receive.Branch.Name,
			&Receive.Branch.Address,
			&Receive.Branch.Type,
		)

		if err != nil {
			return list, err
		}

		list = append(list, Receive)
	}

	return list, rows.Err()
}

// Get Receive by id
func (u *Receive) Get(ctx context.Context, tx *sql.Tx) error {
	query := `
	SELECT 	good_receivings.id, 
		good_receivings.code, 
		good_receivings.date,
		purchases.id,
		purchases.code,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type,
		JSON_ARRAYAGG(good_receiving_details.id),
		JSON_ARRAYAGG(good_receiving_details.code),
		JSON_ARRAYAGG(good_receiving_details.shelve_id),
		JSON_ARRAYAGG(good_receiving_details.qty),
		JSON_ARRAYAGG(products.id),
		JSON_ARRAYAGG(products.code),
		JSON_ARRAYAGG(products.name),
		JSON_ARRAYAGG(products.sale_price)
	FROM good_receivings
	JOIN companies ON good_receivings.company_id = companies.id
	JOIN purchases ON good_receivings.purchase_id = purchases.id AND good_receivings.company_id = purchases.company_id AND good_receivings.branch_id = purchases.branch_id
	JOIN branches ON good_receivings.branch_id = branches.id
	JOIN good_receiving_details ON good_receivings.id = good_receiving_details.good_receiving_id
	JOIN products ON good_receiving_details.product_id = products.id 
	WHERE good_receivings.id=? AND companies.id=?
	`
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	params := []interface{}{u.ID, userLogin.Company.ID}

	switch {
	case userLogin.Region.ID > 0:
		branches, err := userLogin.Region.GetIDBranches(ctx, tx)
		if err != nil {
			return err
		}

		var orWhere []string
		for _, b := range branches {
			orWhere = append(orWhere, "branches.id=?")
			params = append(params, b)
		}

		query += " AND (" + strings.Join(orWhere, " OR ") + ")"

	case userLogin.Branch.ID > 0:
		query += " AND branches.id=?"
		params = append(params, userLogin.Branch.ID)
	}

	var detailID, detailCode, detailShelveID, detailQty, productID, productCode, productName, productPrice string
	err := tx.QueryRowContext(ctx, query+" GROUP BY good_receivings.id", params...).Scan(
		&u.ID,
		&u.Code,
		&u.Date,
		&u.Purchase.ID,
		&u.Purchase.Code,
		&u.Company.ID,
		&u.Company.Code,
		&u.Company.Name,
		&u.Company.Address,
		&u.Branch.ID,
		&u.Branch.Code,
		&u.Branch.Name,
		&u.Branch.Address,
		&u.Branch.Type,
		&detailID,
		&detailCode,
		&detailShelveID,
		&detailQty,
		&productID,
		&productCode,
		&productName,
		&productPrice,
	)

	if err != nil {
		return err
	}

	if len(detailID) > 0 {
		var detailIDs []uint64
		err = json.Unmarshal([]byte(detailID), &detailIDs)
		if err != nil {
			return err
		}

		var detailCodes []string
		err = json.Unmarshal([]byte(detailCode), &detailCodes)
		if err != nil {
			return err
		}

		var detailShelveIDs []uint64
		err = json.Unmarshal([]byte(detailShelveID), &detailShelveIDs)
		if err != nil {
			return err
		}

		var detailQtys []uint
		err = json.Unmarshal([]byte(detailQty), &detailQtys)
		if err != nil {
			return err
		}

		var productIDs []uint64
		err = json.Unmarshal([]byte(productID), &productIDs)
		if err != nil {
			return err
		}

		var productCodes []string
		err = json.Unmarshal([]byte(productCode), &productCodes)
		if err != nil {
			return err
		}

		var productNames []string
		err = json.Unmarshal([]byte(productName), &productNames)
		if err != nil {
			return err
		}

		var productPrices []float64
		err = json.Unmarshal([]byte(productPrice), &productPrices)
		if err != nil {
			return err
		}

		for i, v := range detailIDs {
			u.ReceiveDetails = append(u.ReceiveDetails, ReceiveDetail{
				ID:   uint64(v),
				Code: detailCodes[i],
				Shelve: Shelve{
					ID: detailShelveIDs[i],
				},
				Qty: detailQtys[i],
				Product: Product{
					ID:        productIDs[i],
					Code:      productCodes[i],
					Name:      productNames[i],
					SalePrice: productPrices[i],
					Company:   u.Company,
				},
			})
		}
	}

	return nil
}

// Create new Receive
func (u *Receive) Create(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Branch.ID <= 0 {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	err := u.Purchase.Get(ctx, tx)
	if err != nil {
		return err
	}

	if u.Purchase.Branch.ID != userLogin.Branch.ID || u.Purchase.Company.ID != userLogin.Company.ID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		INSERT INTO good_receivings (code, date, remark, purchase_id, company_id, branch_id, created_by, updated_by, created, updated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	u.Code, err = api.GetCode(ctx, tx, "GR", "good_receivings", userLogin.Company.ID)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, u.Code, u.Date, u.Remark, u.Purchase.ID, userLogin.Company.ID, userLogin.Branch.ID, userLogin.ID, userLogin.ID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = uint64(id)
	u.Company = userLogin.Company
	u.Branch = userLogin.Branch
	u.Branch.Company = u.Company

	for i, d := range u.ReceiveDetails {
		err := u.storeDetail(ctx, tx, d, i)
		if err != nil {
			return err
		}
		u.ReceiveDetails[i].Product.Get(ctx, tx)
	}

	return nil
}

// Update Receive
func (u *Receive) Update(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Company.ID != u.Company.ID || u.Branch.ID <= 0 || userLogin.Branch.ID != u.Branch.ID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	err := u.Purchase.Get(ctx, tx)
	if err != nil {
		return err
	}

	if u.Purchase.Branch.ID != userLogin.Branch.ID || u.Purchase.Company.ID != userLogin.Company.ID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		UPDATE good_receivings 
		SET date = ?, 
			remark = ?,
			purchase_id = ?, 
			updated_by = ?, 
			updated = NOW()
		WHERE id = ?
		AND company_id = ?
		AND branch_id = ?
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Date, u.Remark, u.Purchase.ID, userLogin.ID, u.ID, userLogin.Company.ID, userLogin.Branch.ID)
	if err != nil {
		return err
	}

	existingDetails, err := u.GetExistingDetails(ctx, tx)
	if err != nil {
		return err
	}

	for i, d := range u.ReceiveDetails {
		if d.ID <= 0 {
			err := u.storeDetail(ctx, tx, d, i)
			if err != nil {
				return err
			}
		} else {
			detail, err := u.getDetail(ctx, tx, d.ID)
			if err != nil {
				return err
			}

			d.Code = detail.Code

			err = u.updateDetail(ctx, tx, d)
			if err != nil {
				return err
			}

			var arrUint64 array.ArrUint64
			existingDetails = arrUint64.Remove(existingDetails, d.ID)
		}

		u.ReceiveDetails[i].Product.Get(ctx, tx)
	}

	for _, e := range existingDetails {
		err = u.removeDetail(ctx, tx, e)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetExistingDetails return array of existing receive_details id
func (u *Receive) GetExistingDetails(ctx context.Context, tx *sql.Tx) ([]uint64, error) {
	var list []uint64
	var err error

	rows, err := tx.QueryContext(ctx, "SELECT id FROM good_receiving_details WHERE good_receiving_id=?", u.ID)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var temp uint64
		err = rows.Scan(&temp)
		if err != nil {
			return list, err
		}

		list = append(list, temp)
	}

	return list, rows.Err()
}

func (u *Receive) getDetail(ctx context.Context, tx *sql.Tx, e uint64) (*ReceiveDetail, error) {
	detail := new(ReceiveDetail)
	err := tx.QueryRowContext(ctx,
		`SELECT id, product_id, qty, code, shelve_id FROM good_receiving_details WHERE id = ? AND good_receiving_id = ?`,
		e, u.ID).Scan(&detail.ID, &detail.Product.ID, &detail.Qty, &detail.Code, &detail.Shelve.ID)

	return detail, err
}

func (u *Receive) storeDetail(ctx context.Context, tx *sql.Tx, d ReceiveDetail, i int) error {
	var err error

	const queryDetail = `
		INSERT INTO good_receiving_details (good_receiving_id, product_id, qty, code, shelve_id)
		VALUES (?, ?, ?, ?, ?)
	`

	d.Code, err = u.getProductCode(ctx, tx, d.Product.ID)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, u.ID, d.Product.ID, d.Qty, d.Code, d.Shelve.ID)
	if err != nil {
		return err
	}

	detailID, err := res.LastInsertId()
	u.ReceiveDetails[i].ID = uint64(detailID)
	u.ReceiveDetails[i].Code = d.Code

	inventory := new(Inventory)
	inventory.CompanyID = ctx.Value(api.Ctx("auth")).(User).Company.ID
	inventory.BranchID = ctx.Value(api.Ctx("auth")).(User).Branch.ID
	inventory.ShelveID = d.Shelve.ID
	inventory.ProductID = d.Product.ID
	inventory.ProductCode = d.Code
	inventory.TransactionID = u.ID
	inventory.Code = u.Code
	inventory.TransactionDate = u.Date
	inventory.Type = "GR"
	inventory.InOut = true
	inventory.Qty = 1
	return inventory.Create(ctx, tx)
}

func (u *Receive) updateDetail(ctx context.Context, tx *sql.Tx, d ReceiveDetail) error {
	const queryDetail = `
		UPDATE good_receiving_details 
		SET product_id = ?, 
			code = ?,
			shelve_id = ?
		WHERE id = ?
		AND good_receiving_id = ?
	`

	inventory := new(Inventory)
	inventory.ProductID = d.Product.ID
	inventory.ProductCode = d.Code
	inventory.TransactionID = u.ID
	inventory.Type = "GR"
	err := inventory.GetByComposit(ctx, tx)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, d.Product.ID, d.Code, d.Shelve.ID, d.ID, u.ID)
	if err != nil {
		return err
	}

	inventory.ProductID = d.Product.ID
	inventory.ProductCode = d.Code
	inventory.ShelveID = d.Shelve.ID

	return inventory.Update(ctx, tx)
}

func (u *Receive) removeDetail(ctx context.Context, tx *sql.Tx, e uint64) error {
	const queryDetail = `DELETE FROM  good_receiving_details WHERE id = ? AND good_receiving_id = ?`
	detail, err := u.getDetail(ctx, tx, e)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, e, u.ID)
	if err != nil {
		return err
	}

	inventory := new(Inventory)
	inventory.ProductID = detail.Product.ID
	inventory.ProductCode = detail.Code
	inventory.TransactionID = u.ID
	inventory.Type = "GR"
	inventory.CompanyID = ctx.Value(api.Ctx("auth")).(User).Company.ID
	inventory.BranchID = ctx.Value(api.Ctx("auth")).(User).Branch.ID
	return inventory.DeleteByComposit(ctx, tx)
}

func (u *Receive) getProductCode(ctx context.Context, tx *sql.Tx, productID uint64) (string, error) {
	var code string
	var err error
	var codeInt int

	prefix := time.Now().Format("200601")

	query := `
		SELECT good_receiving_details.code 
		FROM good_receiving_details 
		JOIN good_receivings ON good_receiving_details.good_receiving_id = good_receivings.id
		WHERE good_receivings.company_id = ? AND good_receiving_details.code LIKE ? AND good_receiving_details.product_id = ? 
		ORDER BY good_receiving_details.code DESC LIMIT 1
	`
	err = tx.QueryRowContext(ctx, query, ctx.Value(api.Ctx("auth")).(User).Company.ID, prefix+"%", productID).Scan(&code)

	if err != nil && err != sql.ErrNoRows {
		return code, err
	}

	if len(code) > 0 {
		runes := []rune(code)
		codeInt, err = strconv.Atoi(string(runes[6:]))
		if err != nil {
			return code, err
		}
	}

	return prefix + fmt.Sprintf("%014d", codeInt+1), nil
}
