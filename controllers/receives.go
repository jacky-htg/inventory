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
	var Receive models.Receive
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	list, err := Receive.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting Receives list: %v", err))
		return
	}

	tx.Commit()

	var listResponse []*response.ReceiveListResponse
	for _, Receive := range list {
		var ReceiveResponse response.ReceiveListResponse
		ReceiveResponse.Transform(&Receive)
		listResponse = append(listResponse, &ReceiveResponse)
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

	var Receive models.Receive
	Receive.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = Receive.Get(ctx, tx)

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
	response.Transform(&Receive)
	api.ResponseOK(w, response, http.StatusOK)
}

// Create : http handler for create new Receive
func (u *Receives) Create(w http.ResponseWriter, r *http.Request) {
	var ReceiveRequest request.NewReceiveRequest
	err := api.Decode(r, &ReceiveRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("decode Receive: %v", err))
		return
	}

	Receive := ReceiveRequest.Transform()
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("begin transaction: %v", err))
		return
	}

	err = Receive.Create(r.Context(), tx)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		tx.Rollback()
		api.ResponseError(w, fmt.Errorf("Create Receive: %v", err))
		return
	}

	tx.Commit()

	var response response.ReceiveResponse
	response.Transform(Receive)
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

	var Receive models.Receive
	Receive.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = Receive.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get Receive: %v", err))
		return
	}

	var ReceiveRequest request.ReceiveRequest
	err = api.Decode(r, &ReceiveRequest)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Decode Receive: %v", err))
		return
	}

	if ReceiveRequest.ID <= 0 {
		ReceiveRequest.ID = Receive.ID
	}
	ReceiveUpdate := ReceiveRequest.Transform(&Receive)
	err = ReceiveUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Update Receive: %v", err))
		return
	}

	tx.Commit()

	var response response.ReceiveResponse
	response.Transform(ReceiveUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}
