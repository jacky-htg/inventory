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

// ProductCategories : struct for set ProductCategories Dependency Injection
type ProductCategories struct {
	App   http.Handler
	Token string
}

// Run : http handler for run ProductCategories testing
func (u *ProductCategories) Run(t *testing.T) {
	created := u.Create(t)
	id := created["data"].(map[string]interface{})["id"].(float64)
	u.List(t)
	u.View(t, id)
	u.Update(t, id)
	u.Delete(t, id)
}

// List : http handler for returning list of ProductCategories
func (u *ProductCategories) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/product_categories", nil)
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
				"name": "Lemari",
				"company": map[string]interface{}{
					"id":      float64(1),
					"code":    "DM",
					"name":    "Dummy",
					"address": "",
				},
				"category": map[string]interface{}{
					"id":   float64(1),
					"name": "Accesories",
				},
			},
		},
	}

	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

// Create : http handler for create new product category
func (u *ProductCategories) Create(t *testing.T) map[string]interface{} {
	var created map[string]interface{}
	jsonBody := `
		{
			"name": "Lemari",
			"category": 1
		}
	`
	body := strings.NewReader(jsonBody)

	req := httptest.NewRequest("POST", "/product_categories", body)
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
		t.Fatal("expected non-empty product category id")
	}

	want := map[string]interface{}{
		"status_code":    "REBEL-200",
		"status_message": "OK",
		"data": map[string]interface{}{
			"id":   c["id"],
			"name": "Lemari",
			"company": map[string]interface{}{
				"id":      float64(1),
				"code":    "DM",
				"name":    "Dummy",
				"address": "",
			},
			"category": map[string]interface{}{
				"id":   float64(1),
				"name": "Accesories",
			},
		},
	}

	if diff := cmp.Diff(want, created); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}

	return created
}

// View : http handler for retrieve product category by id
func (u *ProductCategories) View(t *testing.T, id float64) {
	req := httptest.NewRequest("GET", "/product_categories/"+fmt.Sprintf("%d", int(id)), nil)
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
			"name": "Lemari",
			"company": map[string]interface{}{
				"id":      float64(1),
				"code":    "DM",
				"name":    "Dummy",
				"address": "",
			},
			"category": map[string]interface{}{
				"id":   float64(1),
				"name": "Accesories",
			},
		},
	}

	// Fetched product category should match the one we created.
	if diff := cmp.Diff(want, fetched); diff != "" {
		t.Fatalf("Retrieved product category should match created. Diff:\n%s", diff)
	}
}

// Update : http handler for update product category by id
func (u *ProductCategories) Update(t *testing.T, id float64) {
	var updated map[string]interface{}
	jsonBody := `
		{
			"id": %s,
			"name": "Bed"
		}
	`
	body := strings.NewReader(fmt.Sprintf(jsonBody, fmt.Sprintf("%d", int(id))))

	req := httptest.NewRequest("PUT", "/product_categories/"+fmt.Sprintf("%d", int(id)), body)
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
			"name": "Bed",
			"company": map[string]interface{}{
				"id":      float64(1),
				"code":    "DM",
				"name":    "Dummy",
				"address": "",
			},
			"category": map[string]interface{}{
				"id":   float64(1),
				"name": "Accesories",
			},
		},
	}

	if diff := cmp.Diff(want, updated); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

// Delete product category
func (u *ProductCategories) Delete(t *testing.T, id float64) {
	req := httptest.NewRequest("DELETE", "/product_categories/"+fmt.Sprintf("%d", int(id)), nil)
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

	// Fetched product category should match the one we created.
	if diff := cmp.Diff(want, deleted); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}
