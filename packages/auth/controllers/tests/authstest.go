package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//Auths : struct for set AUths Dependency Injection
type Auths struct {
	App   http.Handler
	Token string
}

//List : http handler for returning list of users
func (u *Auths) Login(t *testing.T) {
	jsonBody := `
		{
			"username": "jackyhtg", 
			"password": "12345678"
		}
	`
	body := strings.NewReader(jsonBody)

	req := httptest.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	u.App.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	var list map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	u.Token = list["data"].(map[string]interface{})["token"].(string)

}
