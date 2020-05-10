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

// Delivery : struct of Delivery
type Delivery struct {
	ID              uint64
	Code            string
	Date            time.Time
	Remark          string
	SalesOrder      SalesOrder
	Company         Company
	Branch          Branch
	DeliveryDetails []DeliveryDetail
}

// DeliveryDetail struct
type DeliveryDetail struct {
	ID      uint64
	Product Product
	Qty     uint
	Code    string
	Shelve  Shelve
}

// List Deliveries
func (u *Delivery) List(ctx context.Context, tx *sql.Tx) ([]Delivery, error) {
	var list []Delivery
	var err error

	query := `
	SELECT 	deliveries.id, 
		deliveries.code, 
		deliveries.date,
		sales_orders.id,
		sales_orders.code,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type
	FROM deliveries
	JOIN companies ON deliveries.company_id = companies.id
	JOIN sales_orders ON deliveries.sales_order_id = sales_orders.id AND deliveries.company_id = sales_orders.company_id AND deliveries.branch_id = sales_orders.branch_id 
	JOIN branches ON deliveries.branch_id = branches.id
	JOIN delivery_details ON deliveries.id = delivery_details.delivery_id
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

	query += " GROUP BY deliveries.id"

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var Delivery Delivery
		err = rows.Scan(
			&Delivery.ID,
			&Delivery.Code,
			&Delivery.Date,
			&Delivery.SalesOrder.ID,
			&Delivery.SalesOrder.Code,
			&Delivery.Company.ID,
			&Delivery.Company.Code,
			&Delivery.Company.Name,
			&Delivery.Company.Address,
			&Delivery.Branch.ID,
			&Delivery.Branch.Code,
			&Delivery.Branch.Name,
			&Delivery.Branch.Address,
			&Delivery.Branch.Type,
		)

		if err != nil {
			return list, err
		}

		list = append(list, Delivery)
	}

	return list, rows.Err()
}

// Get Delivery by id
func (u *Delivery) Get(ctx context.Context, tx *sql.Tx) error {
	query := `
	SELECT 	deliveries.id, 
		deliveries.code, 
		deliveries.date,
		sales_orders.id,
		sales_orders.code,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type,
		JSON_ARRAYAGG(delivery_details.id),
		JSON_ARRAYAGG(delivery_details.code),
		JSON_ARRAYAGG(delivery_details.shelve_id),
		JSON_ARRAYAGG(delivery_details.qty),
		JSON_ARRAYAGG(products.id),
		JSON_ARRAYAGG(products.code),
		JSON_ARRAYAGG(products.name),
		JSON_ARRAYAGG(products.sale_price)
	FROM deliveries
	JOIN companies ON deliveries.company_id = companies.id
	JOIN sales_orders ON deliveries.sales_order_id = sales_orders.id AND deliveries.company_id = sales_orders.company_id AND deliveries.branch_id = sales_orders.branch_id
	JOIN branches ON deliveries.branch_id = branches.id
	JOIN delivery_details ON deliveries.id = delivery_details.delivery_id
	JOIN products ON delivery_details.product_id = products.id 
	WHERE deliveries.id=? AND companies.id=?
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
	err := tx.QueryRowContext(ctx, query+" GROUP BY deliveries.id", params...).Scan(
		&u.ID,
		&u.Code,
		&u.Date,
		&u.SalesOrder.ID,
		&u.SalesOrder.Code,
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
			u.DeliveryDetails = append(u.DeliveryDetails, DeliveryDetail{
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

// Create new Delivery
func (u *Delivery) Create(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Branch.ID <= 0 {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	err := u.SalesOrder.Get(ctx, tx)
	if err != nil {
		return err
	}

	if u.SalesOrder.Branch.ID != userLogin.Branch.ID || u.SalesOrder.Company.ID != userLogin.Company.ID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		INSERT INTO deliveries (code, date, remark, sales_order_id, company_id, branch_id, created_by, updated_by, created, updated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	u.Code, err = api.GetCode(ctx, tx, "DO", "deliveries", userLogin.Company.ID)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, u.Code, u.Date, u.Remark, u.SalesOrder.ID, userLogin.Company.ID, userLogin.Branch.ID, userLogin.ID, userLogin.ID)
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

	for i, d := range u.DeliveryDetails {
		err := u.storeDetail(ctx, tx, d, i)
		if err != nil {
			return err
		}
		u.DeliveryDetails[i].Product.Get(ctx, tx)
	}

	return nil
}

// Update Delivery
func (u *Delivery) Update(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Company.ID != u.Company.ID || u.Branch.ID <= 0 || userLogin.Branch.ID != u.Branch.ID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	err := u.SalesOrder.Get(ctx, tx)
	if err != nil {
		return err
	}

	if u.SalesOrder.Branch.ID != userLogin.Branch.ID || u.SalesOrder.Company.ID != userLogin.Company.ID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		UPDATE deliveries 
		SET date = ?, 
			remark = ?,
			sales_order_id = ?, 
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

	_, err = stmt.ExecContext(ctx, u.Date, u.Remark, u.SalesOrder.ID, userLogin.ID, u.ID, userLogin.Company.ID, userLogin.Branch.ID)
	if err != nil {
		return err
	}

	existingDetails, err := u.GetExistingDetails(ctx, tx)
	if err != nil {
		return err
	}

	for i, d := range u.DeliveryDetails {
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

		u.DeliveryDetails[i].Product.Get(ctx, tx)
	}

	for _, e := range existingDetails {
		err = u.removeDetail(ctx, tx, e)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetExistingDetails return array of existing delivery_details id
func (u *Delivery) GetExistingDetails(ctx context.Context, tx *sql.Tx) ([]uint64, error) {
	var list []uint64
	var err error

	rows, err := tx.QueryContext(ctx, "SELECT id FROM delivery_details WHERE delivery_id=?", u.ID)
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

func (u *Delivery) getDetail(ctx context.Context, tx *sql.Tx, e uint64) (*DeliveryDetail, error) {
	detail := new(DeliveryDetail)
	err := tx.QueryRowContext(ctx,
		`SELECT id, product_id, qty, code, shelve_id FROM delivery_details WHERE id = ? AND delivery_id = ?`,
		e, u.ID).Scan(&detail.ID, &detail.Product.ID, &detail.Qty, &detail.Code, &detail.Shelve.ID)

	return detail, err
}

func (u *Delivery) storeDetail(ctx context.Context, tx *sql.Tx, d DeliveryDetail, i int) error {
	var err error

	const queryDetail = `
		INSERT INTO delivery_details (delivery_id, product_id, qty, code, shelve_id)
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
	u.DeliveryDetails[i].ID = uint64(detailID)
	u.DeliveryDetails[i].Code = d.Code

	inventory := new(Inventory)
	inventory.CompanyID = ctx.Value(api.Ctx("auth")).(User).Company.ID
	inventory.BranchID = ctx.Value(api.Ctx("auth")).(User).Branch.ID
	inventory.ShelveID = d.Shelve.ID
	inventory.ProductID = d.Product.ID
	inventory.ProductCode = d.Code
	inventory.TransactionID = u.ID
	inventory.Code = u.Code
	inventory.TransactionDate = u.Date
	inventory.Type = "DO"
	inventory.InOut = false
	inventory.Qty = 1
	return inventory.Create(ctx, tx)
}

func (u *Delivery) updateDetail(ctx context.Context, tx *sql.Tx, d DeliveryDetail) error {
	const queryDetail = `
		UPDATE delivery_details 
		SET product_id = ?, 
			code = ?,
			shelve_id = ?
		WHERE id = ?
		AND delivery_id = ?
	`

	inventory := new(Inventory)
	inventory.ProductID = d.Product.ID
	inventory.ProductCode = d.Code
	inventory.TransactionID = u.ID
	inventory.Type = "DO"
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

func (u *Delivery) removeDetail(ctx context.Context, tx *sql.Tx, e uint64) error {
	const queryDetail = `DELETE FROM  delivery_details WHERE id = ? AND delivery_id = ?`
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
	inventory.Type = "DO"
	inventory.CompanyID = ctx.Value(api.Ctx("auth")).(User).Company.ID
	inventory.BranchID = ctx.Value(api.Ctx("auth")).(User).Branch.ID
	return inventory.DeleteByComposit(ctx, tx)
}

func (u *Delivery) getProductCode(ctx context.Context, tx *sql.Tx, productID uint64) (string, error) {
	var code string
	var err error
	var codeInt int

	prefix := time.Now().Format("200601")

	query := `
		SELECT delivery_details.code 
		FROM delivery_details 
		JOIN deliveries ON delivery_details.delivery_id = deliveries.id
		WHERE deliveries.company_id = ? AND delivery_details.code LIKE ? AND delivery_details.product_id = ? 
		ORDER BY delivery_details.code DESC LIMIT 1
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
