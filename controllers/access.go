package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/models"
	"github.com/jacky-htg/inventory/payloads/response"
)

//Access : struct for set Access Dependency Injection
type Access struct {
	Db  *sql.DB
	Log *log.Logger
}

//List : http handler for returning list of access
func (u *Access) List(w http.ResponseWriter, r *http.Request) {
	var access models.Access
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("Begin tx : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting access list: %v", err))
		return
	}
	list, err := access.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("getting access list: %v", err))
		return
	}

	var listResponse []*response.AccessResponse
	for _, a := range list {
		var accessResponse response.AccessResponse
		accessResponse.Transform(&a)
		listResponse = append(listResponse, &accessResponse)
	}

	api.ResponseOK(w, listResponse, http.StatusOK)
}
