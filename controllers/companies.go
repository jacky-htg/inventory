package controllers

import (
	"database/sql"
	"errors"
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

//Companies : struct for set Companies Dependency Injection
type Companies struct {
	Db  *sql.DB
	Log *log.Logger
}

//View : http handler for retrieve company by id
func (u *Companies) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting: %v", err))
		return
	}

	var company models.Company
	company.ID = uint32(id)
	if company.ID != ctx.Value(api.Ctx("auth")).(models.User).Company.ID {
		err = api.ErrForbidden(errors.New("Forbidden data owner"), "")
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, err)
		return
	}

	err = company.Get(ctx, u.Db)
	if err == sql.ErrNoRows {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrNotFound(err, ""))
		return
	}

	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get company: %v", err))
		return
	}

	var response response.CompanyResponse
	response.Transform(&company)
	api.ResponseOK(w, response, http.StatusOK)
}

//Create : http handler for create new company
func (u *Companies) Create(w http.ResponseWriter, r *http.Request) {
	var companyRequest request.NewCompanyRequest
	err := api.Decode(r, &companyRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("decode role: %v", err))
		return
	}

	company := companyRequest.Transform()
	err = company.Create(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Create company: %v", err))
		return
	}

	var response response.CompanyResponse
	response.Transform(company)
	api.ResponseOK(w, response, http.StatusCreated)
}

//Update : http handler for update company by id
func (u *Companies) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var company models.Company
	company.ID = uint32(id)
	if company.ID != ctx.Value(api.Ctx("auth")).(models.User).Company.ID {
		err = api.ErrForbidden(errors.New("Forbidden data owner"), "")
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, err)
		return
	}

	err = company.Get(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get company: %v", err))
		return
	}

	var companyRequest request.CompanyRequest
	err = api.Decode(r, &companyRequest)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Decode company: %v", err))
		return
	}

	if companyRequest.ID <= 0 {
		companyRequest.ID = company.ID
	}
	companyUpdate := companyRequest.Transform(&company)
	err = companyUpdate.Update(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Update company: %v", err))
		return
	}

	var response response.CompanyResponse
	response.Transform(companyUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}

// Delete : http handler for delete company by id
func (u *Companies) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("type casting paramID: %v", err))
		return
	}

	var company models.Company
	company.ID = uint32(id)
	if company.ID != ctx.Value(api.Ctx("auth")).(models.User).Company.ID {
		err = api.ErrForbidden(errors.New("Forbidden data owner"), "")
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, err)
		return
	}

	err = company.Get(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Get company: %v", err))
		return
	}

	err = company.Delete(ctx, u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("Delete company: %v", err))
		return
	}

	api.ResponseOK(w, nil, http.StatusNoContent)
}
