package main

import (
	"acme/api"
	"acme/config"
	"acme/service"
	"fmt"
	"net/http"
)

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello, World!")
}
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		// Continue with the next handler
		next.ServeHTTP(writer, request)
	})
}
func main() {

	// Initialize the database connection
	config := config.LoadDatabaseConfig(".env")
	dbRepo, err := initializeDatabase(config)
	if err != nil {
		fmt.Println("Error initializing the database:", err)
		return
	}
	defer dbRepo.Close()

	userService := service.NewUserService(dbRepo)
	userAPI := api.NewUserAPI(userService)

	// Use mutliplexer to allow different methods for same PATH but different
	// method VERB.
	router := http.NewServeMux()

	router.HandleFunc("GET /", rootHandler)
	router.HandleFunc("GET /api/users/{id}", userAPI.GetSingleUser)
	router.HandleFunc("GET /api/users", userAPI.GetUsers)
	router.HandleFunc("DELETE /api/users/{id}", userAPI.DeleteUser)
	router.HandleFunc("PUT /api/users/{id}", userAPI.UpdateUser)
	router.HandleFunc("POST /api/users", userAPI.CreateUser)

	// Starting the HTTP server on port 8080
	fmt.Println("Server listening on port 8080...")
	err = http.ListenAndServe(":8080", CorsMiddleware(router))
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
