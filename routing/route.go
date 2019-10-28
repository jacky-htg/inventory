package routing

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/jacky-htg/inventory/controllers"
	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/middleware"
)

//API : hanlder api
func API(db *sql.DB, log *log.Logger) http.Handler {
	app := api.NewApp(
		log,
		middleware.Auths(db, log, []string{"/login", "/health"}),
	)

	// Health Routing
	{
		check := controllers.Checks{Db: db}
		app.Handle(http.MethodGet, "/health", check.Health)
	}

	// Auth Routing
	{
		auth := controllers.Auths{Db: db, Log: log}
		app.Handle(http.MethodPost, "/login", auth.Login)
	}

	// Companies Routing
	{
		companies := controllers.Companies{Db: db, Log: log}
		app.Handle(http.MethodGet, "/companies/:id", companies.View)
		//app.Handle(http.MethodPost, "/companies", companies.Create)
		app.Handle(http.MethodPut, "/companies/:id", companies.Update)
		app.Handle(http.MethodDelete, "/companies/:id", companies.Delete)
	}

	// Users Routing
	{
		user := controllers.Users{Db: db, Log: log}
		app.Handle(http.MethodGet, "/users", user.List)
		app.Handle(http.MethodGet, "/users/:id", user.View)
		app.Handle(http.MethodPost, "/users", user.Create)
		app.Handle(http.MethodPut, "/users/:id", user.Update)
		app.Handle(http.MethodDelete, "/users/:id", user.Delete)
	}

	// Roles Routing
	{
		roles := controllers.Roles{Db: db, Log: log}
		app.Handle(http.MethodGet, "/roles", roles.List)
		app.Handle(http.MethodGet, "/roles/:id", roles.View)
		app.Handle(http.MethodPost, "/roles", roles.Create)
		app.Handle(http.MethodPut, "/roles/:id", roles.Update)
		app.Handle(http.MethodDelete, "/roles/:id", roles.Delete)
		app.Handle(http.MethodPost, "/roles/:id/access/:access_id", roles.Grant)
		app.Handle(http.MethodDelete, "/roles/:id/access/:access_id", roles.Revoke)
	}

	// Access Routing
	{
		access := controllers.Access{Db: db, Log: log}
		app.Handle(http.MethodGet, "/access", access.List)
	}

	// Regions Routing
	{
		regions := controllers.Regions{Db: db, Log: log}
		app.Handle(http.MethodGet, "/regions", regions.List)
		app.Handle(http.MethodGet, "/regions/:id", regions.View)
		app.Handle(http.MethodPost, "/regions", regions.Create)
		app.Handle(http.MethodPut, "/regions/:id", regions.Update)
		app.Handle(http.MethodPost, "/regions/:id/branches/:branch_id", regions.AddBranch)
		app.Handle(http.MethodDelete, "/regions/:id/branches/:branch_id", regions.DeleteBranch)
		app.Handle(http.MethodDelete, "/regions/:id", regions.Delete)
	}

	// Products Routing
	{
		products := controllers.Products{Db: db, Log: log}
		app.Handle(http.MethodGet, "/products", products.List)
		app.Handle(http.MethodGet, "/products/:id", products.View)
		app.Handle(http.MethodPost, "/products", products.Create)
		app.Handle(http.MethodPut, "/products/:id", products.Update)
		app.Handle(http.MethodDelete, "/products/:id", products.Delete)
	}

	// Purchases Routing
	{
		purchases := controllers.Purchases{Db: db, Log: log}
		app.Handle(http.MethodGet, "/purchases", purchases.List)
		app.Handle(http.MethodGet, "/purchases/:id", purchases.View)
		app.Handle(http.MethodPost, "/purchases", purchases.Create)
		app.Handle(http.MethodPut, "/purchases/:id", purchases.Update)
	}

	// Purchase Returns Routing
	{
		purchaseReturns := controllers.PurchaseReturns{Db: db, Log: log}
		app.Handle(http.MethodGet, "/purchase_returns", purchaseReturns.List)
		app.Handle(http.MethodGet, "/purchase_returns/:id", purchaseReturns.View)
		app.Handle(http.MethodPost, "/purchase_returns", purchaseReturns.Create)
		app.Handle(http.MethodPut, "/purchase_returns/:id", purchaseReturns.Update)
	}

	// Closing Stock Routing
	{
		closingStock := controllers.ClosingStocks{Db: db, Log: log}
		app.Handle(http.MethodPost, "/closing_stocks", closingStock.Closing)
	}

	// Customers Routing
	{
		customers := controllers.Customers{Db: db, Log: log}
		app.Handle(http.MethodGet, "/customers", customers.List)
		app.Handle(http.MethodPost, "/customers", customers.Create)
		app.Handle(http.MethodGet, "/customers/:id", customers.View)
		app.Handle(http.MethodPut, "/customers/:id", customers.Update)
		app.Handle(http.MethodDelete, "/customers/:id", customers.Delete)
	}

	// Branches Routing
	{
		branches := controllers.Branches{Db: db, Log: log}
		app.Handle(http.MethodGet, "/branches", branches.List)
		app.Handle(http.MethodPost, "/branches", branches.Create)
		app.Handle(http.MethodGet, "/branches/:id", branches.View)
		app.Handle(http.MethodPut, "/branches/:id", branches.Update)
		app.Handle(http.MethodDelete, "/branches/:id", branches.Delete)
	}

	// Brands Routing
	{
		brands := controllers.Brands{Db: db, Log: log}
		app.Handle(http.MethodGet, "/brands", brands.List)
		app.Handle(http.MethodPost, "/brands", brands.Create)
		app.Handle(http.MethodGet, "/brands/:id", brands.View)
		app.Handle(http.MethodPut, "/brands/:id", brands.Update)
		app.Handle(http.MethodDelete, "/brands/:id", brands.Delete)
	}

	// ProductCategories Routing
	{
		productCategories := controllers.ProductCategories{Db: db, Log: log}
		app.Handle(http.MethodGet, "/product_categories", productCategories.List)
		app.Handle(http.MethodPost, "/product_categories", productCategories.Create)
		app.Handle(http.MethodGet, "/product_categories/:id", productCategories.View)
		app.Handle(http.MethodPut, "/product_categories/:id", productCategories.Update)
		app.Handle(http.MethodDelete, "/product_categories/:id", productCategories.Delete)
	}

	// Receives Routing
	{
		receives := controllers.Receives{Db: db, Log: log}
		app.Handle(http.MethodGet, "/receives", receives.List)
		app.Handle(http.MethodGet, "/receives/:id", receives.View)
		app.Handle(http.MethodPost, "/receives", receives.Create)
		app.Handle(http.MethodPut, "/receives/:id", receives.Update)
	}

	// Shelves Routing
	{
		shelves := controllers.Shelves{Db: db, Log: log}
		app.Handle(http.MethodGet, "/shelves", shelves.List)
		app.Handle(http.MethodPost, "/shelves", shelves.Create)
		app.Handle(http.MethodGet, "/shelves/:id", shelves.View)
		app.Handle(http.MethodPut, "/shelves/:id", shelves.Update)
		app.Handle(http.MethodDelete, "/shelves/:id", shelves.Delete)
	}

	return app
}
