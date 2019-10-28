package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/libraries/array"
)

// ReceiveReturn : struct of ReceiveReturn
type ReceiveReturn struct {
	ID                   uint64
	Code                 string
	Date                 time.Time
	Remark               string
	Receive              Receive
	Company              Company
	Branch               Branch
	ReceiveReturnDetails []ReceiveReturnDetail
}

// ReceiveReturnDetail struct
type ReceiveReturnDetail struct {
	ID      uint64
	Product Product
	Code    string
	Qty     uint
}

// List Receive returns
func (u *ReceiveReturn) List(ctx context.Context, tx *sql.Tx) ([]ReceiveReturn, error) {
	var list []ReceiveReturn
	var err error

	query := `
	SELECT 	receiving_returns.id, 
		receiving_returns.code, 
		receiving_returns.date,
		receiving_returns.remark,
		good_receivings.id,
		good_receivings.code,
		good_receivings.date,
		good_receivings.remark,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type
	FROM receiving_returns
	JOIN good_receivings ON receiving_returns.good_receiving_id = good_receivings.id
	JOIN companies ON receiving_returns.company_id = companies.id
	JOIN branches ON receiving_returns.branch_id = branches.id
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

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var receiveReturn ReceiveReturn
		err = rows.Scan(
			&receiveReturn.ID,
			&receiveReturn.Code,
			&receiveReturn.Date,
			&receiveReturn.Remark,
			&receiveReturn.Receive.ID,
			&receiveReturn.Receive.Code,
			&receiveReturn.Receive.Date,
			&receiveReturn.Receive.Remark,
			&receiveReturn.Company.ID,
			&receiveReturn.Company.Code,
			&receiveReturn.Company.Name,
			&receiveReturn.Company.Address,
			&receiveReturn.Branch.ID,
			&receiveReturn.Branch.Code,
			&receiveReturn.Branch.Name,
			&receiveReturn.Branch.Address,
			&receiveReturn.Branch.Type,
		)

		if err != nil {
			return list, err
		}

		receiveReturn.Receive.Company = receiveReturn.Company
		receiveReturn.Branch.Company = receiveReturn.Company

		list = append(list, receiveReturn)
	}

	return list, rows.Err()
}

// Get Receive return by id
func (u *ReceiveReturn) Get(ctx context.Context, tx *sql.Tx) error {
	query := `
	SELECT 	receiving_returns.id, 
		receiving_returns.code, 
		receiving_returns.date,
		receiving_returns.remark,
		good_receivings.id,
		good_receivings.code,
		good_receivings.date,
		good_receivings.remark,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type,
		JSON_ARRAYAGG(receiving_return_details.id),
		JSON_ARRAYAGG(receiving_return_details.code),
		JSON_ARRAYAGG(receiving_return_details.qty),
		JSON_ARRAYAGG(products.id),
		JSON_ARRAYAGG(products.code),
		JSON_ARRAYAGG(products.name),
		JSON_ARRAYAGG(products.sale_price)
	FROM receiving_returns
	JOIN receiving_return_details ON receiving_returns.id = receiving_return_details.receiving_return_id
	JOIN good_receivings ON receiving_returns.good_receiving_id = good_receivings.id
	JOIN companies ON receiving_returns.company_id = companies.id
	JOIN branches ON receiving_returns.branch_id = branches.id
	JOIN products ON receiving_return_details.product_id = products.id 
	WHERE receiving_returns.id=? AND companies.id=?
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

	var detailID, detailCode, detailQty, productID, productCode, productName, productPrice string
	err := tx.QueryRowContext(ctx, query+" GROUP BY receiving_returns.id", params...).Scan(
		&u.ID,
		&u.Code,
		&u.Date,
		&u.Remark,
		&u.Receive.ID,
		&u.Receive.Code,
		&u.Receive.Date,
		&u.Receive.Remark,
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
			u.ReceiveReturnDetails = append(u.ReceiveReturnDetails, ReceiveReturnDetail{
				ID:   uint64(v),
				Code: detailCodes[i],
				Qty:  detailQtys[i],
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

	u.Receive.Company = u.Company
	u.Branch.Company = u.Company

	return nil
}

// Create new receive return
func (u *ReceiveReturn) Create(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Branch.ID <= 0 {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		INSERT INTO receiving_returns (code, date, remark, good_receiving_id, company_id, branch_id, created_by, updated_by, created, updated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	u.Code, err = api.GetCode(ctx, tx, "RR", "receiving_returns", userLogin.Company.ID)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, u.Code, u.Date, u.Remark, u.Receive.ID, userLogin.Company.ID, userLogin.Branch.ID, userLogin.ID, userLogin.ID)
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
	u.Receive.Get(ctx, tx)

	for i, d := range u.ReceiveReturnDetails {
		err = u.storeDetail(ctx, tx, d, u.Receive.ID, i)
		if err != nil {
			return err
		}

		u.ReceiveReturnDetails[i].Product.Get(ctx, tx)
	}

	return nil
}

// Update Receive return
func (u *ReceiveReturn) Update(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Company.ID != u.Company.ID || u.Branch.ID <= 0 || userLogin.Branch.ID != u.Branch.ID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		UPDATE receiving_returns 
		SET date = ?, 
			remark = ?,
			good_receiving_id = ?, 
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

	_, err = stmt.ExecContext(ctx, u.Date, u.Remark, u.Receive.ID, userLogin.ID, u.ID, userLogin.Company.ID, userLogin.Branch.ID)
	if err != nil {
		return err
	}

	existingDetails, err := u.GetExistingDetails(ctx, tx)
	if err != nil {
		return err
	}

	for i, d := range u.ReceiveReturnDetails {
		if d.ID <= 0 {
			err = u.storeDetail(ctx, tx, d, u.Receive.ID, i)
			if err != nil {
				return err
			}
		} else {
			/*err = u.updateDetail(ctx, tx, d)
			if err != nil {
				return err
			}*/

			var arrUint64 array.ArrUint64
			existingDetails = arrUint64.Remove(existingDetails, d.ID)
		}

	}

	for _, e := range existingDetails {
		err = u.removeDetail(ctx, tx, e)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetExistingDetails return array of existing Receive_return_details id
func (u *ReceiveReturn) GetExistingDetails(ctx context.Context, tx *sql.Tx) ([]uint64, error) {
	var list []uint64
	var err error

	rows, err := tx.QueryContext(ctx, "SELECT id FROM receiving_return_details WHERE receiving_return_id=?", u.ID)
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

func (u *ReceiveReturn) getDetail(ctx context.Context, tx *sql.Tx, e uint64) (*ReceiveReturnDetail, error) {
	detail := new(ReceiveReturnDetail)
	err := tx.QueryRowContext(ctx,
		`SELECT id, product_id, qty, code FROM receiving_return_details WHERE id = ? AND receiving_return_id = ?`,
		e, u.ID).Scan(&detail.ID, &detail.Product.ID, &detail.Qty, &detail.Code)

	return detail, err
}

func (u *ReceiveReturn) storeDetail(ctx context.Context, tx *sql.Tx, d ReceiveReturnDetail, receiveID uint64, i int) error {
	// check :
	// 1. valid detail Receive return is only product in Receive detail list.
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	var code string
	var shelveID uint64
	err := tx.QueryRowContext(ctx, `
		SELECT good_receiving_details.code, good_receiving_details.shelve_id 
		FROM good_receiving_details
		JOIN good_receivings ON good_receiving_details.good_receiving_id = good_receivings.id AND good_receivings.company_id = ? AND good_receivings.branch_id = ?
		LEFT JOIN receiving_returns ON receiving_returns.company_id = good_receivings.company_id AND receiving_returns.branch_id = good_receivings.branch_id AND receiving_returns.good_receiving_id = good_receivings.id
		LEFT JOIN receiving_return_details ON receiving_returns.id = receiving_return_details.receiving_return_id AND receiving_return_details.product_id = good_receiving_details.product_id AND receiving_return_details.code = good_receiving_details.code
		WHERE good_receiving_details.product_id = ? AND good_receiving_details.code = ? AND receiving_return_details.code is null 
	`, userLogin.Company.ID, userLogin.Branch.ID, d.Product.ID, d.Code).Scan(&code, &shelveID)

	if err != nil {
		return err
	}

	const queryDetail = `
		INSERT INTO receiving_return_details (receiving_return_id, product_id, code, qty)
		VALUES (?, ?, ?, ?)
	`
	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, u.ID, d.Product.ID, d.Code, d.Qty)
	if err != nil {
		return err
	}

	detailID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ReceiveReturnDetails[i].ID = uint64(detailID)

	inventory := new(Inventory)
	inventory.CompanyID = userLogin.Company.ID
	inventory.BranchID = userLogin.Branch.ID
	inventory.ProductID = d.Product.ID
	inventory.ProductCode = d.Code
	inventory.TransactionID = u.ID
	inventory.Code = u.Code
	inventory.TransactionDate = u.Date
	inventory.Type = "RR"
	inventory.InOut = false
	inventory.Qty = 1
	inventory.ShelveID = shelveID
	return inventory.Create(ctx, tx)
}

/*func (u *ReceiveReturn) updateDetail(ctx context.Context, tx *sql.Tx, d ReceiveReturnDetail) error {
	const queryDetail = `
		UPDATE receiving_return_details
		SET product_id = ?,
			code = ?,
			qty = ?
		WHERE id = ?
		AND receiving_return_id = ?
	`

	inventory := new(Inventory)
	inventory.ProductID = d.Product.ID
	inventory.ProductCode = d.Code
	inventory.TransactionID = u.ID
	inventory.Type = "RR"
	err := inventory.GetByComposit(ctx, tx)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, d.Product.ID, d.Code, d.Qty, d.ID, u.ID)
	if err != nil {
		return err
	}

	inventory.ProductID = d.Product.ID
	inventory.ProductCode = d.Code
	return inventory.Update(ctx, tx)
}*/

func (u *ReceiveReturn) removeDetail(ctx context.Context, tx *sql.Tx, e uint64) error {
	const queryDetail = `DELETE FROM  receiving_return_details WHERE id = ? AND receiving_return_id = ?`

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
	inventory.Type = "RR"
	inventory.CompanyID = ctx.Value(api.Ctx("auth")).(User).Company.ID
	inventory.BranchID = ctx.Value(api.Ctx("auth")).(User).Branch.ID
	return inventory.DeleteByComposit(ctx, tx)
}
