package main

import (
	"fmt"
	"net/http"
)

const portNumber = ":8080"



func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Println(fmt.Sprintf("Starting application on port %v", portNumber))
	_ = http.ListenAndServe(portNumber, nil) // starts the HTTP server on port 8080. If there's an error, it is printed to the console.
}
