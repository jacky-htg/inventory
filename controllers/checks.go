package controllers

import (
	"database/sql"
	"net/http"

	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/libraries/database"
)

//Checks : struct for set Checks Dependency Injection
type Checks struct {
	Db *sql.DB
}

//Health : http handler for health checking
func (u *Checks) Health(w http.ResponseWriter, r *http.Request) {
	var health struct {
		Status string `json:"status"`
	}

	// Check if the database is ready.
	err := database.StatusCheck(r.Context(), u.Db)
	if err != nil {
		// If the database is not ready we will tell the client and use a 500
		// status. Do not respond by just returning an error because further up in
		// the call stack will interpret that as an unhandled error.
		api.ResponseError(w, err)
		return
	}

	health.Status = "ok"
	api.ResponseOK(w, health, http.StatusOK)
}
