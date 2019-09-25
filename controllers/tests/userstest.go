package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//Users : struct for set Users Dependency Injection
type Users struct {
	App   http.Handler
	Token string
}

//List : http handler for returning list of users
func (u *Users) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/users", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", u.Token)
	resp := httptest.NewRecorder()

	u.App.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	var list map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	want := map[string]interface{}{
		"status_code":    string("REBEL-200"),
		"status_message": string("OK"),
		"data": []interface{}{
			map[string]interface{}{
				"id":        float64(1),
				"username":  string("jackyhtg"),
				"is_active": bool(true),
				"roles":     []interface{}{map[string]interface{}{"ID": float64(1), "Name": string("superadmin")}},
			},
		},
	}

	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

//Crud : http handler for users crud
func (u *Users) Crud(t *testing.T) {
	created := u.Create(t)
	id := created["data"].(map[string]interface{})["id"].(float64)
	u.View(t, id)
	u.Update(t, id)
	u.Delete(t, id)
}

//Create : http handler for create new user
func (u *Users) Create(t *testing.T) map[string]interface{} {
	var created map[string]interface{}
	jsonBody := `
		{
			"username": "peterpan",
			"email": "peterpan@gmail.com", 
			"password": "1234", 
			"re_password": "1234", 
			"is_active": true
		}
	`
	body := strings.NewReader(jsonBody)

	req := httptest.NewRequest("POST", "/users", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", u.Token)
	resp := httptest.NewRecorder()

	u.App.ServeHTTP(resp, req)

	if http.StatusCreated != resp.Code {
		t.Fatalf("posting: expected status code %v, got %v", http.StatusCreated, resp.Code)
	}

	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	c := created["data"].(map[string]interface{})

	if c["id"] == "" || c["id"] == nil {
		t.Fatal("expected non-empty product id")
	}

	want := map[string]interface{}{
		"status_code":    "REBEL-200",
		"status_message": "OK",
		"data": map[string]interface{}{
			"id":        c["id"],
			"username":  "peterpan",
			"is_active": true,
			"roles":     nil,
		},
	}

	if diff := cmp.Diff(want, created); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}

	return created
}

//View : http handler for retrieve user by id
func (u *Users) View(t *testing.T, id float64) {
	req := httptest.NewRequest("GET", "/users/"+fmt.Sprintf("%d", int(id)), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", u.Token)
	resp := httptest.NewRecorder()

	u.App.ServeHTTP(resp, req)

	if http.StatusOK != resp.Code {
		t.Fatalf("retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	var fetched map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	want := map[string]interface{}{
		"status_code":    "REBEL-200",
		"status_message": "OK",
		"data": map[string]interface{}{
			"id":        id,
			"username":  "peterpan",
			"is_active": true,
			"roles":     []interface{}{map[string]interface{}{"ID": float64(0), "Name": string("")}},
		},
	}

	// Fetched product should match the one we created.
	if diff := cmp.Diff(want, fetched); diff != "" {
		t.Fatalf("Retrieved user should match created. Diff:\n%s", diff)
	}
}

//Update : http handler for update user by id
func (u *Users) Update(t *testing.T, id float64) {
	var updated map[string]interface{}
	jsonBody := `
		{
			"id": %s,
			"username": "gatholoco",
			"email": "gatholoco@gmail.com", 
			"password": "1234", 
			"re_password": "1234", 
			"is_active": false
		}
	`
	body := strings.NewReader(fmt.Sprintf(jsonBody, fmt.Sprintf("%d", int(id))))

	req := httptest.NewRequest("PUT", "/users/"+fmt.Sprintf("%d", int(id)), body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", u.Token)
	resp := httptest.NewRecorder()

	u.App.ServeHTTP(resp, req)

	if http.StatusOK != resp.Code {
		t.Fatalf("posting: expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	want := map[string]interface{}{
		"status_code":    "REBEL-200",
		"status_message": "OK",
		"data": map[string]interface{}{
			"id":        id,
			"username":  "gatholoco",
			"is_active": false,
			"roles":     []interface{}{map[string]interface{}{"ID": float64(0), "Name": string("")}},
		},
	}

	if diff := cmp.Diff(want, updated); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

//Delete : http handler for delete user by id
func (u *Users) Delete(t *testing.T, id float64) {
	req := httptest.NewRequest("DELETE", "/users/"+fmt.Sprintf("%d", int(id)), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", u.Token)
	resp := httptest.NewRecorder()

	u.App.ServeHTTP(resp, req)

	if http.StatusNoContent != resp.Code {
		t.Fatalf("retrieving: expected status code %v, got %v", http.StatusNoContent, resp.Code)
	}

	var deleted map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&deleted); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	want := map[string]interface{}{
		"status_code":    "REBEL-200",
		"status_message": "OK",
		"data":           nil,
	}

	// Fetched product should match the one we created.
	if diff := cmp.Diff(want, deleted); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}
