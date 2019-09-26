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

//Regions : struct for set Regions Dependency Injection
type Regions struct {
	Db  *sql.DB
	Log *log.Logger
}

//List : http handler for returning list of regions
func (u *Regions) List(w http.ResponseWriter, r *http.Request) {
	var region models.Region
	list, err := region.List(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting regions list: %v", err))
		return
	}

	var listResponse []*response.RegionResponse
	for _, region := range list {
		var regionResponse response.RegionResponse
		regionResponse.Transform(&region)
		listResponse = append(listResponse, &regionResponse)
	}

	api.ResponseOK(w, listResponse, http.StatusOK)
}

//View : http handler for retrieve region by id
func (u *Regions) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting: %v", err))
		return
	}

	var region models.Region
	region.ID = uint32(id)
	err = region.Get(ctx, u.Db)

	if err == sql.ErrNoRows {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrNotFound(err, ""))
		return
	}

	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get region: %v", err))
		return
	}

	var response response.RegionResponse
	response.Transform(&region)
	api.ResponseOK(w, response, http.StatusOK)
}

//Create : http handler for create new region
func (u *Regions) Create(w http.ResponseWriter, r *http.Request) {
	var regionRequest request.NewRegionRequest
	err := api.Decode(r, &regionRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("decode role: %v", err))
		return
	}

	region := regionRequest.Transform()
	err = region.Create(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Create region: %v", err))
		return
	}

	var response response.RegionResponse
	response.Transform(region)
	api.ResponseOK(w, response, http.StatusCreated)
}

//Update : http handler for update region by id
func (u *Regions) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var region models.Region
	region.ID = uint32(id)
	err = region.Get(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get region: %v", err))
		return
	}

	var regionRequest request.RegionRequest
	err = api.Decode(r, &regionRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Decode region: %v", err))
		return
	}

	if regionRequest.ID <= 0 {
		regionRequest.ID = region.ID
	}
	regionUpdate := regionRequest.Transform(&region)
	err = regionUpdate.Update(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Update region: %v", err))
		return
	}

	var response response.RegionResponse
	response.Transform(regionUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}

// Delete : http handler for delete role by id
func (u *Regions) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var region models.Region
	region.ID = uint32(id)
	err = region.Get(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get region: %v", err))
		return
	}

	err = region.Delete(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Delete region: %v", err))
		return
	}

	api.ResponseOK(w, nil, http.StatusNoContent)
}

// AddBranch : http handler for add branch to region
func (u *Regions) AddBranch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ps := ctx.Value(api.Ctx("ps")).(httprouter.Params)
	paramID := ps.ByName("id")
	paramBranchID := ps.ByName("branch_id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	branchID, err := strconv.Atoi(paramBranchID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramBranchID: %v", err))
		return
	}

	var region models.Region
	region.ID = uint32(id)
	err = region.Get(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get region: %v", err))
		return
	}

	var branch models.Branch
	branch.ID = uint32(branchID)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = branch.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get branch: %v", err))
		return
	}

	err = region.AddBranch(ctx, tx, branch.ID)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("add branch to region: %v", err))
		return
	}

	tx.Commit()

	api.ResponseOK(w, nil, http.StatusOK)
}

// DeleteBranch : http handler for delete branch from region
func (u *Regions) DeleteBranch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ps := ctx.Value(api.Ctx("ps")).(httprouter.Params)
	paramID := ps.ByName("id")
	paramBranchID := ps.ByName("branch_id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	branchID, err := strconv.Atoi(paramBranchID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramBranchID: %v", err))
		return
	}

	var region models.Region
	region.ID = uint32(id)
	err = region.Get(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get region: %v", err))
		return
	}

	var branch models.Branch
	branch.ID = uint32(branchID)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = branch.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get branch: %v", err))
		return
	}

	err = region.DeleteBranch(ctx, tx, branch.ID)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("delete branch from reegion: %v", err))
		return
	}

	tx.Commit()

	api.ResponseOK(w, nil, http.StatusNoContent)
}
