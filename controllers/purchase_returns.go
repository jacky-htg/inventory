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

// PurchaseReturns : struct for set PurchaseReturns Dependency Injection
type PurchaseReturns struct {
	Db  *sql.DB
	Log *log.Logger
}

// List : http handler for returning list of purchases
func (u *PurchaseReturns) List(w http.ResponseWriter, r *http.Request) {
	var purchaseReturn models.PurchaseReturn
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	list, err := purchaseReturn.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting purchase returns list: %v", err))
		return
	}

	tx.Commit()

	var listResponse []*response.PurchaseReturnListResponse
	for _, purchaseReturn := range list {
		var purchaseReturnResponse response.PurchaseReturnListResponse
		purchaseReturnResponse.Transform(&purchaseReturn)
		listResponse = append(listResponse, &purchaseReturnResponse)
	}

	api.ResponseOK(w, listResponse, http.StatusOK)
}

// View : http handler for retrieve purchase return by id
func (u *PurchaseReturns) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting: %v", err))
		return
	}

	var purchaseReturn models.PurchaseReturn
	purchaseReturn.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = purchaseReturn.Get(ctx, tx)

	if err == sql.ErrNoRows {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrNotFound(err, ""))
		return
	}

	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get purchase return: %v", err))
		return
	}

	tx.Commit()

	var response response.PurchaseReturnResponse
	response.Transform(&purchaseReturn)
	api.ResponseOK(w, response, http.StatusOK)
}

// Create : http handler for create new purchase return
func (u *PurchaseReturns) Create(w http.ResponseWriter, r *http.Request) {
	var purchaseReturnRequest request.NewPurchaseReturnRequest
	err := api.Decode(r, &purchaseReturnRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("decode purchase return: %v", err))
		return
	}

	purchaseReturn := purchaseReturnRequest.Transform()
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("begin transaction: %v", err))
		return
	}

	// todo check valid owner of PurchaseID
	// todo check PurchaseID isvalid open to return

	err = purchaseReturn.Create(r.Context(), tx)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		tx.Rollback()
		api.ResponseError(w, fmt.Errorf("Create purchase return: %v", err))
		return
	}

	tx.Commit()

	var response response.PurchaseReturnResponse
	response.Transform(purchaseReturn)
	api.ResponseOK(w, response, http.StatusCreated)
}

// Update : http handler for update purchase return by id
func (u *PurchaseReturns) Update(w http.ResponseWriter, r *http.Request) {
	// TODO : untuk dikerjakan jika modul return pembayaran sudah selesai
	// Edit purchase return hanya boleh dilakukan jika :
	// 1. belum ada pembayaran

	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var purchaseReturn models.PurchaseReturn
	purchaseReturn.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = purchaseReturn.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get purchase return: %v", err))
		return
	}

	var purchaseReturnRequest request.PurchaseReturnRequest
	err = api.Decode(r, &purchaseReturnRequest)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Decode purchase return: %v", err))
		return
	}

	if purchaseReturnRequest.ID <= 0 {
		purchaseReturnRequest.ID = purchaseReturn.ID
	}
	purchaseReturnUpdate := purchaseReturnRequest.Transform(&purchaseReturn)
	err = purchaseReturnUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Update purchase return: %v", err))
		return
	}

	tx.Commit()

	var response response.PurchaseReturnResponse
	response.Transform(purchaseReturnUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}
