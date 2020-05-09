package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/models"
	"github.com/jacky-htg/inventory/payloads/request"
	"github.com/jacky-htg/inventory/payloads/response"
	"github.com/julienschmidt/httprouter"
)

// SalesOrders : struct for set SalesOrders Dependency Injection
type SalesOrders struct {
	Db  *sql.DB
	Log *log.Logger
}

// List : http handler for returning list of SalesOrders
func (u *SalesOrders) List(w http.ResponseWriter, r *http.Request) {
	var salesOrder models.SalesOrder
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	list, err := salesOrder.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting SalesOrders list: %v", err))
		return
	}

	tx.Commit()

	var listResponse []*response.SalesOrderListResponse
	for _, so := range list {
		var salesOrderResponse response.SalesOrderListResponse
		salesOrderResponse.Transform(&so)
		listResponse = append(listResponse, &salesOrderResponse)
	}

	api.ResponseOK(w, listResponse, http.StatusOK)
}

// View : http handler for retrieve SalesOrder by id
func (u *SalesOrders) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting: %v", err))
		return
	}

	var salesOrder models.SalesOrder
	salesOrder.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = salesOrder.Get(ctx, tx)

	if err == sql.ErrNoRows {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrNotFound(err, ""))
		return
	}

	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get sales order: %v", err))
		return
	}

	tx.Commit()

	var res response.SalesOrderResponse
	res.Transform(&salesOrder)
	api.ResponseOK(w, res, http.StatusOK)
}

// Create : http handler for create new SalesOrder
func (u *SalesOrders) Create(w http.ResponseWriter, r *http.Request) {
	var salesOrderRequest request.NewSalesOrderRequest
	err := api.Decode(r, &salesOrderRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("decode SalesOrder: %v", err))
		return
	}

	salesOrder := salesOrderRequest.Transform()
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("begin transaction: %v", err))
		return
	}

	err = salesOrder.Create(r.Context(), tx)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		tx.Rollback()
		api.ResponseError(w, fmt.Errorf("Create SalesOrder: %v", err))
		return
	}

	tx.Commit()

	var res response.SalesOrderResponse
	res.Transform(salesOrder)
	api.ResponseOK(w, res, http.StatusCreated)
}

// Update : http handler for update SalesOrder by id
func (u *SalesOrders) Update(w http.ResponseWriter, r *http.Request) {
	// TODO : untuk dikerjakan jika modul return so, dan delivery order sudah selesai
	// Edit sales order hanya boleh dilakukan jika :
	// 1. belum ada pembayaran
	// 2. belum ada return so
	// 3. belum ada do

	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var salesOrder models.SalesOrder
	salesOrder.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = salesOrder.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get SalesOrder: %v", err))
		return
	}

	var salesOrderRequest request.SalesOrderRequest
	err = api.Decode(r, &salesOrderRequest)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Decode SalesOrder: %v", err))
		return
	}

	if salesOrderRequest.ID <= 0 {
		salesOrderRequest.ID = salesOrder.ID
	}
	salesOrderUpdate := salesOrderRequest.Transform(&salesOrder)
	err = salesOrderUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Update SalesOrder: %v", err))
		return
	}

	tx.Commit()

	var res response.SalesOrderResponse
	res.Transform(salesOrderUpdate)
	api.ResponseOK(w, res, http.StatusOK)
}
