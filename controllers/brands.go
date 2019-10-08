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

// Brands type for handling dependency injection
type Brands struct {
	Db  *sql.DB
	Log *log.Logger
}

// List of Brands
func (u *Brands) List(w http.ResponseWriter, r *http.Request) {
	var Brand models.Brand
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("Error begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	list, err := Brand.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("get Brands list : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx.Commit()

	var BrandResponse []response.BrandResponse
	for _, r := range list {
		var res response.BrandResponse
		res.Transform(&r)
		BrandResponse = append(BrandResponse, res)
	}

	api.ResponseOK(w, BrandResponse, http.StatusOK)
}

// Create new Brand
func (u *Brands) Create(w http.ResponseWriter, r *http.Request) {
	var BrandRequest request.NewBrandRequest
	err := api.Decode(r, &BrandRequest)
	if err != nil {
		u.Log.Printf("Error : %v", err)
		api.ResponseError(w, err)
		return
	}

	Brand := BrandRequest.Transform()

	tx, err := u.Db.Begin()
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Error begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}
	err = Brand.Create(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("create new Brand tx : %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	var response response.BrandResponse
	response.Transform(&Brand)
	api.ResponseOK(w, response, http.StatusCreated)
}

// View of Brand by id
func (u *Brands) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var Brand models.Brand
	Brand.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	err = Brand.View(ctx, tx)
	if err != nil {
		u.Log.Printf("Get Brand: %v", err)
		api.ResponseError(w, err)
		return
	}

	var response response.BrandResponse
	response.Transform(&Brand)
	api.ResponseOK(w, response, http.StatusOK)
}

// Update Brand by id
func (u *Brands) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var BrandRequest request.BrandRequest
	err = api.Decode(r, &BrandRequest)
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

	var Brand models.Brand
	Brand.ID = uint64(id)
	err = Brand.View(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Get Brand: %v", err)
		api.ResponseError(w, err)
		return
	}

	BrandUpdate := BrandRequest.Transform(&Brand)

	err = BrandUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Update Brand: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	var response response.BrandResponse
	response.Transform(BrandUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}

// Delete Brand by id
func (u *Brands) Delete(w http.ResponseWriter, r *http.Request) {
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

	var Brand models.Brand
	Brand.ID = uint64(id)
	err = Brand.View(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Get Brand: %v", err)
		api.ResponseError(w, err)
		return
	}

	err = Brand.Delete(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Update Brand: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	api.ResponseOK(w, nil, http.StatusNoContent)
}
