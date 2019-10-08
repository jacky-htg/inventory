package main

import (
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	apiTest "github.com/jacky-htg/inventory/controllers/tests"
	"github.com/jacky-htg/inventory/libraries/config"
	"github.com/jacky-htg/inventory/routing"
	"github.com/jacky-htg/inventory/schema"
	"github.com/jacky-htg/inventory/tests"
)

var token string

func TestMain(t *testing.T) {
	_, ok := os.LookupEnv("APP_ENV")
	if !ok {
		config.Setup(".env")
	}

	db, teardown := tests.NewUnit(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	log := log.New(os.Stderr, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// api test for auths
	{
		auths := apiTest.Auths{App: routing.API(db, log)}
		t.Run("ApiLogin", auths.Login)
		token = auths.Token
	}

	// api test for users
	{
		users := apiTest.Users{App: routing.API(db, log), Token: token}
		t.Run("APiUsersList", users.List)
		t.Run("APiUsersCrud", users.Crud)
	}

	// api test for access
	{
		access := apiTest.Access{App: routing.API(db, log), Token: token}
		t.Run("APiAccessList", access.List)
	}

	// api test for roles
	{
		roles := apiTest.Roles{App: routing.API(db, log), Token: token}
		t.Run("APiRolesList", roles.List)
		t.Run("APiRolesCrud", roles.Crud)
	}

	// api test for regions
	{
		regions := apiTest.Regions{App: routing.API(db, log), Token: token}
		t.Run("APiRegionsTest", regions.Run)
	}

	// api test for brands
	{
		brands := apiTest.Brands{App: routing.API(db, log), Token: token}
		t.Run("APiBrandsCrud", brands.Run)
	}

	// api test for product categories
	{
		productCategories := apiTest.ProductCategories{App: routing.API(db, log), Token: token}
		t.Run("APiProductCategoriesCrud", productCategories.Run)
	}

	// api test for products
	{
		products := apiTest.Products{App: routing.API(db, log), Token: token}
		t.Run("APiProductsCrud", products.Run)
	}
}
