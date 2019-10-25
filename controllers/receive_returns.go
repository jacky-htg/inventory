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

// ReceiveReturns : struct for set ReceiveReturns Dependency Injection
type ReceiveReturns struct {
	Db  *sql.DB
	Log *log.Logger
}

// List : http handler for returning list of ReceiveReturns
func (u *ReceiveReturns) List(w http.ResponseWriter, r *http.Request) {
	var receiveReturn models.ReceiveReturn
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	list, err := receiveReturn.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting ReceiveReturns list: %v", err))
		return
	}

	tx.Commit()

	var listResponse []*response.ReceiveReturnListResponse
	for _, r := range list {
		var receiveReturnResponse response.ReceiveReturnListResponse
		receiveReturnResponse.Transform(&r)
		listResponse = append(listResponse, &receiveReturnResponse)
	}

	api.ResponseOK(w, listResponse, http.StatusOK)
}

// View : http handler for retrieve ReceiveReturn by id
func (u *ReceiveReturns) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting: %v", err))
		return
	}

	var receiveReturn models.ReceiveReturn
	receiveReturn.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = receiveReturn.Get(ctx, tx)

	if err == sql.ErrNoRows {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrNotFound(err, ""))
		return
	}

	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get ReceiveReturn: %v", err))
		return
	}

	tx.Commit()

	var response response.ReceiveReturnResponse
	response.Transform(&receiveReturn)
	api.ResponseOK(w, response, http.StatusOK)
}

// Create : http handler for create new ReceiveReturn
func (u *ReceiveReturns) Create(w http.ResponseWriter, r *http.Request) {
	var receiveReturnRequest request.NewReceiveReturnRequest
	err := api.Decode(r, &receiveReturnRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("decode Receive Return: %v", err))
		return
	}

	receiveReturn := receiveReturnRequest.Transform()
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("begin transaction: %v", err))
		return
	}

	err = receiveReturn.Create(r.Context(), tx)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		tx.Rollback()
		api.ResponseError(w, fmt.Errorf("Create Receive Return: %v", err))
		return
	}

	tx.Commit()

	var response response.ReceiveReturnResponse
	response.Transform(receiveReturn)
	api.ResponseOK(w, response, http.StatusCreated)
}

// Update : http handler for update ReceiveReturn by id
func (u *ReceiveReturns) Update(w http.ResponseWriter, r *http.Request) {
	// TODO :
	// Edit Receive Return hanya boleh dilakukan jika :
	// belum ada transaksi lain atas good receiving seperti delivery

	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var receiveReturn models.ReceiveReturn
	receiveReturn.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = receiveReturn.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get ReceiveReturn: %v", err))
		return
	}

	var receiveReturnRequest request.ReceiveReturnRequest
	err = api.Decode(r, &receiveReturnRequest)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Decode Receive Return: %v", err))
		return
	}

	if receiveReturnRequest.ID <= 0 {
		receiveReturnRequest.ID = receiveReturn.ID
	}
	receiveReturnUpdate := receiveReturnRequest.Transform(&receiveReturn)
	err = receiveReturnUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Update Receive Return: %v", err))
		return
	}

	tx.Commit()

	var response response.ReceiveReturnResponse
	response.Transform(receiveReturnUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}
