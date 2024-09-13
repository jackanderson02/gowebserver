package main

import (
	"acme/model"
	"acme/api"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestRootHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(rootHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Check the response body
    expected := "Hello, World!"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }

}


func TestRootHandlerWithServer(t *testing.T){
	server := httptest.NewServer(http.HandlerFunc(rootHandler))
    defer server.Close()

	resp, err := http.Get(server.URL + "/")
    if err != nil {
        t.Fatalf("Failed to send GET request: %v", err)
    }
    defer resp.Body.Close()

	// Check the status code
    if status := resp.StatusCode; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Check the response body
    expected := "Hello, World!"
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        t.Fatalf("Failed to read response body: %v", err)
    }
    if string(bodyBytes) != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", string(bodyBytes), expected)
    }

}

func TestGetUsersHandlerWithServer(t *testing.T) {
    // ARRANGE
	server := httptest.NewServer(http.HandlerFunc(api.GetUsers))
    defer server.Close()

	//Arrange our expected response
	expected := []model.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	// Make request GET request to server, expecting JSON response
	resp, err := http.Get(server.URL + "/api/users")
    if err != nil {
        t.Fatalf("Failed to send GET request: %v", err)
    }
    defer resp.Body.Close()

	// Check the status code
    if status := resp.StatusCode; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    bodyBytes, err := io.ReadAll(resp.Body)
	var actual []model.User
	if err := json.Unmarshal(bodyBytes, &actual); err != nil {
        t.Fatalf("Failed to unmarshal response body: %v", err)
    }

    if !reflect.DeepEqual(actual, expected) {
        t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
    }

	if(err != nil){
		t.Fatalf("Failed to read response body")
	}

}

func TestCreateUserWithServer(t *testing.T) {

    // ARRANGE
	server := httptest.NewServer(http.HandlerFunc(api.GetUsers))
    defer server.Close()

	resp, _:= http.Get(server.URL + "/api/users")

	// Check the status code
    if status := resp.StatusCode; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    bodyBytes, _:= io.ReadAll(resp.Body)
	var actual []model.User
	// Get data from model to see how many users we currently have
	json.Unmarshal(bodyBytes, &actual)
	fmt.Println(actual)

	numUsers := len(actual)

	server = httptest.NewServer(http.HandlerFunc(api.CreateUser))

	newUser, _ := json.Marshal(model.User{Name:"John"})
	resp, err := http.Post(server.URL + "/api/users", "JSON", bytes.NewReader(newUser))
    if err != nil {
        t.Fatalf("Failed to send GET request: %v", err)
    }
    defer resp.Body.Close()

	bodyBytes, err = io.ReadAll(resp.Body)

	if(err != nil){
		t.Fatalf("Failed to read response body")
	}

	resp_string := string(bodyBytes)
	numUsersNew, err := strconv.Atoi(resp_string[len(resp_string)-1:]) 
	if(err != nil){
		t.Error("A non integer number of users was returned.")
	}
	
	if numUsers == numUsersNew{
		t.Errorf("The number of users reported has not increased, previous count was %d and new count is %d.", numUsers, numUsersNew)
	}
}

func TestDeleteUserWithServer(t *testing.T){

	// First create a user
	server := httptest.NewServer(http.HandlerFunc(api.CreateUser))

	newUser, _ := json.Marshal(model.User{Name:"John"})
	resp, err := http.Post(server.URL + "/api/users", "JSON", bytes.NewReader(newUser))
    if err != nil {
        t.Fatalf("Failed to send GET request: %v", err)
    }
    defer resp.Body.Close()

	bodyBytes, _:= io.ReadAll(resp.Body)

	idStr := string(bodyBytes)
	idStr = idStr[len(idStr)-1:]
	// Then delete the user
	server = httptest.NewServer(http.HandlerFunc(api.DeleteUser))

	req, _ :=http.NewRequest("DELETE", server.URL + "/api/users/" + idStr, nil)
	server.Client().Do(req)

	// We do not expect a response, so just need to assert that the user we previously created does not exist anymore

	server = httptest.NewServer(http.HandlerFunc(api.GetSingleUser))
    defer server.Close()


	resp, _= http.Get(server.URL + "/api/users/" + idStr)

	var respJson model.User
	bodyBytes, _ = io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &respJson)
	fmt.Println(string(bodyBytes))
	// No response so no user found
	if (respJson != model.User{}){
		t.Fatalf("User with id %s still exists in the database", idStr)
	}
}


func TestUpdateUserWithServer(t *testing.T){

	// First create a user
	server := httptest.NewServer(http.HandlerFunc(api.CreateUser))

	newUser, _ := json.Marshal(model.User{Name:"John"})
	resp, err := http.Post(server.URL + "/api/users", "JSON", bytes.NewReader(newUser))
    if err != nil {
        t.Fatalf("Failed to send GET request: %v", err)
    }
    defer resp.Body.Close()


	bodyBytes, _:= io.ReadAll(resp.Body)

	idStr := string(bodyBytes)
	idStr = idStr[len(idStr)-1:]

	// Then update the name of that user
	server = httptest.NewServer(http.HandlerFunc(api.UpdateUser))
	updatedName, _ := json.Marshal(model.User{Name:"James"})

	req, _ :=http.NewRequest("PUT", server.URL + "/api/users/" + idStr, bytes.NewReader(updatedName))
	server.Client().Do(req)

	// Then get that user and assert that the names are not equal
	server = httptest.NewServer(http.HandlerFunc(api.GetSingleUser))
    defer server.Close()

	resp, _= http.Get(server.URL + "/api/users/" + idStr)

	var respJson model.User
	bodyBytes, _ = io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &respJson)

	// assert that the name of the user is now James
	fmt.Println(respJson)
	if (respJson.Name != "James"){
		t.Fatalf("Username for user with id %s was not upated", idStr)
	}

}