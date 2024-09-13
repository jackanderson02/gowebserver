package api

import (
	"acme/service"
	"acme/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"io"
)

func decodeUser(body io.ReadCloser) (user model.User, err error) {

	err = json.NewDecoder(body).Decode(&user)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		return model.User{}, err
	}

	return user, nil
}


func parseID(request *http.Request) (int, error){
	// Should just be able to extract from fields
	// idStr := request.PathValue("id")
	idStr := request.URL.Path
	idStr = idStr[len(idStr)-1:]
	id, err:= strconv.Atoi(idStr)
	return id, err

}
func CreateUser(writer http.ResponseWriter, request *http.Request) {
	var user model.User
	user, err := decodeUser(request.Body)
	if err != nil {
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	id, err := service.CreateUserService(user)

	if err != nil{
		http.Error(writer, "Error creating user.", http.StatusInternalServerError)
	}
	writer.WriteHeader(http.StatusCreated)
	fmt.Fprintf(writer, "User created successfully: %d", id)

}

func GetSingleUser(writer http.ResponseWriter, request *http.Request) {
	id, err := parseID(request)
	if err != nil {
		http.Error(writer, "Non-integer user number provided", http.StatusInternalServerError)
	}
	fmt.Printf("got /api/users/{%d} request\n", id)
	user, err := service.GetUserByIDService(id)

	if err != nil{
		http.Error(writer, "Error retrieving user.",  http.StatusInternalServerError)
	}

	userJSON, errMarshal := json.Marshal(user)
	if errMarshal != nil {
		// Handle error if marshalling fails
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	//3. Finally, we write the JSON response string directly to the
	_, err = writer.Write(userJSON)
	if err != nil {
		// Handle error if writing response fails
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
func GetUsers(writer http.ResponseWriter, _ *http.Request) {

	users := service.GetUsersService()
	json.NewEncoder(writer).Encode(users)

}

func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	id, err := parseID(request)
	if err != nil {
		http.Error(writer, "Non-integer user number provided", http.StatusInternalServerError)
	}
	err = service.DeleteUserService(id)
	if err != nil {
		http.Error(writer, "Failed to delete user", http.StatusInternalServerError)
	}
}

func UpdateUser(writer http.ResponseWriter, request *http.Request) {
	id, err := parseID(request)
	if err != nil {
		http.Error(writer, "Non-integer user number provided", http.StatusInternalServerError)
	}

	var user model.User
	user, err = decodeUser(request.Body)

	if err != nil {
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	err = service.UpdateUsernameService(id, user)
	if err != nil {
		http.Error(writer, "Failed to update user", http.StatusInternalServerError)
	}
}