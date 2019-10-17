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

// Receives : struct for set Receives Dependency Injection
type Receives struct {
	Db  *sql.DB
	Log *log.Logger
}

// List : http handler for returning list of Receives
func (u *Receives) List(w http.ResponseWriter, r *http.Request) {
	var receive models.Receive
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	list, err := receive.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting Receives list: %v", err))
		return
	}

	tx.Commit()

	var listResponse []*response.ReceiveListResponse
	for _, r := range list {
		var receiveResponse response.ReceiveListResponse
		receiveResponse.Transform(&r)
		listResponse = append(listResponse, &receiveResponse)
	}

	api.ResponseOK(w, listResponse, http.StatusOK)
}

// View : http handler for retrieve Receive by id
func (u *Receives) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting: %v", err))
		return
	}

	var receive models.Receive
	receive.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = receive.Get(ctx, tx)

	if err == sql.ErrNoRows {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrNotFound(err, ""))
		return
	}

	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get Receive: %v", err))
		return
	}

	tx.Commit()

	var response response.ReceiveResponse
	response.Transform(&receive)
	api.ResponseOK(w, response, http.StatusOK)
}

// Create : http handler for create new Receive
func (u *Receives) Create(w http.ResponseWriter, r *http.Request) {
	var receiveRequest request.NewReceiveRequest
	err := api.Decode(r, &receiveRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("decode Receive: %v", err))
		return
	}

	receive := receiveRequest.Transform()
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("begin transaction: %v", err))
		return
	}

	err = receive.Create(r.Context(), tx)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		tx.Rollback()
		api.ResponseError(w, fmt.Errorf("Create Receive: %v", err))
		return
	}

	tx.Commit()

	var response response.ReceiveResponse
	response.Transform(receive)
	api.ResponseOK(w, response, http.StatusCreated)
}

// Update : http handler for update Receive by id
func (u *Receives) Update(w http.ResponseWriter, r *http.Request) {
	// TODO :
	// Edit Receive hanya boleh dilakukan jika :
	// belum ada transaksi lain atas good receiving seperti return receiving, mutations, atau sales/delivery

	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var receive models.Receive
	receive.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = receive.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get Receive: %v", err))
		return
	}

	var receiveRequest request.ReceiveRequest
	err = api.Decode(r, &receiveRequest)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Decode Receive: %v", err))
		return
	}

	if receiveRequest.ID <= 0 {
		receiveRequest.ID = receive.ID
	}
	receiveUpdate := receiveRequest.Transform(&receive)
	err = receiveUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Update Receive: %v", err))
		return
	}

	tx.Commit()

	var response response.ReceiveResponse
	response.Transform(receiveUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}
