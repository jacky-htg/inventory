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

// DeliveryReturn : struct of DeliveryReturn
type DeliveryReturn struct {
	ID                    uint64
	Code                  string
	Date                  time.Time
	Remark                string
	Delivery              Delivery
	Company               Company
	Branch                Branch
	DeliveryReturnDetails []DeliveryReturnDetail
}

// DeliveryReturnDetail struct
type DeliveryReturnDetail struct {
	ID      uint64
	Product Product
	Code    string
	Qty     uint
}

// List Delivery returns
func (u *DeliveryReturn) List(ctx context.Context, tx *sql.Tx) ([]DeliveryReturn, error) {
	var list []DeliveryReturn
	var err error

	query := `
	SELECT 	delivery_returns.id, 
		delivery_returns.code, 
		delivery_returns.date,
		delivery_returns.remark,
		deliveries.id,
		deliveries.code,
		deliveries.date,
		deliveries.remark,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type
	FROM delivery_returns
	JOIN deliveries ON delivery_returns.delivery_id = deliveries.id
	JOIN companies ON delivery_returns.company_id = companies.id
	JOIN branches ON delivery_returns.branch_id = branches.id
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
		var deliveryReturn DeliveryReturn
		err = rows.Scan(
			&deliveryReturn.ID,
			&deliveryReturn.Code,
			&deliveryReturn.Date,
			&deliveryReturn.Remark,
			&deliveryReturn.Delivery.ID,
			&deliveryReturn.Delivery.Code,
			&deliveryReturn.Delivery.Date,
			&deliveryReturn.Delivery.Remark,
			&deliveryReturn.Company.ID,
			&deliveryReturn.Company.Code,
			&deliveryReturn.Company.Name,
			&deliveryReturn.Company.Address,
			&deliveryReturn.Branch.ID,
			&deliveryReturn.Branch.Code,
			&deliveryReturn.Branch.Name,
			&deliveryReturn.Branch.Address,
			&deliveryReturn.Branch.Type,
		)

		if err != nil {
			return list, err
		}

		deliveryReturn.Delivery.Company = deliveryReturn.Company
		deliveryReturn.Branch.Company = deliveryReturn.Company

		list = append(list, deliveryReturn)
	}

	return list, rows.Err()
}

// Get Delivery return by id
func (u *DeliveryReturn) Get(ctx context.Context, tx *sql.Tx) error {
	query := `
	SELECT 	delivery_returns.id, 
		delivery_returns.code, 
		delivery_returns.date,
		delivery_returns.remark,
		deliveries.id,
		deliveries.code,
		deliveries.date,
		deliveries.remark,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type,
		JSON_ARRAYAGG(delivery_return_details.id),
		JSON_ARRAYAGG(delivery_return_details.code),
		JSON_ARRAYAGG(delivery_return_details.qty),
		JSON_ARRAYAGG(products.id),
		JSON_ARRAYAGG(products.code),
		JSON_ARRAYAGG(products.name),
		JSON_ARRAYAGG(products.sale_price)
	FROM delivery_returns
	JOIN delivery_return_details ON delivery_returns.id = delivery_return_details.delivery_return_id
	JOIN deliveries ON delivery_returns.delivery_id = deliveries.id
	JOIN companies ON delivery_returns.company_id = companies.id
	JOIN branches ON delivery_returns.branch_id = branches.id
	JOIN products ON delivery_return_details.product_id = products.id 
	WHERE delivery_returns.id=? AND companies.id=?
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
	err := tx.QueryRowContext(ctx, query+" GROUP BY delivery_returns.id", params...).Scan(
		&u.ID,
		&u.Code,
		&u.Date,
		&u.Remark,
		&u.Delivery.ID,
		&u.Delivery.Code,
		&u.Delivery.Date,
		&u.Delivery.Remark,
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
			u.DeliveryReturnDetails = append(u.DeliveryReturnDetails, DeliveryReturnDetail{
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

	u.Delivery.Company = u.Company
	u.Branch.Company = u.Company

	return nil
}

// Create new delivery return
func (u *DeliveryReturn) Create(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Branch.ID <= 0 {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		INSERT INTO delivery_returns (code, date, remark, delivery_id, company_id, branch_id, created_by, updated_by, created, updated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	u.Code, err = api.GetCode(ctx, tx, "DR", "delivery_returns", userLogin.Company.ID)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, u.Code, u.Date, u.Remark, u.Delivery.ID, userLogin.Company.ID, userLogin.Branch.ID, userLogin.ID, userLogin.ID)
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
	u.Delivery.Get(ctx, tx)

	for i, d := range u.DeliveryReturnDetails {
		err = u.storeDetail(ctx, tx, d, u.Delivery.ID, i)
		if err != nil {
			return err
		}

		u.DeliveryReturnDetails[i].Product.Get(ctx, tx)
	}

	return nil
}

// Update Delivery return
func (u *DeliveryReturn) Update(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Company.ID != u.Company.ID || u.Branch.ID <= 0 || userLogin.Branch.ID != u.Branch.ID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		UPDATE delivery_returns 
		SET date = ?, 
			remark = ?,
			delivery_id = ?, 
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

	_, err = stmt.ExecContext(ctx, u.Date, u.Remark, u.Delivery.ID, userLogin.ID, u.ID, userLogin.Company.ID, userLogin.Branch.ID)
	if err != nil {
		return err
	}

	existingDetails, err := u.GetExistingDetails(ctx, tx)
	if err != nil {
		return err
	}

	for i, d := range u.DeliveryReturnDetails {
		if d.ID <= 0 {
			err = u.storeDetail(ctx, tx, d, u.Delivery.ID, i)
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

// GetExistingDetails return array of existing Delivery_return_details id
func (u *DeliveryReturn) GetExistingDetails(ctx context.Context, tx *sql.Tx) ([]uint64, error) {
	var list []uint64
	var err error

	rows, err := tx.QueryContext(ctx, "SELECT id FROM delivery_return_details WHERE delivery_return_id=?", u.ID)
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

func (u *DeliveryReturn) getDetail(ctx context.Context, tx *sql.Tx, e uint64) (*DeliveryReturnDetail, error) {
	detail := new(DeliveryReturnDetail)
	err := tx.QueryRowContext(ctx,
		`SELECT id, product_id, qty, code FROM delivery_return_details WHERE id = ? AND delivery_return_id = ?`,
		e, u.ID).Scan(&detail.ID, &detail.Product.ID, &detail.Qty, &detail.Code)

	return detail, err
}

func (u *DeliveryReturn) storeDetail(ctx context.Context, tx *sql.Tx, d DeliveryReturnDetail, deliveryID uint64, i int) error {
	// check :
	// 1. valid detail Delivery return is only product in Delivery detail list.
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	var code string
	var shelveID uint64
	err := tx.QueryRowContext(ctx, `
		SELECT delivery_details.code, delivery_details.shelve_id 
		FROM delivery_details
		JOIN deliveries ON delivery_details.delivery_id = deliveries.id AND deliveries.company_id = ? AND deliveries.branch_id = ?
		LEFT JOIN delivery_returns ON delivery_returns.company_id = deliveries.company_id AND delivery_returns.branch_id = deliveries.branch_id AND delivery_returns.delivery_id = deliveries.id
		LEFT JOIN delivery_return_details ON delivery_returns.id = delivery_return_details.delivery_return_id AND delivery_return_details.product_id = delivery_details.product_id AND delivery_return_details.code = delivery_details.code
		WHERE delivery_details.product_id = ? AND delivery_details.code = ? AND delivery_return_details.code is null 
	`, userLogin.Company.ID, userLogin.Branch.ID, d.Product.ID, d.Code).Scan(&code, &shelveID)

	if err != nil {
		return err
	}

	const queryDetail = `
		INSERT INTO delivery_return_details (delivery_return_id, product_id, code, qty)
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

	u.DeliveryReturnDetails[i].ID = uint64(detailID)

	inventory := new(Inventory)
	inventory.CompanyID = userLogin.Company.ID
	inventory.BranchID = userLogin.Branch.ID
	inventory.ProductID = d.Product.ID
	inventory.ProductCode = d.Code
	inventory.TransactionID = u.ID
	inventory.Code = u.Code
	inventory.TransactionDate = u.Date
	inventory.Type = "DR"
	inventory.InOut = true
	inventory.Qty = 1
	inventory.ShelveID = shelveID
	return inventory.Create(ctx, tx)
}

/*func (u *DeliveryReturn) updateDetail(ctx context.Context, tx *sql.Tx, d DeliveryReturnDetail) error {
	const queryDetail = `
		UPDATE delivery_return_details
		SET product_id = ?,
			code = ?,
			qty = ?
		WHERE id = ?
		AND delivery_return_id = ?
	`

	inventory := new(Inventory)
	inventory.ProductID = d.Product.ID
	inventory.ProductCode = d.Code
	inventory.TransactionID = u.ID
	inventory.Type = "DR"
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

func (u *DeliveryReturn) removeDetail(ctx context.Context, tx *sql.Tx, e uint64) error {
	const queryDetail = `DELETE FROM  delivery_return_details WHERE id = ? AND delivery_return_id = ?`

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
	inventory.Type = "DR"
	inventory.CompanyID = ctx.Value(api.Ctx("auth")).(User).Company.ID
	inventory.BranchID = ctx.Value(api.Ctx("auth")).(User).Branch.ID
	return inventory.DeleteByComposit(ctx, tx)
}
