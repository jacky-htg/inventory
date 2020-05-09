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

// SalesOrderReturn : struct of SalesOrderReturn
type SalesOrderReturn struct {
	ID                      uint64
	Code                    string
	Date                    time.Time
	Price                   float64
	Disc                    float64
	AdditionalDisc          float64
	Total                   float64
	SalesOrder              SalesOrder
	Company                 Company
	Branch                  Branch
	SalesOrderReturnDetails []SalesOrderReturnDetail
}

// SalesOrderReturnDetail struct
type SalesOrderReturnDetail struct {
	ID      uint64
	Product Product
	Price   float64
	Disc    float64
	Qty     uint
}

// List sales order returns
func (u *SalesOrderReturn) List(ctx context.Context, tx *sql.Tx) ([]SalesOrderReturn, error) {
	var list []SalesOrderReturn
	var err error

	query := `
	SELECT 	sales_order_returns.id, 
		sales_order_returns.code, 
		sales_order_returns.date,
		sales_orders.id,
		sales_orders.code,
		sales_orders.date,
		sales_orders.disc,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type,
		SUM(sales_order_return_details.price),
		SUM(sales_order_return_details.disc),
		sales_order_returns.disc
	FROM sales_order_returns
	JOIN companies ON sales_order_returns.company_id = companies.id
	JOIN sales_orders ON sales_order_returns.sales_order_id = sales_orders.id
	JOIN branches ON sales_orders.branch_id = branches.id
	JOIN sales_order_return_details ON sales_order_returns.id = sales_order_return_details.sales_order_return_id
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

	query += " GROUP BY sales_order_returns.id"

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var salesOrderReturn SalesOrderReturn
		err = rows.Scan(
			&salesOrderReturn.ID,
			&salesOrderReturn.Code,
			&salesOrderReturn.Date,
			&salesOrderReturn.SalesOrder.ID,
			&salesOrderReturn.SalesOrder.Code,
			&salesOrderReturn.SalesOrder.Date,
			&salesOrderReturn.SalesOrder.Disc,
			&salesOrderReturn.Company.ID,
			&salesOrderReturn.Company.Code,
			&salesOrderReturn.Company.Name,
			&salesOrderReturn.Company.Address,
			&salesOrderReturn.Branch.ID,
			&salesOrderReturn.Branch.Code,
			&salesOrderReturn.Branch.Name,
			&salesOrderReturn.Branch.Address,
			&salesOrderReturn.Branch.Type,
			&salesOrderReturn.Price,
			&salesOrderReturn.Disc,
			&salesOrderReturn.AdditionalDisc,
		)

		if err != nil {
			return list, err
		}

		salesOrderReturn.Total = salesOrderReturn.Price - salesOrderReturn.Disc - salesOrderReturn.AdditionalDisc
		salesOrderReturn.SalesOrder.Company = salesOrderReturn.Company
		salesOrderReturn.Branch.Company = salesOrderReturn.Company

		list = append(list, salesOrderReturn)
	}

	return list, rows.Err()
}

// Get salesOrder return by id
func (u *SalesOrderReturn) Get(ctx context.Context, tx *sql.Tx) error {
	query := `
	SELECT 	sales_order_returns.id, 
		sales_order_returns.code, 
		sales_order_returns.date,
		sales_orders.id,
		sales_orders.code,
		sales_orders.date,
		sales_orders.disc,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type,
		SUM(sales_order_return_details.price),
		SUM(sales_order_return_details.disc),
		JSON_ARRAYAGG(sales_order_return_details.id),
		JSON_ARRAYAGG(sales_order_return_details.price),
		JSON_ARRAYAGG(sales_order_return_details.disc),
		JSON_ARRAYAGG(sales_order_return_details.qty),
		JSON_ARRAYAGG(products.id),
		JSON_ARRAYAGG(products.code),
		JSON_ARRAYAGG(products.name),
		JSON_ARRAYAGG(products.sale_price),
		sales_order_returns.disc
	FROM sales_order_returns
	JOIN companies ON sales_order_returns.company_id = companies.id
	JOIN sales_orders ON sales_order_returns.sales_order_id = sales_orders.id
	JOIN branches ON sales_order_returns.branch_id = branches.id
	JOIN sales_order_return_details ON sales_order_returns.id = sales_order_return_details.sales_order_return_id
	JOIN products ON sales_order_return_details.product_id = products.id 
	WHERE sales_order_returns.id=? AND companies.id=?
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
	err := tx.QueryRowContext(ctx, query+" GROUP BY sales_order_returns.id", params...).Scan(
		&u.ID,
		&u.Code,
		&u.Date,
		&u.SalesOrder.ID,
		&u.SalesOrder.Code,
		&u.SalesOrder.Date,
		&u.SalesOrder.Disc,
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
			u.SalesOrderReturnDetails = append(u.SalesOrderReturnDetails, SalesOrderReturnDetail{
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

	u.SalesOrder.Company = u.Company
	u.Branch.Company = u.Company

	return nil
}

// Create new salesOrder return
func (u *SalesOrderReturn) Create(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Branch.ID <= 0 {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		INSERT INTO sales_order_returns (code, date, disc, sales_order_id, company_id, branch_id, created_by, updated_by, created, updated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	u.Code, err = u.getCode(ctx, tx)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, u.Code, u.Date, u.AdditionalDisc, u.SalesOrder.ID, userLogin.Company.ID, userLogin.Branch.ID, userLogin.ID, userLogin.ID)
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
	u.SalesOrder.Get(ctx, tx)

	for i, d := range u.SalesOrderReturnDetails {
		// TODO :
		// VALIDATE that detail is open
		// 1. Detail belum pernah direturn
		// 2. Detail belum pernah didelivery
		detailID, err := u.storeDetail(ctx, tx, d, u.SalesOrder.ID)
		if err != nil {
			return err
		}

		u.SalesOrderReturnDetails[i].ID = detailID
		u.Price += d.Price
		u.Disc += d.Disc
		u.SalesOrderReturnDetails[i].Product.Get(ctx, tx)
	}
	u.Total = u.Price - u.Disc - u.AdditionalDisc

	return nil
}

// Update salesOrder return
func (u *SalesOrderReturn) Update(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Company.ID != u.Company.ID || u.Branch.ID <= 0 || userLogin.Branch.ID != u.Branch.ID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		UPDATE sales_order_returns 
		SET date = ?, 
			disc = ?,
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

	_, err = stmt.ExecContext(ctx, u.Date, u.AdditionalDisc, u.SalesOrder.ID, userLogin.ID, u.ID, userLogin.Company.ID, userLogin.Branch.ID)
	if err != nil {
		return err
	}

	existingDetails, err := u.GetExistingDetails(ctx, tx)
	if err != nil {
		return err
	}

	for i, d := range u.SalesOrderReturnDetails {
		if d.ID <= 0 {
			detailID, err := u.storeDetail(ctx, tx, d, u.SalesOrder.ID)
			if err != nil {
				return err
			}
			u.SalesOrderReturnDetails[i].ID = detailID
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

// GetExistingDetails return array of existing sales_order_return_details id
func (u *SalesOrderReturn) GetExistingDetails(ctx context.Context, tx *sql.Tx) ([]uint64, error) {
	var list []uint64
	var err error

	rows, err := tx.QueryContext(ctx, "SELECT id FROM sales_order_return_details WHERE sales_order_return_id=?", u.ID)
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

func (u *SalesOrderReturn) getCode(ctx context.Context, tx *sql.Tx) (string, error) {
	var code, prefix string
	var codeInt int
	prefix = "SR" + time.Now().Format("200601")

	query := `SELECT code FROM sales_order_returns WHERE company_id = ? AND code LIKE ? ORDER BY code DESC LIMIT 1`
	err := tx.QueryRowContext(ctx, query, ctx.Value(api.Ctx("auth")).(User).Company.ID, prefix+"%").Scan(&code)

	if err != nil && err != sql.ErrNoRows {
		return code, err
	}

	if len(code) > 0 {
		runes := []rune(code)
		codeInt, err = strconv.Atoi(string(runes[8:]))
		if err != nil {
			return code, err
		}
	}

	return prefix + fmt.Sprintf("%05d", codeInt+1), nil
}

func (u *SalesOrderReturn) storeDetail(ctx context.Context, tx *sql.Tx, d SalesOrderReturnDetail, salesOrderID uint64) (uint64, error) {
	var id uint64

	// check :
	// 1. valid detail sales order return is only product in sales order detail list.
	// 2. Max qty of product detail sales order return = qty product detail sales order - qty existing return
	// 3. Jika modul delivery sudah selesai, tambahkan validasi max qty = qty sales order - qty return - qty delivery
	var qty uint
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	err := tx.QueryRowContext(ctx, `
		SELECT (MAX(sales_order_details.qty) - IFNULL(SUM(sales_order_return_details.qty), 0)) as qty
		FROM sales_order_details 
		JOIN sales_orders ON sales_order_details.sales_order_id = sales_orders.id
		LEFT JOIN sales_order_returns ON sales_orders.id = sales_order_returns.sales_order_id
		LEFT JOIN sales_order_return_details ON sales_order_returns.id = sales_order_return_details.sales_order_return_id AND sales_order_details.product_id = sales_order_return_details.product_id
		WHERE sales_order_details.sales_order_id=? AND sales_order_details.product_id=?
		AND sales_orders.company_id=? AND sales_orders.branch_id=?
		GROUP BY sales_order_details.product_id
	`, salesOrderID, d.Product.ID, userLogin.Company.ID, userLogin.Branch.ID).Scan(&qty)

	if err != nil {
		return id, err
	}

	if d.Qty > qty {
		return id, api.ErrBadRequest(errors.New("Invalid quantity"), "")
	}

	const queryDetail = `
		INSERT INTO sales_order_return_details (sales_order_return_id, product_id, price, disc, qty)
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

func (u *SalesOrderReturn) updateDetail(ctx context.Context, tx *sql.Tx, d SalesOrderReturnDetail) error {
	const queryDetail = `
		UPDATE sales_order_return_details 
		SET product_id = ?, 
			price = ?,
			disc = ?,
			qty = ?
		WHERE id = ?
		AND sales_order_return_id = ?
	`
	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, d.Product.ID, d.Price, d.Disc, d.Qty, d.ID, u.ID)
	return err
}

func (u *SalesOrderReturn) removeDetail(ctx context.Context, tx *sql.Tx, e uint64) error {
	const queryDetail = `DELETE FROM  sales_order_return_details WHERE id = ? AND sales_order_return_id = ?`
	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, e, u.ID)
	return err
}
