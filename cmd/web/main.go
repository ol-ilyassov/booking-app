package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ol-ilyassov/booking-app/pkg/config"
	"github.com/ol-ilyassov/booking-app/pkg/handlers"
	"github.com/ol-ilyassov/booking-app/pkg/render"
)

const portNumber string = ":8081"

func main() {
	var app config.AppConfig

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
	// _ = http.ListenAndServe(portNumber, nil)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
