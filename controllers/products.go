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

// Products : struct for set Products Dependency Injection
type Products struct {
	Db  *sql.DB
	Log *log.Logger
}

//List : http handler for returning list of products
func (u *Products) List(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	list, err := product.List(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting products list: %v", err))
		return
	}

	var listResponse []*response.ProductResponse
	for _, product := range list {
		var productResponse response.ProductResponse
		productResponse.Transform(&product)
		listResponse = append(listResponse, &productResponse)
	}

	api.ResponseOK(w, listResponse, http.StatusOK)
}

//View : http handler for retrieve product by id
func (u *Products) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting: %v", err))
		return
	}

	var product models.Product
	product.ID = uint64(id)
	err = product.Get(ctx, u.Db)

	if err == sql.ErrNoRows {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrNotFound(err, ""))
		return
	}

	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get product: %v", err))
		return
	}

	var response response.ProductResponse
	response.Transform(&product)
	api.ResponseOK(w, response, http.StatusOK)
}

//Create : http handler for create new product
func (u *Products) Create(w http.ResponseWriter, r *http.Request) {
	var productRequest request.NewProductRequest
	err := api.Decode(r, &productRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("decode product: %v", err))
		return
	}

	product := productRequest.Transform()
	err = product.Create(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Create product: %v", err))
		return
	}

	var response response.ProductResponse
	response.Transform(product)
	api.ResponseOK(w, response, http.StatusCreated)
}

//Update : http handler for update product by id
func (u *Products) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var product models.Product
	product.ID = uint64(id)
	err = product.Get(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get product: %v", err))
		return
	}

	var productRequest request.ProductRequest
	err = api.Decode(r, &productRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Decode product: %v", err))
		return
	}

	if productRequest.ID <= 0 {
		productRequest.ID = product.ID
	}
	productUpdate := productRequest.Transform(&product)
	err = productUpdate.Update(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Update product: %v", err))
		return
	}

	var response response.ProductResponse
	response.Transform(productUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}

// Delete : http handler for delete product by id
func (u *Products) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var product models.Product
	product.ID = uint64(id)
	err = product.Get(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get product: %v", err))
		return
	}

	err = product.Delete(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Delete product: %v", err))
		return
	}

	api.ResponseOK(w, nil, http.StatusNoContent)
}
