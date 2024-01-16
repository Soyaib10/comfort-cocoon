package main

import (
	"fmt"
	"net/http"
	"github.com/Soyaib10/comfort-cocoon/pkg/handlers"
)

const portNumber = ":8080"

// main is the main application function
func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Println(fmt.Sprintf("Starting application on port %v", portNumber))
	_ = http.ListenAndServe(portNumber, nil) // starts the HTTP server on port 8080. If there's an error, it is printed to the console.
}
