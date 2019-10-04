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

// Branches type for handling dependency injection
type Branches struct {
	Db  *sql.DB
	Log *log.Logger
}

// List of branches
func (b *Branches) List(w http.ResponseWriter, r *http.Request) {
	var branch models.Branch
	tx, err := b.Db.Begin()
	if err != nil {
		b.Log.Printf("Error begin tx: %v", err)
		api.ResponseError(w, err)
		return
	}

	list, err := branch.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		b.Log.Printf("get branches list: %v", err)
		api.ResponseError(w, err)
		return
	}

	tx.Commit()

	var branchResponse []response.BranchResponse
	for _, r := range list {
		var res response.BranchResponse
		res.Transform(&r)
		branchResponse = append(branchResponse, res)
	}

	api.ResponseOK(w, branchResponse, http.StatusOK)

}

// Create new branch
func (b *Branches) Create(w http.ResponseWriter, r *http.Request) {
	var branchRequest request.NewBranchRequest
	err := api.Decode(r, &branchRequest)
	if err != nil {
		b.Log.Printf("Error : %v", err)
		api.ResponseError(w, err)
		return
	}

	branch := branchRequest.Transform()

	tx, err := b.Db.Begin()
	if err != nil {
		tx.Rollback()
		b.Log.Printf("Error begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}
	err = branch.Create(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		b.Log.Printf("create new branch tx : %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	var response response.BranchResponse
	response.Transform(branch)
	api.ResponseOK(w, response, http.StatusCreated)
}

// View branches by id
func (b *Branches) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		b.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var branch models.Branch
	branch.ID = uint32(id)
	tx, err := b.Db.Begin()
	if err != nil {
		b.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	err = branch.Get(ctx, tx)
	if err != nil {
		b.Log.Printf("Get branch: %v", err)
		api.ResponseError(w, err)
		return
	}

	var response response.BranchResponse
	response.Transform(&branch)
	api.ResponseOK(w, response, http.StatusOK)
}

// Update branch by id
func (b *Branches) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		b.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var branchRequest request.BranchRequest
	err = api.Decode(r, &branchRequest)
	if err != nil {
		b.Log.Printf("Error : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx, err := b.Db.Begin()
	if err != nil {
		tx.Rollback()
		b.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	var branch models.Branch
	branch.ID = uint32(id)
	err = branch.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		b.Log.Printf("Get branch: %v", err)
		api.ResponseError(w, err)
		return
	}

	branchUpdate := branchRequest.Transform(&branch)

	err = branchUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		b.Log.Printf("Update branch: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	var response response.BranchResponse
	response.Transform(branchUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}

// Delete branch by id
func (b *Branches) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		b.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx, err := b.Db.Begin()
	if err != nil {
		tx.Rollback()
		b.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	var branch models.Branch
	branch.ID = uint32(id)
	err = branch.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		b.Log.Printf("Get branch: %v", err)
		api.ResponseError(w, err)
		return
	}

	err = branch.Delete(ctx, tx)
	if err != nil {
		tx.Rollback()
		b.Log.Printf("Delete branch: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	api.ResponseOK(w, nil, http.StatusNoContent)
}
