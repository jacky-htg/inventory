package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/models"
	"github.com/jacky-htg/inventory/payloads/request"
	"github.com/jacky-htg/inventory/payloads/response"
	"github.com/julienschmidt/httprouter"
)

// Suppliers type for handling dependency injection
type Suppliers struct {
	Db  *sql.DB
	Log *log.Logger
}

// List of suppliers
func (u *Suppliers) List(w http.ResponseWriter, r *http.Request) {
	var supplier models.Supplier
	list, err := supplier.List(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("get supplier list : %v", err)
		api.ResponseError(w, err)
		return
	}

	var supplierResponse []response.SupplierResponse
	for _, r := range list {
		var res response.SupplierResponse
		res.Transform(&r)
		supplierResponse = append(supplierResponse, res)
	}

	api.ResponseOK(w, supplierResponse, http.StatusOK)
}

// Create new supplier
func (u *Suppliers) Create(w http.ResponseWriter, r *http.Request) {
	var supplierRequest request.NewSupplierRequest
	err := api.Decode(r, &supplierRequest)
	if err != nil {
		u.Log.Printf("Error : %v", err)
		api.ResponseError(w, err)
		return
	}

	supplier := supplierRequest.Transform()

	err = supplier.Create(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("create new supplier tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	var res response.SupplierResponse
	res.Transform(&supplier)
	api.ResponseOK(w, supplier, http.StatusCreated)
}

// View of supplier by id
func (u *Suppliers) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var supplier models.Supplier
	supplier.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	err = supplier.Get(ctx, tx)
	if err != nil {
		u.Log.Printf("Get supplier: %v", err)
		api.ResponseError(w, err)
		return
	}

	var res response.SupplierResponse
	res.Transform(&supplier)
	api.ResponseOK(w, res, http.StatusOK)
}

// Update supplier by id
func (u *Suppliers) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var supplierRequest request.SupplierRequest
	err = api.Decode(r, &supplierRequest)
	if err != nil {
		u.Log.Printf("Error : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx, err := u.Db.Begin()
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	var supplier models.Supplier
	supplier.ID = uint64(id)
	err = supplier.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Get supplier: %v", err)
		api.ResponseError(w, err)
		return
	}

	tx.Commit()

	supplierUpdate := supplierRequest.Transform(&supplier)

	err = supplierUpdate.Update(ctx, u.Db)
	if err != nil {
		u.Log.Printf("Update supplier : %v", err)
		api.ResponseError(w, err)
		return
	}

	var res response.SupplierResponse
	res.Transform(supplierUpdate)
	api.ResponseOK(w, res, http.StatusOK)
}

// Delete supplier by id
func (u *Suppliers) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx, err := u.Db.Begin()
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	var supplier models.Supplier
	supplier.ID = uint64(id)
	err = supplier.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Get supplier: %v", err)
		api.ResponseError(w, err)
		return
	}

	tx.Commit()

	err = supplier.Delete(ctx, u.Db)
	if err != nil {
		u.Log.Printf("Update supplier: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	api.ResponseOK(w, nil, http.StatusNoContent)
}
