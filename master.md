# Add New Module (Master)

## Design Table
- Open file schema/migrate.go
- Add design of table in schema migration
```
	{
	Version:     24,
	Description: "Add Closing Stocks",
	Script: `
CREATE TABLE customers (
id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
company_id	INT(10) UNSIGNED NOT NULL,
name VARCHAR(100) NOT NULL,
email VARCHAR(100) NOT NULL UNIQUE,
address VARCHAR(255) NOT NULL,
hp CHAR(15) NOT NULL,
PRIMARY KEY (id),
KEY customers_company_id (company_id),
CONSTRAINT fk_customers_to_companies FOREIGN KEY (company_id) REFERENCES companies(id)
);`,
	},
```
- go run cmd/main.go migrate

## Model
- Create new file models/customer.go
```
package models

import (
	"context"
	"database/sql"

	"github.com/jacky-htg/inventory/libraries/api"
)

// Customer : struct of customer
type Customer struct {
	ID      uint64
	Company Company
	Name    string
	Email   string
	Address string
	Hp      string
}

const qCustomers = `SELECT id, name, email, address, hp FROM customers`

// List of customers
func (u *Customer) List(ctx context.Context, tx *sql.Tx) ([]Customer, error) {
	var list []Customer

	rows, err := tx.QueryContext(
		ctx, 
		qCustomers+" WHERE company_id=?", 
		ctx.Value(api.Ctx("auth")).(User).Company.ID,
	)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Customer
		err = rows.Scan(&c.ID, &c.Name, &c.Email, &c.Address, &c.Hp)
		if err != nil {
			return list, err
		}

		list = append(list, c)
	}

	return list, rows.Err()
}

// Create new customer
func (u *Customer) Create(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO customers (company_id, name, email, address, hp) 
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	userLogin := ctx.Value(api.Ctx("auth")).(User)
	res, err := stmt.ExecContext(ctx, userLogin.Company.ID, u.Name, u.Email, u.Address, u.Hp)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	u.ID = uint64(id)
	u.Company = userLogin.Company

	return err
}

// View customer by id
func (u *Customer) View(ctx context.Context, tx *sql.Tx) error {
	return tx.QueryRowContext(
		ctx,
		qCustomers+" WHERE id=? AND company_id=?",
		u.ID,
		ctx.Value(api.Ctx("auth")).(User).Company.ID,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Address, &u.Hp)
}

// Update customer by id
func (u *Customer) Update(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `
		UPDATE customers  
		SET name = ?, 
			email = ?, 
			address = ?,
			hp = ?
		WHERE id = ? AND company_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	userLogin := ctx.Value(api.Ctx("auth")).(User)
	u.Company = userLogin.Company
	_, err = stmt.ExecContext(ctx, u.Name, u.Email, u.Address, u.Hp, u.ID, userLogin.Company.ID)

	return err
}

// Delete customer by id
func (u *Customer) Delete(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM customers WHERE id = ? AND company_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID, ctx.Value(api.Ctx("auth")).(User).Company.ID)

	return err
}

```

## Payload Response
- Create new file payloads/response/customer_response.go
```
package response

import "github.com/jacky-htg/inventory/models"

// CustomerResponse json
type CustomerResponse struct {
	ID      uint64          `json:"id"`
	Company CompanyResponse `json:"company"`
	Name    string          `json:"name"`
	Email   string          `json:"email"`
	Address string          `json:"address"`
	Hp      string          `json:"hp"`
}

// Transform Customer models to customer response
func (u *CustomerResponse) Transform(c models.Customer) {
	u.ID = c.ID
	u.Name = c.Name
	u.Email = c.Email
	u.Address = c.Address
	u.Hp = c.Hp
	u.Company.Transform(&c.Company)
}
```

## Payload Request
- Create new file payloads/request/customer_request.go
```
package request

import "github.com/jacky-htg/inventory/models"

// NewCustomerRequest is json request for new customer and validation
type NewCustomerRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required"`
	Address string `json:"address" validate:"required"`
	Hp      string `json:"hp" validate:"required"`
}

// Transform NewCustomerRequest to Customer model
func (u *NewCustomerRequest) Transform() models.Customer {
	var c models.Customer
	c.Name = u.Name
	c.Email = u.Email
	c.Address = u.Address
	c.Hp = u.Hp

	return c
}

// CustomerRequest is json request for update customer and validation
type CustomerRequest struct {
	ID      uint64 `json:"id" validate:"required"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Hp      string `json:"hp"`
}

// Transform CustomerRequest to Customer model
func (u *CustomerRequest) Transform(c *models.Customer) *models.Customer {
	if c.ID == u.ID {
		if len(u.Name) > 0 {
			c.Name = u.Name
		}
		if len(u.Email) > 0 {
			c.Email = u.Email
		}
		if len(u.Address) > 0 {
			c.Address = u.Address
		}
		if len(u.Hp) > 0 {
			c.Hp = u.Hp
		}
	}
	return c
}
```

## Controller
- Create new file controllers/customers.go
```
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

// Customers type for handling dependency injection
type Customers struct {
	Db  *sql.DB
	Log *log.Logger
}

// List of customers
func (u *Customers) List(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("Error begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	list, err := customer.List(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("get customers list : %v", err)
		api.ResponseError(w, err)
		return
	}

	tx.Commit()

	var customerResponse []response.CustomerResponse
	for _, r := range list {
		var res response.CustomerResponse
		res.Transform(&r)
		customerResponse = append(customerResponse, res)
	}

	api.ResponseOK(w, customerResponse, http.StatusOK)
}

// Create new customer
func (u *Customers) Create(w http.ResponseWriter, r *http.Request) {
	var customerRequest request.NewCustomerRequest
	err := api.Decode(r, &customerRequest)
	if err != nil {
		u.Log.Printf("Error : %v", err)
		api.ResponseError(w, err)
		return
	}

	customer := customerRequest.Transform()

	tx, err := u.Db.Begin()
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Error begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}
	err = customer.Create(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("create new customer tx : %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	var response response.CustomerResponse
	response.Transform(&customer)
	api.ResponseOK(w, response, http.StatusCreated)
}

// View of customer by id
func (u *Customers) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var customer models.Customer
	customer.ID = uint64(id)
	tx, err := u.Db.Begin()
	if err != nil {
		u.Log.Printf("Begin tx : %v", err)
		api.ResponseError(w, err)
		return
	}

	err = customer.View(ctx, tx)
	if err != nil {
		u.Log.Printf("Get customer: %v", err)
		api.ResponseError(w, err)
		return
	}

	var response response.CustomerResponse
	response.Transform(&customer)
	api.ResponseOK(w, response, http.StatusOK)
}

// Update customer by id
func (u *Customers) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := ctx.Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Printf("casting paramID : %v", err)
		api.ResponseError(w, err)
		return
	}

	var customerRequest request.CustomerRequest
	err = api.Decode(r, &customerRequest)
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

	var customer models.Customer
	customer.ID = uint64(id)
	err = customer.View(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Get customer: %v", err)
		api.ResponseError(w, err)
		return
	}

	customerUpdate := customerRequest.Transform(&customer)

	err = customerUpdate.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Update customer: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	var response response.CustomerResponse
	response.Transform(customerUpdate)
	api.ResponseOK(w, response, http.StatusOK)
}

// Delete customer by id
func (u *Customers) Delete(w http.ResponseWriter, r *http.Request) {
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

	var customer models.Customer
	customer.ID = uint64(id)
	err = customer.View(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Get customer: %v", err)
		api.ResponseError(w, err)
		return
	}

	err = customer.Delete(ctx, tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("Update customer: %v", err)
		api.ResponseError(w, err)
		return
	}
	tx.Commit()

	api.ResponseOK(w, nil, http.StatusNoContent)
}

```

## Routing
- Open routing file on routing/route.go
- Add customer routing
```
// Customers Routing
{
    customers := controllers.Customers{Db: db, Log: log}
    app.Handle(http.MethodGet, "/customers", customers.List)
    app.Handle(http.MethodPost, "/customers", customers.Create)
    app.Handle(http.MethodGet, "/customers/:id", customers.View)
    app.Handle(http.MethodPut, "/customers/:id", customers.Update)
    app.Handle(http.MethodDelete, "/customers/:id", customers.Delete)
}
```