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
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %v", portNumber))
	
	srv := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
