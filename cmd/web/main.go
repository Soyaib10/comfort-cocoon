package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Soyaib10/comfort-cocoon/internal/config"
	"github.com/Soyaib10/comfort-cocoon/internal/driver"
	"github.com/Soyaib10/comfort-cocoon/internal/handlers"
	"github.com/Soyaib10/comfort-cocoon/internal/helpers"
	"github.com/Soyaib10/comfort-cocoon/internal/models"
	"github.com/Soyaib10/comfort-cocoon/internal/render"
	"github.com/Soyaib10/comfort-cocoon/repoisitory"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	_ = db
	fmt.Println(fmt.Sprintf("Starting application on port %v", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	dsn := fmt.Sprintf("root:@tcp(localhost:3306)/cocoon")
	db, err := repoisitory.ConnectionDB(dsn)

	// putting things in season
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("root:@tcp(localhost:3306)/cocoon")

	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to Database!")

	tc, err := render.CreateTemplateCache()

	// if err != nil {
	// 	log.Fatal("can't create template cache")
	// 	return err
	// }
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return nil, db
}
