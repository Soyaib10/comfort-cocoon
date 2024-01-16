package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // line registers a handler function for the root ("/") route
		n, err := fmt.Fprintf(w, "Hello world!")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(fmt.Sprintf("Number of bytes written: %v", n))
	})

	_ = http.ListenAndServe(":8080", nil) // starts the HTTP server on port 8080. If there's an error, it is printed to the console.
}