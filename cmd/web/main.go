package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ol-ilyassov/booking-app/internal/config"
	"github.com/ol-ilyassov/booking-app/internal/driver"
	"github.com/ol-ilyassov/booking-app/internal/handlers"
	"github.com/ol-ilyassov/booking-app/internal/helpers"
	"github.com/ol-ilyassov/booking-app/internal/models"
	"github.com/ol-ilyassov/booking-app/internal/render"
)

// Further Grow:
// Alex Edwards' session Manager and Redis storage
// to store sessions and do not lose them after re-launch app.

// Not good/bad approach to use global variables.
const portNumber string = ":8081"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)
	listenForMail()

	app.InfoLog.Printf("Starting application on port%v\n", portNumber)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	// Data that could be stored in session.
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// Application working mode (development|production).
	inProduction := flag.Bool("production", false, "Application production mode (true=production | false=development)")
	useCache := flag.Bool("cache", true, "Use caching on templates")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbName := flag.String("dbname", "", "Database name")
	dbUsername := flag.String("dbuser", "", "Database user")
	dbPassword := flag.String("dbpass", "", "Database password")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable|prefer|require)")
	flag.Parse()

	if *dbName == "" || *dbUsername == "" {
		app.InfoLog.Println("missing required flags...")
		os.Exit(1)
	}

	app.InProduction = *inProduction

	app.UseCache = *useCache

	// Session management.
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // on browser tab close, saves cookie.
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false // app.InProduction
	app.Session = session

	// DB connection.
	app.InfoLog.Println("Connection to database ...")

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUsername, *dbPassword, *dbSSL)
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		log.Fatal("Cannot cinnect to database! Dying...")
	}
	app.InfoLog.Println("Database connection succeded.")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, err
}
