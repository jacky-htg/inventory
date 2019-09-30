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

// Purchases : struct for set Purchases Dependency Injection
type Purchases struct {
	Db  *sql.DB
	Log *log.Logger
}

// List : http handler for returning list of purchases
func (u *Purchases) List(w http.ResponseWriter, r *http.Request) {
	var purchase models.Purchase
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	list, err := purchase.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting purchases list: %v", err))
		return
	}

	tx.Commit()

	var listResponse []*response.PurchaseListResponse
	for _, purchase := range list {
		var purchaseResponse response.PurchaseListResponse
		purchaseResponse.Transform(&purchase)
		listResponse = append(listResponse, &purchaseResponse)
	}

	api.ResponseOK(w, listResponse, http.StatusOK)
}

// View : http handler for retrieve purchase by id
func (u *Purchases) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting: %v", err))
		return
	}

	var purchase models.Purchase
	purchase.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = purchase.Get(ctx, tx)

	if err == sql.ErrNoRows {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrNotFound(err, ""))
		return
	}

	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get purchase: %v", err))
		return
	}

	tx.Commit()

	var response response.PurchaseResponse
	response.Transform(&purchase)
	api.ResponseOK(w, response, http.StatusOK)
}

// Create : http handler for create new purchase
func (u *Purchases) Create(w http.ResponseWriter, r *http.Request) {
	var purchaseRequest request.NewPurchaseRequest
	err := api.Decode(r, &purchaseRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("decode purchase: %v", err))
		return
	}

	purchase := purchaseRequest.Transform()
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("begin transaction: %v", err))
		return
	}

	err = purchase.Create(r.Context(), tx)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		tx.Rollback()
		api.ResponseError(w, fmt.Errorf("Create purchase: %v", err))
		return
	}

	tx.Commit()

	var response response.PurchaseResponse
	response.Transform(purchase)
	api.ResponseOK(w, response, http.StatusCreated)
}

// Update : http handler for update purchase by id
func (u *Purchases) Update(w http.ResponseWriter, r *http.Request) {
	// TODO : untuk dikerjakan jika modul return purchasing, dan good receiving sudah selesai
	// Edit purchase hanya boleh dilakukan jika :
	// 1. belum ada pembayaran
	// 2. belum ada return purchasing
	// 3. belum ada good receiving

	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var purchase models.Purchase
	purchase.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = purchase.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get purchase: %v", err))
		return
	}

	var purchaseRequest request.PurchaseRequest
	err = api.Decode(r, &purchaseRequest)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Decode purchase: %v", err))
		return
	}

	if purchaseRequest.ID <= 0 {
		purchaseRequest.ID = purchase.ID
	}
	purchaseUpdate := purchaseRequest.Transform(&purchase)
	err = purchaseUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Update purchase: %v", err))
		return
	}

	tx.Commit()

	var response response.PurchaseResponse
	response.Transform(purchaseUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}
