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

// SalesOrder : struct of SalesOrder
type SalesOrder struct {
	ID                uint64
	Code              string
	Date              time.Time
	Price             float64
	Disc              float64
	AdditionalDisc    float64
	Total             float64
	Salesman          Salesman
	Customer          Customer
	Company           Company
	Branch            Branch
	SalesOrderDetails []SalesOrderDetail
}

// SalesOrderDetail struct
type SalesOrderDetail struct {
	ID      uint64
	Product Product
	Price   float64
	Disc    float64
	Qty     uint
}

// List sales orders
func (u *SalesOrder) List(ctx context.Context, tx *sql.Tx) ([]SalesOrder, error) {
	var list []SalesOrder
	var err error

	query := `
	SELECT 	sales_orders.id, 
		sales_orders.code, 
		sales_orders.date,
		salesmen.id,
		salesmen.code,
		salesmen.name,
		salesmen.address,
		customers.id,
		customers.name,
		customers.email,
		customers.address,
		customers.hp,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type,
		SUM(sales_order_details.price),
		SUM(sales_order_details.disc),
		sales_orders.disc
	FROM sales_orders
	JOIN customers ON sales_orders.customer_id = customers.id
	JOIN companies ON sales_orders.company_id = companies.id
	JOIN salesmen ON sales_orders.salesman_id = salesmen.id
	JOIN branches ON sales_orders.branch_id = branches.id
	JOIN sales_order_details ON sales_orders.id = sales_order_details.sales_order_id
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

	query += " GROUP BY sales_orders.id"

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var salesOrder SalesOrder
		err = rows.Scan(
			&salesOrder.ID,
			&salesOrder.Code,
			&salesOrder.Date,
			&salesOrder.Salesman.ID,
			&salesOrder.Salesman.Code,
			&salesOrder.Salesman.Name,
			&salesOrder.Salesman.Address,
			&salesOrder.Customer.ID,
			&salesOrder.Customer.Name,
			&salesOrder.Customer.Email,
			&salesOrder.Customer.Address,
			&salesOrder.Customer.Hp,
			&salesOrder.Company.ID,
			&salesOrder.Company.Code,
			&salesOrder.Company.Name,
			&salesOrder.Company.Address,
			&salesOrder.Branch.ID,
			&salesOrder.Branch.Code,
			&salesOrder.Branch.Name,
			&salesOrder.Branch.Address,
			&salesOrder.Branch.Type,
			&salesOrder.Price,
			&salesOrder.Disc,
			&salesOrder.AdditionalDisc,
		)

		if err != nil {
			return list, err
		}

		salesOrder.Total = salesOrder.Price - salesOrder.Disc - salesOrder.AdditionalDisc
		salesOrder.Salesman.Company = salesOrder.Company
		salesOrder.Branch.Company = salesOrder.Company

		list = append(list, salesOrder)
	}

	return list, rows.Err()
}

// Get sales order by id
func (u *SalesOrder) Get(ctx context.Context, tx *sql.Tx) error {
	query := `
	SELECT 	sales_orders.id, 
	sales_orders.code, 
	sales_orders.date,
		salesmen.id,
		salesmen.code,
		salesmen.name,
		salesmen.address,
		customers.id,
		customers.name,
		customers.email,
		customers.address,
		customers.hp,		
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type,
		SUM(sales_order_details.price),
		SUM(sales_order_details.disc),
		JSON_ARRAYAGG(sales_order_details.id),
		JSON_ARRAYAGG(sales_order_details.price),
		JSON_ARRAYAGG(sales_order_details.disc),
		JSON_ARRAYAGG(sales_order_details.qty),
		JSON_ARRAYAGG(products.id),
		JSON_ARRAYAGG(products.code),
		JSON_ARRAYAGG(products.name),
		JSON_ARRAYAGG(products.sale_price),
		sales_orders.disc
	FROM sales_orders
	JOIN companies ON sales_orders.company_id = companies.id
	JOIN salesmen ON sales_orders.salesman_id = salesmen.id
	JOIN customers ON sales_orders.customer_id = customers.id
	JOIN branches ON sales_orders.branch_id = branches.id
	JOIN sales_order_details ON sales_orders.id = sales_order_details.sales_order_id
	JOIN products ON sales_order_details.product_id = products.id 
	WHERE sales_orders.id=? AND companies.id=?
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

	var detailID, detailPrice, detailDisc, detailQty, productID, productCode, productName, productPrice string
	err := tx.QueryRowContext(ctx, query+" GROUP BY sales_orders.id", params...).Scan(
		&u.ID,
		&u.Code,
		&u.Date,
		&u.Salesman.ID,
		&u.Salesman.Code,
		&u.Salesman.Name,
		&u.Salesman.Address,
		&u.Customer.ID,
		&u.Customer.Name,
		&u.Customer.Email,
		&u.Customer.Address,
		&u.Customer.Hp,
		&u.Company.ID,
		&u.Company.Code,
		&u.Company.Name,
		&u.Company.Address,
		&u.Branch.ID,
		&u.Branch.Code,
		&u.Branch.Name,
		&u.Branch.Address,
		&u.Branch.Type,
		&u.Price,
		&u.Disc,
		&detailID,
		&detailPrice,
		&detailDisc,
		&detailQty,
		&productID,
		&productCode,
		&productName,
		&productPrice,
		&u.AdditionalDisc,
	)

	if err != nil {
		return err
	}

	u.Total = u.Price - u.Disc - u.AdditionalDisc

	if len(detailID) > 0 {
		var detailIDs []uint64
		err = json.Unmarshal([]byte(detailID), &detailIDs)
		if err != nil {
			return err
		}

		var detailPrices []float64
		err = json.Unmarshal([]byte(detailPrice), &detailPrices)
		if err != nil {
			return err
		}

		var detailDiscs []float64
		err = json.Unmarshal([]byte(detailDisc), &detailDiscs)
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
			u.SalesOrderDetails = append(u.SalesOrderDetails, SalesOrderDetail{
				ID:    uint64(v),
				Price: detailPrices[i],
				Disc:  detailDiscs[i],
				Qty:   detailQtys[i],
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

	u.Salesman.Company = u.Company
	u.Branch.Company = u.Company

	return nil
}

// Create new sales order
func (u *SalesOrder) Create(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Branch.ID <= 0 {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		INSERT INTO sales_orders (code, date, disc, salesman_id, customer_id, company_id, branch_id, created_by, updated_by, created, updated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	u.Code, err = api.GetCode(ctx, tx, "SO", "sales_orders", userLogin.Company.ID)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, u.Code, u.Date, u.AdditionalDisc, u.Salesman.ID, u.Customer.ID, userLogin.Company.ID, userLogin.Branch.ID, userLogin.ID, userLogin.ID)
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
	u.Salesman.Get(ctx, tx)

	for i, d := range u.SalesOrderDetails {
		detailID, err := u.storeDetail(ctx, tx, d)
		if err != nil {
			return err
		}
		u.SalesOrderDetails[i].ID = detailID
		u.Price += d.Price
		u.Disc += d.Disc
		u.SalesOrderDetails[i].Product.Get(ctx, tx)
	}
	u.Total = u.Price - u.Disc - u.AdditionalDisc

	return nil
}

// Update sales order
func (u *SalesOrder) Update(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Company.ID != u.Company.ID || u.Branch.ID <= 0 || userLogin.Branch.ID != u.Branch.ID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		UPDATE sales_orders 
		SET date = ?, 
			disc = ?,
			salesman_id = ?,
			customer_id = ?, 
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

	_, err = stmt.ExecContext(ctx, u.Date, u.AdditionalDisc, u.Salesman.ID, u.Customer.ID, userLogin.ID, u.ID, userLogin.Company.ID, userLogin.Branch.ID)
	if err != nil {
		return err
	}

	existingDetails, err := u.GetExistingDetails(ctx, tx)
	if err != nil {
		return err
	}

	for i, d := range u.SalesOrderDetails {
		if d.ID <= 0 {
			detailID, err := u.storeDetail(ctx, tx, d)
			if err != nil {
				return err
			}
			u.SalesOrderDetails[i].ID = detailID
		} else {
			err = u.updateDetail(ctx, tx, d)
			if err != nil {
				return err
			}

			var arrUint64 array.ArrUint64
			existingDetails = arrUint64.Remove(existingDetails, d.ID)
		}

		u.Price += d.Price
		u.Disc += d.Disc
		u.SalesOrderDetails[i].Product.Get(ctx, tx)
	}

	u.Total = u.Price - u.Disc - u.AdditionalDisc

	for _, e := range existingDetails {
		err = u.removeDetail(ctx, tx, e)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetExistingDetails return array of existing sales_order_details id
func (u *SalesOrder) GetExistingDetails(ctx context.Context, tx *sql.Tx) ([]uint64, error) {
	var list []uint64
	var err error

	rows, err := tx.QueryContext(ctx, "SELECT id FROM sales_order_details WHERE sales_order_id=?", u.ID)
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

func (u *SalesOrder) storeDetail(ctx context.Context, tx *sql.Tx, d SalesOrderDetail) (uint64, error) {
	var id uint64
	const queryDetail = `
		INSERT INTO sales_order_details (sales_order_id, product_id, price, disc, qty)
		VALUES (?, ?, ?, ?, ?)
	`
	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return id, err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, u.ID, d.Product.ID, d.Price, d.Disc, d.Qty)
	if err != nil {
		return id, err
	}

	detailID, err := res.LastInsertId()
	if err != nil {
		return id, err
	}

	return uint64(detailID), nil
}

func (u *SalesOrder) updateDetail(ctx context.Context, tx *sql.Tx, d SalesOrderDetail) error {
	const queryDetail = `
		UPDATE sales_order_details 
		SET product_id = ?, 
			price = ?,
			disc = ?,
			qty = ?
		WHERE id = ?
		AND sales_order_id = ?
	`
	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, d.Product.ID, d.Price, d.Disc, d.Qty, d.ID, u.ID)
	return err
}

func (u *SalesOrder) removeDetail(ctx context.Context, tx *sql.Tx, e uint64) error {
	const queryDetail = `DELETE FROM  sales_order_details WHERE id = ? AND sales_order_id = ?`
	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, e, u.ID)
	return err
}
