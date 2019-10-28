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

// Shelve type for handling dependency injection
type Shelves struct {
	Db  *sql.DB
	Log *log.Logger
}

// List of Shelves by branch id
func (s *Shelves) List(w http.ResponseWriter, r *http.Request) {
	var shelve models.Shelve
	tx, err := s.Db.Begin()
	if err != nil {
		s.Log.Printf("Error begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	list, err := shelve.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		s.Log.Printf("get shelves list : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx.Commit()

	var shelveResponse []response.ShelveResponse
	for _, r := range list {
		var res response.ShelveResponse
		res.Transform(&r)
		shelveResponse = append(shelveResponse, res)
	}

	api.ResponseOK(w, shelveResponse, http.StatusOK)
}

// Create new Shelve
func (s *Shelves) Create(w http.ResponseWriter, r *http.Request) {
	var shelveRequest request.NewShelveRequest
	err := api.Decode(r, &shelveRequest)
	if err != nil {
		s.Log.Printf("Error Decode : %v", err)
		api.ResponseError(w, err)
		return
	}

	shelve := shelveRequest.Transform()

	tx, err := s.Db.Begin()
	if err != nil {
		tx.Rollback()
		s.Log.Printf("Error begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	err = shelve.Create(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		s.Log.Printf("Error create new Shelve tx : %v", err)
		api.ResponseError(w, err)
	}

	tx.Commit()

	var response response.ShelveResponse
	response.Transform(&shelve)
	api.ResponseOK(w, response, http.StatusCreated)
}

// View Shelve by id
func (s *Shelves) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		s.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var shelve models.Shelve
	shelve.ID = uint64(id)
	tx, err := s.Db.Begin()
	if err != nil {
		s.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	err = shelve.View(ctx, tx)
	if err != nil {
		s.Log.Printf("Get Shelve Error: %v", err)
		api.ResponseError(w, err)
		return
	}

	var response response.ShelveResponse
	response.Transform(&shelve)
	api.ResponseOK(w, response, http.StatusOK)
}

// Update Shelve by id
func (s *Shelves) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		s.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var shelveRequest request.ShelveRequest
	err = api.Decode(r, &shelveRequest)
	if err != nil {
		s.Log.Printf("Error : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx, err := s.Db.Begin()
	if err != nil {
		tx.Rollback()
		s.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	var shelve models.Shelve
	shelve.ID = uint64(id)
	err = shelve.View(ctx, tx)
	if err != nil {
		tx.Rollback()
		s.Log.Printf("Get shelve : %v", err)
		api.ResponseError(w, err)
		return
	}

	shelveUpdate := shelveRequest.Transform(&shelve)

	err = shelveUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		s.Log.Printf("Update shelve : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx.Commit()

	var response response.ShelveResponse
	response.Transform(shelveUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}

// Delete Shelve by id
func (s *Shelves) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramId := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		s.Log.Printf("casting paramId : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx, err := s.Db.Begin()
	if err != nil {
		tx.Rollback()
		s.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	var shelve models.Shelve
	shelve.ID = uint64(id)
	err = shelve.View(ctx, tx)
	if err != nil {
		tx.Rollback()
		s.Log.Printf("Get shelve : %v", err)
		api.ResponseError(w, err)
		return
	}

	err = shelve.Delete(ctx, tx)
	if err != nil {
		tx.Rollback()
		s.Log.Printf("Delete shelve : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx.Commit()
	api.ResponseOK(w, nil, http.StatusNoContent)
}
