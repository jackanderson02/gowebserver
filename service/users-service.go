package service

import (
	"acme/db"
	"acme/model"
	"fmt"
	"errors"
)


func GetUsersService() []model.User {
	ret, _ := db.GetUsers()
	return ret
}


func CreateUserService(user model.User) (int, error){
	id, err := db.AddUser(user)

	if err != nil{
		fmt.Println("Error creating new use:", err)
		return 0, errors.New("Error retrieving user.")
	}

	return id, nil


}

func DeleteUserService(id int) error {
	err := db.DeleteUserByID(id)

	if err != nil{
		fmt.Println("Error deleting user from DB:", err)
		return errors.New("Could not delete user")
	}

	return nil
}


func UpdateUsernameService(id int, user model.User) error {
	err := db.UpdateUsernameByID(id, user)

	if err != nil{
		fmt.Println("Error updating user in DB:", err)
		return errors.New("Could not update username")
	}

	return nil

}


func GetUserByIDService(id int ) (model.User, error){
	user, err := db.GetUserByID(id)


	if err != nil{
		fmt.Println("Error retrieving user by ID.", err)
		return model.User{}, errors.New("Could not Find user with that ID.")
	}

	return user, nil

}