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

// Deliveries : struct for set Deliveries Dependency Injection
type Deliveries struct {
	Db  *sql.DB
	Log *log.Logger
}

// List : http handler for returning list of Deliveries
func (u *Deliveries) List(w http.ResponseWriter, r *http.Request) {
	var delivery models.Delivery
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	list, err := delivery.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting Deliveries list: %v", err))
		return
	}

	tx.Commit()

	var listResponse []*response.DeliveryListResponse
	for _, r := range list {
		var deliveryResponse response.DeliveryListResponse
		deliveryResponse.Transform(&r)
		listResponse = append(listResponse, &deliveryResponse)
	}

	api.ResponseOK(w, listResponse, http.StatusOK)
}

// View : http handler for retrieve Delivery by id
func (u *Deliveries) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting: %v", err))
		return
	}

	var delivery models.Delivery
	delivery.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = delivery.Get(ctx, tx)

	if err == sql.ErrNoRows {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrNotFound(err, ""))
		return
	}

	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get Delivery: %v", err))
		return
	}

	tx.Commit()

	var response response.DeliveryResponse
	response.Transform(&delivery)
	api.ResponseOK(w, response, http.StatusOK)
}

// Create : http handler for create new Delivery
func (u *Deliveries) Create(w http.ResponseWriter, r *http.Request) {
	var deliveryRequest request.NewDeliveryRequest
	err := api.Decode(r, &deliveryRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("decode Delivery: %v", err))
		return
	}

	delivery := deliveryRequest.Transform()
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("begin transaction: %v", err))
		return
	}

	err = delivery.Create(r.Context(), tx)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		tx.Rollback()
		api.ResponseError(w, fmt.Errorf("Create Delivery: %v", err))
		return
	}

	tx.Commit()

	var response response.DeliveryResponse
	response.Transform(delivery)
	api.ResponseOK(w, response, http.StatusCreated)
}

// Update : http handler for update Delivery by id
func (u *Deliveries) Update(w http.ResponseWriter, r *http.Request) {
	// TODO :
	// Edit Delivery hanya boleh dilakukan jika :
	// belum ada transaksi lain atas delivery seperti return delivery, mutations, atau sales/delivery

	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var delivery models.Delivery
	delivery.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = delivery.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get Delivery: %v", err))
		return
	}

	var deliveryRequest request.DeliveryRequest
	err = api.Decode(r, &deliveryRequest)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Decode Delivery: %v", err))
		return
	}

	if deliveryRequest.ID <= 0 {
		deliveryRequest.ID = delivery.ID
	}
	deliveryUpdate := deliveryRequest.Transform(&delivery)
	err = deliveryUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Update Delivery: %v", err))
		return
	}

	tx.Commit()

	var response response.DeliveryResponse
	response.Transform(deliveryUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}
