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

// Salesmen type for handling dependency injection
type Salesmen struct {
	Db  *sql.DB
	Log *log.Logger
}

// List of salesmen
func (u *Salesmen) List(w http.ResponseWriter, r *http.Request) {
	var salesman models.Salesman

	list, err := salesman.List(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("get salesmen list : %v", err)
		api.ResponseError(w, err)
		return
	}

	var salesmanResponse []response.SalesmanResponse
	for _, r := range list {
		var res response.SalesmanResponse
		res.Transform(&r)
		salesmanResponse = append(salesmanResponse, res)
	}

	api.ResponseOK(w, salesmanResponse, http.StatusOK)
}

// Create new salesman
func (u *Salesmen) Create(w http.ResponseWriter, r *http.Request) {
	var salesmanRequest request.NewSalesmanRequest
	err := api.Decode(r, &salesmanRequest)
	if err != nil {
		u.Log.Printf("Error : %v", err)
		api.ResponseError(w, err)
		return
	}

	salesman := salesmanRequest.Transform()

	err = salesman.Create(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("create new salesman: %v", err)
		api.ResponseError(w, err)
		return
	}

	var res response.SalesmanResponse
	res.Transform(&salesman)
	api.ResponseOK(w, res, http.StatusCreated)
}

// View of salesman by id
func (u *Salesmen) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var salesman models.Salesman
	salesman.ID = uint64(id)

	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = salesman.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Get salesman: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	var res response.SalesmanResponse
	res.Transform(&salesman)
	api.ResponseOK(w, res, http.StatusOK)
}

// Update salesman by id
func (u *Salesmen) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var salesmanRequest request.SalesmanRequest
	err = api.Decode(r, &salesmanRequest)
	if err != nil {
		u.Log.Printf("Error : %v", err)
		api.ResponseError(w, err)
		return
	}

	var salesman models.Salesman
	salesman.ID = uint64(id)

	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = salesman.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Get salesman: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	salesmanUpdate := salesmanRequest.Transform(&salesman)

	err = salesmanUpdate.Update(ctx, u.Db)
	if err != nil {
		u.Log.Printf("Update salesman: %v", err)
		api.ResponseError(w, err)
		return
	}

	var res response.SalesmanResponse
	res.Transform(salesmanUpdate)
	api.ResponseOK(w, res, http.StatusOK)
}

// Delete salesman by id
func (u *Salesmen) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var salesman models.Salesman
	salesman.ID = uint64(id)

	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = salesman.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Get salesman: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	err = salesman.Delete(ctx, u.Db)
	if err != nil {
		u.Log.Printf("Delete salesman: %v", err)
		api.ResponseError(w, err)
		return
	}

	api.ResponseOK(w, nil, http.StatusNoContent)
}
