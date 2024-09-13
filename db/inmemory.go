package db

import (
	"errors"
	"fmt"
	"slices"
	"acme/model"

)
var users []model.User

func init() {
    // Initialize the in-memory database with some sample data
    users = []model.User{
        {ID: 1, Name: "Alice"},
        {ID: 2, Name: "Bob"},
    }
}

func GetUsers() ([]model.User, error){
	return users, nil
}


func AddUser(user model.User) (int, error){
	lastID := users[len(users)-1].ID
	newID := lastID + 1
	user.ID = newID
	users = append(users, user)

	return newID, nil
}

func GetUserByID(id int ) (model.User, error){
	fmt.Println("Getting user:")
	fmt.Println(id)
	for _, user := range users{
		if user.ID == id{
			return user, nil
		}
	
	}

	fmt.Println("returning empty user")
	return model.User{}, nil

}

func DeleteUserByID(id int) error{
	fmt.Println("Attempting to delete user with ID:")
	fmt.Println(id)
	for idx, user := range users{
		if user.ID == id{
			users = slices.Delete(users, idx, idx+1)
			return nil
		}
	}
	return errors.New("No user exists with this id")
}

func UpdateUsernameByID(id int, newUser model.User) error{
	// Create deep copy of users
	userCopy := make([]model.User, len(users))
	copy(userCopy, users)

	for idx, user := range userCopy{
		if user.ID == id{
			users[idx].Name = newUser.Name
			return nil
		}
	}

	return errors.New("No user exists with this id")
}