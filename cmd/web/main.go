package main

import (
	"fmt"
	"net/http"

	"github.com/ol-ilyassov/booking-app/pkg/handlers"
)

const portNumber string = ":8081"

func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.AboutUs)

	fmt.Println("Starting application on port", portNumber)
	_ = http.ListenAndServe(portNumber, nil)
}
