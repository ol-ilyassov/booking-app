package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ol-ilyassov/booking-app/internal/config"
	"github.com/ol-ilyassov/booking-app/internal/helpers"
	"github.com/ol-ilyassov/booking-app/internal/models"
	"github.com/ol-ilyassov/booking-app/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{
	"humanDate":  render.HumanDate,
	"formatDate": render.FormatDate,
	"iterate":    render.Iterate,
	"add":        render.Add,
}

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	app.InProduction = false

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false
	app.Session = session

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(app.MailChan)

	listenForMail()

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = true

	repo := NewTestingRepo(&app)
	NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	os.Exit(m.Run())
}

func listenForMail() {
	go func() {
		for {
			_ = <-app.MailChan
		}
	}()

}

func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.AboutUs)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.SearchAvailability)
	mux.Post("/search-availability", Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", Repo.PostSearchAvailabilityJSON)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/user/login", Repo.ShowLogin)
	mux.Post("/user/login", Repo.PostShowLogin)
	mux.Get("/user/logout", Repo.Logout)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		// mux.Use(Auth)
		mux.Get("/dashboard", Repo.ShowAdminDashboard)

		mux.Get("/reservations-new", Repo.ShowAdminNewReservations)
		mux.Get("/reservations-all", Repo.ShowAdminAllReservations)

		mux.Get("/reservations-calendar", Repo.ShowAdminCalendarReservations)
		mux.Post("/reservations-calendar", Repo.PostAdminCalendarReservations)

		mux.Get("/process-reservation/{src}/{id}/do", Repo.AdminProcessReservation)
		mux.Get("/delete-reservation/{src}/{id}/do", Repo.AdminDeleteReservation)

		mux.Get("/reservations/{src}/{id}/show", Repo.ShowAdminReservation)
		mux.Post("/reservations/{src}/{id}", Repo.PostAdminReservation)
	})

	return mux
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return cache, err
	}

	matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}
