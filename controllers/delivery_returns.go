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

// DeliveryReturns : struct for set DeliveryReturns Dependency Injection
type DeliveryReturns struct {
	Db  *sql.DB
	Log *log.Logger
}

// List : http handler for returning list of DeliveryReturns
func (u *DeliveryReturns) List(w http.ResponseWriter, r *http.Request) {
	var deliveryReturn models.DeliveryReturn
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	list, err := deliveryReturn.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting DeliveryReturns list: %v", err))
		return
	}

	tx.Commit()

	var listResponse []*response.DeliveryReturnListResponse
	for _, r := range list {
		var deliveryReturnResponse response.DeliveryReturnListResponse
		deliveryReturnResponse.Transform(&r)
		listResponse = append(listResponse, &deliveryReturnResponse)
	}

	api.ResponseOK(w, listResponse, http.StatusOK)
}

// View : http handler for retrieve DeliveryReturn by id
func (u *DeliveryReturns) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting: %v", err))
		return
	}

	var deliveryReturn models.DeliveryReturn
	deliveryReturn.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = deliveryReturn.Get(ctx, tx)

	if err == sql.ErrNoRows {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrNotFound(err, ""))
		return
	}

	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get DeliveryReturn: %v", err))
		return
	}

	tx.Commit()

	var response response.DeliveryReturnResponse
	response.Transform(&deliveryReturn)
	api.ResponseOK(w, response, http.StatusOK)
}

// Create : http handler for create new DeliveryReturn
func (u *DeliveryReturns) Create(w http.ResponseWriter, r *http.Request) {
	var deliveryReturnRequest request.NewDeliveryReturnRequest
	err := api.Decode(r, &deliveryReturnRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("decode Delivery Return: %v", err))
		return
	}

	deliveryReturn := deliveryReturnRequest.Transform()
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("begin transaction: %v", err))
		return
	}

	err = deliveryReturn.Create(r.Context(), tx)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		tx.Rollback()
		api.ResponseError(w, fmt.Errorf("Create Delivery Return: %v", err))
		return
	}

	tx.Commit()

	var response response.DeliveryReturnResponse
	response.Transform(deliveryReturn)
	api.ResponseOK(w, response, http.StatusCreated)
}

// Update : http handler for update DeliveryReturn by id
func (u *DeliveryReturns) Update(w http.ResponseWriter, r *http.Request) {
	// TODO :
	// Edit Delivery Return hanya boleh dilakukan jika :

	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var deliveryReturn models.DeliveryReturn
	deliveryReturn.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = deliveryReturn.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get DeliveryReturn: %v", err))
		return
	}

	var deliveryReturnRequest request.DeliveryReturnRequest
	err = api.Decode(r, &deliveryReturnRequest)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Decode Delivery Return: %v", err))
		return
	}

	if deliveryReturnRequest.ID <= 0 {
		deliveryReturnRequest.ID = deliveryReturn.ID
	}
	deliveryReturnUpdate := deliveryReturnRequest.Transform(&deliveryReturn)
	err = deliveryReturnUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Update Delivery Return: %v", err))
		return
	}

	tx.Commit()

	var response response.DeliveryReturnResponse
	response.Transform(deliveryReturnUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}
