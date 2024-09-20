package main 

import(
	"acme/config"
	"acme/repository/user"
	"fmt"
	"acme/db"
)

func initializeDatabase(config config.DatabaseConfig) (user.Repository, error) {
	var userRepo user.Repository
    switch config.Type {
    case "postgres":
        connectionString := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=%s", config.User, config.DBName, config.Password, config.Host, config.SSLMode)
		db, err := postgres.PostgresConnection(connectionString)
		userRepo := user.NewPostgresUserRepository(db.DB)
		return userRepo, err

    case "inmemory":
		userRepo = user.NewInMemoryUserRepository()
		return userRepo, nil
    default:
        return nil, fmt.Errorf("unsupported database type: %s", config.Type)
    }
}


