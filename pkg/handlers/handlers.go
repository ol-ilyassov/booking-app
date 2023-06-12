package handlers

import (
	"fmt"
	"net/http"

	"github.com/ol-ilyassov/booking-app/pkg/render"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.tmpl")
}

func AboutUs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About page!")
}
