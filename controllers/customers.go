package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/models"
	"github.com/jacky-htg/inventory/payloads/request"
	"github.com/jacky-htg/inventory/payloads/response"
	"github.com/julienschmidt/httprouter"
)

// Customers type for handling dependency injection
type Customers struct {
	Db  *sql.DB
	Log *log.Logger
}

// List of customers
func (u *Customers) List(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("Error begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	list, err := customer.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("get customers list : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx.Commit()

	var customerResponse []response.CustomerResponse
	for _, r := range list {
		var res response.CustomerResponse
		res.Transform(&r)
		customerResponse = append(customerResponse, res)
	}

	api.ResponseOK(w, customerResponse, http.StatusOK)
}

// Create new customer
func (u *Customers) Create(w http.ResponseWriter, r *http.Request) {
	var customerRequest request.NewCustomerRequest
	err := api.Decode(r, &customerRequest)
	if err != nil {
		u.Log.Printf("Error : %v", err)
		api.ResponseError(w, err)
		return
	}

	customer := customerRequest.Transform()

	tx, err := u.Db.Begin()
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Error begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}
	err = customer.Create(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("create new customer tx : %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	var response response.CustomerResponse
	response.Transform(&customer)
	api.ResponseOK(w, response, http.StatusCreated)
}

// View of customer by id
func (u *Customers) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var customer models.Customer
	customer.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	err = customer.View(ctx, tx)
	if err != nil {
		u.Log.Printf("Get customer: %v", err)
		api.ResponseError(w, err)
		return
	}

	var response response.CustomerResponse
	response.Transform(&customer)
	api.ResponseOK(w, response, http.StatusOK)
}

// Update customer by id
func (u *Customers) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var customerRequest request.CustomerRequest
	err = api.Decode(r, &customerRequest)
	if err != nil {
		u.Log.Printf("Error : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx, err := u.Db.Begin()
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	var customer models.Customer
	customer.ID = uint64(id)
	err = customer.View(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Get customer: %v", err)
		api.ResponseError(w, err)
		return
	}

	customerUpdate := customerRequest.Transform(&customer)

	err = customerUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Update customer: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	var response response.CustomerResponse
	response.Transform(customerUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}

// Delete customer by id
func (u *Customers) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx, err := u.Db.Begin()
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	var customer models.Customer
	customer.ID = uint64(id)
	err = customer.View(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Get customer: %v", err)
		api.ResponseError(w, err)
		return
	}

	err = customer.Delete(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Update customer: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	api.ResponseOK(w, nil, http.StatusNoContent)
}
