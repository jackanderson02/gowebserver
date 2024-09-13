package main

import(
    "acme/api"
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

	// Use mutliplexer to allow different methods for same PATH but different
	// method VERB.
	router := http.NewServeMux()

	router.HandleFunc("GET /", rootHandler)
	router.HandleFunc("GET /api/users/{id}", api.GetSingleUser)
	router.HandleFunc("GET /api/users", api.GetUsers)
	router.HandleFunc("DELETE /api/users/{id}", api.DeleteUser)
	router.HandleFunc("PUT /api/users/{id}", api.UpdateUser)
	router.HandleFunc("POST /api/users", api.CreateUser)

	// Starting the HTTP server on port 8080
	fmt.Println("Server listening on port 8080...")
	err := http.ListenAndServe(":8080", CorsMiddleware(router))
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
