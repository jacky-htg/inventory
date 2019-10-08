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

// Products : struct for set Products Dependency Injection
type Products struct {
	App   http.Handler
	Token string
}

// Run : http handler for run products testing
func (u *Products) Run(t *testing.T) {
	created := u.Create(t)
	id := created["data"].(map[string]interface{})["id"].(float64)
	u.List(t)
	u.View(t, id)
	u.Update(t, id)
	u.Delete(t, id)
}

// List : http handler for returning list of products
func (u *Products) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/products", nil)
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
				"id":            float64(1),
				"code":          string("PROD-1"),
				"name":          "Tes",
				"price":         float64(1),
				"minimum_stock": float64(25),
				"company": map[string]interface{}{
					"id":      float64(1),
					"code":    "DM",
					"name":    "Dummy",
					"address": "",
				},
				"brand": map[string]interface{}{
					"id":   float64(1),
					"code": "BRAND-01",
					"name": "TOP",
				},
				"product_category": map[string]interface{}{
					"id":   float64(1),
					"name": "Lemari",
				},
			},
		},
	}

	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

// Create : http handler for create new product
func (u *Products) Create(t *testing.T) map[string]interface{} {
	var created map[string]interface{}
	jsonBody := `
		{
			"code": "PROD-200",
			"name": "Tes",
			"price": 1,
			"minimum_stock" : "25",
			"brand": "1",
			"product_category": "1"
		}
	`
	body := strings.NewReader(jsonBody)

	req := httptest.NewRequest("POST", "/products", body)
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
			"id":            c["id"],
			"code":          "PROD-200",
			"name":          "Tes",
			"price":         float64(1),
			"minimum_stock": float64(25),
			"company": map[string]interface{}{
				"id":      float64(1),
				"code":    "DM",
				"name":    "Dummy",
				"address": "",
			},
			"brand": map[string]interface{}{
				"id":   float64(1),
				"code": "BRAND-01",
				"name": "TOP",
			},
			"product_category": map[string]interface{}{
				"id":   float64(1),
				"name": "Lemari",
			},
		},
	}

	if diff := cmp.Diff(want, created); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}

	return created
}

// View : http handler for retrieve product by id
func (u *Products) View(t *testing.T, id float64) {
	req := httptest.NewRequest("GET", "/products/"+fmt.Sprintf("%d", int(id)), nil)
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
			"id":            id,
			"code":          "PROD-200",
			"name":          "Tes",
			"price":         float64(1),
			"minimum_stock": float64(25),
			"company": map[string]interface{}{
				"id":      float64(1),
				"code":    "DM",
				"name":    "Dummy",
				"address": "",
			},
			"brand": map[string]interface{}{
				"id":   float64(1),
				"code": "BRAND-01",
				"name": "TOP",
			},
			"product_category": map[string]interface{}{
				"id":   float64(1),
				"name": "Lemari",
			},
		},
	}

	// Fetched product should match the one we created.
	if diff := cmp.Diff(want, fetched); diff != "" {
		t.Fatalf("Retrieved user should match created. Diff:\n%s", diff)
	}
}

// Update : http handler for update product by id
func (u *Products) Update(t *testing.T, id float64) {
	var updated map[string]interface{}
	jsonBody := `
		{
			"id": %s,
			"name": "Test",
			"price": 2,
			"minimum_stock": "50",
			"brand":"1",
			"product_category": "1"
		}
	`
	body := strings.NewReader(fmt.Sprintf(jsonBody, fmt.Sprintf("%d", int(id))))

	req := httptest.NewRequest("PUT", "/products/"+fmt.Sprintf("%d", int(id)), body)
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
			"id":            id,
			"code":          "PROD-200",
			"name":          "Test",
			"price":         float64(2),
			"minimum_stock": float64(50),
			"company": map[string]interface{}{
				"id":      float64(1),
				"code":    "DM",
				"name":    "Dummy",
				"address": "",
			},
			"brand": map[string]interface{}{
				"id":   float64(1),
				"code": "BRAND-01",
				"name": "TOP",
			},
			"product_category": map[string]interface{}{
				"id":   float64(1),
				"name": "Lemari",
			},
		},
	}

	if diff := cmp.Diff(want, updated); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

// Delete product
func (u *Products) Delete(t *testing.T, id float64) {
	req := httptest.NewRequest("DELETE", "/products/"+fmt.Sprintf("%d", int(id)), nil)
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
