package tests

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/libraries/config"
	"github.com/jacky-htg/inventory/models"
	"github.com/jacky-htg/inventory/schema"
	"github.com/jacky-htg/inventory/tests"
)

// User struct for test users
type User struct {
	Db        *sql.DB
	UserLogin models.User
}

func TestMain(t *testing.T) {
	_, ok := os.LookupEnv("APP_ENV")
	if !ok {
		config.Setup("../../.env")
	}

	db, teardown := tests.NewUnit(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	userLogin := models.User{Username: "jackyhtg"}
	err := userLogin.GetByUsername(context.Background(), db)
	if err != nil {
		t.Fatal(err)
	}

	user := User{Db: db, UserLogin: userLogin}
	t.Run("List", user.List)
	t.Run("Crud", user.Crud)
}

//Crud : unit test  for create get and delete user function
func (u *User) Crud(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, api.Ctx("auth"), u.UserLogin)

	u0 := models.User{
		Username: "Aladin",
		Email:    "aladin@gmail.com",
		Password: "1234",
		IsActive: true,
	}

	tx, err := u.Db.Begin()
	if err != nil {
		t.Fatalf("begin transaction: %s", err)
	}

	err = u0.Create(ctx, tx)
	if err != nil {
		t.Fatalf("creating user u0: %s", err)
	}

	u1 := models.User{
		ID: u0.ID,
	}

	err = u1.Get(ctx, u.Db)
	if err != nil {
		t.Fatalf("getting user u1: %s", err)
	}

	if diff := cmp.Diff(u1, u0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}

	u1.IsActive = false
	err = u1.Update(ctx, tx)
	if err != nil {
		t.Fatalf("update user u1: %s", err)
	}

	u2 := models.User{
		ID: u1.ID,
	}

	err = u2.Get(ctx, u.Db)
	if err != nil {
		t.Fatalf("getting user u2: %s", err)
	}

	if diff := cmp.Diff(u1, u2); diff != "" {
		t.Fatalf("fetched != updated:\n%s", diff)
	}

	err = u2.Delete(ctx, u.Db)
	if err != nil {
		t.Fatalf("delete user u2: %s", err)
	}

	u3 := models.User{
		ID: u2.ID,
	}

	err = u3.Get(ctx, u.Db)
	if err != sql.ErrNoRows {
		t.Fatalf("getting user u3: %s", err)
	}
}

//List : unit test for user list function
func (u *User) List(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, api.Ctx("auth"), u.UserLogin)
	var user models.User
	users, err := user.List(ctx, u.Db)
	if err != nil {
		t.Fatalf("listing users: %s", err)
	}
	if exp, got := 1, len(users); exp != got {
		t.Fatalf("expected users list size %v, got %v", exp, got)
	}
}
