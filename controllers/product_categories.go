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

// ProductCategories type for handling dependency injection
type ProductCategories struct {
	Db  *sql.DB
	Log *log.Logger
}

// List of ProductCategories
func (u *ProductCategories) List(w http.ResponseWriter, r *http.Request) {
	var ProductCategory models.ProductCategory
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("Error begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	list, err := ProductCategory.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("get ProductCategories list : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx.Commit()

	var ProductCategoryResponse []response.ProductCategoryResponse
	for _, r := range list {
		var res response.ProductCategoryResponse
		res.Transform(&r)
		ProductCategoryResponse = append(ProductCategoryResponse, res)
	}

	api.ResponseOK(w, ProductCategoryResponse, http.StatusOK)
}

// Create new ProductCategory
func (u *ProductCategories) Create(w http.ResponseWriter, r *http.Request) {
	var ProductCategoryRequest request.NewProductCategoryRequest
	err := api.Decode(r, &ProductCategoryRequest)
	if err != nil {
		u.Log.Printf("Error : %v", err)
		api.ResponseError(w, err)
		return
	}

	ProductCategory := ProductCategoryRequest.Transform()

	tx, err := u.Db.Begin()
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Error begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}
	err = ProductCategory.Create(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("create new ProductCategory tx : %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	var response response.ProductCategoryResponse
	response.Transform(&ProductCategory)
	api.ResponseOK(w, response, http.StatusCreated)
}

// View of ProductCategory by id
func (u *ProductCategories) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var ProductCategory models.ProductCategory
	ProductCategory.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	err = ProductCategory.View(ctx, tx)
	if err != nil {
		u.Log.Printf("Get ProductCategory: %v", err)
		api.ResponseError(w, err)
		return
	}

	var response response.ProductCategoryResponse
	response.Transform(&ProductCategory)
	api.ResponseOK(w, response, http.StatusOK)
}

// Update ProductCategory by id
func (u *ProductCategories) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var ProductCategoryRequest request.ProductCategoryRequest
	err = api.Decode(r, &ProductCategoryRequest)
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

	var ProductCategory models.ProductCategory
	ProductCategory.ID = uint64(id)
	err = ProductCategory.View(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Get ProductCategory: %v", err)
		api.ResponseError(w, err)
		return
	}

	ProductCategoryUpdate := ProductCategoryRequest.Transform(&ProductCategory)

	err = ProductCategoryUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Update ProductCategory: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	var response response.ProductCategoryResponse
	response.Transform(ProductCategoryUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}

// Delete ProductCategory by id
func (u *ProductCategories) Delete(w http.ResponseWriter, r *http.Request) {
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

	var ProductCategory models.ProductCategory
	ProductCategory.ID = uint64(id)
	err = ProductCategory.View(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Get ProductCategory: %v", err)
		api.ResponseError(w, err)
		return
	}

	err = ProductCategory.Delete(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Update ProductCategory: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	api.ResponseOK(w, nil, http.StatusNoContent)
}
