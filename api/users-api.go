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

type UserAPI struct {
    userService *service.UserService
}

func NewUserAPI(userService *service.UserService) *UserAPI {
    return &UserAPI{
        userService: userService,
    }
}

func (api *UserAPI) decodeUser(body io.ReadCloser) (user model.User, err error) {

	err = json.NewDecoder(body).Decode(&user)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		return model.User{}, err
	}

	return user, nil
}


func (api *UserAPI) parseID(request *http.Request) (int, error){
	// Should just be able to extract from fields
	// idStr := request.PathValue("id")
	idStr := request.URL.Path
	idStr = idStr[len(idStr)-1:]
	id, err:= strconv.Atoi(idStr)
	return id, err

}
func (api *UserAPI) CreateUser(writer http.ResponseWriter, request *http.Request) {
	service := api.userService
	var user model.User
	user, err := api.decodeUser(request.Body)
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

func (api *UserAPI) GetSingleUser(writer http.ResponseWriter, request *http.Request) {
	service := api.userService
	id, err := api.parseID(request)
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
func (api *UserAPI) GetUsers(writer http.ResponseWriter, _ *http.Request) {
	service := api.userService

	users, err := service.GetUsersService()
	json.NewEncoder(writer).Encode(users)

	if (err != nil){
		http.Error(writer, "Failed to retrieve users.", http.StatusInternalServerError)
	}

}

func (api *UserAPI) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	service := api.userService

	id, err := api.parseID(request)
	if err != nil {
		http.Error(writer, "Non-integer user number provided", http.StatusInternalServerError)
	}
	err = service.DeleteUserService(id)
	if err != nil {
		http.Error(writer, "Failed to delete user", http.StatusInternalServerError)
	}
}

func (api *UserAPI) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	service := api.userService

	id, err := api.parseID(request)
	if err != nil {
		http.Error(writer, "Non-integer user number provided", http.StatusInternalServerError)
	}

	var user model.User
	user, err = api.decodeUser(request.Body)

	if err != nil {
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	err = service.UpdateUsernameService(id, user)
	if err != nil {
		http.Error(writer, "Failed to update user", http.StatusInternalServerError)
	}
}