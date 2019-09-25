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

//Roles : struct for set Roles Dependency Injection
type Roles struct {
	Db  *sql.DB
	Log *log.Logger
}

//List : http handler for returning list of roles
func (u *Roles) List(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	list, err := role.List(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting roles list: %v", err))
		return
	}

	var listResponse []*response.RoleResponse
	for _, role := range list {
		var roleResponse response.RoleResponse
		roleResponse.Transform(&role)
		listResponse = append(listResponse, &roleResponse)
	}

	api.ResponseOK(w, listResponse, http.StatusOK)
}

//View : http handler for retrieve role by id
func (u *Roles) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting: %v", err))
		return
	}

	var role models.Role
	role.ID = uint32(id)
	err = role.Get(ctx, u.Db)

	if err == sql.ErrNoRows {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrNotFound(err, ""))
		return
	}

	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get Role: %v", err))
		return
	}

	var response response.RoleResponse
	response.Transform(&role)
	api.ResponseOK(w, response, http.StatusOK)
}

//Create : http handler for create new role
func (u *Roles) Create(w http.ResponseWriter, r *http.Request) {
	var roleRequest request.NewRoleRequest
	err := api.Decode(r, &roleRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("decode role: %v", err))
		return
	}

	role := roleRequest.Transform()
	err = role.Create(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Create Role: %v", err))
		return
	}

	var response response.RoleResponse
	response.Transform(role)
	api.ResponseOK(w, response, http.StatusCreated)
}

//Update : http handler for update role by id
func (u *Roles) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var role models.Role
	role.ID = uint32(id)
	err = role.Get(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get Role: %v", err))
		return
	}

	var roleRequest request.RoleRequest
	err = api.Decode(r, &roleRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Decode Role: %v", err))
		return
	}

	if roleRequest.ID <= 0 {
		roleRequest.ID = role.ID
	}
	roleUpdate := roleRequest.Transform(&role)
	err = roleUpdate.Update(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Update role: %v", err))
		return
	}

	var response response.RoleResponse
	response.Transform(roleUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}

// Delete : http handler for delete role by id
func (u *Roles) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var role models.Role
	role.ID = uint32(id)
	err = role.Get(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get role: %v", err))
		return
	}

	err = role.Delete(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Delete role: %v", err))
		return
	}

	api.ResponseOK(w, nil, http.StatusNoContent)
}

//Grant : http handler for grant access to role
func (u *Roles) Grant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ps := ctx.Value(api.Ctx("ps")).(httprouter.Params)
	paramID := ps.ByName("id")
	paramAccessID := ps.ByName("access_id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	accessID, err := strconv.Atoi(paramAccessID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramAccessID: %v", err))
		return
	}

	var role models.Role
	role.ID = uint32(id)
	err = role.Get(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get role: %v", err))
		return
	}

	var access models.Access
	access.ID = uint32(accessID)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = access.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get access: %v", err))
		return
	}

	err = role.Grant(ctx, u.Db, access.ID)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Grant role: %v", err))
		return
	}

	tx.Commit()

	api.ResponseOK(w, nil, http.StatusOK)
}

//Revoke : http handler for revoke access from role
func (u *Roles) Revoke(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ps := ctx.Value(api.Ctx("ps")).(httprouter.Params)
	paramID := ps.ByName("id")
	paramAccessID := ps.ByName("access_id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	accessID, err := strconv.Atoi(paramAccessID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramAccessID: %v", err))
		return
	}

	var role models.Role
	role.ID = uint32(id)
	err = role.Get(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get role: %v", err))
		return
	}

	var access models.Access
	access.ID = uint32(accessID)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Begin tx: %v", err))
		return
	}

	err = access.Get(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get access: %v", err))
		return
	}

	err = role.Revoke(ctx, u.Db, access.ID)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Revoke role: %v", err))
		return
	}

	tx.Commit()

	api.ResponseOK(w, nil, http.StatusNoContent)
}
