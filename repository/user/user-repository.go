package user

import (
	"acme/model"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresUserRepository struct {
	DB *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func (repo *PostgresUserRepository) GetUsers() ([]model.User, error) {

	DB := repo.DB

	users := []model.User{}

	err := sqlx.Select(DB, &users, "SELECT * FROM users")
	if err != nil {
		fmt.Println("Error querying the database:", err)
		return []model.User{}, errors.New("Database could not be queried")
	}

	return users, nil
}

func (repo *PostgresUserRepository) AddUser(user model.User) (id int, err error) {
	DB := repo.DB
	err = DB.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id", user.Name).Scan(&id)
	if err != nil {
		fmt.Println("Error inserting user into the database:", err)
		return 0, errors.New("Could not insert user")
	}

	return id, nil
}

func (repo *PostgresUserRepository) GetUser(id int) (user model.User, err error) {
	DB := repo.DB
	usr := []model.User{}
	err = sqlx.Select(DB, &usr, "SELECT * FROM users WHERE id=$1", id)

	fmt.Println(usr)

	if err != nil || len(usr) == 0{
		fmt.Println("Could not retrieve a user with this ID", err)
		return model.User{}, errors.New("No such user exists")
	}
	return usr[0], nil
}

func (repo *PostgresUserRepository) DeleteUser(id int) error {
	DB := repo.DB
	_, err := DB.Query("DELETE FROM users WHERE id=$1", id)

	if err != nil {
		fmt.Println("Could not deleted this user.", err)
		return errors.New("No such user exists")
	}

	return nil
}

func (repo *PostgresUserRepository) UpdateUser(id int, user model.User) error {
	DB := repo.DB
	_, err := DB.Query("UPDATE users SET name=$1 WHERE id=$2", user.Name, id)

	if err != nil {
		fmt.Println("Could not update user.", err)
		return errors.New("No such user exists.")
	}

	return nil
}

func (repo *PostgresUserRepository) Close() {

}
