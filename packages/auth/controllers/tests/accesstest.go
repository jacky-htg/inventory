package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Access : struct for set Access Dependency Injection
type Access struct {
	App   http.Handler
	Token string
}

//List : http handler for returning list of access
func (u *Access) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/access", nil)
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
}
