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

// Brands : : struct for set Brands Dependency Injection
type Brands struct {
	App   http.Handler
	Token string
}

// Run : http handler for run brands testing
func (u *Brands) Run(t *testing.T) {
	created := u.Create(t)
	id := created["data"].(map[string]interface{})["id"].(float64)
	u.List(t)
	u.View(t, id)
	u.Update(t, id)
	u.Delete(t, id)
}

// List : http handler for returning list of brands
func (u *Brands) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/brands", nil)
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
				"code": string("BRAND-1"),
				"name": "Tes",
				"company": map[string]interface{}{
					"id":      float64(1),
					"code":    "DM",
					"name":    "Dummy",
					"address": "",
				},
			},
		},
	}

	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

// Create : http handler for create new brand
func (u *Brands) Create(t *testing.T) map[string]interface{} {
	var created map[string]interface{}
	jsonBody := `
		{
			"code": "BRAND-1",
			"name": "Tes"
		}
	`
	body := strings.NewReader(jsonBody)

	req := httptest.NewRequest("POST", "/brands", body)
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
		t.Fatal("expected non-empty brand id")
	}

	want := map[string]interface{}{
		"status_code":    "REBEL-200",
		"status_message": "OK",
		"data": map[string]interface{}{
			"id":   c["id"],
			"code": "BRAND-1",
			"name": "Tes",
			"company": map[string]interface{}{
				"id":      float64(1),
				"code":    "DM",
				"name":    "Dummy",
				"address": "",
			},
		},
	}

	if diff := cmp.Diff(want, created); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}

	return created
}

// View : http handler for retrieve brand by id
func (u *Brands) View(t *testing.T, id float64) {
	req := httptest.NewRequest("GET", "/brands/"+fmt.Sprintf("%d", int(id)), nil)
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
			"code": "BRAND-1",
			"name": "Tes",
			"company": map[string]interface{}{
				"id":      float64(1),
				"code":    "DM",
				"name":    "Dummy",
				"address": "",
			},
		},
	}

	// Fetched brand should match the one we created.
	if diff := cmp.Diff(want, fetched); diff != "" {
		t.Fatalf("Retrieved brand should match created. Diff:\n%s", diff)
	}
}

// Update : http handler for update brand by id
func (u *Brands) Update(t *testing.T, id float64) {
	var updated map[string]interface{}
	jsonBody := `
		{
			"id": %s,
			"name": "Test"
		}
	`
	body := strings.NewReader(fmt.Sprintf(jsonBody, fmt.Sprintf("%d", int(id))))

	req := httptest.NewRequest("PUT", "/brands/"+fmt.Sprintf("%d", int(id)), body)
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
			"code": "BRAND-1",
			"name": "Test",
			"company": map[string]interface{}{
				"id":      float64(1),
				"code":    "DM",
				"name":    "Dummy",
				"address": "",
			},
		},
	}

	if diff := cmp.Diff(want, updated); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

// Delete brand
func (u *Brands) Delete(t *testing.T, id float64) {
	req := httptest.NewRequest("DELETE", "/brands/"+fmt.Sprintf("%d", int(id)), nil)
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

	// Fetched brand should match the one we created.
	if diff := cmp.Diff(want, deleted); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}
