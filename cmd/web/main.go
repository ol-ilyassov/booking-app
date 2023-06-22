package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ol-ilyassov/booking-app/internal/config"
	"github.com/ol-ilyassov/booking-app/internal/handlers"
	"github.com/ol-ilyassov/booking-app/internal/render"
)

// Not good/bad approach to use global variables.
const portNumber string = ":8081"

var app config.AppConfig

func main() {
	// Application working mode (development|production).
	app.InProduction = false

	// Session management.
	app.Session = scs.New()
	app.Session.Lifetime = 24 * time.Hour
	app.Session.Cookie.Persist = true // on browser tab close, saves cookie.
	app.Session.Cookie.SameSite = http.SameSiteLaxMode
	app.Session.Cookie.Secure = false // app.InProduction

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc

	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println("Starting application on port", portNumber)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
