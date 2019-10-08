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

// PurchaseReturn : struct of PurchaseReturn
type PurchaseReturn struct {
	ID                    uint64
	Code                  string
	Date                  time.Time
	Price                 float64
	Disc                  float64
	AdditionalDisc        float64
	Total                 float64
	Purchase              Purchase
	Company               Company
	Branch                Branch
	PurchaseReturnDetails []PurchaseReturnDetail
}

// PurchaseReturnDetail struct
type PurchaseReturnDetail struct {
	ID      uint64
	Product Product
	Price   float64
	Disc    float64
	Qty     uint
}

// List purchase returns
func (u *PurchaseReturn) List(ctx context.Context, tx *sql.Tx) ([]PurchaseReturn, error) {
	var list []PurchaseReturn
	var err error

	query := `
	SELECT 	purchase_returns.id, 
		purchase_returns.code, 
		purchase_returns.date,
		purchases.id,
		purchases.code,
		purchases.date,
		purchases.disc,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type,
		SUM(purchase_return_details.price),
		SUM(purchase_return_details.disc),
		purchase_returns.disc
	FROM purchase_returns
	JOIN companies ON purchase_returns.company_id = companies.id
	JOIN purchases ON purchase_returns.purchase_id = purchases.id
	JOIN branches ON purchases.branch_id = branches.id
	JOIN purchase_return_details ON purchase_returns.id = purchase_return_details.purchase_return_id
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

	query += " GROUP BY purchase_returns.id"

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var purchaseReturn PurchaseReturn
		err = rows.Scan(
			&purchaseReturn.ID,
			&purchaseReturn.Code,
			&purchaseReturn.Date,
			&purchaseReturn.Purchase.ID,
			&purchaseReturn.Purchase.Code,
			&purchaseReturn.Purchase.Date,
			&purchaseReturn.Purchase.Disc,
			&purchaseReturn.Company.ID,
			&purchaseReturn.Company.Code,
			&purchaseReturn.Company.Name,
			&purchaseReturn.Company.Address,
			&purchaseReturn.Branch.ID,
			&purchaseReturn.Branch.Code,
			&purchaseReturn.Branch.Name,
			&purchaseReturn.Branch.Address,
			&purchaseReturn.Branch.Type,
			&purchaseReturn.Price,
			&purchaseReturn.Disc,
			&purchaseReturn.AdditionalDisc,
		)

		if err != nil {
			return list, err
		}

		purchaseReturn.Total = purchaseReturn.Price - purchaseReturn.Disc - purchaseReturn.AdditionalDisc
		purchaseReturn.Purchase.Company = purchaseReturn.Company
		purchaseReturn.Branch.Company = purchaseReturn.Company

		list = append(list, purchaseReturn)
	}

	return list, rows.Err()
}

// Get purchase return by id
func (u *PurchaseReturn) Get(ctx context.Context, tx *sql.Tx) error {
	query := `
	SELECT 	purchase_returns.id, 
		purchase_returns.code, 
		purchase_returns.date,
		purchases.id,
		purchases.code,
		purchases.date,
		purchases.disc,
		companies.id, 
		companies.code, 
		companies.name,
		companies.address,
		branches.id,
		branches.code,
		branches.name,
		branches.address,
		branches.type,
		SUM(purchase_return_details.price),
		SUM(purchase_return_details.disc),
		JSON_ARRAYAGG(purchase_return_details.id),
		JSON_ARRAYAGG(purchase_return_details.price),
		JSON_ARRAYAGG(purchase_return_details.disc),
		JSON_ARRAYAGG(purchase_return_details.qty),
		JSON_ARRAYAGG(products.id),
		JSON_ARRAYAGG(products.code),
		JSON_ARRAYAGG(products.name),
		JSON_ARRAYAGG(products.sale_price),
		purchase_returns.disc
	FROM purchase_returns
	JOIN companies ON purchase_returns.company_id = companies.id
	JOIN purchases ON purchase_returns.purchase_id = purchases.id
	JOIN branches ON purchase_returns.branch_id = branches.id
	JOIN purchase_return_details ON purchase_returns.id = purchase_return_details.purchase_return_id
	JOIN products ON purchase_return_details.product_id = products.id 
	WHERE purchase_returns.id=? AND companies.id=?
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
	err := tx.QueryRowContext(ctx, query+" GROUP BY purchase_returns.id", params...).Scan(
		&u.ID,
		&u.Code,
		&u.Date,
		&u.Purchase.ID,
		&u.Purchase.Code,
		&u.Purchase.Date,
		&u.Purchase.Disc,
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
			u.PurchaseReturnDetails = append(u.PurchaseReturnDetails, PurchaseReturnDetail{
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

	u.Purchase.Company = u.Company
	u.Branch.Company = u.Company

	return nil
}

// Create new purchase return
func (u *PurchaseReturn) Create(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Branch.ID <= 0 {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		INSERT INTO purchase_returns (code, date, disc, purchase_id, company_id, branch_id, created_by, updated_by, created, updated)
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

	res, err := stmt.ExecContext(ctx, u.Code, u.Date, u.AdditionalDisc, u.Purchase.ID, userLogin.Company.ID, userLogin.Branch.ID, userLogin.ID, userLogin.ID)
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
	u.Purchase.Get(ctx, tx)

	for i, d := range u.PurchaseReturnDetails {
		// TODO :
		// VALIDATE that detail is open
		// 1. Detail belum pernah direturn
		// 2. Detail belum pernah direceiving
		detailID, err := u.storeDetail(ctx, tx, d, u.Purchase.ID)
		if err != nil {
			return err
		}

		u.PurchaseReturnDetails[i].ID = detailID
		u.Price += d.Price
		u.Disc += d.Disc
		u.PurchaseReturnDetails[i].Product.Get(ctx, tx)
	}
	u.Total = u.Price - u.Disc - u.AdditionalDisc

	return nil
}

// Update purchase return
func (u *PurchaseReturn) Update(ctx context.Context, tx *sql.Tx) error {
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	if userLogin.Company.ID != u.Company.ID || u.Branch.ID <= 0 || userLogin.Branch.ID != u.Branch.ID {
		return api.ErrForbidden(errors.New("Forbidden data owner"), "")
	}

	const query = `
		UPDATE purchase_returns 
		SET date = ?, 
			disc = ?,
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

	_, err = stmt.ExecContext(ctx, u.Date, u.AdditionalDisc, u.Purchase.ID, userLogin.ID, u.ID, userLogin.Company.ID, userLogin.Branch.ID)
	if err != nil {
		return err
	}

	existingDetails, err := u.GetExistingDetails(ctx, tx)
	if err != nil {
		return err
	}

	for i, d := range u.PurchaseReturnDetails {
		if d.ID <= 0 {
			detailID, err := u.storeDetail(ctx, tx, d, u.Purchase.ID)
			if err != nil {
				return err
			}
			u.PurchaseReturnDetails[i].ID = detailID
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

// GetExistingDetails return array of existing purchase_return_details id
func (u *PurchaseReturn) GetExistingDetails(ctx context.Context, tx *sql.Tx) ([]uint64, error) {
	var list []uint64
	var err error

	rows, err := tx.QueryContext(ctx, "SELECT id FROM purchase_return_details WHERE purchase_return_id=?", u.ID)
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

func (u *PurchaseReturn) getCode(ctx context.Context, tx *sql.Tx) (string, error) {
	var code, prefix string
	var codeInt int
	prefix = "PR" + time.Now().Format("200601")

	query := `SELECT code FROM purchase_returns WHERE company_id = ? AND code LIKE ? ORDER BY code DESC LIMIT 1`
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

func (u *PurchaseReturn) storeDetail(ctx context.Context, tx *sql.Tx, d PurchaseReturnDetail, purchaseID uint64) (uint64, error) {
	var id uint64

	// check :
	// 1. valid detail purchase return is only product in purchase detail list.
	// 2. Max qty of product detail purchase return = qty product detail purchase - qty existing return
	// 3. Jika modul receiving sudah selesai, tambahkan validasi max qty = qty purchase - qty return - qty receiving
	var qty uint
	userLogin := ctx.Value(api.Ctx("auth")).(User)
	err := tx.QueryRowContext(ctx, `
		SELECT (MAX(purchase_details.qty) - IFNULL(SUM(purchase_return_details.qty), 0)) as qty
		FROM purchase_details 
		JOIN purchases ON purchase_details.purchase_id = purchases.id
		LEFT JOIN purchase_returns ON purchases.id = purchase_returns.purchase_id
		LEFT JOIN purchase_return_details ON purchase_returns.id = purchase_return_details.purchase_return_id AND purchase_details.product_id = purchase_return_details.product_id
		WHERE purchase_details.purchase_id=? AND purchase_details.product_id=?
		AND purchases.company_id=? AND purchases.branch_id=?
		GROUP BY purchase_details.product_id
	`, purchaseID, d.Product.ID, userLogin.Company.ID, userLogin.Branch.ID).Scan(&qty)

	if err != nil {
		return id, err
	}

	if d.Qty > qty {
		return id, api.ErrBadRequest(errors.New("Invalid quantity"), "")
	}

	const queryDetail = `
		INSERT INTO purchase_return_details (purchase_return_id, product_id, price, disc, qty)
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

func (u *PurchaseReturn) updateDetail(ctx context.Context, tx *sql.Tx, d PurchaseReturnDetail) error {
	const queryDetail = `
		UPDATE purchase_return_details 
		SET product_id = ?, 
			price = ?,
			disc = ?,
			qty = ?
		WHERE id = ?
		AND purchase_return_id = ?
	`
	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, d.Product.ID, d.Price, d.Disc, d.Qty, d.ID, u.ID)
	return err
}

func (u *PurchaseReturn) removeDetail(ctx context.Context, tx *sql.Tx, e uint64) error {
	const queryDetail = `DELETE FROM  purchase_return_details WHERE id = ? AND purchase_return_id = ?`
	stmt, err := tx.PrepareContext(ctx, queryDetail)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, e, u.ID)
	return err
}
