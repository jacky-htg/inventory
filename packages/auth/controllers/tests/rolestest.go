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

//Roles : struct for set Roles Dependency Injection
type Roles struct {
	App   http.Handler
	Token string
}

//List : http handler for returning list of roles
func (u *Roles) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/roles", nil)
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
				"id":   float64(1),
				"name": string("superadmin"),
			},
		},
	}

	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

//Crud : http handler for roles crud
func (u *Roles) Crud(t *testing.T) {
	created := u.Create(t)
	id := created["data"].(map[string]interface{})["id"].(float64)
	u.View(t, id)
	u.Update(t, id)
	u.Grant(t, id, 1)
	u.Revoke(t, id, 1)
	u.Delete(t, id)
}

//Create : http handler for create new role
func (u *Roles) Create(t *testing.T) map[string]interface{} {
	var created map[string]interface{}
	jsonBody := `
		{
			"name": "Tes"
		}
	`
	body := strings.NewReader(jsonBody)

	req := httptest.NewRequest("POST", "/roles", body)
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
			"id":   c["id"],
			"name": "Tes",
		},
	}

	if diff := cmp.Diff(want, created); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}

	return created
}

//View : http handler for retrieve roles by id
func (u *Roles) View(t *testing.T, id float64) {
	req := httptest.NewRequest("GET", "/roles/"+fmt.Sprintf("%d", int(id)), nil)
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
			"id":   id,
			"name": "Tes",
		},
	}

	// Fetched product should match the one we created.
	if diff := cmp.Diff(want, fetched); diff != "" {
		t.Fatalf("Retrieved user should match created. Diff:\n%s", diff)
	}
}

//Update : http handler for update roles by id
func (u *Roles) Update(t *testing.T, id float64) {
	var updated map[string]interface{}
	jsonBody := `
		{
			"id": %s,
			"name": "Test"
		}
	`
	body := strings.NewReader(fmt.Sprintf(jsonBody, fmt.Sprintf("%d", int(id))))

	req := httptest.NewRequest("PUT", "/roles/"+fmt.Sprintf("%d", int(id)), body)
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
			"id":   id,
			"name": "Test",
		},
	}

	if diff := cmp.Diff(want, updated); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

//Grant : http handler for grant access
func (u *Roles) Grant(t *testing.T, id float64, accessID float64) {
	req := httptest.NewRequest("POST", "/roles/"+fmt.Sprintf("%d", int(id))+"/access/"+fmt.Sprintf("%d", int(accessID)), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", u.Token)
	resp := httptest.NewRecorder()

	u.App.ServeHTTP(resp, req)

	if http.StatusOK != resp.Code {
		t.Fatalf("retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
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

//Revoke : http handler for revoke access
func (u *Roles) Revoke(t *testing.T, id float64, accessID float64) {
	req := httptest.NewRequest("DELETE", "/roles/"+fmt.Sprintf("%d", int(id))+"/access/"+fmt.Sprintf("%d", int(accessID)), nil)
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

//Delete : http handler for delete roles by id
func (u *Roles) Delete(t *testing.T, id float64) {
	req := httptest.NewRequest("DELETE", "/roles/"+fmt.Sprintf("%d", int(id)), nil)
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
