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

// Purchase : struct of Purchase
type Purchase struct {
	ID              uint64
	Code            string
	Date            time.Time
	Price           float64
	Disc            float64
	AdditionalDisc  float64
	Total           float64
	Supplier        Supplier
	Company         Company
	Branch          Branch
	PurchaseDetails []PurchaseDetail
}

// PurchaseDetail struct
type PurchaseDetail struct {
	ID      uint64
	Product Product
	Price   float64
	Disc    float64
	Qty     uint
}

// List purchases
func (u *Purchase) List(ctx context.Context, tx *sql.Tx) ([]Purchase, error) {
	var list []Purchase
	var err error

	query := `
	SELECT 	purchases.id, 
		purchases.code, 
		purchases.date,
		suppliers.id,
		suppliers.code,
		suppliers.name,
		suppliers.address,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type,
		SUM(purchase_details.price),
		SUM(purchase_details.disc),
		purchases.disc
	FROM purchases
	JOIN companies ON purchases.company_id = companies.id
	JOIN suppliers ON purchases.supplier_id = suppliers.id
	JOIN branches ON purchases.branch_id = branches.id
	JOIN purchase_details ON purchases.id = purchase_details.purchase_id
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

	query += " GROUP BY purchases.id"

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var purchase Purchase
		err = rows.Scan(
			&purchase.ID,
			&purchase.Code,
			&purchase.Date,
			&purchase.Supplier.ID,
			&purchase.Supplier.Code,
			&purchase.Supplier.Name,
			&purchase.Supplier.Address,
			&purchase.Company.ID,
			&purchase.Company.Code,
			&purchase.Company.Name,
			&purchase.Company.Address,
			&purchase.Branch.ID,
			&purchase.Branch.Code,
			&purchase.Branch.Name,
			&purchase.Branch.Address,
			&purchase.Branch.Type,
			&purchase.Price,
			&purchase.Disc,
			&purchase.AdditionalDisc,
		)

		if err != nil {
			return list, err
		}

		purchase.Total = purchase.Price - purchase.Disc - purchase.AdditionalDisc
		purchase.Supplier.Company = purchase.Company
		purchase.Branch.Company = purchase.Company

		list = append(list, purchase)
	}

	return list, rows.Err()
}

// Get purchase by id
func (u *Purchase) Get(ctx context.Context, tx *sql.Tx) error {
	query := `
	SELECT 	purchases.id, 
		purchases.code, 
		purchases.date,
		suppliers.id,
		suppliers.code,
		suppliers.name,
		suppliers.address,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type,
		SUM(purchase_details.price),
		SUM(purchase_details.disc),
		JSON_ARRAYAGG(purchase_details.id),
		JSON_ARRAYAGG(purchase_details.price),
		JSON_ARRAYAGG(purchase_details.disc),
		JSON_ARRAYAGG(purchase_details.qty),
		JSON_ARRAYAGG(products.id),
		JSON_ARRAYAGG(products.code),
		JSON_ARRAYAGG(products.name),
		JSON_ARRAYAGG(products.sale_price),
		purchases.disc
	FROM purchases
	JOIN companies ON purchases.company_id = companies.id
	JOIN suppliers ON purchases.supplier_id = suppliers.id
	JOIN branches ON purchases.branch_id = branches.id
	JOIN purchase_details ON purchases.id = purchase_details.purchase_id
	JOIN products ON purchase_details.product_id = products.id 
	WHERE purchases.id=? AND companies.id=?
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
	err := tx.QueryRowContext(ctx, query+" GROUP BY purchases.id", params...).Scan(
		&u.ID,
		&u.Code,
		&u.Date,
		&u.Supplier.ID,
		&u.Supplier.Code,
		&u.Supplier.Name,
		&u.Supplier.Address,
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
			u.PurchaseDetails = append(u.PurchaseDetails, PurchaseDetail{
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

	u.Supplier.Company = u.Company
	u.Branch.Company = u.Company

	return nil
}

// Create new purchase
func (u *Purchase) Create(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Branch.ID <= 0 {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		INSERT INTO purchases (code, date, disc, supplier_id, company_id, branch_id, created_by, updated_by, created, updated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	u.Code, err = api.GetCode(ctx, tx, "POR", "purchases", userLogin.Company.ID)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, u.Code, u.Date, u.AdditionalDisc, u.Supplier.ID, userLogin.Company.ID, userLogin.Branch.ID, userLogin.ID, userLogin.ID)
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
	u.Supplier.Get(ctx, tx)

	for i, d := range u.PurchaseDetails {
		detailID, err := u.storeDetail(ctx, tx, d)
		if err != nil {
			return err
		}
		u.PurchaseDetails[i].ID = detailID
		u.Price += d.Price
		u.Disc += d.Disc
		u.PurchaseDetails[i].Product.Get(ctx, tx)
	}
	u.Total = u.Price - u.Disc - u.AdditionalDisc

	return nil
}

// Update purchase
func (u *Purchase) Update(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Company.ID != u.Company.ID || u.Branch.ID <= 0 || userLogin.Branch.ID != u.Branch.ID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		UPDATE purchases 
		SET date = ?, 
			disc = ?,
			supplier_id = ?, 
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

	_, err = stmt.ExecContext(ctx, u.Date, u.AdditionalDisc, u.Supplier.ID, userLogin.ID, u.ID, userLogin.Company.ID, userLogin.Branch.ID)
	if err != nil {
		return err
	}

	existingDetails, err := u.GetExistingDetails(ctx, tx)
	if err != nil {
		return err
	}

	for i, d := range u.PurchaseDetails {
		if d.ID <= 0 {
			detailID, err := u.storeDetail(ctx, tx, d)
			if err != nil {
				return err
			}
			u.PurchaseDetails[i].ID = detailID
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
		u.PurchaseDetails[i].Product.Get(ctx, tx)
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

// GetExistingDetails return array of existing purchase_details id
func (u *Purchase) GetExistingDetails(ctx context.Context, tx *sql.Tx) ([]uint64, error) {
	var list []uint64
	var err error

	rows, err := tx.QueryContext(ctx, "SELECT id FROM purchase_details WHERE purchase_id=?", u.ID)
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

func (u *Purchase) storeDetail(ctx context.Context, tx *sql.Tx, d PurchaseDetail) (uint64, error) {
	var id uint64
	const queryDetail = `
		INSERT INTO purchase_details (purchase_id, product_id, price, disc, qty)
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

func (u *Purchase) updateDetail(ctx context.Context, tx *sql.Tx, d PurchaseDetail) error {
	const queryDetail = `
		UPDATE purchase_details 
		SET product_id = ?, 
			price = ?,
			disc = ?,
			qty = ?
		WHERE id = ?
		AND purchase_id = ?
	`
	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, d.Product.ID, d.Price, d.Disc, d.Qty, d.ID, u.ID)
	return err
}

func (u *Purchase) removeDetail(ctx context.Context, tx *sql.Tx, e uint64) error {
	const queryDetail = `DELETE FROM  purchase_details WHERE id = ? AND purchase_id = ?`
	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, e, u.ID)
	return err
}
