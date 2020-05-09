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

// SalesOrderReturns : struct for set SalesOrderReturns Dependency Injection
type SalesOrderReturns struct {
	Db  *sql.DB
	Log *log.Logger
}

// List : http handler for returning list of salesOrders
func (u *SalesOrderReturns) List(w http.ResponseWriter, r *http.Request) {
	var salesOrderReturn models.SalesOrderReturn
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	list, err := salesOrderReturn.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting salesOrder returns list: %v", err))
		return
	}

	tx.Commit()

	var listResponse []*response.SalesOrderReturnListResponse
	for _, salesOrderReturn := range list {
		var salesOrderReturnResponse response.SalesOrderReturnListResponse
		salesOrderReturnResponse.Transform(&salesOrderReturn)
		listResponse = append(listResponse, &salesOrderReturnResponse)
	}

	api.ResponseOK(w, listResponse, http.StatusOK)
}

// View : http handler for retrieve salesOrder return by id
func (u *SalesOrderReturns) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting: %v", err))
		return
	}

	var salesOrderReturn models.SalesOrderReturn
	salesOrderReturn.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = salesOrderReturn.Get(ctx, tx)

	if err == sql.ErrNoRows {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrNotFound(err, ""))
		return
	}

	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get salesOrder return: %v", err))
		return
	}

	tx.Commit()

	var response response.SalesOrderReturnResponse
	response.Transform(&salesOrderReturn)
	api.ResponseOK(w, response, http.StatusOK)
}

// Create : http handler for create new salesOrder return
func (u *SalesOrderReturns) Create(w http.ResponseWriter, r *http.Request) {
	var salesOrderReturnRequest request.NewSalesOrderReturnRequest
	err := api.Decode(r, &salesOrderReturnRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("decode salesOrder return: %v", err))
		return
	}

	salesOrderReturn := salesOrderReturnRequest.Transform()
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("begin transaction: %v", err))
		return
	}

	// todo check valid owner of SalesOrderID
	// todo check SalesOrderID isvalid open to return

	err = salesOrderReturn.Create(r.Context(), tx)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		tx.Rollback()
		api.ResponseError(w, fmt.Errorf("Create salesOrder return: %v", err))
		return
	}

	tx.Commit()

	var response response.SalesOrderReturnResponse
	response.Transform(salesOrderReturn)
	api.ResponseOK(w, response, http.StatusCreated)
}

// Update : http handler for update salesOrder return by id
func (u *SalesOrderReturns) Update(w http.ResponseWriter, r *http.Request) {
	// TODO : untuk dikerjakan jika modul return pembayaran sudah selesai
	// Edit salesOrder return hanya boleh dilakukan jika :
	// 1. belum ada pembayaran

	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var salesOrderReturn models.SalesOrderReturn
	salesOrderReturn.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = salesOrderReturn.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get salesOrder return: %v", err))
		return
	}

	var salesOrderReturnRequest request.SalesOrderReturnRequest
	err = api.Decode(r, &salesOrderReturnRequest)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Decode salesOrder return: %v", err))
		return
	}

	if salesOrderReturnRequest.ID <= 0 {
		salesOrderReturnRequest.ID = salesOrderReturn.ID
	}
	salesOrderReturnUpdate := salesOrderReturnRequest.Transform(&salesOrderReturn)
	err = salesOrderReturnUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Update salesOrder return: %v", err))
		return
	}

	tx.Commit()

	var response response.SalesOrderReturnResponse
	response.Transform(salesOrderReturnUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}
