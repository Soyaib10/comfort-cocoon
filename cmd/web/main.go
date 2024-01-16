package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Soyaib10/comfort-cocoon/pkg/config"
	"github.com/Soyaib10/comfort-cocoon/pkg/handlers"
	"github.com/Soyaib10/comfort-cocoon/pkg/render"
)

const portNumber = ":8080"

// main is the main application function
func main() {
	var app config.AppConfig
	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("can't create template cache")
	}
	app.TemplateCache = tc
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Println(fmt.Sprintf("Starting application on port %v", portNumber))
	_ = http.ListenAndServe(portNumber, nil) 
}
