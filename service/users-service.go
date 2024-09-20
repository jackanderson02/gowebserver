package service

import (
	"acme/model"
	"acme/repository/user"
	"fmt"
	"errors"
)

type UserService struct {
	repository user.Repository
}

func NewUserService(repo user.Repository) *UserService {
	return &UserService{
		repository: repo,
	}
}

func (service *UserService) GetUsersService() ([]model.User, error) {
	repo := service.repository
	users, err := repo.GetUsers() //used to be db.GetUsers()

	if err != nil {
		fmt.Println("Error getting users from DB:", err)
		return nil, errors.New("There was an error getting the users from the database.")
	}

	return users, nil
}


func (service *UserService) CreateUserService(user model.User) (int, error){
	repo := service.repository
	id, err := repo.AddUser(user)

	if err != nil{
		fmt.Println("Error creating new use:", err)
		return 0, errors.New("Error retrieving user.")
	}

	return id, nil


}

func (service *UserService) DeleteUserService(id int) error {
	repo := service.repository
	err := repo.DeleteUser(id)

	if err != nil{
		fmt.Println("Error deleting user from DB:", err)
		return errors.New("Could not delete user")
	}

	return nil
}


func (service *UserService) UpdateUsernameService(id int, user model.User) error {
	repo := service.repository
	err := repo.UpdateUser(id, user)

	if err != nil{
		fmt.Println("Error updating user in DB:", err)
		return errors.New("Could not update username")
	}

	return nil

}


func (service *UserService) GetUserByIDService(id int ) (model.User, error){
	repo := service.repository
	user, err := repo.GetUser(id)

	if err != nil{
		fmt.Println("Error retrieving user by ID.", err)
		return model.User{}, errors.New("Could not Find user with that ID.")
	}

	return user, nil

}