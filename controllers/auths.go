package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/libraries/token"
	"github.com/jacky-htg/inventory/models"
	"github.com/jacky-htg/inventory/payloads/request"
	"github.com/jacky-htg/inventory/payloads/response"
	"golang.org/x/crypto/bcrypt"
)

// Auths struct
type Auths struct {
	Db  *sql.DB
	Log *log.Logger
}

// Login http handler
func (u *Auths) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest request.LoginRequest
	err := api.Decode(r, &loginRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, err)
		return
	}

	uLogin := models.User{Username: loginRequest.Username}
	err = uLogin.GetByUsername(r.Context(), u.Db)
	if err != nil {
		err = fmt.Errorf("call login: %v", err)
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(uLogin.Password), []byte(loginRequest.Password))
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrBadRequest(fmt.Errorf("compare password: %v", err), ""))
		return
	}

	token, err := token.ClaimToken(uLogin.Username)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("claim token: %v", err))
		return
	}

	var response response.TokenResponse
	response.Token = token

	api.ResponseOK(w, response, http.StatusOK)
}
