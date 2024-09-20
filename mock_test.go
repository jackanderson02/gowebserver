package main

import (
	"acme/api"
	"acme/repository/mock"
	// "acme/config"
	"acme/model"
	"acme/service"
	// "bytes"
	"encoding/json"
	// "fmt"
	// "io"
	"net/http"
	"net/http/httptest"
	// "strconv"
	"reflect"
	"testing"
)
type MockEndToEnd struct{
	a *api.UserAPI
	s *service.UserService
	m *mock.MockRepository
}

func setupMock(r *mock.MockRepository) *MockEndToEnd {
	userService := service.NewUserService(r)
	userAPI := api.NewUserAPI(userService)
	return &MockEndToEnd{
		s: userService,
		a: userAPI,
		m: r,
	}

}
func TestGetUsersHandlerWithMock(t *testing.T) {
    req, err := http.NewRequest("GET", "/api/users", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    expected := []model.User{
        {ID: 1, Name: "Alice"},
        {ID: 2, Name: "Bob"},
        {ID: 3, Name: "Terry"},
    }
    mock := setupMock(&mock.MockRepository{
		MockGetUsers: func() ([]model.User, error) {
			return expected, nil
		},
	})
    
    handler := http.HandlerFunc(mock.a.GetUsers)
    handler.ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
	var actual []model.User
    if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
        t.Fatalf("Failed to unmarshal response body: %v", err)
    }
    if !reflect.DeepEqual(actual, expected) {
        t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
    }
}


